package calendar

import (
	"context"
	"encoding/base64"
	"fmt"

	// "log"

	"github.com/rs/zerolog"
	apperr "github.com/vapusdata-ecosystem/vapusai/core/app/errors"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	"github.com/vapusdata-ecosystem/vapusai/core/options"
	dmerrors "github.com/vapusdata-ecosystem/vapusai/core/pkgs/errors"
	"github.com/vapusdata-ecosystem/vapusai/core/thirdparty/gcp"
	"github.com/vapusdata-ecosystem/vapusai/core/types"
)

type Calendar interface {
	// CreateEvent(ctx context.Context, summary, location, description string, startTime, endTime string) (string, error)
	// UpdateEvent(ctx context.Context, eventId, newSummary , newLocation, newDescription, newStart, newEnd string) error
	// ListEvent(ctx context.Context) error
	// GetEvent(ctx context.Context,eventId string) (*calendar.Event, error)
	// DeleteEvent(ctx context.Context, eventId string) error
	CreateEvent(ctx context.Context, req *options.CreateEventRequest) (*options.CreateEventResponse, error)
	UpdateEvent(ctx context.Context, req *options.UpdateEventRequest) error
	ListEvent(ctx context.Context, userEmail string) (*options.ListEventsResponse, error)
	GetEvent(ctx context.Context, req *options.GetEventRequest) (*options.GetEventResponse, error)
	DeleteEvent(ctx context.Context, req *options.DeleteEventRequest) error
}
type CalendarClient struct {
	logger zerolog.Logger
	Calendar
}

func New(ctx context.Context, service string, netOps *models.PluginNetworkParams,
	ops []*models.Mapper, logger zerolog.Logger) (Calendar, error) {
	calendar := &CalendarClient{
		logger: logger,
	}
	// log.Println(netOps.Credentials.GcpCreds.ServiceAccountKey)
	switch service {
	case types.GOOGLE_CALENDAR.String():
		serviceAccountKey, err := base64.StdEncoding.DecodeString(netOps.Credentials.GcpCreds.ServiceAccountKey)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to decode service account key")
			return nil, err
		}
		cfg := &gcp.GcpConfig{
			ServiceAccountKey: serviceAccountKey,
			ProjectID:         netOps.Credentials.GcpCreds.ProjectId,
			IsDomainScopeApp:  netOps.Credentials.GcpCreds.IsDomainScopeApp,
		}

		client, err := gcp.NewGoogleCalendar(ctx, cfg, logger)

		if err != nil {
			logger.Error().Err(err).Msg("Failed to create Google Calendar client")
			return nil, err
		}

		calendar.Calendar = client
	default:
		logger.Error().Msg("Invalid calendar service")
		return nil, dmerrors.DMError(apperr.ErrInvalidCalenderService, apperr.ErrInvalidCalenderConn)
	}

	return calendar, nil
}
