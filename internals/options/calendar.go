package options

type CalendarBase struct {
	UserEmail string
}

type CreateEventRequest struct {
	CalendarBase
	Summary     string `json:"summary"`
	Location    string `json:"location"`
	Description string `json:"description"`
	StartTime   string `json:"starttime"`
	EndTime     string `json:"endtime"`
	Attendees   []string
}

type CreateEventResponse struct {
	EventId string `json:"eventId"`
	Status  string `json:"status"`
	Link    string `json:"link"`
}

type UpdateEventRequest struct {
	CalendarBase
	EventId     string `json:"eventid"`
	Summary     string `json:"summary"`
	Location    string `json:"location"`
	Description string `json:"description"`
	StartTime   string `json:"starttime"`
	EndTime     string `json:"endtime"`
}

type DeleteEventRequest struct {
	CalendarBase
	EventId string `json:"eventid"`
}

type GetEventRequest struct {
	CalendarBase
	EventId string `json:"eventid"`
}

type GetEventResponse struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
	StartTime   string `json:"starttime"`
	EndTime     string `json:"endtime"`
}

type EventSummary struct {
	CalendarBase
	EventID   string `json:"eventId"`
	Summary   string `json:"summary"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type ListEventsResponse struct {
	Events []EventSummary `json:"events"`
	Status string         `json:"status"`
}
