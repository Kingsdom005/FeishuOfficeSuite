package feishu

import (
	"context"
	"fmt"

	"github.com/feishu/feishu-office-suite/api/feishu/v1"
	"github.com/feishu/feishu-office-suite/internal/data"
	"github.com/feishu/feishu-office-suite/internal/domain/entity"
)

func GetApproval(ctx context.Context, d *data.Data, approvalID string) (*entity.Approval, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/approval/v4/approvals/%s", approvalID)

	result, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	approvalData, ok := data["approval"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid approval data")
	}

	return parseApproval(approvalData), nil
}

func ListApprovals(ctx context.Context, d *data.Data, pageSize int) ([]*entity.Approval, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/approval/v4/approvals?page_size=%d", pageSize)

	result, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	items, ok := data["list"].([]interface{})
	if !ok {
		return nil, nil
	}

	var approvals []*entity.Approval
	for _, item := range items {
		if approvalData, ok := item.(map[string]interface{}); ok {
			approvals = append(approvals, parseApproval(approvalData))
		}
	}

	return approvals, nil
}

func CreateApproval(ctx context.Context, d *data.Data, req *v1.CreateApprovalRequest) (*entity.Approval, error) {
	client := GetClient(d)

	url := "https://open.feishu.cn/open-apis/approval/v4/approvals"

	payload := map[string]interface{}{
		"name":          req.Name,
		"description":   req.Description,
		"form_content":  req.FormContent,
		"approver_ids":  req.ApproverIds,
	}

	result, err := client.DoRequest(ctx, "POST", url, payload)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	approvalData, ok := data["approval"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid approval data")
	}

	return parseApproval(approvalData), nil
}

func GetApprovalInstance(ctx context.Context, d *data.Data, instanceID string) (*entity.ApprovalInstance, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/approval/v4/approval_instances/%s", instanceID)

	result, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	instanceData, ok := data["instance"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid instance data")
	}

	return parseApprovalInstance(instanceData), nil
}

func ListApprovalInstances(ctx context.Context, d *data.Data, approvalID string, pageSize int) ([]*entity.ApprovalInstance, error) {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/approval/v4/approval_instances?approval_id=%s&page_size=%d", approvalID, pageSize)

	result, err := client.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	items, ok := data["list"].([]interface{})
	if !ok {
		return nil, nil
	}

	var instances []*entity.ApprovalInstance
	for _, item := range items {
		if instanceData, ok := item.(map[string]interface{}); ok {
			instances = append(instances, parseApprovalInstance(instanceData))
		}
	}

	return instances, nil
}

func ApproveInstance(ctx context.Context, d *data.Data, instanceID, comment string) error {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/approval/v4/approval_instances/%s/approve", instanceID)

	payload := map[string]interface{}{
		"comment": comment,
	}

	_, err := client.DoRequest(ctx, "POST", url, payload)
	return err
}

func RejectInstance(ctx context.Context, d *data.Data, instanceID, comment string) error {
	client := GetClient(d)

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/approval/v4/approval_instances/%s/reject", instanceID)

	payload := map[string]interface{}{
		"comment": comment,
	}

	_, err := client.DoRequest(ctx, "POST", url, payload)
	return err
}

func parseApproval(data map[string]interface{}) *entity.Approval {
	approval := &entity.Approval{}

	if v, ok := data["approval_id"].(string); ok {
		approval.ApprovalID = v
	}
	if v, ok := data["name"].(string); ok {
		approval.Name = v
	}
	if v, ok := data["description"].(string); ok {
		approval.Description = v
	}
	if v, ok := data["form_content"].(string); ok {
		approval.FormContent = v
	}

	return approval
}

func parseApprovalInstance(data map[string]interface{}) *entity.ApprovalInstance {
	instance := &entity.ApprovalInstance{}

	if v, ok := data["instance_id"].(string); ok {
		instance.InstanceID = v
	}
	if v, ok := data["approval_id"].(string); ok {
		instance.ApprovalID = v
	}
	if v, ok := data["title"].(string); ok {
		instance.Title = v
	}
	if v, ok := data["status"].(string); ok {
		instance.Status = v
	}
	if v, ok := data["initiator_id"].(string); ok {
		instance.InitiatorID = v
	}

	return instance
}