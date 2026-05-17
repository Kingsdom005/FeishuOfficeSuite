package handler

import (
	"time"

	"github.com/feishu/feishu-office-suite/api/feishu/v1"
	"github.com/feishu/feishu-office-suite/internal/domain/entity"
)

func convertUserToProto(user *entity.User) *v1.User {
	if user == nil {
		return nil
	}
	return &v1.User{
		Id:             user.ID,
		UnionId:        user.UnionID,
		OpenId:         user.OpenID,
		Name:           user.Name,
		EnName:         user.EnName,
		Email:          user.Email,
		Phone:          user.Phone,
		AvatarUrl:      user.AvatarURL,
		AvatarThumb:    user.AvatarThumb,
		AvatarMiddle:   user.AvatarMiddle,
		Status:         user.Status,
		IsActivated:   user.IsActivated,
		IsTenantAccess: user.IsTenantAccess,
		CreatedAt:      timestampToProto(user.CreatedAt),
		UpdatedAt:      timestampToProto(user.UpdatedAt),
	}
}

func convertDepartmentToProto(dept *entity.Department) *v1.Department {
	if dept == nil {
		return nil
	}
	return &v1.Department{
		DepartmentId: dept.DepartmentID,
		Name:         dept.Name,
		NameEn:       dept.NameEn,
		ParentId:     dept.ParentID,
		Order:        int32(dept.Order),
		IsRoot:       dept.IsRoot,
		MemberCount:  int32(dept.MemberCount),
		CreatedAt:    timestampToProto(dept.CreatedAt),
		UpdatedAt:    timestampToProto(dept.UpdatedAt),
	}
}

func convertMessageToProto(msg *entity.Message) *v1.Message {
	if msg == nil {
		return nil
	}
	return &v1.Message{
		MessageId:   msg.MessageID,
		ChatId:      msg.ChatID,
		Sender:     msg.Sender,
		SenderId:   msg.SenderID,
		SenderType: msg.SenderType,
		Content:    msg.Content,
		MsgType:    msg.MsgType,
		IsDeleted: msg.IsDeleted,
		CreateTime: timestampToProto(msg.CreateTime),
		UpdateTime: timestampToProto(msg.UpdateTime),
	}
}

func convertChatToProto(chat *entity.Chat) *v1.Chat {
	if chat == nil {
		return nil
	}
	return &v1.Chat{
		ChatId:       chat.ChatID,
		Name:         chat.Name,
		Description:  chat.Description,
		MemberCount:  int32(chat.MemberCount),
		OwnerId:      chat.OwnerID,
		OwnerIdType:  chat.OwnerIDType,
		IsActivated: chat.IsActivated,
		CreatedAt:   timestampToProto(chat.CreatedAt),
	}
}

func convertCalendarToProto(cal *entity.Calendar) *v1.Calendar {
	if cal == nil {
		return nil
	}
	return &v1.Calendar{
		CalendarId: cal.CalendarID,
		Summary:    cal.Summary,
		Description: cal.Description,
		Type:       cal.Type,
		Timezone:   cal.Timezone,
		IsPrimary: cal.IsPrimary,
		CreatedAt: timestampToProto(cal.CreatedAt),
		UpdatedAt: timestampToProto(cal.UpdatedAt),
	}
}

func convertEventToProto(event *entity.Event) *v1.Event {
	if event == nil {
		return nil
	}
	var attendees []*v1.Attendee
	for _, a := range event.Attendees {
		attendees = append(attendees, &v1.Attendee{
			UserId:     a.UserID,
			UserName:   a.UserName,
			Email:      a.Email,
			Status:     a.Status,
			IsOrganizer: a.IsOrganizer,
		})
	}
	return &v1.Event{
		EventId:     event.EventID,
		CalendarId: event.CalendarID,
		Summary:     event.Summary,
		Description: event.Description,
		IsAllDay:   event.IsAllDay,
		StartTime:  timestampToProto(event.StartTime),
		EndTime:    timestampToProto(event.EndTime),
		Timezone:   event.Timezone,
		Attendees:  attendees,
		Status:     event.Status,
		CreatedAt: timestampToProto(event.CreatedAt),
		UpdatedAt: timestampToProto(event.UpdatedAt),
	}
}

func convertApprovalToProto(approval *entity.Approval) *v1.Approval {
	if approval == nil {
		return nil
	}
	return &v1.Approval{
		ApprovalId:   approval.ApprovalID,
		Name:         approval.Name,
		Description:  approval.Description,
		FormContent:  approval.FormContent,
		InstanceCount: int32(approval.InstanceCount),
		CreatedAt:   timestampToProto(approval.CreatedAt),
		UpdatedAt:   timestampToProto(approval.UpdatedAt),
	}
}

func convertApprovalInstanceToProto(instance *entity.ApprovalInstance) *v1.ApprovalInstance {
	if instance == nil {
		return nil
	}
	return &v1.ApprovalInstance{
		InstanceId: instance.InstanceID,
		ApprovalId: instance.ApprovalID,
		Title:      instance.Title,
		Status:     instance.Status,
		InitiatorId: instance.InitiatorID,
		StartTime: timestampToProto(instance.StartTime),
		EndTime:   timestampToProto(instance.EndTime),
	}
}

func timestampToProto(t time.Time) *timestamp {
	return &timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	}
}

type timestamp struct {
	Seconds int64 `json:"seconds"`
	Nanos   int32 `json:"nanos"`
}