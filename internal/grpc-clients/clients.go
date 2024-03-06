package clients

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	itemgrpc "go-server/internal/clients/item/grpc"
	"go-server/internal/config"
)

func CreateItemClient(ctx context.Context, cfg *config.Config) (*itemgrpc.Client, error) { // составить на русском языке что делает этот метод
	loggerSlog := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	itemClient, err := itemgrpc.New(ctx, loggerSlog, cfg.Clients.Item.Address, cfg.Clients.Item.Timeout, cfg.Clients.Item.RetriesCount)
	if err != nil {
		return nil, fmt.Errorf("failed to create item client: %w", err)
	}
	return itemClient, nil
}
