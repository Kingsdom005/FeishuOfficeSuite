package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/feishu/feishu-office-suite/api/feishu/v1"
	"github.com/feishu/feishu-office-suite/internal/data"
	"github.com/feishu/feishu-office-suite/pkg/feishu"
)

type CalendarHandler struct {
	v1.UnimplementedFeishuCalendarServer
	data *data.Data
}

func NewCalendarHandler(data *data.Data) *CalendarHandler {
	return &CalendarHandler{data: data}
}

func (h *CalendarHandler) GetCalendar(ctx context.Context, req *v1.GetCalendarRequest) (*v1.GetCalendarResponse, error) {
	calendar, err := feishu.GetCalendar(ctx, h.data, req.CalendarId)
	if err != nil {
		return nil, err
	}

	return &v1.GetCalendarResponse{
		Calendar: convertCalendarToProto(calendar),
	}, nil
}

func (h *CalendarHandler) ListCalendars(ctx context.Context, req *v1.ListCalendarsRequest) (*v1.ListCalendarsResponse, error) {
	calendars, err := feishu.ListCalendars(ctx, h.data, int(req.PageSize))
	if err != nil {
		return nil, err
	}

	var protoCalendars []*v1.Calendar
	for _, cal := range calendars {
		protoCalendars = append(protoCalendars, convertCalendarToProto(cal))
	}

	return &v1.ListCalendarsResponse{
		Calendars:    protoCalendars,
		HasMore:      false,
		NextPageToken: "",
	}, nil
}

func (h *CalendarHandler) CreateEvent(ctx context.Context, req *v1.CreateEventRequest) (*v1.CreateEventResponse, error) {
	event, err := feishu.CreateEvent(ctx, h.data, req)
	if err != nil {
		return nil, err
	}

	return &v1.CreateEventResponse{
		Event: convertEventToProto(event),
	}, nil
}

func (h *CalendarHandler) GetEvent(ctx context.Context, req *v1.GetEventRequest) (*v1.GetEventResponse, error) {
	event, err := feishu.GetEvent(ctx, h.data, req.CalendarId, req.EventId)
	if err != nil {
		return nil, err
	}

	return &v1.GetEventResponse{
		Event: convertEventToProto(event),
	}, nil
}

func (h *CalendarHandler) ListEvents(ctx context.Context, req *v1.ListEventsRequest) (*v1.ListEventsResponse, error) {
	events, err := feishu.ListEvents(ctx, h.data, req.CalendarId, int(req.PageSize))
	if err != nil {
		return nil, err
	}

	var protoEvents []*v1.Event
	for _, event := range events {
		protoEvents = append(protoEvents, convertEventToProto(event))
	}

	return &v1.ListEventsResponse{
		Events:        protoEvents,
		HasMore:      false,
		NextPageToken: "",
	}, nil
}

func (h *CalendarHandler) UpdateEvent(ctx context.Context, req *v1.UpdateEventRequest) (*emptypb.Empty, error) {
	err := feishu.UpdateEvent(ctx, h.data, req)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *CalendarHandler) DeleteEvent(ctx context.Context, req *v1.DeleteEventRequest) (*emptypb.Empty, error) {
	err := feishu.DeleteEvent(ctx, h.data, req.CalendarId, req.EventId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

type ApprovalHandler struct {
	v1.UnimplementedFeishuApprovalServer
	data *data.Data
}

func NewApprovalHandler(data *data.Data) *ApprovalHandler {
	return &ApprovalHandler{data: data}
}

func (h *ApprovalHandler) GetApproval(ctx context.Context, req *v1.GetApprovalRequest) (*v1.GetApprovalResponse, error) {
	approval, err := feishu.GetApproval(ctx, h.data, req.ApprovalId)
	if err != nil {
		return nil, err
	}

	return &v1.GetApprovalResponse{
		Approval: convertApprovalToProto(approval),
	}, nil
}

func (h *ApprovalHandler) ListApprovals(ctx context.Context, req *v1.ListApprovalsRequest) (*v1.ListApprovalsResponse, error) {
	approvals, err := feishu.ListApprovals(ctx, h.data, int(req.PageSize))
	if err != nil {
		return nil, err
	}

	var protoApprovals []*v1.Approval
	for _, approval := range approvals {
		protoApprovals = append(protoApprovals, convertApprovalToProto(approval))
	}

	return &v1.ListApprovalsResponse{
		Approvals:    protoApprovals,
		HasMore:      false,
		NextPageToken: "",
	}, nil
}

func (h *ApprovalHandler) CreateApproval(ctx context.Context, req *v1.CreateApprovalRequest) (*v1.CreateApprovalResponse, error) {
	approval, err := feishu.CreateApproval(ctx, h.data, req)
	if err != nil {
		return nil, err
	}

	return &v1.CreateApprovalResponse{
		Approval: convertApprovalToProto(approval),
	}, nil
}

func (h *ApprovalHandler) GetApprovalInstance(ctx context.Context, req *v1.GetApprovalInstanceRequest) (*v1.GetApprovalInstanceResponse, error) {
	instance, err := feishu.GetApprovalInstance(ctx, h.data, req.InstanceId)
	if err != nil {
		return nil, err
	}

	return &v1.GetApprovalInstanceResponse{
		Instance: convertApprovalInstanceToProto(instance),
	}, nil
}

func (h *ApprovalHandler) ListApprovalInstances(ctx context.Context, req *v1.ListApprovalInstancesRequest) (*v1.ListApprovalInstancesResponse, error) {
	instances, err := feishu.ListApprovalInstances(ctx, h.data, req.ApprovalId, int(req.PageSize))
	if err != nil {
		return nil, err
	}

	var protoInstances []*v1.ApprovalInstance
	for _, instance := range instances {
		protoInstances = append(protoInstances, convertApprovalInstanceToProto(instance))
	}

	return &v1.ListApprovalInstancesResponse{
		Instances:    protoInstances,
		HasMore:      false,
		NextPageToken: "",
	}, nil
}

func (h *ApprovalHandler) ApproveInstance(ctx context.Context, req *v1.ApproveInstanceRequest) (*emptypb.Empty, error) {
	err := feishu.ApproveInstance(ctx, h.data, req.InstanceId, req.Comment)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *ApprovalHandler) RejectInstance(ctx context.Context, req *v1.RejectInstanceRequest) (*emptypb.Empty, error) {
	err := feishu.RejectInstance(ctx, h.data, req.InstanceId, req.Comment)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}