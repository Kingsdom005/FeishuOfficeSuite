package feishu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/feishu/feishu-office-suite/internal/data"
	"github.com/feishu/feishu-office-suite/internal/domain/entity"
)

type FeishuClient struct {
	appID     string
	appSecret string
	token     string
	tokenExpiry time.Time
}

func NewFeishuClient(appID, appSecret, token string) *FeishuClient {
	return &FeishuClient{
		appID:     appID,
		appSecret: appSecret,
		token:     token,
	}
}

func (c *FeishuClient) GetAccessToken(ctx context.Context) (string, error) {
	if c.token != "" && time.Now().Before(c.tokenExpiry) {
		return c.token, nil
	}

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal?app_id=%s&app_secret=%s", c.appID, c.appSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if code, ok := result["code"].(float64); ok && code != 0 {
		return "", fmt.Errorf("failed to get access token: %v", result["msg"])
	}

	c.token = result["tenant_access_token"].(string)
	c.tokenExpiry = time.Now().Add(2 * time.Hour)

	return c.token, nil
}

func (c *FeishuClient) DoRequest(ctx context.Context, method, url string, body interface{}) (map[string]interface{}, error) {
	token, err := c.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if code, ok := result["code"].(float64); ok && code != 0 {
		return nil, fmt.Errorf("feishu API error: code=%v, msg=%v", code, result["msg"])
	}

	return result, nil
}

func GetClient(d *data.Data) *FeishuClient {
	cfg := &data.FeishuConfig{}
	return NewFeishuClient(cfg.AppID, cfg.AppSecret, cfg.Token)
}

func GetUser(ctx context.Context, d *data.Data, userID, userIDType string) (*entity.User, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/contact/v3/users/%s?user_id_type=%s", userID, userIDType)

	result, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	userData, ok := data["user"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid user data")
	}

	return parseUser(userData), nil
}

func GetUserByEmail(ctx context.Context, d *data.Data, email string) (*entity.User, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/contact/v3/users?emails=%s&user_id_type=open_id", email)

	result, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	items, ok := data["items"].([]interface{})
	if !ok || len(items) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	userData, ok := items[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid user data")
	}

	return parseUser(userData), nil
}

func GetUserByPhone(ctx context.Context, d *data.Data, phone string) (*entity.User, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/contact/v3/users?phones=%s&user_id_type=open_id", phone)

	result, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	items, ok := data["items"].([]interface{})
	if !ok || len(items) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	userData, ok := items[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid user data")
	}

	return parseUser(userData), nil
}

func CreateUser(ctx context.Context, d *data.Data, req interface{}) (*entity.User, error) {
	client := GetClient(d)

	url := "https://open.feishu.cn/open-apis/contact/v3/users?user_id_type=open_id"

	result, err := client.DoRequest(ctx, "POST", url, req)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	userData, ok := data["user"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid user data")
	}

	return parseUser(userData), nil
}

func UpdateUser(ctx context.Context, d *data.Data, req interface{}) (*entity.User, error) {
	client := GetClient(d)

	url := "https://open.feishu.cn/open-apis/contact/v3/users?user_id_type=open_id"

	result, err := client.DoRequest(ctx, "PATCH", url, req)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	userData, ok := data["user"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid user data")
	}

	return parseUser(userData), nil
}

func DeleteUser(ctx context.Context, d *data.Data, userID, userIDType string) error {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/contact/v3/users/%s?user_id_type=%s", userID, userIDType)

	_, err := client.DoRequest(ctx, "DELETE", url, nil)
	return err
}

func GetUserDepartments(ctx context.Context, d *data.Data, userID, userIDType string) ([]*entity.Department, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/contact/v3/users/%s/departments?user_id_type=%s", userID, userIDType)

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

	var departments []*entity.Department
	for _, item := range items {
		if deptData, ok := item.(map[string]interface{}); ok {
			departments = append(departments, parseDepartment(deptData))
		}
	}

	return departments, nil
}

func parseUser(data map[string]interface{}) *entity.User {
	user := &entity.User{}

	if v, ok := data["user_id"].(string); ok {
		user.ID = v
	}
	if v, ok := data["union_id"].(string); ok {
		user.UnionID = v
	}
	if v, ok := data["open_id"].(string); ok {
		user.OpenID = v
	}
	if v, ok := data["name"].(string); ok {
		user.Name = v
	}
	if v, ok := data["en_name"].(string); ok {
		user.EnName = v
	}
	if v, ok := data["email"].(string); ok {
		user.Email = v
	}
	if v, ok := data["phone"].(string); ok {
		user.Phone = v
	}
	if v, ok := data["avatar"]["avatar_72"].(string); ok {
		user.AvatarURL = v
	}
	if v, ok := data["status"].(float64); ok {
		user.IsActivated = v == 1
	}

	return user
}

func parseDepartment(data map[string]interface{}) *entity.Department {
	dept := &entity.Department{}

	if v, ok := data["department_id"].(string); ok {
		dept.DepartmentID = v
	}
	if v, ok := data["name"].(string); ok {
		dept.Name = v
	}
	if v, ok := data["name_en"].(string); ok {
		dept.NameEn = v
	}
	if v, ok := data["parent_id"].(string); ok {
		dept.ParentID = v
	}
	if v, ok := data["order"].(float64); ok {
		dept.Order = int(v)
	}

	return dept
}