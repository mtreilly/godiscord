package gateway

import "encoding/json"

// OpCode defines the gateway message operation codes used by Discord.
type OpCode int

const (
	OpCodeDispatch OpCode = iota
	OpCodeHeartbeat
	OpCodeIdentify
	OpCodePresenceUpdate
	OpCodeVoiceStateUpdate
	OpCodeResume OpCode = 6
	OpCodeReconnect
	OpCodeRequestGuildMembers
	OpCodeInvalidSession
	OpCodeHello
	OpCodeHeartbeatAck
)

// Payload represents the generic envelope sent over the gateway.
type Payload struct {
	Op OpCode          `json:"op"`
	D  json.RawMessage `json:"d,omitempty"`
	S  int             `json:"s,omitempty"`
	T  string          `json:"t,omitempty"`
}

// IdentifyProperties describes the identifying metadata Discord expects.
type IdentifyProperties struct {
	OS      string `json:"os"`
	Browser string `json:"browser"`
	Device  string `json:"device"`
}

// IdentifyPayload is sent when the client first connects to the gateway.
type IdentifyPayload struct {
	Token      string             `json:"token"`
	Properties IdentifyProperties `json:"properties"`
	Compress   bool               `json:"compress,omitempty"`
	Intents    int                `json:"intents"`
	Shard      []int              `json:"shard,omitempty"`
}

// ResumePayload is used when resuming an existing session.
type ResumePayload struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int    `json:"seq"`
}
