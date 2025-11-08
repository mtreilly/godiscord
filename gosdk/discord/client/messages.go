package client

import (
	"context"
	"fmt"

	"github.com/yourusername/agent-discord/gosdk/discord/types"
)

// MessageService provides helpers for channel message operations.
type MessageService struct {
	client *Client
}

// Messages exposes message helpers bound to the client.
func (c *Client) Messages() *MessageService {
	return &MessageService{client: c}
}

// CreateMessage sends a message to a channel.
func (m *MessageService) CreateMessage(ctx context.Context, channelID string, params *types.MessageCreateParams) (*types.Message, error) {
	if err := validateID("channelID", channelID); err != nil {
		return nil, err
	}
	if params == nil {
		return nil, &types.ValidationError{Field: "params", Message: "message create params required"}
	}

	var msg types.Message
	if err := m.client.Post(ctx, fmt.Sprintf("/channels/%s/messages", channelID), params, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// GetMessage fetches a single message.
func (m *MessageService) GetMessage(ctx context.Context, channelID, messageID string) (*types.Message, error) {
	if err := validateID("channelID", channelID); err != nil {
		return nil, err
	}
	if err := validateID("messageID", messageID); err != nil {
		return nil, err
	}

	var msg types.Message
	if err := m.client.Get(ctx, fmt.Sprintf("/channels/%s/messages/%s", channelID, messageID), &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// EditMessage updates message content/embeds.
func (m *MessageService) EditMessage(ctx context.Context, channelID, messageID string, params *types.MessageEditParams) (*types.Message, error) {
	if err := validateID("channelID", channelID); err != nil {
		return nil, err
	}
	if err := validateID("messageID", messageID); err != nil {
		return nil, err
	}
	if params == nil {
		return nil, &types.ValidationError{Field: "params", Message: "message edit params required"}
	}

	var msg types.Message
	if err := m.client.Patch(ctx, fmt.Sprintf("/channels/%s/messages/%s", channelID, messageID), params, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// DeleteMessage removes a single message.
func (m *MessageService) DeleteMessage(ctx context.Context, channelID, messageID string) error {
	if err := validateID("channelID", channelID); err != nil {
		return err
	}
	if err := validateID("messageID", messageID); err != nil {
		return err
	}
	return m.client.Delete(ctx, fmt.Sprintf("/channels/%s/messages/%s", channelID, messageID))
}

// BulkDeleteMessages deletes multiple messages in one call.
func (m *MessageService) BulkDeleteMessages(ctx context.Context, channelID string, messageIDs []string) error {
	if err := validateID("channelID", channelID); err != nil {
		return err
	}
	if len(messageIDs) == 0 {
		return &types.ValidationError{Field: "messages", Message: "at least one message ID required"}
	}
	if len(messageIDs) > 100 {
		return &types.ValidationError{Field: "messages", Message: "maximum 100 messages per bulk delete"}
	}

	payload := struct {
		Messages []string `json:"messages"`
	}{
		Messages: messageIDs,
	}
	return m.client.Post(ctx, fmt.Sprintf("/channels/%s/messages/bulk-delete", channelID), payload, nil)
}
