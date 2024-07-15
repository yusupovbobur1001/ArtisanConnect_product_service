package pkg

import (
	"errors"
	"log"
	"product_service/config"
	pbu "product_service/genproto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateAuthClient(cfg *config.Config) pbu.AuthClient {
	conn, err := grpc.NewClient(cfg.AUTH_SERVICE_PORT,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(errors.New("failed to connect to the address: " + err.Error()))
		return nil
	}

	return pbu.NewAuthClient(conn)
}
