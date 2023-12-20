package db

import (
	"go-server/pkg/client/postgresql"
	"go-server/pkg/logging"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}
