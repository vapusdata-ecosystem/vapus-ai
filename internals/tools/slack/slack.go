package slack

import (
	slack "github.com/slack-go/slack"
)

type Options func(*SlackParams)

type SlackParams struct {
	WebhookURL string
	Token      string
}

func WithWebhookURL(url string) Options {
	return func(p *SlackParams) {
		p.WebhookURL = url
	}
}

func WithToken(token string) Options {
	return func(p *SlackParams) {
		p.Token = token
	}
}

type Slack struct {
	client *slack.Client
}

func NewSlack(opts ...Options) *Slack {
	params := &SlackParams{}
	for _, opt := range opts {
		opt(params)
	}
	return &Slack{
		client: slack.New(params.Token),
	}
}
