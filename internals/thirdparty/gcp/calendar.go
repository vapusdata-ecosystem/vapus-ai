package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	options "github.com/vapusdata-ecosystem/vapusai/core/options"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

type CalendarClient struct {
	client         *calendar.Service
	logger         zerolog.Logger
	clientEmail    string
	saKey          string
	isImpersonator bool
}

func NewGoogleCalendar(ctx context.Context, opts *GcpConfig, logger zerolog.Logger) (*CalendarClient, error) {
	if opts.IsDomainScopeApp {
		creds, err := google.CredentialsFromJSON(ctx, opts.ServiceAccountKey)
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
		return &CalendarClient{
			logger:         logger,
			clientEmail:    clEmail,
			isImpersonator: true,
			saKey: string(opts.ServiceAccountKey),
		}, nil
	}

	return nil, nil
}
func (c *CalendarClient) getClient(ctx context.Context, userEmail string) *calendar.Service {
	if !c.isImpersonator {
		return c.client
	}
	tokenSource, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
		TargetPrincipal: c.clientEmail,
		Subject:         userEmail,
		Scopes:          []string{"https://www.googleapis.com/auth/calendar"},
	}, option.WithCredentialsJSON([]byte(c.saKey)))
	if err != nil {
		c.logger.Err(err).Msgf("Error while impersonating the user -- %v", err)
		return nil
	}
	client, err := calendar.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil
	}

	return client
}

func (c *CalendarClient) CreateEvent(ctx context.Context, req *options.CreateEventRequest) (*options.CreateEventResponse, error) {
	client := c.getClient(ctx, req.UserEmail)
	if client == nil {
		return nil, fmt.Errorf("failed to get calendar client")

	}
	var attendees []*calendar.EventAttendee
	for _, email := range req.Attendees {
		attendees = append(attendees, &calendar.EventAttendee{
			Email: email,
		})
	}
	event := &calendar.Event{
		Summary:     req.Summary,
		Location:    req.Location,
		Description: req.Description,
		Start: &calendar.EventDateTime{
			DateTime: req.StartTime,
			TimeZone: "Asia/Kolkata",
		},
		End: &calendar.EventDateTime{
			DateTime: req.EndTime,
			TimeZone: "Asia/Kolkata",
		},
		Attendees: attendees,
	}

	createdEvent, err := client.Events.Insert("primary", event).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	response := &options.CreateEventResponse{
		EventId: createdEvent.Id,
		Status:  "success",
		Link:    createdEvent.HtmlLink,
	}
	return response, nil

}

func (c *CalendarClient) UpdateEvent(ctx context.Context, req *options.UpdateEventRequest) error {
	client := c.getClient(ctx, req.UserEmail)
	if client == nil {
		return fmt.Errorf("failed to get calendar client")
	}
	event, err := client.Events.Get("primary", req.EventId).Do()
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}
	// Update fields if provided
	if req.Summary != "" {
		event.Summary = req.Summary
	}
	if req.Description != "" {
		event.Description = req.Description
	}
	if req.Location != "" {
		event.Location = req.Location
	}
	if req.StartTime != "" {
		event.Start = &calendar.EventDateTime{DateTime: req.StartTime, TimeZone: "Asia/Kolkata"}
	}
	if req.EndTime != "" {
		event.End = &calendar.EventDateTime{DateTime: req.EndTime, TimeZone: "Asia/Kolkata"}
	}

	_, err = client.Events.Update("primary", req.EventId, event).Do()
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}
func (c *CalendarClient) GetEvent(ctx context.Context, req *options.GetEventRequest) (*options.GetEventResponse, error) {
	client := c.getClient(ctx, req.UserEmail)
	if client == nil {
		return nil, fmt.Errorf("failed to get calendar client")
	}
	event, err := client.Events.Get("primary", req.EventId).Do()
	if err != nil {
		return nil, err
	}
	resp := &options.GetEventResponse{
		Summary:     event.Summary,
		Description: event.Description,
		StartTime:   event.Start.DateTime,
		EndTime:     event.End.DateTime,
	}

	fmt.Println("These are the event details ", resp.Summary, resp.StartTime, resp.EndTime)
	return resp, nil
}

func (c *CalendarClient) DeleteEvent(ctx context.Context, req *options.DeleteEventRequest) error {
	client := c.getClient(ctx, req.UserEmail)
	if client == nil {
		return fmt.Errorf("failed to get calendar client")
	}
	err := client.Events.Delete("primary", req.EventId).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	return nil
}

func (c *CalendarClient) ListEvent(ctx context.Context, userEmail string) (*options.ListEventsResponse, error) {
	client := c.getClient(ctx, userEmail)
	if client == nil {
		return nil, fmt.Errorf("failed to get calendar client")
	}
	events, err := client.Events.List("primary").Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	var eventList []options.EventSummary
	for _, item := range events.Items {
		eventList = append(eventList, options.EventSummary{
			EventID:   item.Id,
			Summary:   item.Summary,
			StartTime: item.Start.DateTime,
			EndTime:   item.End.DateTime,
		})
	}

	return &options.ListEventsResponse{
		Events: eventList,
		Status: "success",
	}, nil
}
