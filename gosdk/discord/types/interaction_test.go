package types

import "testing"

func TestInteractionValidate(t *testing.T) {
	inter := &Interaction{ID: "123", Token: "token"}
	if err := inter.Validate(); err != nil {
		t.Fatalf("expected valid interaction: %v", err)
	}

	inter.ID = ""
	if err := inter.Validate(); err == nil {
		t.Fatalf("expected error for missing ID")
	}
}

func TestApplicationCommandValidate(t *testing.T) {
	cmd := &ApplicationCommand{Name: "hello", Description: "desc"}
	if err := cmd.Validate(); err != nil {
		t.Fatalf("expected valid command: %v", err)
	}

	cmd.Name = ""
	if err := cmd.Validate(); err == nil {
		t.Fatalf("expected error for empty name")
	}
}

func TestApplicationCommandOptionValidate(t *testing.T) {
	opt := ApplicationCommandOption{
		Name:        "option",
		Description: "desc",
	}
	if err := opt.Validate(); err != nil {
		t.Fatalf("expected valid option: %v", err)
	}

	opt.Name = ""
	if err := opt.Validate(); err == nil {
		t.Fatalf("expected error for name")
	}
}
