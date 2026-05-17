package entity

import (
	"time"
)

type User struct {
	ID            string    `json:"id"`
	UnionID       string    `json:"union_id"`
	OpenID        string    `json:"open_id"`
	Name          string    `json:"name"`
	EnName        string    `json:"en_name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	AvatarURL     string    `json:"avatar_url"`
	AvatarThumb   string    `json:"avatar_thumb"`
	AvatarMiddle  string    `json:"avatar_middle"`
	Status        string    `json:"status"`
	IsActivated   bool      `json:"is_activated"`
	IsTenantAccess bool     `json:"is_tenant_access"`
	DepartmentID  string    `json:"department_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Department struct {
	DepartmentID string    `json:"department_id"`
	Name        string    `json:"name"`
	NameEn      string    `json:"name_en"`
	ParentID    string    `json:"parent_id"`
	Order       int       `json:"order"`
	IsRoot      bool      `json:"is_root"`
	MemberCount int       `json:"member_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Message struct {
	MessageID   string    `json:"message_id"`
	ChatID      string    `json:"chat_id"`
	Sender     string    `json:"sender"`
	SenderID   string    `json:"sender_id"`
	SenderType  string    `json:"sender_type"`
	Content     string    `json:"content"`
	MsgType     string    `json:"msg_type"`
	IsDeleted   bool      `json:"is_deleted"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

type Chat struct {
	ChatID       string    `json:"chat_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	MemberCount  int       `json:"member_count"`
	OwnerID      string    `json:"owner_id"`
	OwnerIDType  string    `json:"owner_id_type"`
	IsActivated  bool      `json:"is_activated"`
	CreatedAt    time.Time `json:"created_at"`
}

type Calendar struct {
	CalendarID  string    `json:"calendar_id"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Timezone    string    `json:"timezone"`
	IsPrimary   bool      `json:"is_primary"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Event struct {
	EventID     string     `json:"event_id"`
	CalendarID  string     `json:"calendar_id"`
	Summary     string     `json:"summary"`
	Description string     `json:"description"`
	IsAllDay    bool       `json:"is_all_day"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     time.Time  `json:"end_time"`
	Timezone    string     `json:"timezone"`
	Attendees   []Attendee `json:"attendees"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Attendee struct {
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Email       string `json:"email"`
	Status      string `json:"status"`
	IsOrganizer bool   `json:"is_organizer"`
}

type Approval struct {
	ApprovalID    string    `json:"approval_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	FormContent    string    `json:"form_content"`
	InstanceCount  int       `json:"instance_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ApprovalInstance struct {
	InstanceID string    `json:"instance_id"`
	ApprovalID string    `json:"approval_id"`
	Title      string    `json:"title"`
	Status     string    `json:"status"`
	InitiatorID string   `json:"initiator_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}