package slack

import (
	"fmt"

	"github.com/slack-go/slack"
)

type SlackDmChannel struct {
	ChannelID string
	UserID    string
	Message   string
	ThreadTs  string
}

func (s *Slack) NewDMChannel(params *SlackMessageOpts) (*slack.Channel, error) {
	if params.UserID == "" {
		return nil, fmt.Errorf("user id is required")
	}
	user, err := s.client.GetUserByEmail(params.UserID)
	if err != nil {
		params.Logger.Error().Err(err).Msg("failed to get user by email")
		return nil, err
	}
	channel, _, _, err := s.client.OpenConversation(&slack.OpenConversationParameters{
		Users: []string{user.ID},
	})
	if err != nil {
		params.Logger.Error().Err(err).Msg("failed to open conversation")
		return nil, err
	}
	return channel, nil
}
