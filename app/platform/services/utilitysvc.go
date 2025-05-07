package services

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	pb "github.com/vapusdata-ecosystem/apis/protos/vapusdata/v1alpha1"
	appdrepo "github.com/vapusdata-ecosystem/vapusdata/core/app/datarepo"
	apperr "github.com/vapusdata-ecosystem/vapusdata/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusdata/core/models"
	"github.com/vapusdata-ecosystem/vapusdata/core/options"
	encryption "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/encryption"
	dmerrors "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/errors"
	dmutils "github.com/vapusdata-ecosystem/vapusdata/core/pkgs/utils"
	processes "github.com/vapusdata-ecosystem/vapusdata/core/process"
	filetools "github.com/vapusdata-ecosystem/vapusdata/core/tools/files"
	types "github.com/vapusdata-ecosystem/vapusdata/core/types"
	pkgs "github.com/vapusdata-ecosystem/vapusdata/platform/pkgs"
)

type BlobAgent struct {
	*processes.VapusInterfaceBase
	method        string
	uploadRequest *pb.UploadRequest
	uploadStream  pb.UtilityService_UploadStreamServer
	organization  *models.Organization
	uploadResult  *pb.UploadResponse
	*DMServices
}

func (s *DMServices) NewUtilityAgent(ctx context.Context, uploadRequest *pb.UploadRequest, uploadStream pb.UtilityService_UploadStreamServer) (*BlobAgent, error) {
	vapusPlatformClaim, ok := encryption.GetCtxClaim(ctx)
	if !ok {
		s.logger.Error().Ctx(ctx).Msg("error while getting claim metadata from context")
		return nil, encryption.ErrInvalidJWTClaims
	}
	organization, err := s.DMStore.GetOrganization(ctx, vapusPlatformClaim[encryption.ClaimOrganizationKey], vapusPlatformClaim)
	if err != nil {
		s.logger.Error().Err(err).Ctx(ctx).Msg("error while getting organization from datastore")
		return nil, dmerrors.DMError(apperr.ErrOrganization404, err)
	}
	agent := &BlobAgent{
		uploadRequest: uploadRequest,
		uploadStream:  uploadStream,
		DMServices:    s,
		organization:  organization,
		VapusInterfaceBase: &processes.VapusInterfaceBase{
			CtxClaim: vapusPlatformClaim,
			// Ctx:      ctx,
			InitAt: dmutils.GetEpochTime(),
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

func (v *BlobAgent) Act(ctx context.Context) error {
	if v.uploadRequest != nil {
		v.uploadResult = &pb.UploadResponse{
			Output: []*pb.UploadResponse_ObjectUploadResult{},
		}
		return v.uploadFile(ctx)
	} else if v.uploadStream != nil {
		return v.uploadFileStream()
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
		counter, exists, _ := appdrepo.ValidateFileCache(ctx, v.DMStore.VapusStore, checksum, fileData.Name, filepath.Join(v.uploadRequest.Resource, v.uploadRequest.ResourceId))

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
			v.logger.Info().Msgf("uploading file to blob storage with key path: %s", keyPath)
			err := v.DMStore.BlobStore.UploadObject(ctx, &options.BlobOpsParams{
				BucketName: v.CtxClaim[encryption.ClaimOrganizationKey],
				ObjectName: keyPath,
				Data:       fileData.Data,
			})
			if err != nil {
				v.Logger.Error().Err(err).Msg("error while uploading file to blob storage")
				v.uploadResult.Output = append(v.uploadResult.Output, &pb.UploadResponse_ObjectUploadResult{
					Error:  err.Error(),
					Object: fileData,
				})
				continue
			}
			_ = appdrepo.LogFileStoreLog(ctx, v.DMStore.VapusStore, logFileObj, v.CtxClaim)
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
