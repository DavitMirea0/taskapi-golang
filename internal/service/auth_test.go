package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"taskapi/internal/domain"
)

type fakeUsers struct {
	byEmail map[string]domain.User
}

func (f *fakeUsers) Create(_ context.Context, email, hash string) (domain.User, error) {
	if _, ok := f.byEmail[email]; ok {
		return domain.User{}, domain.ErrEmailTaken
	}
	u := domain.User{ID: int64(len(f.byEmail) + 1), Email: email, PasswordHash: hash}
	f.byEmail[email] = u
	return u, nil
}

func (f *fakeUsers) GetByEmail(_ context.Context, email string) (domain.User, error) {
	u, ok := f.byEmail[email]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return u, nil
}

func TestAuth_RegisterLogin(t *testing.T) {
	ctx := context.Background()
	svc := NewAuthService(&fakeUsers{byEmail: map[string]domain.User{}}, []byte("secret"), time.Hour)

	if _, err := svc.Register(ctx, "A@Example.com", "password123"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if _, err := svc.Register(ctx, "a@example.com", "password123"); !errors.Is(err, domain.ErrEmailTaken) {
		t.Fatalf("want ErrEmailTaken, got %v", err)
	}
	if _, err := svc.Register(ctx, "b@example.com", "short"); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("want ErrValidation, got %v", err)
	}

	token, err := svc.Login(ctx, "a@example.com", "password123")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	id, err := svc.ParseToken(token)
	if err != nil || id != 1 {
		t.Fatalf("parse token: id=%d err=%v", id, err)
	}

	if _, err := svc.Login(ctx, "a@example.com", "wrong-password"); !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("want ErrInvalidCredentials, got %v", err)
	}
	if _, err := svc.ParseToken("garbage"); !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("want ErrInvalidCredentials for bad token, got %v", err)
	}
}
