package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sexy_backend/internal/model"
	"sexy_backend/pkg/supabase"
)

type UserService struct {
	supabaseClient *supabase.Client
}

func NewUserService(client *supabase.Client) *UserService {
	return &UserService{
		supabaseClient: client,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
	var result []model.User
	data := map[string]interface{}{
		"email":    user.Email,
		"username": user.Username,
	}

	resp, _, err := s.supabaseClient.Client.From("users").
		Insert(data, false, "", "", "").
		Execute()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("no user created")
	}
	*user = result[0]
	return nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*model.User, error) {
	var result []model.User

	resp, _, err := s.supabaseClient.Client.From("users").
		Select("*", "", false).
		Eq("id", id).
		Execute()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return &result[0], nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, user *model.User) error {
	var result []model.User
	data := map[string]interface{}{
		"email":    user.Email,
		"username": user.Username,
	}

	resp, _, err := s.supabaseClient.Client.From("users").
		Update(data, "", "").
		Eq("id", id).
		Execute()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("no user updated")
	}
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	var result []model.User

	resp, _, err := s.supabaseClient.Client.From("users").
		Delete("", "").
		Eq("id", id).
		Execute()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("no user deleted")
	}
	return nil
}
