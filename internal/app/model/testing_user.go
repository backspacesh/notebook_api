package model

import "testing"

func TestUser(t *testing.T) *User {
	t.Helper()

	return &User{
		Name: "testUser",
		Email: "test@example.com",
		Password: "password",
	}
}
