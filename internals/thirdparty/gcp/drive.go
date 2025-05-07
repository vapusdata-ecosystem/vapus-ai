package gcp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"sync"

	"github.com/rs/zerolog"
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	options "github.com/vapusdata-ecosystem/vapusai/core/options"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

type GDriveOpts struct {
	Client *drive.Service
	logger zerolog.Logger
}

func NewImpersonatorGDrive(ctx context.Context, opts *GcpConfig, userEmailAddr string, logger zerolog.Logger) (*GDriveOpts, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(string(opts.ServiceAccountKey))
	if err != nil {
		logger.Err(err).Msgf("Error while decoding the GCP KEY -- %v", err)
		return nil, err
	}
	creds, err := google.CredentialsFromJSON(ctx, decodedKey)
	if err != nil || creds == nil {
		logger.Err(err).Msgf("Error while creating credentials from json for GCP drive plugin-- %v", err)
		return nil, err
	}
	keyJson := map[string]any{}
	err = json.Unmarshal(creds.JSON, &keyJson)
	if err != nil {
		logger.Err(err).Msgf("Error while unmarshalling the GCP KEY json -- %v", err)
		return nil, err
	}
	clEmail, ok := keyJson["client_email"].(string)
	if !ok {
		logger.Err(err).Msgf("Error while getting the client_email from the GCP KEY json -- %v", err)
		return nil, err
	}
	log.Println("gcp-SVC-KEY-EMAIL", string(decodedKey))
	log.Println("gcp-SVC-KEY-EMAIL", clEmail)
	log.Println("gcp-SVC-KEY-EMAIL", userEmailAddr)
	log.Println("gcp-SVC-KEY-TokenSource", creds.TokenSource)

	// Now impersonate the ORGANIZATION user
	tokenSource, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
		TargetPrincipal: clEmail,
		Subject:         userEmailAddr,
		Scopes:          []string{"https://www.googleapis.com/auth/drive"},
	}, option.WithCredentialsJSON(decodedKey))
	if err != nil {
		logger.Err(err).Msgf("Error while impersonating the user -- %v", err)
		return nil, err
	}
	client, err := drive.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		logger.Err(err).Msgf("Error while creating credentials from json -- %v", err)
		return nil, err
	}
	return &GDriveOpts{
		Client: client,
		logger: logger,
	}, nil
}

func NewGDrive(ctx context.Context, opts *GcpConfig, userEmailAddr string, logger zerolog.Logger) (*GDriveOpts, error) {
	log.Println("GCP-DRIVE-KEY", string(opts.ServiceAccountKey))
	decodedKey, err := base64.StdEncoding.DecodeString(string(opts.ServiceAccountKey))
	if err != nil {
		logger.Err(err).Msgf("Error while decoding the GCP KEY -- %v", err)
		return nil, err
	}
	log.Println("gcp-SVC-KEY", string(decodedKey))
	// creds, err := google.CredentialsFromJSON(ctx, decodedKey)
	// if err != nil || creds == nil {
	// 	logger.Err(err).Msgf("Error while creating credentials from json for GCP drive plugin-- %v", err)
	// 	return nil, err
	// }
	// keyJson := map[string]any{}
	// err = json.Unmarshal(creds.JSON, &keyJson)
	// if err != nil {
	// 	logger.Err(err).Msgf("Error while unmarshalling the GCP KEY json -- %v", err)
	// 	return nil, err
	// }
	// clEmail, ok := keyJson["client_email"].(string)
	// if !ok {
	// 	logger.Err(err).Msgf("Error while getting the client_email from the GCP KEY json -- %v", err)
	// 	return nil, err
	// }
	// log.Println("gcp-SVC-KEY-EMAIL", string(decodedKey))
	// log.Println("gcp-SVC-KEY-EMAIL", clEmail)
	// log.Println("gcp-SVC-KEY-EMAIL", userEmailAddr)
	// log.Println("gcp-SVC-KEY-TokenSource", creds.TokenSource)

	// // Now impersonate the ORGANIZATION user
	// tokenSource, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
	// 	TargetPrincipal: clEmail,
	// 	Subject:         userEmailAddr,
	// 	Scopes:          []string{"https://www.googleapis.com/auth/drive"},
	// }, option.WithCredentialsJSON(decodedKey))
	// if err != nil {
	// 	logger.Err(err).Msgf("Error while impersonating the user -- %v", err)
	// 	return nil, err
	// }
	client, err := drive.NewService(ctx, option.WithCredentialsJSON(decodedKey))
	if err != nil {
		logger.Err(err).Msgf("Error while creating credentials from json -- %v", err)
		return nil, err
	}
	return &GDriveOpts{
		Client: client,
		logger: logger,
	}, nil
}

func (g *GDriveOpts) UploadFiles(ctx context.Context, opts *options.FileStoreUploadRequest) (*options.FileStoreUploadResponse, error) {
	var err error
	log.Println("gcp-DRIVE-UPLOAD", opts)
	response := &options.UploadResponse{}
	reader := bytes.NewReader(opts.Data)
	fileMetadata := &drive.File{
		Name: opts.FileName,
		// OwnedByMe: opts.OwnedByMe,
	}
	if opts.DirectoryId != "" {
		fileMetadata.Parents = []string{opts.DirectoryId}
	}
	res, err := g.Client.Files.Create(fileMetadata).Media(reader).Do()
	if err != nil {
		g.logger.Err(err).Msgf("Error while uploading the file -- %v", err)
		return nil, err
	}
	response.FileId = res.Id
	log.Println("File ID =========================>>>>>>>>>>>>>>>>>>>>>", res.Id)
	if opts.SharingParams != nil {
		shareRequest := &options.FileStoreShareRequest{
			FileId:        res.Id,
			SharingParams: opts.SharingParams,
		}
		// if opts.AllowSearch {
		// 	shareRequest.AllowSearch = opts.AllowSearch
		// }
		shareResponse, err := g.Share(ctx, shareRequest)
		if err != nil {
			g.logger.Err(err).Msgf("Error while sharing the file -- %v", err)
			return nil, err
		} else {
			response.PermissionId = shareResponse.PermissionId
		}
	}
	return &options.FileStoreUploadResponse{
		Results: []*options.UploadResponse{response},
	}, nil
}

func (g *GDriveOpts) ListFiles(ctx context.Context, opts *options.FileStoreListRequest) (*options.FileStoreListResponse, error) {
	var err error
	response := &options.FileStoreListResponse{}
	resp, err := g.Client.Files.List().Do()
	if err != nil {
		g.logger.Err(err).Msgf("Error while listing the files -- %v", err)
		return nil, err
	}
	for _, file := range resp.Files {
		log.Println("gcp-DRIVE-FILE", file.Name)
	}
	return response, nil
}

func (g *GDriveOpts) DownloadFiles(ctx context.Context, request *options.FileStoreDownloadRequest) (*options.FileStoreDownloadResponse, error) {
	var err error
	response := &options.FileStoreDownloadResponse{
		Files: []*mpb.FileData{},
	}
	respFiles, err := g.getFiles(ctx, request.FileIds)
	if err != nil {
		g.logger.Err(err).Msgf("Error while getting the files -- %v", err)
		return nil, err
	}
	response.Files = respFiles

	return response, nil
}

func (g *GDriveOpts) DeleteFiles(ctx context.Context, opts *options.FileStoreDeleteRequest) error {
	var err error
	err = g.Client.Files.Delete(opts.FileId).Do()
	if err != nil {
		g.logger.Err(err).Msgf("Error while deleting the file -- %v", err)
		return err
	}
	return nil
}

func (g *GDriveOpts) Share(ctx context.Context, opts *options.FileStoreShareRequest) (*options.FileStoreShareResponse, error) {
	var err error
	var resId string
	response := &options.FileStoreShareResponse{}
	if opts.DirectoryId != "" {
		resId = opts.DirectoryId
		response.DirectoryId = opts.DirectoryId
	} else if opts.FileId != "" {
		resId = opts.FileId
		response.FileId = opts.FileId
	} else {
		g.logger.Err(err).Msgf("Error while deleting the folder -- %v", err)
		return nil, apperr.ErrInvalidFileSharingResource
	}
	for _, shareRequest := range opts.SharingParams {
		switch shareRequest.Type {
		case options.FileStoreSharingTypeUser:
			for _, email := range shareRequest.Emails {
				res, err := g.Client.Permissions.Create(resId, &drive.Permission{
					Type:         shareRequest.Type.String(),
					Role:         shareRequest.Role,
					EmailAddress: email,
					// AllowFileDiscovery: opts.AllowSearch,
				}).Do()
				if err != nil {
					g.logger.Err(err).Msgf("Error while sharing the file -- %v", err)
				} else {
					response.PermissionId = res.Id
				}
			}
		case options.FileStoreSharingTypeOrganization:
			res, err := g.Client.Permissions.Create(resId, &drive.Permission{
				Type:               shareRequest.Type.String(),
				Role:               shareRequest.Role,
				Domain:             shareRequest.Organization,
				AllowFileDiscovery: opts.AllowSearch,
			}).Do()
			if err != nil {
				g.logger.Err(err).Msgf("Error while sharing the file -- %v", err)
			} else {
				response.PermissionId = res.Id
			}
		case options.FileStoreSharingTypeGroup:
			res, err := g.Client.Permissions.Create(resId, &drive.Permission{
				Type:               shareRequest.Type.String(),
				Role:               shareRequest.Role,
				AllowFileDiscovery: opts.AllowSearch,
				// Group: shareRequest.Group, // TO:DO Revisit this parameters
			}).Do()
			if err != nil {
				g.logger.Err(err).Msgf("Error while sharing the file -- %v", err)
			} else {
				response.PermissionId = res.Id
			}
		default:
			g.logger.Err(err).Msgf("Error while sharing the file -- %v", err)
			return response, apperr.ErrInvalidFileSharingType
		}
	}
	return response, nil
}

func (g *GDriveOpts) UnShare(ctx context.Context, opts *options.FileStoreShareRequest) (*options.FileStoreShareResponse, error) {
	var err error
	var resId string
	response := &options.FileStoreShareResponse{}
	if opts.DirectoryId != "" {
		resId = opts.DirectoryId
		response.DirectoryId = opts.DirectoryId
	} else if opts.FileId != "" {
		resId = opts.FileId
		response.FileId = opts.FileId
	} else {
		g.logger.Err(err).Msgf("Error while deleting the folder -- %v", err)
		return nil, apperr.ErrInvalidFileSharingResource
	}
	for _, shareRequest := range opts.SharingParams {
		switch shareRequest.Type {
		case options.FileStoreSharingTypeUser:
			for _, email := range shareRequest.Emails {
				res, err := g.Client.Permissions.Update(resId, opts.PermissionId, &drive.Permission{
					Type:         shareRequest.Type.String(),
					Role:         shareRequest.Role,
					EmailAddress: email,
				}).Do()
				if err != nil {
					g.logger.Err(err).Msgf("Error while sharing the file -- %v", err)
				} else {
					response.PermissionId = res.Id
				}
			}
		case options.FileStoreSharingTypeOrganization:
			res, err := g.Client.Permissions.Update(resId, opts.PermissionId, &drive.Permission{
				Type:   shareRequest.Type.String(),
				Role:   shareRequest.Role,
				Domain: shareRequest.Organization,
			}).Do()
			if err != nil {
				g.logger.Err(err).Msgf("Error while sharing the file -- %v", err)
			} else {
				response.PermissionId = res.Id
			}
		case options.FileStoreSharingTypeGroup:
			res, err := g.Client.Permissions.Update(resId, opts.PermissionId, &drive.Permission{
				Type: shareRequest.Type.String(),
				Role: shareRequest.Role,
				// Group: shareRequest.Group, // TO:DO Revisit this parameters
			}).Do()
			if err != nil {
				g.logger.Err(err).Msgf("Error while sharing the file -- %v", err)
			} else {
				response.PermissionId = res.Id
			}
		default:
			g.logger.Err(err).Msgf("Error while sharing the file -- %v", err)
			return response, apperr.ErrInvalidFileSharingType
		}
	}
	return response, nil
}

func (g *GDriveOpts) CreateFolder(ctx context.Context, opts *options.FileStoreCreateDirectoryRequest) (*options.FileStoreCreateDirectoryResponse, error) {
	var err error
	folderMetadata := &drive.File{
		Name:     opts.DirectoryName,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{opts.ParentId},
	}
	res, err := g.Client.Files.Create(folderMetadata).Do()
	if err != nil {
		g.logger.Err(err).Msgf("Error while creating the folder -- %v", err)
		return nil, err
	}

	return &options.FileStoreCreateDirectoryResponse{
		DirectoryId: res.Id,
	}, nil
}
func (g *GDriveOpts) MoveFiles(ctx context.Context, opts *options.FileStoreMoveFilesRequest) error {
	var err error
	response := &options.FileStoreMoveFilesResponse{
		FailedFileIds:   map[string][]string{},
		FailedFileNames: map[string][]string{},
	}
	for _, obj := range opts.Files {
		_, err = g.Client.Files.Update(obj, &drive.File{
			Id:      obj,
			Parents: []string{opts.ExistingDirectoryId},
		}).AddParents(opts.NewDirectoryId).Do()
		if err != nil {
			g.logger.Err(err).Msgf("Error while moving the file -- %v", err)
			exists, ok := response.FailedFileIds[obj]
			if !ok {
				response.FailedFileIds[obj] = []string{}
			} else {
				response.FailedFileIds[obj] = append(exists, obj)
			}
		}
	}
	for _, obj := range opts.Filenames {
		_, err = g.Client.Files.Update(obj, &drive.File{
			Name:    obj,
			Parents: []string{opts.ExistingDirectoryId},
		}).AddParents(opts.NewDirectoryId).Do()
		if err != nil {
			g.logger.Err(err).Msgf("Error while moving the file -- %v", err)
			exists, ok := response.FailedFileNames[obj]
			if !ok {
				response.FailedFileNames[obj] = []string{}
			} else {
				response.FailedFileNames[obj] = append(exists, obj)
			}
		}
	}
	return nil
}

func (g *GDriveOpts) DownloadFolder(ctx context.Context, request *options.FileStoreDownloadRequest) (*options.FileStoreDownloadResponse, error) {
	var err error
	response := &options.FileStoreDownloadResponse{
		Files: []*mpb.FileData{},
	}
	respFiles, err := g.getFiles(ctx, request.DirectoryIds)
	if err != nil {
		g.logger.Err(err).Msgf("Error while getting the files -- %v", err)
		return nil, err
	}
	response.Files = respFiles

	return response, nil
}

func (g *GDriveOpts) ShareFileLink(ctx context.Context) error {

	return nil
}

func (g *GDriveOpts) DeleteFolder(ctx context.Context, opts *options.FileStoreDeleteRequest) error {
	var err error
	var resId string
	if opts.DirectoryId != "" {
		resId = opts.DirectoryId
	} else if opts.FileId != "" {
		resId = opts.FileId
	} else {
		g.logger.Err(err).Msgf("Error while deleting the folder -- %v", err)
		resId = opts.FileName
	}
	err = g.Client.Files.Delete(resId).Do()
	if err != nil {
		g.logger.Err(err).Msgf("Error while deleting the folder -- %v", err)
		return err
	}
	return nil
}

func (g *GDriveOpts) convertToFile(ctx context.Context, resp *drive.File, fileId string) (*mpb.FileData, error) {
	var err error
	file := &mpb.FileData{}

	res, err := g.Client.Files.Get(fileId).Download()
	if err != nil {
		g.logger.Err(err).Msgf("Error while downloading the file -- %v", err)
		return nil, err
	}
	defer res.Body.Close()
	fileData, err := io.ReadAll(res.Body)
	if err != nil {
		g.logger.Err(err).Msgf("Error while reading file data -- %v", err)
		return nil, err
	}
	file.Data = fileData
	file.Name = resp.Name
	paramBytes, err := resp.MarshalJSON()
	if err != nil {
		g.logger.Err(err).Msgf("Error while marshalling the file Request -- %v", err)
	}
	if err == nil {
		mapObj := map[string]any{}
		err = json.Unmarshal(paramBytes, &mapObj)
		if err != nil {
			g.logger.Err(err).Msgf("Error while unmarshalling the file Request -- %v", err)
		} else {
			structpbObj, err := structpb.NewStruct(mapObj)
			if err != nil {
				g.logger.Err(err).Msgf("Error while creating the structpb object -- %v", err)
			} else {
				file.Params = structpbObj
			}
		}
	}
	return file, nil
}

func (g *GDriveOpts) getFile(ctx context.Context, fileId string) (*mpb.FileData, error) {
	var err error
	resp, err := g.Client.Files.Get(fileId).Do()

	if err != nil {
		g.logger.Err(err).Msgf("Error while getting the file -- %v", err)
		return nil, err
	}
	return g.convertToFile(ctx, resp, fileId)
}

func (g *GDriveOpts) getFiles(ctx context.Context, fileIds []string) ([]*mpb.FileData, error) {
	var wg sync.WaitGroup
	var errChan = make(chan error, len(fileIds))
	result := make([]*mpb.FileData, len(fileIds))
	for i, fileId := range fileIds {
		wg.Add(1)
		go func(i int, fileId string) {
			defer wg.Done()
			file, err := g.getFile(ctx, fileId)
			if err != nil {
				errChan <- err
			} else {
				result[i] = file
			}
		}(i, fileId)
	}
	wg.Wait()
	close(errChan)
	if err := <-errChan; err != nil {
		g.logger.Err(err).Msgf("Error while getting the files -- %v", err)
		return result, err
	}
	return result, nil
}
