package interactions

import "testing"

func TestCommandBuilder(t *testing.T) {
	builder := NewSlashCommand("hello", "Description").AddStringOption("name", "Who to greet", true)
	cmd, err := builder.Build()
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	if cmd.Name != "hello" || len(cmd.Options) != 1 {
		t.Fatalf("unexpected command: %+v", cmd)
	}
}
