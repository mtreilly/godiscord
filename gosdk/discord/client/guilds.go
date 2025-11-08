package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/yourusername/agent-discord/gosdk/discord/types"
)

// Guilds exposes guild-related REST helpers.
type Guilds struct {
	client *Client
}

// Guilds returns a guild service bound to the Client.
func (c *Client) Guilds() *Guilds {
	return &Guilds{client: c}
}

// GetGuild retrieves a guild by ID (optionally including approximate counts).
func (g *Guilds) GetGuild(ctx context.Context, guildID string, withCounts bool) (*types.Guild, error) {
	if err := validateID("guildID", guildID); err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/guilds/%s", guildID)
	if withCounts {
		path += "?with_counts=true"
	}

	var guild types.Guild
	if err := g.client.Get(ctx, path, &guild); err != nil {
		return nil, err
	}
	return &guild, nil
}

// GetGuildPreview fetches a preview for a discoverable guild.
func (g *Guilds) GetGuildPreview(ctx context.Context, guildID string) (*types.GuildPreview, error) {
	if err := validateID("guildID", guildID); err != nil {
		return nil, err
	}

	var preview types.GuildPreview
	if err := g.client.Get(ctx, fmt.Sprintf("/guilds/%s/preview", guildID), &preview); err != nil {
		return nil, err
	}
	return &preview, nil
}

// ModifyGuild updates mutable guild properties.
func (g *Guilds) ModifyGuild(ctx context.Context, guildID string, params *types.GuildModifyParams) (*types.Guild, error) {
	if err := validateID("guildID", guildID); err != nil {
		return nil, err
	}
	if err := params.Validate(); err != nil {
		return nil, err
	}

	headers := http.Header{}
	if params.AuditLogReason != "" {
		headers.Set("X-Audit-Log-Reason", url.QueryEscape(params.AuditLogReason))
	}

	var guild types.Guild
	if err := g.client.do(ctx, http.MethodPatch, fmt.Sprintf("/guilds/%s", guildID), params, &guild, headers); err != nil {
		return nil, err
	}
	return &guild, nil
}

// GetGuildChannels lists a guild's channels.
func (g *Guilds) GetGuildChannels(ctx context.Context, guildID string) ([]*types.Channel, error) {
	if err := validateID("guildID", guildID); err != nil {
		return nil, err
	}

	var channels []*types.Channel
	if err := g.client.Get(ctx, fmt.Sprintf("/guilds/%s/channels", guildID), &channels); err != nil {
		return nil, err
	}
	return channels, nil
}

// CreateGuildChannel creates a channel within a guild.
func (g *Guilds) CreateGuildChannel(ctx context.Context, guildID string, params *types.ChannelCreateParams) (*types.Channel, error) {
	if err := validateID("guildID", guildID); err != nil {
		return nil, err
	}
	if params == nil {
		return nil, &types.ValidationError{Field: "params", Message: "channel params required"}
	}
	if err := params.Validate(); err != nil {
		return nil, err
	}

	headers := http.Header{}
	if params.AuditLogReason != "" {
		headers.Set("X-Audit-Log-Reason", url.QueryEscape(params.AuditLogReason))
	}

	var channel types.Channel
	if err := g.client.do(ctx, http.MethodPost, fmt.Sprintf("/guilds/%s/channels", guildID), params, &channel, headers); err != nil {
		return nil, err
	}
	return &channel, nil
}
