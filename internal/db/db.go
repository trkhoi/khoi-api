package db

import (
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/trkhoi/khoi-api/config"
)

type DB struct {
	Store IStore
}

func New(cfg config.View, logger *logrus.Entry) *DB {
	db := NewPostgresStore(cfg, logger)
	return &DB{
		Store: db,
	}
}
