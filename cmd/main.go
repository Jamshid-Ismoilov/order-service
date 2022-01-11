package main

import (
	"net"

	"google.golang.org/grpc"

	"github.com/Jamshid-Ismoilov/order-service/config"
	pb "github.com/Jamshid-Ismoilov/order-service/genproto/order_service"
	"github.com/Jamshid-Ismoilov/order-service/pkg/db"
	"github.com/Jamshid-Ismoilov/order-service/pkg/logger"
	"github.com/Jamshid-Ismoilov/order-service/service"
	"github.com/Jamshid-Ismoilov/order-service/storage"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "template-service")
	defer func(l logger.Logger) {
		err := logger.Cleanup(l)
		if err != nil {
			log.Fatal("failed cleanup logger", logger.Error(err))
		}
	}(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	pgStorage := storage.NewStoragePg(connDB)

	orderService := service.NewOrderService(pgStorage, log)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, orderService)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
