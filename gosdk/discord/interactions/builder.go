package interactions

import (
	"github.com/yourusername/agent-discord/gosdk/discord/types"
)

// CommandBuilder builds application commands fluently.
type CommandBuilder struct {
	cmd *types.ApplicationCommand
}

// NewSlashCommand initializes a chat-input command builder.
func NewSlashCommand(name, description string) *CommandBuilder {
	return &CommandBuilder{
		cmd: &types.ApplicationCommand{
			Name:        name,
			Description: description,
			Type:        types.ApplicationCommandTypeChatInput,
		},
	}
}

// AddStringOption appends a string option to the command.
func (b *CommandBuilder) AddStringOption(name, description string, required bool) *CommandBuilder {
	b.cmd.Options = append(b.cmd.Options, types.ApplicationCommandOption{
		Type:        types.CommandOptionString,
		Name:        name,
		Description: description,
		Required:    required,
	})
	return b
}

// Build validates and returns the command.
func (b *CommandBuilder) Build() (*types.ApplicationCommand, error) {
	if err := b.cmd.Validate(); err != nil {
		return nil, err
	}
	return b.cmd, nil
}
