package main

import (
	"log"
	"net"
	"product_service/config"
	pbu "product_service/genproto/orders"
	pb "product_service/genproto/products"
	"product_service/pkg/logger"
	"product_service/service"
	"product_service/storage/postgres"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	logger, err := logger.New("debug", "develop", "app.log")
	if err != nil {
		log.Fatal(err)
		return
	}

	db, err := postgres.ConnDB()
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	lisstener, err := net.Listen("tcp", cfg.PRODUCT_SERVICE_PORT)
	if err != nil {
		logger.Fatal("error while making tcp connection")
		return
	}

	server := grpc.NewServer()	

	pb.RegisterProductServiceServer(server, service.NewProductService(db, cfg))
	pbu.RegisterOrderServiceServer(server, service.NewOrderService(db, cfg))

	err = server.Serve(lisstener)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}
}
