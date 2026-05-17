package handler

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/feishu/feishu-office-suite/api/feishu/v1"
	"github.com/feishu/feishu-office-suite/internal/data"
	"github.com/feishu/feishu-office-suite/pkg/feishu"
)

type UserHandler struct {
	v1.UnimplementedFeishuUserServer
	data *data.Data
}

func NewUserHandler(data *data.Data) *UserHandler {
	return &UserHandler{data: data}
}

func (h *UserHandler) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	userID := req.UserId
	userIDType := req.UserIdType

	if userIDType == "" {
		userIDType = "open_id"
	}

	user, err := h.data.EntClient().User.Query().
		Where(userIDEQ(userID)).
		First(ctx)

	if err != nil {
		feishuUser, feishuErr := feishu.GetUser(ctx, h.data, userID, userIDType)
		if feishuErr != nil {
			return nil, feishuErr
		}
		user = feishuUser
	}

	return &v1.GetUserResponse{
		User: convertUserToProto(user),
	}, nil
}

func (h *UserHandler) GetUserByEmail(ctx context.Context, req *v1.GetUserByEmailRequest) (*v1.GetUserByEmailResponse, error) {
	user, err := h.data.EntClient().User.Query().
		Where(emailEQ(req.Email)).
		First(ctx)

	if err != nil {
		feishuUser, feishuErr := feishu.GetUserByEmail(ctx, h.data, req.Email)
		if feishuErr != nil {
			return nil, feishuErr
		}
		user = feishuUser
	}

	return &v1.GetUserByEmailResponse{
		User: convertUserToProto(user),
	}, nil
}

func (h *UserHandler) GetUserByPhone(ctx context.Context, req *v1.GetUserByPhoneRequest) (*v1.GetUserByPhoneResponse, error) {
	user, err := h.data.EntClient().User.Query().
		Where(phoneEQ(req.Phone)).
		First(ctx)

	if err != nil {
		feishuUser, feishuErr := feishu.GetUserByPhone(ctx, h.data, req.Phone)
		if feishuErr != nil {
			return nil, feishuErr
		}
		user = feishuUser
	}

	return &v1.GetUserByPhoneResponse{
		User: convertUserToProto(user),
	}, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, req *v1.ListUsersRequest) (*v1.ListUsersResponse, error) {
	query := h.data.EntClient().User.Query()

	if req.DepartmentId != "" {
		query = query.Where(departmentIDEQ(req.DepartmentId))
	}

	if !req.IncludeResigned {
		query = query.Where(statusNEQ("resigned"))
	}

	users, err := query.Limit(int(req.PageSize)).Offset(0).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	var protoUsers []*v1.User
	for _, user := range users {
		protoUsers = append(protoUsers, convertUserToProto(user))
	}

	return &v1.ListUsersResponse{
		Users:         protoUsers,
		TotalCount:    int32(len(protoUsers)),
		HasMore:       false,
		NextPageToken: "",
	}, nil
}

func (h *UserHandler) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	feishuUser, err := feishu.CreateUser(ctx, h.data, req)
	if err != nil {
		return nil, err
	}

	return &v1.CreateUserResponse{
		User: convertUserToProto(feishuUser),
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	userID := req.UserId
	userIDType := req.UserIdType

	if userIDType == "" {
		userIDType = "open_id"
	}

	feishuUser, err := feishu.UpdateUser(ctx, h.data, req)
	if err != nil {
		return nil, err
	}

	_ = userID

	return &v1.UpdateUserResponse{
		User: convertUserToProto(feishuUser),
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*emptypb.Empty, error) {
	userID := req.UserId
	userIDType := req.UserIdType

	if userIDType == "" {
		userIDType = "open_id"
	}

	_, err := feishu.DeleteUser(ctx, h.data, userID, userIDType)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) GetDepartments(ctx context.Context, req *v1.GetDepartmentsRequest) (*v1.GetDepartmentsResponse, error) {
	departments, err := h.data.EntClient().Department.Query().
		Where(parentIDEQ(req.ParentId)).
		Limit(int(req.PageSize)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get departments: %w", err)
	}

	var protoDepts []*v1.Department
	for _, dept := range departments {
		protoDepts = append(protoDepts, convertDepartmentToProto(dept))
	}

	return &v1.GetDepartmentsResponse{
		Departments:  protoDepts,
		TotalCount:   int32(len(protoDepts)),
		HasMore:      false,
		NextPageToken: "",
	}, nil
}

func (h *UserHandler) GetUserDepartments(ctx context.Context, req *v1.GetUserDepartmentsRequest) (*v1.GetUserDepartmentsResponse, error) {
	userID := req.UserId
	userIDType := req.UserIdType

	if userIDType == "" {
		userIDType = "open_id"
	}

	departments, err := feishu.GetUserDepartments(ctx, h.data, userID, userIDType)
	if err != nil {
		return nil, err
	}

	var protoDepts []*v1.Department
	for _, dept := range departments {
		protoDepts = append(protoDepts, convertDepartmentToProto(dept))
	}

	return &v1.GetUserDepartmentsResponse{
		Departments: protoDepts,
	}, nil
}

type userEntity struct{}

func (u *userEntity) IDEQ(id string) func(*ent.UserWhereInput) {
	return nil
}