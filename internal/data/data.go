package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/feishu/feishu-office-suite/internal/data/ent"
	"github.com/feishu/feishu-office-suite/internal/kafka"
	"github.com/feishu/feishu-office-suite/internal/cache"
)

type OpenTelemetryConfig struct {
	Endpoint string `json:"endpoint"`
}

type BootstrapConfig struct {
	Server      ServerConfig      `json:"server"`
	MySQL       MySQLConfig       `json:"mysql"`
	Redis       RedisConfig       `json:"redis"`
	Kafka       KafkaConfig       `json:"kafka"`
	Asynq       AsynqConfig       `json:"asynq"`
	Feishu      FeishuConfig      `json:"feishu"`
	OpenTelemetry OpenTelemetryConfig `json:"open_telemetry"`
}

type ServerConfig struct {
	HTTP HTTPConfig `json:"http"`
	GRPC GRPCConfig `json:"grpc"`
}

type HTTPConfig struct {
	Network string `json:"network"`
	Timeout int    `json:"timeout"`
}

type GRPCConfig struct {
	Network string `json:"network"`
	Timeout int    `json:"timeout"`
}

type MySQLConfig struct {
	DSN          string `json:"dsn"`
	MaxOpenConns int    `json:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns"`
	MaxLifetime  int    `json:"max_lifetime"`
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	PoolSize int    `json:"pool_size"`
}

type KafkaConfig struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
}

type AsynqConfig struct {
	Redis RedisConfig `json:"redis"`
}

type FeishuConfig struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	Token     string `json:"token"`
}

type Data struct {
	db        *gorm.DB
	rdb       *redis.Client
	entClient *ent.Client
	kafka     *kafka.KafkaProducer
	cache     *cache.Cache
	asynq     *AsynqClient
}

type AsynqClient struct {
	client *redis.Client
}

func NewData(logger log.Logger, cfg *BootstrapConfig, tp *trace.TracerProvider) (*Data, func(), error) {
	d := &Data{}
	var cleanup cleanupFunc

	gormLogger := logger.DefaultLogger

	db, err := NewMySQL(cfg.MySQL, gormLogger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create mysql: %w", err)
	}
	d.db = db

	entClient, err := ent.NewClient(ent.Driver(db))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ent client: %w", err)
	}
	d.entClient = entClient

	rdb := NewRedis(cfg.Redis)
	d.rdb = rdb

	d.cache = cache.NewCache(rdb)

	if cfg.Kafka.Brokers != nil && len(cfg.Kafka.Brokers) > 0 {
		producer, err := kafka.NewProducer(cfg.Kafka.Brokers)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create kafka producer: %w", err)
		}
		d.kafka = producer
	}

	asynqClient := NewAsynqClient(cfg.Asynq.Redis)
	d.asynq = asynqClient

	cleanup = func() {
		if entClient != nil {
			entClient.Close()
		}
		if rdb != nil {
			rdb.Close()
		}
		if d.kafka != nil {
			d.kafka.Close()
		}
	}

	return d, cleanup, nil
}

func NewMySQL(cfg MySQLConfig, logger log.Logger) (*gorm.DB, error) {
	gormLogger := logger.Default()

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second)

	return db, nil
}

func NewRedis(cfg RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
}

func (d *Data) DB() *gorm.DB {
	return d.db
}

func (d *Data) EntClient() *ent.Client {
	return d.entClient
}

func (d *Data) Redis() *redis.Client {
	return d.rdb
}

func (d *Data) Cache() *cache.Cache {
	return d.cache
}

func (d *Data) Kafka() *kafka.KafkaProducer {
	return d.kafka
}

func (d *Data) Close() error {
	if d.entClient != nil {
		d.entClient.Close()
	}
	if d.rdb != nil {
		d.rdb.Close()
	}
	if d.kafka != nil {
		d.kafka.Close()
	}
	return nil
}

type cleanupFunc func()