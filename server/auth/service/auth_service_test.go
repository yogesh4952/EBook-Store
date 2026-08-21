package service

import (
	"context"
	"errors"
	"testing"

	"github.com/yogesh4952/ebookstore/user/models"
)

type mockUserLookup struct {
	user *models.User
	err  error
}

func (m mockUserLookup) FindByEmail(email string) (*models.User, error) {
	return m.user, m.err
}

func TestLogin_UserFound(t *testing.T) {
	mock := mockUserLookup{user: &models.User{Email: "a@b.com"}}
	s := &authService{userStore: mock}

	token, err := s.Login(context.Background(), "a@b.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token != "" {
		t.Fatalf("expected empty placeholder token, got %q", token)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	mock := mockUserLookup{err: errors.New("record not found")}
	s := &authService{userStore: mock}

	_, err := s.Login(context.Background(), "missing@b.com")
	if err == nil {
		t.Fatal("expected error for missing user, got nil")
	}
}
