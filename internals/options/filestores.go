package options

import (
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
)

type FileStoreUploadRequest struct {
	// Optional field, schema generator will create an array schema for this.
	// Added omitempty as slices are often optional.
	SharingParams []*FileStoreAssetSharingParams `json:"sharingParams,omitempty" jsonschema:"description=Optional sharing parameters for the asset."`

	// Required field based on the target schema's "required" array.
	FileName string `json:"fileName" jsonschema:"description=The name of the file to be uploaded. This is a required field.,required"`

	// Required field based on the target schema's "required" array.
	DirectoryId string `json:"directoryId" jsonschema:"description=The ID of the directory where the file will be uploaded. This is a required field.,required"`

	// NOTE: Not marked 'required' here because it wasn't in the target schema's explicit "required" array.
	// []byte usually maps to string with contentEncoding: base64
	// Added omitempty as it's effectively optional per the target schema's required list.
	Data []byte `json:"data,omitempty" jsonschema:"description=The base64 encoded data of the file to be uploaded."`

	// Required field with enum constraint based on the target schema.
	Format string `json:"format" jsonschema:"description=The format of the file to be uploaded. This is a required field.,required,enum=json,enum=csv,enum=xlsx,enum=pdf,enum=png,enum=jpeg,enum=jpg,enum=gif,enum=bmp,enum=tiff,enum=svg,enum=webp,enum=ico,enum=raw,enum=psd"`

	// NOTE: Not marked 'required' here because it wasn't in the target schema's explicit "required" array.
	// Added default value. Added omitempty assuming false (the zero value) doesn't need to be sent if default is true.
	OwnedByMe bool `json:"ownedByMe,omitempty" jsonschema:"description=Indicates whether the file is owned by the user. Should typically be true.,default=true"`

	// Required field based on the target schema's "required" array. Added default value.
	AllowSearch bool `json:"allowSearch" jsonschema:"description=Indicates whether the file should be searchable. Should typically be true.,required,default=true"`
}

type FileStoreUploadResponse struct {
	Results []*UploadResponse `json:"results"`
}

type UploadResponse struct {
	FileId       string `json:"fileId"`
	FileName     string `json:"fileName"`
	DirectoryId  string `json:"directoryId"`
	PermissionId string `json:"permissionId"`
	Error        string `json:"error"`
}

type FileStoreDownloadRequest struct {
	// Field name mismatch: Schema uses 'fileIds', struct uses 'FileIds'. Assuming map. Marked required.
	FileIds []string `json:"fileIds,omitempty" jsonschema:"description=The IDs of the files to be downloaded. This is a required field and can be ampty array if the user is not mentioning the file id.,required"`
	// Field name mismatch: Schema uses 'directoryIds', struct uses 'DirectoryIds'. Assuming map. Marked required.
	DirectoryIds []string `json:"directoryIds,omitempty" jsonschema:"description=The IDs of the directories where the files are located. This is a required field and can be ampty array if the user is not mentioning the directory id.,required"`
	// Marked required based on schema.
	FileNames []string `json:"fileNames,omitempty" jsonschema:"description=The names of the files to be downloaded. This is a required field and can be ampty array if the user is not mentioning the file name.,required"`
	// Marked required based on schema.
	DirectoryNames []string `json:"directoryNames,omitempty" jsonschema:"description=The names of the directories where the files are located. This is a required field and can be ampty array if the user is not mentioning the directory name.,required"`
}

type FileStoreDownloadResponse struct {
	Files []*mpb.FileData `json:"files"`
}

type FileStoreDeleteRequest struct {
	// Marked required based on schema.
	FileName string `json:"fileName" jsonschema:"description=The name of the file to be deleted. This is a required field and can be empty string if the user is not mentioning the file name.,required"`
	// Marked required based on schema.
	DirectoryId string `json:"directoryId" jsonschema:"description=The ID of the directory where the file is located. This is a required field and can be empty string if the user is not mentioning the directory ID.,required"`
	// Marked required based on schema.
	FileId string `json:"fileId" jsonschema:"description=The ID of the file to be deleted. This is a required field and can be empty string if the user is not mentioning the file ID.,required"`
}

type FileStoreListRequest struct {
	// JSON name 'patterns' vs schema 'filePatterns'. Assuming map. Marked required.
	FilePatterns []string `json:"filePatterns,omitempty" jsonschema:"description=The patterns of the files to be listed. This is a required field and can be empty array if the user is not mentioning the file patterns.,required"`
	// Marked required based on schema.
	DirectoryPatterns []string `json:"directoryPatterns,omitempty" jsonschema:"description=The patterns of the directories to be listed. This is a required field and can be empty array if the user is not mentioning the directory patterns.,required"`
	// Marked required based on schema.
	Owners []string `json:"owners,omitempty" jsonschema:"description=The owners of the files to be listed. This is a required field and can be empty array if the user is not mentioning the owners.,required"`
	// Marked required based on schema.
	OwnedByMe bool `json:"ownedByMe" jsonschema:"description=Indicates whether the files should be owned by the user. This is a required field and should be set to true.,required"`
	// Marked required based on schema.
	SharedWithMe bool `json:"sharedWithMe" jsonschema:"description=Indicates whether the files should be shared with the user. This is a required field and should be set to true.,required"`
}

type FileStoreListResponse struct {
	FileIds []string `json:"fileIds"`
}

type FileStoreShareRequest struct {
	// Using FileStoreAssetSharingParams which is already tagged. Marked required based on schema.
	SharingParams []*FileStoreAssetSharingParams `json:"sharingParams,omitempty" jsonschema:"description=The parameters for the sharing request. This is a required field.,required"`
	// Marked required based on schema.
	FileId string `json:"fileId" jsonschema:"description=The ID of the file to be shared. This is a required field.,required"`
	// Marked required based on schema.
	DirectoryId string `json:"directoryId" jsonschema:"description=The ID of the directory where the file is located. This is a required field.,required"`
	// Marked required based on schema.
	PermissionId string `json:"permissionId" jsonschema:"description=The ID of the permission to be applied. This is a required field.,required"`
	// Marked required based on schema. Consider adding default=true if appropriate.
	AllowSearch bool `json:"allowSearch" jsonschema:"description=Indicates whether the file should be searchable. This is a required field and should be set to true.,required"`
}

type FileStoreShareResponse struct {
	PermissionId string `json:"permissionId"`
	FileId       string `json:"fileId"`
	DirectoryId  string `json:"directoryId"`
}

type FileStoreMoveFilesRequest struct {
	Files               []string `json:"files"`
	ExistingDirectoryId string   `json:"existingDirectoryId"`
	NewDirectoryId      string   `json:"newDirectoryId"`
	Filenames           []string `json:"filenames"`
}

type FileStoreMoveFilesResponse struct {
	FailedFileIds   map[string][]string `json:"failedFileIds"`
	FailedFileNames map[string][]string `json:"failedFileNames"`
}

type FileStoreCreateDirectoryRequest struct {
	DirectoryName string `json:"directoryName"`
	ParentId      string `json:"parentId"`
}
type FileStoreCreateDirectoryResponse struct {
	DirectoryId string `json:"directoryId"`
}

type FileStoreSharingType string

func (f FileStoreSharingType) String() string {
	return string(f)
}

const (
	FileStoreSharingTypeUser         FileStoreSharingType = "user"
	FileStoreSharingTypeOrganization FileStoreSharingType = "Organization"
	FileStoreSharingTypeGroup        FileStoreSharingType = "group"
)

type FileStoreAssetSharingParams struct {
	// Schema used 'emailAddress' but struct uses 'Emails'. Assuming mapping. Marked required.
	Emails []string `json:"emailAddress,omitempty" jsonschema:"description=The email address of the user to share the file with. This is a required field.,required"`
	// Marked required and added enum based on schema. Using string type directly for simplicity with tags.
	Role string `json:"role,omitempty" jsonschema:"description=The role of the user to share the file with. This is a required field.,required,enum=owner,enum=organizer,enum=fileOrganizer,enum=writer,enum=commenter,enum=reader"`
	// Marked required based on schema.
	Organization string `json:"Organization,omitempty" jsonschema:"description=The Organization of the user to share the file with. This is a required field.,required"`
	// Marked required based on schema.
	Group string `json:"group,omitempty" jsonschema:"description=The group of the user to share the file with. This is a required field.,required"`
	// Marked required and added enum based on schema. Using string type directly for simplicity with tags.
	Type FileStoreSharingType `json:"type,omitempty" jsonschema:"description=The type of the user to share the file with. This is a required field.,required,enum=user,enum=group,enum=ORGANIZATION,enum=anyone"`
}
