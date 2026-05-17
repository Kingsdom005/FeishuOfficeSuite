package feishu

import (
	"context"
	"fmt"

	"github.com/feishu/feishu-office-suite/api/feishu/v1"
	"github.com/feishu/feishu-office-suite/internal/data"
	"github.com/feishu/feishu-office-suite/internal/domain/entity"
)

func GetCalendar(ctx context.Context, d *data.Data, calendarID string) (*entity.Calendar, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/calendar/v4/calendars/%s", calendarID)

	result, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	calData, ok := data["calendar"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid calendar data")
	}

	return parseCalendar(calData), nil
}

func ListCalendars(ctx context.Context, d *data.Data, pageSize int) ([]*entity.Calendar, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/calendar/v4/calendars?page_size=%d", pageSize)

	result, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	items, ok := data["calendar_list"].([]interface{})
	if !ok {
		return nil, nil
	}

	var calendars []*entity.Calendar
	for _, item := range items {
		if calData, ok := item.(map[string]interface{}); ok {
			calendars = append(calendars, parseCalendar(calData))
		}
	}

	return calendars, nil
}

func CreateEvent(ctx context.Context, d *data.Data, req *v1.CreateEventRequest) (*entity.Event, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/calendar/v4/calendars/%s/events", req.CalendarId)

	payload := map[string]interface{}{
		"summary":     req.Summary,
		"description": req.Description,
		"start_time":  req.StartTime,
		"end_time":    req.EndTime,
		"timezone":    req.Timezone,
		"attendees":   req.AttendeeEmails,
	}

	result, err := client.DoRequest(ctx, "POST", url, payload)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	eventData, ok := data["event"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid event data")
	}

	return parseEvent(eventData), nil
}

func GetEvent(ctx context.Context, d *data.Data, calendarID, eventID string) (*entity.Event, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/calendar/v4/calendars/%s/events/%s", calendarID, eventID)

	result, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	eventData, ok := data["event"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid event data")
	}

	return parseEvent(eventData), nil
}

func ListEvents(ctx context.Context, d *data.Data, calendarID string, pageSize int) ([]*entity.Event, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/calendar/v4/calendars/%s/events?page_size=%d", calendarID, pageSize)

	result, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	items, ok := data["items"].([]interface{})
	if !ok {
		return nil, nil
	}

	var events []*entity.Event
	for _, item := range items {
		if eventData, ok := item.(map[string]interface{}); ok {
			events = append(events, parseEvent(eventData))
		}
	}

	return events, nil
}

func UpdateEvent(ctx context.Context, d *data.Data, req *v1.UpdateEventRequest) error {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/calendar/v4/calendars/%s/events/%s", req.CalendarId, req.EventId)

	payload := map[string]interface{}{
		"summary":     req.Summary,
		"description": req.Description,
		"start_time":  req.StartTime,
		"end_time":    req.EndTime,
		"timezone":    req.Timezone,
	}

	_, err := client.DoRequest(ctx, "PATCH", url, payload)
	return err
}

func DeleteEvent(ctx context.Context, d *data.Data, calendarID, eventID string) error {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/calendar/v4/calendars/%s/events/%s", calendarID, eventID)

	_, err := client.DoRequest(ctx, "DELETE", url, nil)
	return err
}

func parseCalendar(data map[string]interface{}) *entity.Calendar {
	cal := &entity.Calendar{}

	if v, ok := data["calendar_id"].(string); ok {
		cal.CalendarID = v
	}
	if v, ok := data["summary"].(string); ok {
		cal.Summary = v
	}
	if v, ok := data["description"].(string); ok {
		cal.Description = v
	}
	if v, ok := data["type"].(string); ok {
		cal.Type = v
	}
	if v, ok := data["timezone"].(string); ok {
		cal.Timezone = v
	}
	if v, ok := data["summary"].(string); ok {
		cal.Summary = v
	}

	return cal
}

func parseEvent(data map[string]interface{}) *entity.Event {
	event := &entity.Event{}

	if v, ok := data["event_id"].(string); ok {
		event.EventID = v
	}
	if v, ok := data["summary"].(string); ok {
		event.Summary = v
	}
	if v, ok := data["description"].(string); ok {
		event.Description = v
	}
	if v, ok := data["status"].(string); ok {
		event.Status = v
	}
	if v, ok := data["timezone"].(string); ok {
		event.Timezone = v
	}

	if v, ok := data["attendees"].([]interface{}); ok {
		for _, a := range v {
			if attendeeData, ok := a.(map[string]interface{}); ok {
				attendee := entity.Attendee{}
				if userID, ok := attendeeData["user_id"].(string); ok {
					attendee.UserID = userID
				}
				if displayName, ok := attendeeData["display_name"].(string); ok {
					attendee.UserName = displayName
				}
				if email, ok := attendeeData["email"].(string); ok {
					attendee.Email = email
				}
				event.Attendees = append(event.Attendees, attendee)
			}
		}
	}

	return event
}