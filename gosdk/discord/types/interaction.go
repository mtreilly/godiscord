package types

// InteractionType defines the type of interaction from Discord.
type InteractionType int

const (
	InteractionTypePing InteractionType = iota + 1
	InteractionTypeApplicationCommand
	InteractionTypeMessageComponent
	InteractionTypeApplicationCommandAutocomplete
	InteractionTypeModalSubmit
)

// Interaction represents an incoming interaction payload.
type Interaction struct {
	ID            string           `json:"id"`
	ApplicationID string           `json:"application_id"`
	Type          InteractionType  `json:"type"`
	Data          *InteractionData `json:"data,omitempty"`
	GuildID       string           `json:"guild_id,omitempty"`
	ChannelID     string           `json:"channel_id,omitempty"`
	Member        *Member          `json:"member,omitempty"`
	User          *User            `json:"user,omitempty"`
	Token         string           `json:"token"`
	Version       int              `json:"version"`
	Message       *Message         `json:"message,omitempty"`
	Locale        string           `json:"locale,omitempty"`
	GuildLocale   string           `json:"guild_locale,omitempty"`
}

// InteractionData contains payload-specific data (commands/components).
type InteractionData struct {
	ID            string                     `json:"id,omitempty"`
	Name          string                     `json:"name,omitempty"`
	Type          ApplicationCommandType     `json:"type,omitempty"`
	Resolved      *ResolvedData              `json:"resolved,omitempty"`
	Options       []ApplicationCommandOption `json:"options,omitempty"`
	CustomID      string                     `json:"custom_id,omitempty"`
	ComponentType ComponentType              `json:"component_type,omitempty"`
	Values        []string                   `json:"values,omitempty"`
	TargetID      string                     `json:"target_id,omitempty"`
}

// ResolvedData contains hydrated entities referenced in commands.
type ResolvedData struct {
	Users    map[string]User    `json:"users,omitempty"`
	Members  map[string]Member  `json:"members,omitempty"`
	Roles    map[string]Role    `json:"roles,omitempty"`
	Channels map[string]Channel `json:"channels,omitempty"`
	Messages map[string]Message `json:"messages,omitempty"`
}

// ApplicationCommand represents a slash command or user/message command.
type ApplicationCommand struct {
	ID                       string                     `json:"id,omitempty"`
	Type                     ApplicationCommandType     `json:"type,omitempty"`
	ApplicationID            string                     `json:"application_id,omitempty"`
	GuildID                  string                     `json:"guild_id,omitempty"`
	Name                     string                     `json:"name"`
	NameLocalizations        map[string]string          `json:"name_localizations,omitempty"`
	Description              string                     `json:"description"`
	DescriptionLocalizations map[string]string          `json:"description_localizations,omitempty"`
	Options                  []ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions *string                    `json:"default_member_permissions,omitempty"`
	DMPermission             *bool                      `json:"dm_permission,omitempty"`
	NSFW                     bool                       `json:"nsfw,omitempty"`
	Version                  string                     `json:"version,omitempty"`
	AuditLogReason           string                     `json:"-"`
}

// ApplicationCommandType enumerates command kinds.
type ApplicationCommandType int

const (
	ApplicationCommandTypeChatInput ApplicationCommandType = iota + 1
	ApplicationCommandTypeUser
	ApplicationCommandTypeMessage
)

// ApplicationCommandOption describes command options or subcommands.
type ApplicationCommandOption struct {
	Type                     ApplicationCommandOptionType `json:"type"`
	Name                     string                       `json:"name"`
	NameLocalizations        map[string]string            `json:"name_localizations,omitempty"`
	Description              string                       `json:"description"`
	DescriptionLocalizations map[string]string            `json:"description_localizations,omitempty"`
	Required                 bool                         `json:"required,omitempty"`
	Choices                  []ApplicationCommandChoice   `json:"choices,omitempty"`
	Options                  []ApplicationCommandOption   `json:"options,omitempty"`
	ChannelTypes             []ChannelType                `json:"channel_types,omitempty"`
	MinValue                 *float64                     `json:"min_value,omitempty"`
	MaxValue                 *float64                     `json:"max_value,omitempty"`
	MinLength                *int                         `json:"min_length,omitempty"`
	MaxLength                *int                         `json:"max_length,omitempty"`
	Autocomplete             bool                         `json:"autocomplete,omitempty"`
}

// ApplicationCommandOptionType enumerates option types.
type ApplicationCommandOptionType int

const (
	CommandOptionSubCommand ApplicationCommandOptionType = iota + 1
	CommandOptionSubCommandGroup
	CommandOptionString
	CommandOptionInteger
	CommandOptionBoolean
	CommandOptionUser
	CommandOptionChannel
	CommandOptionRole
	CommandOptionMentionable
	CommandOptionNumber
	CommandOptionAttachment
)

// ApplicationCommandChoice represents an option choice.
type ApplicationCommandChoice struct {
	Name              string            `json:"name"`
	NameLocalizations map[string]string `json:"name_localizations,omitempty"`
	Value             interface{}       `json:"value"`
}

// ComponentType enumerates message component types.
type ComponentType int

const (
	ComponentTypeActionRow  ComponentType = 1
	ComponentTypeButton     ComponentType = 2
	ComponentTypeSelectMenu ComponentType = 3
	ComponentTypeTextInput  ComponentType = 4
)

// Validate ensures interactions are well-formed.
func (i *Interaction) Validate() error {
	if i == nil {
		return &ValidationError{Field: "interaction", Message: "interaction is required"}
	}
	if i.ID == "" {
		return &ValidationError{Field: "interaction.id", Message: "interaction ID is required"}
	}
	if i.Token == "" {
		return &ValidationError{Field: "interaction.token", Message: "interaction token is required"}
	}
	return nil
}

// Validate ensures command definitions satisfy Discord constraints.
func (c *ApplicationCommand) Validate() error {
	if c == nil {
		return &ValidationError{Field: "command", Message: "command is required"}
	}
	if len(c.Name) < 1 || len(c.Name) > 32 {
		return &ValidationError{Field: "name", Message: "command name must be 1-32 characters"}
	}
	if len(c.Description) > 100 {
		return &ValidationError{Field: "description", Message: "description must be <=100 characters"}
	}
	for _, opt := range c.Options {
		if err := opt.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate ensures option details respect Discord limits.
func (o *ApplicationCommandOption) Validate() error {
	if len(o.Name) < 1 || len(o.Name) > 32 {
		return &ValidationError{Field: "option.name", Message: "option name must be 1-32 characters"}
	}
	if len(o.Description) < 1 || len(o.Description) > 100 {
		return &ValidationError{Field: "option.description", Message: "description must be 1-100 characters"}
	}
	for _, opt := range o.Options {
		if err := opt.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// InteractionResponseType enumerates response types.
type InteractionResponseType int

const (
	InteractionResponsePong                             InteractionResponseType = 1
	InteractionResponseChannelMessageWithSource         InteractionResponseType = 4
	InteractionResponseDeferredChannelMessageWithSource InteractionResponseType = 5
	InteractionResponseDeferredUpdateMessage            InteractionResponseType = 6
	InteractionResponseUpdateMessage                    InteractionResponseType = 7
	InteractionResponseAutocompleteResult               InteractionResponseType = 8
	InteractionResponseModal                            InteractionResponseType = 9
)

// InteractionResponse represents the payload for responding to interactions.
type InteractionResponse struct {
	Type InteractionResponseType                    `json:"type"`
	Data *InteractionApplicationCommandCallbackData `json:"data,omitempty"`
}

// InteractionApplicationCommandCallbackData describes response data.
type InteractionApplicationCommandCallbackData struct {
	TTS             bool               `json:"tts,omitempty"`
	Content         string             `json:"content,omitempty"`
	Embeds          []Embed            `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions   `json:"allowed_mentions,omitempty"`
	Flags           int                `json:"flags,omitempty"`
	Components      []MessageComponent `json:"components,omitempty"`
	Attachments     []Attachment       `json:"attachments,omitempty"`
}

// MessageComponent represents a generic component.
type MessageComponent struct {
	Type        ComponentType      `json:"type"`
	CustomID    string             `json:"custom_id,omitempty"`
	Disabled    bool               `json:"disabled,omitempty"`
	Style       int                `json:"style,omitempty"`
	Label       string             `json:"label,omitempty"`
	Emoji       *Emoji             `json:"emoji,omitempty"`
	URL         string             `json:"url,omitempty"`
	Options     []SelectOption     `json:"options,omitempty"`
	Placeholder string             `json:"placeholder,omitempty"`
	MinValues   int                `json:"min_values,omitempty"`
	MaxValues   int                `json:"max_values,omitempty"`
	Components  []MessageComponent `json:"components,omitempty"`
}

// SelectOption represents select menu options.
type SelectOption struct {
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Default     bool   `json:"default,omitempty"`
	Emoji       *Emoji `json:"emoji,omitempty"`
}

// AllowedMentions controls mention parsing in responses.
type AllowedMentions struct {
	Parse       []string `json:"parse,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Users       []string `json:"users,omitempty"`
	RepliedUser bool     `json:"replied_user,omitempty"`
}
