package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	itemv1 "github.com/tolseone/protos/gen/go/item"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	model "go-server/internal/models"

)

type Client struct {
	api itemv1.ItemServiceClient
	log *slog.Logger
}

func New(ctx context.Context, log *slog.Logger, addr string, timeout time.Duration, retriesCount int) (*Client, error) {
	const op = "grpc.New"

	// Options for interceptor grpcretry
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	// Options for interceptor grpclog
	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	// Create connection with gRPC-server ITEM for client
	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Create gRPC-client ITEM
	grpcClient := itemv1.NewItemServiceClient(cc)

	return &Client{
		api: grpcClient,
	}, nil
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (c *Client) CreateItem(ctx context.Context, name string, rarity string, description string) (uuid.UUID, error) {
	const op = "grpc.CreateItem"

	resp, err := c.api.CreateItem(ctx, &itemv1.CreateItemRequest{
		Name:        name,
		Rarity:      rarity,
		Description: description,
	})
	if err != nil {
		c.log.Error("%s: %s", op, err)
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	ItemID, err := uuid.Parse(resp.GetItemId())
	if err != nil {
		c.log.Error("%s: %s", op, err)
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return ItemID, nil
}

func (c *Client) GetAllItems(ctx context.Context) ([]*model.Item, error) {
	const op = "grpc.GetItemList"

	resp, err := c.api.GetAllItems(ctx, &itemv1.GetAllItemsRequest{})
	if err != nil {
		return []*model.Item{}, fmt.Errorf("%s: %w", op, err)
	}

	items := resp.GetItems()

	var newItems []*model.Item

	for _, item := range items {
		itemId, err := uuid.Parse(item.GetItemId())
		if err != nil {
			c.log.Info("failed to parse item id", slog.Any("item_id", item.GetItemId))
		}

		newItem := &model.Item{
			ItemId:      itemId,
			Name:        item.GetName(),
			Rarity:      item.GetRarity(),
			Description: item.GetDescription(),
		}

		newItems = append(newItems, newItem)
	}

	return newItems, nil
}

func (c *Client) GetItem(ctx context.Context, itemId string) (*model.Item, error) {
	const op = "grpc.GetItem"

	resp, err := c.api.GetItem(ctx, &itemv1.GetItemRequest{
		ItemId: itemId,
	})
	if err != nil {
		return &model.Item{}, fmt.Errorf("%s: %w", op, err)
	}

	item := resp.GetItem()

	itemID, err := uuid.Parse(item.GetItemId())
	if err != nil {
		c.log.Info("failed to parse item id", slog.Any("item_id", item.GetItemId))
	}

	return &model.Item{
		ItemId:      itemID,
		Name:        item.GetName(),
		Rarity:      item.GetRarity(),
		Description: item.GetDescription(),
	}, nil
}

func (c *Client) DeleteItem(ctx context.Context, itemId string) (*model.Item, error) {
	const op = "grpc.DeleteItem"

	_, err := c.api.DeleteItem(ctx, &itemv1.DeleteItemRequest{
		ItemId: itemId,
	})
	if err != nil {
		return &model.Item{}, fmt.Errorf("%s: %w", op, err)
	}

	return &model.Item{}, nil
}
