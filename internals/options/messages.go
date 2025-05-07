package options

import (
	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
)

type MessageWithAttachment struct {
	Color   string                    `json:"color"`
	Text    string                    `json:"text"`
	Title   string                    `json:"title"`
	Actions []MessageAttachmentAction `json:"actions"`
}

type MessageAttachmentAction struct {
	Name string `json:"name"`
	Text string `json:"text"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

type MessageUploadAttachment struct {
	Files   []*mpb.FileData `json:"files"`
	Channel string          `json:"channel"`
	Comment string          `json:"comment"`
	Title   string          `json:"title"`
}
