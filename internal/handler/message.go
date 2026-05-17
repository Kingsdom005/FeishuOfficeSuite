package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/feishu/feishu-office-suite/api/feishu/v1"
	"github.com/feishu/feishu-office-suite/internal/data"
	"github.com/feishu/feishu-office-suite/pkg/feishu"
)

type MessageHandler struct {
	v1.UnimplementedFeishuMessageServer
	data *data.Data
}

func NewMessageHandler(data *data.Data) *MessageHandler {
	return &MessageHandler{data: data}
}

func (h *MessageHandler) SendMessage(ctx context.Context, req *v1.SendMessageRequest) (*v1.SendMessageResponse, error) {
	msg, err := feishu.SendMessage(ctx, h.data, req.ReceiveId, req.ReceiveIdType, req.MsgType, req.Content)
	if err != nil {
		return nil, err
	}

	return &v1.SendMessageResponse{
		Message: convertMessageToProto(msg),
	}, nil
}

func (h *MessageHandler) SendTextMessage(ctx context.Context, req *v1.SendTextMessageRequest) (*v1.SendMessageResponse, error) {
	msg, err := feishu.SendTextMessage(ctx, h.data, req.ReceiveId, req.ReceiveIdType, req.Text)
	if err != nil {
		return nil, err
	}

	return &v1.SendMessageResponse{
		Message: convertMessageToProto(msg),
	}, nil
}

func (h *MessageHandler) SendImageMessage(ctx context.Context, req *v1.SendImageMessageRequest) (*v1.SendMessageResponse, error) {
	msg, err := feishu.SendImageMessage(ctx, h.data, req.ReceiveId, req.ReceiveIdType, req.ImageKey)
	if err != nil {
		return nil, err
	}

	return &v1.SendMessageResponse{
		Message: convertMessageToProto(msg),
	}, nil
}

func (h *MessageHandler) SendCardMessage(ctx context.Context, req *v1.SendCardMessageRequest) (*v1.SendMessageResponse, error) {
	msg, err := feishu.SendCardMessage(ctx, h.data, req.ReceiveId, req.ReceiveIdType, req.CardContent)
	if err != nil {
		return nil, err
	}

	return &v1.SendMessageResponse{
		Message: convertMessageToProto(msg),
	}, nil
}

func (h *MessageHandler) GetMessage(ctx context.Context, req *v1.GetMessageRequest) (*v1.GetMessageResponse, error) {
	msg, err := feishu.GetMessage(ctx, h.data, req.MessageId)
	if err != nil {
		return nil, err
	}

	return &v1.GetMessageResponse{
		Message: convertMessageToProto(msg),
	}, nil
}

func (h *MessageHandler) ListMessages(ctx context.Context, req *v1.ListMessagesRequest) (*v1.ListMessagesResponse, error) {
	messages, err := feishu.ListMessages(ctx, h.data, req.ChatId, int(req.PageSize))
	if err != nil {
		return nil, err
	}

	var protoMessages []*v1.Message
	for _, msg := range messages {
		protoMessages = append(protoMessages, convertMessageToProto(msg))
	}

	return &v1.ListMessagesResponse{
		Messages: protoMessages,
		HasMore:  false,
	}, nil
}

func (h *MessageHandler) UpdateMessage(ctx context.Context, req *v1.UpdateMessageRequest) (*emptypb.Empty, error) {
	err := feishu.UpdateMessage(ctx, h.data, req.MessageId, req.Content, req.MsgType)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *MessageHandler) DeleteMessage(ctx context.Context, req *v1.DeleteMessageRequest) (*emptypb.Empty, error) {
	err := feishu.DeleteMessage(ctx, h.data, req.MessageId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *MessageHandler) ReplyMessage(ctx context.Context, req *v1.ReplyMessageRequest) (*v1.SendMessageResponse, error) {
	msg, err := feishu.ReplyMessage(ctx, h.data, req.MessageId, req.Content, req.MsgType)
	if err != nil {
		return nil, err
	}

	return &v1.SendMessageResponse{
		Message: convertMessageToProto(msg),
	}, nil
}

type ChatHandler struct {
	v1.UnimplementedFeishuChatServer
	data *data.Data
}

func NewChatHandler(data *data.Data) *ChatHandler {
	return &ChatHandler{data: data}
}

func (h *ChatHandler) GetChat(ctx context.Context, req *v1.GetChatRequest) (*v1.GetChatResponse, error) {
	chat, err := feishu.GetChat(ctx, h.data, req.ChatId)
	if err != nil {
		return nil, err
	}

	return &v1.GetChatResponse{
		Chat: convertChatToProto(chat),
	}, nil
}

func (h *ChatHandler) ListChats(ctx context.Context, req *v1.ListChatsRequest) (*v1.ListChatsResponse, error) {
	chats, err := feishu.ListChats(ctx, h.data, int(req.PageSize))
	if err != nil {
		return nil, err
	}

	var protoChats []*v1.Chat
	for _, chat := range chats {
		protoChats = append(protoChats, convertChatToProto(chat))
	}

	return &v1.ListChatsResponse{
		Chats:        protoChats,
		HasMore:      false,
		NextPageToken: "",
	}, nil
}

func (h *ChatHandler) CreateChat(ctx context.Context, req *v1.CreateChatRequest) (*v1.CreateChatResponse, error) {
	chat, err := feishu.CreateChat(ctx, h.data, req)
	if err != nil {
		return nil, err
	}

	return &v1.CreateChatResponse{
		Chat: convertChatToProto(chat),
	}, nil
}

func (h *ChatHandler) UpdateChat(ctx context.Context, req *v1.UpdateChatRequest) (*emptypb.Empty, error) {
	err := feishu.UpdateChat(ctx, h.data, req.ChatId, req.Name, req.Description)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *ChatHandler) DeleteChat(ctx context.Context, req *v1.DeleteChatRequest) (*emptypb.Empty, error) {
	err := feishu.DeleteChat(ctx, h.data, req.ChatId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *ChatHandler) AddMembers(ctx context.Context, req *v1.AddChatMembersRequest) (*emptypb.Empty, error) {
	err := feishu.AddChatMembers(ctx, h.data, req.ChatId, req.MemberIds)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *ChatHandler) RemoveMembers(ctx context.Context, req *v1.RemoveChatMembersRequest) (*emptypb.Empty, error) {
	err := feishu.RemoveChatMembers(ctx, h.data, req.ChatId, req.MemberIds)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}