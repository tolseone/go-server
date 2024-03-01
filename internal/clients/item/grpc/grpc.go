package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	itemv1 "github.com/tolseone/protos/gen/go/item"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

)

type Client struct {
	api itemv1.ItemServiceClient
	log *slog.Logger
}

func New(ctx context.Context, log *slog.Logger, addr string, timeout time.Duration, retriesCount int) (*Client, error) {
	const op = "grpc.New"

	// Опции для интерсептора grpcretry
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	// Опции для интерсептора grpclog
	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	// Создаём соединение с gRPC-сервером AUTH для клиента
	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Создаём gRPC-клиент SSO/Auth
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

func (c *Client) CreateItem(ctx context.Context, name string, rarity string, description string) (string, error) {
	const op = "grpc.CreateItem"

	resp, err := c.api.CreateItem(ctx, &itemv1.CreateItemRequest{
		Name:        name,
		Rarity:      rarity,
		Description: description,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.ItemId, nil
}

func (c *Client) GetItem(ctx context.Context, itemID string) ()
