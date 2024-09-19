package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"github.com/trkhoi/khoi-api/config"
)

// store is implimentation of repository
type store struct {
	database     *gorm.DB
	shutdownFunc func() error
}

// NewPostgresStore postgres init by gorm
func NewPostgresStore(cfg config.View, logger *logrus.Entry) IStore {
	dbUser := cfg.GetString("DB_USER")
	dbPass := cfg.GetString("DB_PASS")
	dbHost := cfg.GetString("DB_HOST")
	dbPort := cfg.GetString("DB_PORT")
	dbName := cfg.GetString("DB_NAME")

	ds := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass,
		dbHost, dbPort, dbName,
	)
	conn, err := sql.Open("postgres", ds)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.New(
		postgres.Config{Conn: conn}),
		&gorm.Config{
			Logger: NewLogrusLogger(logger),
		},
	)
	if err != nil {
		panic(err)
	}

	var readDialectors []gorm.Dialector
	for _, r := range cfg.GetStringSlice("DB_READ_HOSTS") {
		if r == "" {
			continue
		}

		if strings.Contains(r, ":") {
			p := strings.Split(r, ":")
			r = p[0]
			dbPort = p[1]
		}

		ds := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			dbUser, dbPass,
			r, dbPort, dbName,
		)

		conn, err := sql.Open("postgres", ds)
		if err != nil {
			logger.Logger.Errorf("cannot init replica %v:%v", r, dbPort)
			continue
		}

		readDialectors = append(readDialectors, postgres.New(postgres.Config{Conn: conn}))
	}

	if len(readDialectors) > 0 {
		err := db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: readDialectors,
			// sources/replicas load balancing policy
			Policy: dbresolver.RandomPolicy{},
		}))
		if err != nil {
			panic(err)
		}
	}

	return &store{
		database:     db,
		shutdownFunc: conn.Close,
	}
}

// Shutdown close database connection
func (s *store) Shutdown() error {
	if s.shutdownFunc != nil {
		return s.shutdownFunc()
	}
	return nil
}

// DB database connection
func (s *store) DB() *gorm.DB {
	return s.database
}

// NewTx for database connection
func (s *store) NewTx() (txDB DB, finalFn IFinalFunc) {
	newDB := s.database.Begin()

	fn := FinalFunc{db: newDB}
	store := &store{database: newDB}

	return DB{
		Store: store,
	}, fn
}

type FinalFunc struct {
	db *gorm.DB
}

func (fn FinalFunc) Commit() error {
	return fn.db.Commit().Error
}

func (fn FinalFunc) Rollback(err error) error {
	return fn.db.Rollback().Error
}
