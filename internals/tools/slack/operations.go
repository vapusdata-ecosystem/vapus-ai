package slack

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	options "github.com/vapusdata-ecosystem/vapusdata/core/options"
)

type SlackMessageOpts struct {
	ChannelID   string
	UserID      string
	Message     string
	ThreadTs    string
	Channel     *slack.Channel
	FileId      string
	Files       *options.MessageUploadAttachment
	Attachments []*options.MessageWithAttachment
	Logger      zerolog.Logger
}

func (s *Slack) SendMessageToChannel(params *SlackMessageOpts) error {
	_, _, err := s.client.PostMessage(
		params.ChannelID,
		slack.MsgOptionText(params.Message, false),
	)
	return err
}

func (s *Slack) SendMessageToUser(params *SlackMessageOpts) error {
	if params.Channel == nil {
		channel, err := s.NewDMChannel(params)
		if err != nil {
			params.Logger.Error().Err(err).Msg("failed to create DM channel")
			return err
		}
		params.Channel = channel
	}
	_, _, err := s.client.PostMessage(
		params.Channel.ID,
		slack.MsgOptionText(params.Message, false),
	)
	return err
}

func (s *Slack) SendMessageToThread(params *SlackMessageOpts) error {
	_, _, err := s.client.PostMessage(
		params.ChannelID,
		slack.MsgOptionTS(params.ThreadTs),
	)
	return err
}

func (s *Slack) UploadFile(params *SlackMessageOpts) error {
	if params.Channel == nil {
		channel, err := s.NewDMChannel(params)
		if err != nil {
			params.Logger.Error().Err(err).Msg("failed to create DM channel")
			return err
		}
		params.Channel = channel
	}
	summary, err := s.client.UploadFileV2(slack.UploadFileV2Parameters{
		File:           params.Files.Files[0].Name,
		Channel:        params.Channel.ID,
		InitialComment: params.Files.Comment,
		Title:          params.Files.Title,
		Reader:         strings.NewReader(string(params.Files.Files[0].Data)),
	})
	if err != nil {
		params.Logger.Error().Err(err).Msg("failed to upload file")
		return err
	}
	params.FileId = summary.ID
	return nil
}
