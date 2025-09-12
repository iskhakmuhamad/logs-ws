package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iskhakmuhamad/mylogs-ws/internal/config"
	"github.com/iskhakmuhamad/mylogs-ws/internal/delivery/consumer"
	"github.com/iskhakmuhamad/mylogs-ws/internal/delivery/http"
	"github.com/iskhakmuhamad/mylogs-ws/internal/delivery/ws"
	"github.com/iskhakmuhamad/mylogs-ws/internal/repository"
	"github.com/iskhakmuhamad/mylogs-ws/internal/usecase"
	"github.com/iskhakmuhamad/mylogs-ws/pkg/broker"
	"github.com/iskhakmuhamad/mylogs-ws/pkg/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.LoadConfig()
	database := db.NewPostgresDB(cfg)
	notifRepo := repository.NewNotificationRepository(database)
	hub := ws.NewHub()
	go hub.Run()
	notifUsecase := usecase.NewNotificationUsecase(notifRepo, hub)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	http.NewNotificationHandler(app, notifUsecase)
	ws.NewWsHandler(app, hub, notifRepo)

	kc := broker.NewKafkaConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID, cfg.KafkaTopic)
	nc := consumer.NewNotificationConsumer(kc, notifUsecase)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := nc.Start(ctx); err != nil {
			log.Fatalf("Kafka consumer error: %v", err)
		}
	}()

	log.Printf("Server started at %s", cfg.Port)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", cfg.Port)))
}
