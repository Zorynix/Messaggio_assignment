package app

import (
	"context"
	"fmt"
	"messagio_testsuite/config"
	"messagio_testsuite/internal/repo"
	"messagio_testsuite/internal/repo/pgdb"
	v1 "messagio_testsuite/internal/routes/http/v1"
	"messagio_testsuite/internal/service"
	"messagio_testsuite/pkg/kafka"
	"messagio_testsuite/pkg/postgres"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func Run(configPath string) {

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		logrus.Fatalf("Config error: %s", err)
	}

	SetLogrus(cfg.Log.Level)

	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		logrus.Fatal(fmt.Errorf("app - Run - pgdb.NewServices: %w", err))
	}
	defer pg.Close()

	consumer := kafka.NewKafkaConsumer(cfg.Kafka.Brokers, cfg.Kafka.GroupID, cfg.Kafka.Topic)
	if err != nil {
		logrus.Fatalf("Failed to initialize Kafka consumer: %v", err)
	}
	defer consumer.Close()

	producer := kafka.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	if err != nil {
		logrus.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	defer producer.Close()

	messageRepo := pgdb.NewMessageRepo(pg)

	services := service.NewServices(service.ServicesDependencies{
		Repos:         &repo.Repositories{Message: messageRepo},
		KafkaProducer: producer,
		KafkaConsumer: consumer,
	})

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{validator: validator.New()}

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "all systems operational"})
	})

	v1.NewMessageRoutes(e.Group("/api/v1"), services.Message)

	addr := cfg.Server.Port
	logrus.Infof("Starting server on %s...", addr)
	go func() {
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	}

	logrus.Info("Server exiting")
}
