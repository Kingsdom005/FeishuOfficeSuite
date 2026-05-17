package feishu

import (
	"context"
	"fmt"

	"github.com/feishu/feishu-office-suite/api/feishu/v1"
	"github.com/feishu/feishu-office-suite/internal/data"
	"github.com/feishu/feishu-office-suite/internal/domain/entity"
)

func SendMessage(ctx context.Context, d *data.Data, receiveID, receiveIDType, msgType, content string) (*entity.Message, error) {
	client := GetClient(d)

	url := "https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=" + receiveIDType

	payload := map[string]interface{}{
		"receive_id": receiveID,
		"msg_type":   msgType,
		"content":    content,
	}

	result, err := client.DoRequest(ctx, "POST", url, payload)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	msgData, ok := data["message"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid message data")
	}

	return parseMessage(msgData), nil
}

func SendTextMessage(ctx context.Context, d *data.Data, receiveID, receiveIDType, text string) (*entity.Message, error) {
	content := fmt.Sprintf(`{"text":"%s"}`, text)
	return SendMessage(ctx, d, receiveID, receiveIDType, "text", content)
}

func SendImageMessage(ctx context.Context, d *data.Data, receiveID, receiveIDType, imageKey string) (*entity.Message, error) {
	content := fmt.Sprintf(`{"image_key":"%s"}`, imageKey)
	return SendMessage(ctx, d, receiveID, receiveIDType, "image", content)
}

func SendCardMessage(ctx context.Context, d *data.Data, receiveID, receiveIDType, cardContent string) (*entity.Message, error) {
	return SendMessage(ctx, d, receiveID, receiveIDType, "interactive", cardContent)
}

func GetMessage(ctx context.Context, d *data.Data, messageID string) (*entity.Message, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/messages/%s", messageID)

	result, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	msgData, ok := data["message"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid message data")
	}

	return parseMessage(msgData), nil
}

func ListMessages(ctx context.Context, d *data.Data, chatID string, pageSize int) ([]*entity.Message, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/messages?container_id_type=chat&container_id=%s&page_size=%d", chatID, pageSize)

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

	var messages []*entity.Message
	for _, item := range items {
		if msgData, ok := item.(map[string]interface{}); ok {
			messages = append(messages, parseMessage(msgData))
		}
	}

	return messages, nil
}

func UpdateMessage(ctx context.Context, d *data.Data, messageID, content, msgType string) error {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/messages/%s", messageID)

	payload := map[string]interface{}{
		"msg_type": msgType,
		"content":  content,
	}

	_, err := client.DoRequest(ctx, "PATCH", url, payload)
	return err
}

func DeleteMessage(ctx context.Context, d *data.Data, messageID string) error {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/messages/%s", messageID)

	_, err := client.DoRequest(ctx, "DELETE", url, nil)
	return err
}

func ReplyMessage(ctx context.Context, d *data.Data, messageID, content, msgType string) (*entity.Message, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/messages/%s/reply", messageID)

	payload := map[string]interface{}{
		"msg_type": msgType,
		"content":  content,
	}

	result, err := client.DoRequest(ctx, "POST", url, payload)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	msgData, ok := data["message"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid message data")
	}

	return parseMessage(msgData), nil
}

func GetChat(ctx context.Context, d *data.Data, chatID string) (*entity.Chat, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/chats/%s", chatID)

	result, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	chatData, ok := data["chat"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid chat data")
	}

	return parseChat(chatData), nil
}

func ListChats(ctx context.Context, d *data.Data, pageSize int) ([]*entity.Chat, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/chats?page_size=%d", pageSize)

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

	var chats []*entity.Chat
	for _, item := range items {
		if chatData, ok := item.(map[string]interface{}); ok {
			chats = append(chats, parseChat(chatData))
		}
	}

	return chats, nil
}

func CreateChat(ctx context.Context, d *data.Data, req *v1.CreateChatRequest) (*entity.Chat, error) {
	client := GetClient(d)

	url := "https://open.feishu.cn/open-apis/im/v1/chats"

	payload := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"owner_id":    req.OwnerId,
		"user_id_list": req.MemberIds,
	}

	result, err := client.DoRequest(ctx, "POST", url, payload)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	chatData, ok := data["chat"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid chat data")
	}

	return parseChat(chatData), nil
}

func UpdateChat(ctx context.Context, d *data.Data, chatID, name, description string) error {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/chats/%s", chatID)

	payload := map[string]interface{}{
		"name":        name,
		"description": description,
	}

	_, err := client.DoRequest(ctx, "PATCH", url, payload)
	return err
}

func DeleteChat(ctx context.Context, d *data.Data, chatID string) error {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/chats/%s", chatID)

	_, err := client.DoRequest(ctx, "DELETE", url, nil)
	return err
}

func AddChatMembers(ctx context.Context, d *data.Data, chatID string, memberIDs []string) error {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/chats/%s/members", chatID)

	payload := map[string]interface{}{
		"user_id_list": memberIDs,
	}

	_, err := client.DoRequest(ctx, "POST", url, payload)
	return err
}

func RemoveChatMembers(ctx context.Context, d *data.Data, chatID string, memberIDs []string) error {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/im/v1/chats/%s/members", chatID)

	payload := map[string]interface{}{
		"user_id_list": memberIDs,
	}

	_, err := client.DoRequest(ctx, "DELETE", url, payload)
	return err
}

func parseMessage(data map[string]interface{}) *entity.Message {
	msg := &entity.Message{}

	if v, ok := data["message_id"].(string); ok {
		msg.MessageID = v
	}
	if v, ok := data["chat_id"].(string); ok {
		msg.ChatID = v
	}
	if v, ok := data["sender"].(map[string]interface{}); ok {
		if senderID, ok := v["id"].(string); ok {
			msg.SenderID = senderID
		}
		if senderType, ok := v["sender_type"].(string); ok {
			msg.SenderType = senderType
		}
	}
	if v, ok := data["msg_type"].(string); ok {
		msg.MsgType = v
	}
	if v, ok := data["body"].(map[string]interface{}); ok {
		if content, ok := v["content"].(string); ok {
			msg.Content = content
		}
	}

	return msg
}

func parseChat(data map[string]interface{}) *entity.Chat {
	chat := &entity.Chat{}

	if v, ok := data["chat_id"].(string); ok {
		chat.ChatID = v
	}
	if v, ok := data["name"].(string); ok {
		chat.Name = v
	}
	if v, ok := data["description"].(string); ok {
		chat.Description = v
	}
	if v, ok := data["owner_id"].(string); ok {
		chat.OwnerID = v
	}
	if v, ok := data["member_count"].(float64); ok {
		chat.MemberCount = int(v)
	}

	return chat
}