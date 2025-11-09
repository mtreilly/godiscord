package gateway

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/yourusername/agent-discord/gosdk/logger"
)

// Shard represents a gateway shard (ID + total + client).
type Shard struct {
	id          int
	totalShards int
	client      *Client
}

// ShardManagerOption configures the shard manager.
type ShardManagerOption func(*ShardManager)

// WithShardLogger overrides the logger.
func WithShardLogger(l *logger.Logger) ShardManagerOption {
	return func(sm *ShardManager) {
		if l != nil {
			sm.logger = l
		}
	}
}

// WithShardDispatcher uses a custom dispatcher.
func WithShardDispatcher(d *Dispatcher) ShardManagerOption {
	return func(sm *ShardManager) {
		if d != nil {
			sm.dispatcher = d
		}
	}
}

// WithShardConnectionOptions augments the underlying connections.
func WithShardConnectionOptions(opts ...ConnectionOption) ShardManagerOption {
	return func(sm *ShardManager) {
		sm.connectionOpts = append(sm.connectionOpts, opts...)
	}
}

// ShardManager orchestrates multiple gateway shards.
type ShardManager struct {
	token          string
	intents        int
	shardCount     int
	logger         *logger.Logger
	dispatcher     *Dispatcher
	connectionOpts []ConnectionOption

	shards []*Shard
	mu     sync.Mutex
}

// NewShardManager constructs a shard manager.
func NewShardManager(token string, shardCount int, intents int, opts ...ShardManagerOption) *ShardManager {
	sm := &ShardManager{
		token:      token,
		intents:    intents,
		shardCount: shardCount,
		logger:     logger.Default(),
		dispatcher: NewDispatcher(),
	}
	for _, opt := range opts {
		opt(sm)
	}
	return sm
}

// Connect initializes and starts all shard clients.
func (sm *ShardManager) Connect(ctx context.Context) error {
	sm.mu.Lock()
	if len(sm.shards) > 0 {
		sm.mu.Unlock()
		return errors.New("shard manager already connected")
	}
	sm.mu.Unlock()

	for id := 0; id < sm.shardCount; id++ {
		connOpts := append([]ConnectionOption{}, sm.connectionOpts...)
		shardURL := fmt.Sprintf("%s&shard=%d,%d", defaultGatewayURL, id, sm.shardCount)
		connOpts = append(connOpts, WithGatewayURL(shardURL))

		client, err := NewClient(sm.token, sm.intents,
			WithDispatcher(sm.dispatcher),
			WithGatewayLogger(sm.logger),
			WithConnectionOptions(connOpts...),
		)
		if err != nil {
			return fmt.Errorf("init shard %d: %w", id, err)
		}
		if err := client.Connect(ctx); err != nil {
			return fmt.Errorf("connect shard %d: %w", id, err)
		}

		sm.mu.Lock()
		sm.shards = append(sm.shards, &Shard{id: id, totalShards: sm.shardCount, client: client})
		sm.mu.Unlock()
	}
	return nil
}

// Disconnect closes all shard clients.
func (sm *ShardManager) Disconnect() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	var errs []error
	for _, shard := range sm.shards {
		if err := shard.client.Disconnect(); err != nil {
			errs = append(errs, fmt.Errorf("shard %d: %w", shard.id, err))
		}
	}
	sm.shards = nil
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

// On registers an event handler across all shards.
func (sm *ShardManager) On(eventType string, handler EventHandler) {
	sm.dispatcher.On(eventType, handler)
}

// OnMessageCreate registers a MESSAGE_CREATE handler.
func (sm *ShardManager) OnMessageCreate(handler func(context.Context, *MessageCreateEvent) error) {
	sm.dispatcher.OnMessageCreate(handler)
}

// OnMessageUpdate registers a MESSAGE_UPDATE handler.
func (sm *ShardManager) OnMessageUpdate(handler func(context.Context, *MessageUpdateEvent) error) {
	sm.dispatcher.OnMessageUpdate(handler)
}

// OnInteraction registers an INTERACTION_CREATE handler.
func (sm *ShardManager) OnInteraction(handler func(context.Context, *InteractionCreateEvent) error) {
	sm.dispatcher.OnInteraction(handler)
}

// Broadcast sends the payload to every shard.
func (sm *ShardManager) Broadcast(ctx context.Context, payload *Payload) error {
	sm.mu.Lock()
	shards := append([]*Shard(nil), sm.shards...)
	sm.mu.Unlock()
	var errs []error
	for _, shard := range shards {
		if err := shard.client.Send(ctx, payload); err != nil {
			errs = append(errs, fmt.Errorf("shard %d: %w", shard.id, err))
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}
