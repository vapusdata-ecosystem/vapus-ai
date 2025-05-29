package services

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/databricks/databricks-sql-go/logger"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusai-studio/v1alpha1"
	pkgs "github.com/vapusdata-ecosystem/vapusai/aistudio/pkgs"
	appdrepo "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo"
	aidmstore "github.com/vapusdata-ecosystem/vapusai/core/app/datarepo/aistudio"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	"github.com/vapusdata-ecosystem/vapusai/core/options"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	processes "github.com/vapusdata-ecosystem/vapusai/core/process"
	filetools "github.com/vapusdata-ecosystem/vapusai/core/tools/files"
	types "github.com/vapusdata-ecosystem/vapusai/core/types"
)

type BlobAgent struct {
	*processes.VapusInterfaceBase
	method        string
	uploadRequest *pb.UploadRequest
	uploadStream  pb.UtilityService_UploadStreamServer
	fetchRequest  *pb.DownloadRequest
	fetchResponse *pb.DownloadResponse
	organization  *models.Organization
	uploadResult  *pb.UploadResponse
	dmStore       *aidmstore.AIStudioDMStore
}

func (s *AIStudioServices) NewUtilityAgent(ctx context.Context, uploadRequest *pb.UploadRequest, uploadStream pb.UtilityService_UploadStreamServer, fetchRequest *pb.DownloadRequest) (*BlobAgent, error) {
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.Logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, encryption.ErrInvalidJWTClaims
	}
	organization, err := s.DMStore.GetOrganization(ctx, vapusPlatformClaim[encryption.ClaimOrganizationKey], vapusPlatformClaim)
	if err != nil {
		s.Logger.Error().Err(err).Ctx(ctx).Msg("error while getting organization from datastore")
		return nil, dmerrors.DMError(apperr.ErrOrganization404, err)
	}
	agent := &BlobAgent{
		uploadRequest: uploadRequest,
		uploadStream:  uploadStream,
		fetchRequest:  fetchRequest,
		organization:  organization,
		dmStore:       s.DMStore,
		fetchResponse: &pb.DownloadResponse{},
		uploadResult:  &pb.UploadResponse{},
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			CtxClaim: vapusPlatformClaim,
			InitAt:   dmutils.GetEpochTime(),
		},
	}
	agent.SetAgentId()
	agent.Logger = pkgs.GetSubDMLogger(types.AIPROMPTAGENT.String(), agent.AgentId)
	return agent, nil
}

func (v *BlobAgent) GetUploadedResult() *pb.UploadResponse {
	v.FinishAt = dmutils.GetEpochTime()
	v.FinalLog()
	return v.uploadResult
}

func (v *BlobAgent) GetDownloadResult() *pb.DownloadResponse {
	v.FinishAt = dmutils.GetEpochTime()
	v.FinalLog()
	return v.fetchResponse
}

func (v *BlobAgent) Act(ctx context.Context) error {
	if v.uploadRequest != nil {
		v.uploadResult = &pb.UploadResponse{
			Output: []*pb.UploadResponse_ObjectUploadResult{},
		}
		return v.uploadFile(ctx)
	} else if v.uploadStream != nil {
		return v.uploadFileStream()
	} else if v.fetchRequest != nil {
		return v.downloadFile(ctx)
	}
	return nil
}

func (v *BlobAgent) uploadFile(ctx context.Context) error {
	if v.uploadRequest.Resource == "" {
		v.Logger.Error().Msg("resource is empty")
		return dmerrors.DMError(apperr.ErrMissingUploadResourceName, nil)
	}
	if v.uploadRequest.ResourceId == "" {
		v.uploadRequest.ResourceId = v.CtxClaim[encryption.ClaimUserIdKey]
	}
	for _, fileData := range v.uploadRequest.Objects {
		if fileData.Name == "" && fileData.Data == nil {
			v.Logger.Error().Msg("file data is empty")
			v.uploadResult.Output = append(v.uploadResult.Output, &pb.UploadResponse_ObjectUploadResult{
				Error:  "file data is empty",
				Object: fileData,
			})
			continue
		}
		fType := filetools.GetConfFileType(fileData.Name)
		fileData.Format = mpb.ContentFormats(mpb.ContentFormats_value[strings.ToUpper(fType)])
		checksum := encryption.GenerateSHA3_256FromBytes(fileData.Data, []byte{})
		counter, exists, _ := appdrepo.ValidateFileCache(ctx, v.dmStore.VapusStore, checksum, fileData.Name, filepath.Join(v.uploadRequest.Resource, v.uploadRequest.ResourceId))
		if fileData.Name == "" {
			fileData.Name = dmutils.GetUUID() + "." + strings.ToLower(fileData.Format.String())
		} else {
			fTL := len(fileData.Format.String()) + 1
			if strings.ToLower(fileData.Name[len(fileData.Name)-fTL:]) != "."+strings.ToLower(fileData.Format.String()) {
				fileData.Name = fmt.Sprintf("%s_%d.%s", fileData.Name, counter+1, strings.ToLower(fileData.Format.String()))
			} else {
				fileData.Name = fmt.Sprintf("%d_%s", counter, fileData.Name)
			}
		}
		fileData.Name = filetools.SanitizeFileName(fileData.Name)
		if fileData.Data == nil {
			v.Logger.Error().Msg("file data is empty")
			v.uploadResult.Output = append(v.uploadResult.Output, &pb.UploadResponse_ObjectUploadResult{
				Error:  "file data is empty",
				Object: fileData,
			})
			continue
		}
		keyPath := fmt.Sprintf("%s/%s/%s", v.uploadRequest.Resource, v.uploadRequest.ResourceId, fileData.Name)
		logFileObj := &models.FileStoreLog{
			Name:       fileData.Name,
			Path:       fileData.Path,
			Format:     fileData.Format.String(),
			Size:       int64(len(fileData.Data)),
			Resource:   v.uploadRequest.Resource,
			ResourceId: v.uploadRequest.ResourceId,
			Checksum:   encryption.GenerateSHA3_256FromBytes(fileData.Data, []byte{}),
		}
		if !exists {
			v.Logger.Info().Msgf("uploading file to blob storage with key path: %s", keyPath)
			err := v.dmStore.BlobStore.UploadObject(ctx, &options.BlobOpsParams{
				BucketName: v.CtxClaim[encryption.ClaimOrganizationKey],
				ObjectName: keyPath,
				Data:       fileData.Data,
			})
			if err != nil {
				_ = appdrepo.DeleteRedisKey(ctx, keyPath, v.dmStore.VapusStore, v.Logger) // Deleting the Key from the redis
				v.Logger.Error().Err(err).Msg("error while uploading file to blob storage")
				v.uploadResult.Output = append(v.uploadResult.Output, &pb.UploadResponse_ObjectUploadResult{
					Error:  err.Error(),
					Object: fileData,
				})
				continue
			}
			// Saving the file data in DB
			err = appdrepo.LogFileStoreLog(ctx, v.dmStore.VapusStore, logFileObj, v.CtxClaim)
			if err != nil {
				_ = appdrepo.DeleteRedisKey(ctx, keyPath, v.dmStore.VapusStore, v.Logger)
				v.Logger.Error().Err(err).Msg("error while DB entery")
				return err
			}
		}
		fileData.Data = nil

		v.uploadResult.Output = append(v.uploadResult.Output, &pb.UploadResponse_ObjectUploadResult{
			Object:       fileData,
			ResponsePath: keyPath,
			Fid:          dmutils.GetUUID(),
		})
	}

	return nil
}

func (v *BlobAgent) uploadFileStream() error {
	// upload file stream
	return nil
}

func (v *BlobAgent) downloadFile(ctx context.Context) error {
	if v.fetchRequest.GetPath() == "" {
		v.Logger.Error().Msg("FilePath is empty")
		return dmerrors.DMError(apperr.ErrEmptyFile, nil)
	}
	fileStoreLog, err := appdrepo.GetFile(ctx, v.dmStore.VapusStore, v.fetchRequest.GetPath(), v.CtxClaim)
	if err != nil {
		v.Logger.Err(err).Msgf("Error while fetching the file path")
		return err
	}
	if fileStoreLog.DeletedBy != "" || fileStoreLog.DeletedAt != 0 {
		v.Logger.Err(err).Msgf("Error File is Deleted")
		return err
	}

	data, err := v.dmStore.BlobStore.DownloadObject(ctx, &options.BlobOpsParams{
		BucketName: v.CtxClaim[encryption.ClaimOrganizationKey],
		ObjectName: v.fetchRequest.GetPath(),
	})

	if err != nil {
		logger.Err(err).Msg("error while downloading data file")
		return err
	}

	fmt.Println("Data: ", data[0])

	v.fetchResponse = &pb.DownloadResponse{
		Data:   data,
		Name:   fileStoreLog.Name,
		Format: mpb.ContentFormats(mpb.ContentFormats_value[fileStoreLog.Format]),
	}
	return nil
}
