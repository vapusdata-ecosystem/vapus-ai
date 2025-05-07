package plugins

import "github.com/vapusdata-ecosystem/vapusai/core/models"

type EmailOpts struct {
	From             string
	SenderName       string
	To               []string
	Subject          string
	Body             string
	Attachments      []*Attachment
	HtmlTemplateBody string
	Tags             []*models.Mapper
	RawMessage       []byte
}

type Attachment struct {
	FileName    string
	Data        []byte
	Format      string
	ContentType string
}

type GoogleDriveUploadOpts struct {
	FileName    string
	Data        []byte
	Path        string
	ContentType string
}
