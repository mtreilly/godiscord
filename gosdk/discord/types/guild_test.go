package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestGuildValidate(t *testing.T) {
	guild := &Guild{ID: "123", Name: "Test"}
	if err := guild.Validate(); err != nil {
		t.Fatalf("expected valid guild, got %v", err)
	}

	guild.ID = ""
	if err := guild.Validate(); err == nil {
		t.Fatalf("expected error for missing ID")
	}
}

func TestRoleValidate(t *testing.T) {
	role := &Role{Name: "Admin"}
	if err := role.Validate(); err != nil {
		t.Fatalf("expected role valid, got %v", err)
	}

	role.Name = ""
	if err := role.Validate(); err == nil {
		t.Fatalf("expected error for missing role name")
	}
}

func TestGuildJSONMarshaling(t *testing.T) {
	now := time.Now().UTC()
	guild := &Guild{
		ID:      "1",
		Name:    "Guild",
		Members: []Member{{User: &User{ID: "u1"}, JoinedAt: now}},
	}

	data, err := json.Marshal(guild)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	if len(data) == 0 || data[0] != '{' {
		t.Fatalf("unexpected JSON: %s", data)
	}
}
