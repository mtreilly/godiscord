package types

// WebhookMessage represents a message to be sent via webhook
type WebhookMessage struct {
	Content         string  `json:"content,omitempty"`
	Username        string  `json:"username,omitempty"`
	AvatarURL       string  `json:"avatar_url,omitempty"`
	TTS             bool    `json:"tts,omitempty"`
	Embeds          []Embed `json:"embeds,omitempty"`
	AllowedMentions *struct {
		Parse []string `json:"parse,omitempty"`
	} `json:"allowed_mentions,omitempty"`
}

// Validate checks if the webhook message is valid
func (w *WebhookMessage) Validate() error {
	if w.Content == "" && len(w.Embeds) == 0 {
		return &ValidationError{
			Field:   "content/embeds",
			Message: "at least one of content or embeds is required",
		}
	}

	if len(w.Content) > 2000 {
		return &ValidationError{
			Field:   "content",
			Message: "content exceeds 2000 characters",
		}
	}

	if len(w.Embeds) > 10 {
		return &ValidationError{
			Field:   "embeds",
			Message: "maximum 10 embeds allowed",
		}
	}

	for i, embed := range w.Embeds {
		if err := validateEmbed(&embed); err != nil {
			return err
		}
		_ = i // silence unused variable warning for now
	}

	return nil
}

func validateEmbed(e *Embed) error {
	if len(e.Title) > 256 {
		return &ValidationError{Field: "embed.title", Message: "title exceeds 256 characters"}
	}
	if len(e.Description) > 4096 {
		return &ValidationError{Field: "embed.description", Message: "description exceeds 4096 characters"}
	}
	if len(e.Fields) > 25 {
		return &ValidationError{Field: "embed.fields", Message: "maximum 25 fields allowed"}
	}
	return nil
}
