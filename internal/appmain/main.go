// Package appmain contains the common application initialization code for Social Payment servers.
package appmain

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	// "github.com/consolelabs/mochi-toolkit/config"
	// "github.com/consolelabs/mochi-toolkit/kafka"
	"github.com/consolelabs/mochi-toolkit/kafka"
	"github.com/sirupsen/logrus"

	"github.com/trkhoi/khoi-api/config"
	cache "github.com/trkhoi/khoi-api/internal/cache"
	repo "github.com/trkhoi/khoi-api/internal/db"
	"github.com/trkhoi/khoi-api/logger"
)

var (
	l *logrus.Entry
)

// RunApplication starts and runs the given application forever.  For use in
// main functions to run the full application.
func RunApplication(serviceName string, bindService Bind) {
	l = logger.New(serviceName)
	c := make(chan os.Signal, 1)
	// SIGTERM is signaled by k8s when it wants a pod to stop.
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	readConfig := func() (config.View, error) {
		return config.Read()
	}

	// mpKafka, err := newKafkaConsumer(readConfig)
	// if err != nil {
	// 	l.Fatal(err)
	// }

	// // start producer
	// go func() {
	// 	if err := mpKafka.RunProducer(); err != nil {
	// 		l.Fatal(err)
	// 	}
	// }()

	a, err := newApplication(serviceName, bindService, readConfig, net.Listen, nil)
	if err != nil {
		l.Fatal(err)
	}

	go a.serve()
	l.Info("Application started successfully")

	<-c
	err = a.stop()
	if err != nil {
		l.Fatal(err)
	}
	l.Info("Application stopped successfully.")
}

// Bind is a function which starts an application, and binds it to serving.
type Bind func(p *Params, b *Bindings) IServer

// Params are inputs to starting an application.
type Params struct {
	config      config.View
	logger      *logrus.Entry
	serviceName string
	db          *repo.DB
	// mpKafka     *kafka.Kafka
	Cache cache.Cache
}

// NewParams creates a new Params object.
func NewParams(serviceName string, config config.View, logger *logrus.Entry, db *repo.DB, mpKafka *kafka.Kafka, cache cache.Cache) *Params {
	return &Params{
		config:      config,
		logger:      logger,
		serviceName: serviceName,
		db:          db,
		// mpKafka: mpKafka,
		Cache: cache,
	}
}

// DB provides a database repository for the application.
func (p *Params) DB() *repo.DB {
	return p.db
}

// Config provides the configuration for the application.
func (p *Params) Config() config.View {
	return p.config
}

// Logger provides a logger for the application.
func (p *Params) Logger() *logrus.Entry {
	return p.logger
}

// Kafka provides a kafka consumer for the application.
// func (p *Params) Kafka() *kafka.Kafka {
// 	return p.mpKafka
// }

type IServer interface {
	ListenAndServe() error
}

// ServiceName is a name for the currently running binary specified by
// RunApplication.
func (p *Params) ServiceName() string {
	return p.serviceName
}

// Bindings allows applications to bind various functions to the running servers.
type Bindings struct {
	a *App
}

// AddCloserErr specifies a function to be called when the application is being
// stopped.  Closers are called in reverse order.  The first error returned by
// a closer will be logged.
func (b *Bindings) AddCloserErr(c func() error) {
	b.a.closers = append(b.a.closers, c)
}

// App is used internally, and public only for apptest.  Do not use, and use apptest instead.
type App struct {
	closers []func() error
	srv     IServer
}

// newApplication is used internally, and public only for apptest.  Do not use, and use apptest instead.
func newApplication(serviceName string, bindService Bind, getCfg func() (config.View, error), listen func(network, address string) (net.Listener, error), kafka *kafka.Kafka) (*App, error) {
	a := &App{}

	cfg, err := getCfg()
	if err != nil {
		l.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalf("cannot read configuration.")
	}

	// pgRepo := repo.New(cfg, l)

	// *** cache ***
	// var redisCli *redis.Client
	// if cfg.GetString("REDIS_MASTER_NAME") == "" {
	// 	redisURL := cfg.GetString("REDIS_URL")
	// 	if redisURL == "" {
	// 		l.Fatalf("redis url is not set")
	// 	}
	// 	redisOpt, err := redis.ParseURL(redisURL)
	// 	if err != nil {
	// 		log.Fatal(err, "failed to init redis")
	// 	}
	// 	redisCli = redis.NewClient(redisOpt)
	// } else {
	// 	redisURL := cfg.GetString("REDIS_SENTINEL_URL")
	// 	if redisURL == "" {
	// 		l.Fatalf("redis url is not set")
	// 	}
	// 	redisCli = redis.NewFailoverClient(&redis.FailoverOptions{
	// 		SentinelAddrs: strings.Split(redisURL, ","),
	// 		MasterName:    cfg.GetString("REDIS_MASTER_NAME"),
	// 	})
	// }
	// cache, err := cache.NewRedisCache(nil)
	// if err != nil {
	// 	log.Fatal(err, "failed to init redis cache")
	// }

	p := NewParams(serviceName, cfg, l, nil, kafka, nil)

	b := &Bindings{
		a: a,
	}

	a.srv = bindService(p, b)

	return a, nil
}

// newKafkaConsumer is used internally, and public only for apptest.  Do not use, and use apptest instead.
// func newKafkaConsumer(getCfg func() (config.View, error)) (*kafka.Kafka, error) {
// 	cfg, err := getCfg()
// 	if err != nil {
// 		l.WithFields(logrus.Fields{
// 			"error": err.Error(),
// 		}).Fatalf("cannot read configuration.")
// 	}
// 	kafkaBrokers := cfg.GetString("KAFKA_BROKERS")
// 	kafkaProducer := kafka.New(kafkaBrokers, l)

// 	return kafkaProducer, nil
// }

// stop is used internally, and public only for apptest.  Do not use, and use apptest instead.
func (a *App) stop() error {
	// Use closers in reverse order: Since dependencies are created before
	// their dependants, this helps ensure no dependencies are closed
	// unexpectedly.
	var firstErr error
	for i := len(a.closers) - 1; i >= 0; i-- {
		err := a.closers[i]()
		if firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (a *App) serve() error {
	return a.srv.ListenAndServe()
}

// func NewJob(jobName string, bindService Bind) {
// 	l = logger.New(jobName)
// 	readConfig := func() (config.View, error) {
// 		return config.Read()
// 	}
// 	mpKafka, err := newKafkaConsumer(readConfig)
// 	if err != nil {
// 		l.Fatal(err)
// 	}

// 	// start producer
// 	go func() {
// 		if err := mpKafka.RunProducer(); err != nil {
// 			l.Fatal(err)
// 		}
// 	}()

// 	_, err = newApplication(jobName, bindService, readConfig, nil, mpKafka)
// 	if err != nil {
// 		l.Fatal(err)
// 	}
// 	l.Info("Job initialized successfully")
// }

// func NewConsumer() {
// 	l = logger.New("app.consumer")

// 	readConfig := func() (config.View, error) {
// 		return config.Read()
// 	}

// 	cfg, err := readConfig()
// 	if err != nil {
// 		l.Fatal("Failed to read config")
// 	}

// 	pgRepo := repo.New(cfg, l)
// 	p := NewParams("app.consumer", cfg, l, pgRepo, nil, nil)

// 	appconsumer.NewConsumer(p)
// }
