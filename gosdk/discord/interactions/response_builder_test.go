package interactions

import (
	"testing"

	"github.com/yourusername/agent-discord/gosdk/discord/types"
)

func TestResponseBuilder_Message(t *testing.T) {
	row := types.MessageComponent{
		Type: types.ComponentTypeActionRow,
		Components: []types.MessageComponent{
			{Type: types.ComponentTypeButton, CustomID: "run", Label: "Run"},
		},
	}

	builder := NewMessageResponse("hello").
		SetContent("updated").
		SetTTS(true).
		SetAllowedMentions(&types.AllowedMentions{RepliedUser: true}).
		AddEmbed(types.Embed{Title: "Embed"}).
		AddComponentRow(row).
		SetEphemeral(true)

	resp, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	if resp.Type != types.InteractionResponseChannelMessageWithSource {
		t.Fatalf("unexpected response type %d", resp.Type)
	}
	if resp.Data == nil || resp.Data.Content != "updated" {
		t.Fatalf("expected updated content, got %+v", resp.Data)
	}
	if resp.Data.Flags&interactionResponseFlagEphemeral == 0 {
		t.Fatalf("expected ephemeral flag to be set")
	}
	if len(resp.Data.Components) != 1 {
		t.Fatalf("expected component row")
	}
}

func TestResponseBuilder_Modal(t *testing.T) {
	modal := NewModalResponse("modal", "Title")
	modal.SetModalComponents(types.MessageComponent{
		Type: types.ComponentTypeActionRow,
		Components: []types.MessageComponent{
			{Type: types.ComponentTypeTextInput, CustomID: "input", Label: "Label"},
		},
	})

	resp, err := modal.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if resp.Type != types.InteractionResponseModal {
		t.Fatalf("expected modal type, got %d", resp.Type)
	}
}

func TestResponseBuilder_ComponentValidation(t *testing.T) {
	builder := NewMessageResponse("hello")
	builder.AddComponentRow(types.MessageComponent{
		Type: types.ComponentTypeButton,
	})

	if _, err := builder.Build(); err == nil {
		t.Fatalf("expected error for non action row")
	}

	modal := NewModalResponse("modal", "Title")
	modal.SetModalComponents(types.MessageComponent{
		Type: types.ComponentTypeActionRow,
		Components: []types.MessageComponent{
			{Type: types.ComponentTypeButton, CustomID: "btn"},
		},
	})
	if _, err := modal.Build(); err == nil {
		t.Fatalf("expected error for modal child that is not text input")
	}
}
