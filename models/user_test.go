package models

import (
	"testing"
)

func TestUserNewInstance(t *testing.T) {
	user := &User{
		Username:  "john_doe",
		FirstName: "John",
	}

	if user.Username != "john_doe" {
		t.Errorf("Username = %s, want john_doe", user.Username)
	}
	if user.FirstName != "John" {
		t.Errorf("FirstName = %s, want John", user.FirstName)
	}
}

func TestUserWithEmptyFields(t *testing.T) {
	user := &User{}

	if user.Username != "" {
		t.Errorf("Username = %s, want empty string", user.Username)
	}
	if user.FirstName != "" {
		t.Errorf("FirstName = %s, want empty string", user.FirstName)
	}
}

func TestUserWithSpecialCharactersUsername(t *testing.T) {
	user := &User{
		Username:  "user_123-test",
		FirstName: "Test User",
	}

	if user.Username != "user_123-test" {
		t.Errorf("Username = %s, want user_123-test", user.Username)
	}
}

func TestUserWithLongNames(t *testing.T) {
	longUsername := "a_very_long_username_with_many_characters"
	longFirstName := "Christopher Alexander"

	user := &User{
		Username:  longUsername,
		FirstName: longFirstName,
	}

	if user.Username != longUsername {
		t.Errorf("Username = %s, want %s", user.Username, longUsername)
	}
	if user.FirstName != longFirstName {
		t.Errorf("FirstName = %s, want %s", user.FirstName, longFirstName)
	}
}

func TestUserWithWhitespace(t *testing.T) {
	user := &User{
		Username:  "user@123",
		FirstName: "John Doe",
	}

	if user.Username != "user@123" {
		t.Errorf("Username = %s, want user@123", user.Username)
	}
	if user.FirstName != "John Doe" {
		t.Errorf("FirstName = %s, want John Doe", user.FirstName)
	}
}
