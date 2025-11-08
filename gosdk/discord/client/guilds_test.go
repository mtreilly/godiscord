package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/agent-discord/gosdk/discord/types"
)

func TestGuildsGetGuild(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "with_counts=true" {
			t.Fatalf("expected with_counts query, got %s", r.URL.RawQuery)
		}
		json.NewEncoder(w).Encode(types.Guild{ID: "1", Name: "Guild"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	guild, err := client.Guilds().GetGuild(context.Background(), "1", true)
	if err != nil {
		t.Fatalf("GetGuild error: %v", err)
	}
	if guild.ID != "1" {
		t.Fatalf("expected guild ID 1, got %s", guild.ID)
	}
}

func TestGuildsModifyGuild(t *testing.T) {
	var reason string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reason = r.Header.Get("X-Audit-Log-Reason")
		json.NewEncoder(w).Encode(types.Guild{ID: "1", Name: "Updated"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	params := &types.GuildModifyParams{Name: "Updated", AuditLogReason: "cleanup"}
	guild, err := client.Guilds().ModifyGuild(context.Background(), "1", params)
	if err != nil {
		t.Fatalf("ModifyGuild error: %v", err)
	}
	if guild.Name != "Updated" {
		t.Fatalf("expected updated name")
	}
	if reason == "" {
		t.Fatalf("expected audit log reason header")
	}
}

func TestGuildsCreateChannel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(types.Channel{ID: "10", Name: "new"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	channel, err := client.Guilds().CreateGuildChannel(context.Background(), "1", &types.ChannelCreateParams{Name: "new"})
	if err != nil {
		t.Fatalf("CreateGuildChannel error: %v", err)
	}
	if channel.ID != "10" {
		t.Fatalf("expected channel ID 10")
	}
}
