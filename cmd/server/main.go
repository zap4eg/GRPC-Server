package main

import (
	"GRPCProject/internal/config"
	"GRPCProject/internal/repository/postgres"
	"GRPCProject/internal/server"
	pb "GRPCProject/proto"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	if err := setupViper(); err != nil {
		log.Fatal(err)
	}

	host := viper.GetString("server.host")
	port := viper.GetInt("server.port")
	addr := fmt.Sprintf("%s:%d", host, port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("TCP listener error %v", err)
	}
	defer listener.Close()

	postgresPool, err := config.SetupPostgres(ctx)
	if err != nil {
		log.Fatalf("PostgreSQL setup error %v", err)
	}

	userRepository := postgres.NewUserRepository(postgresPool)

	userServer := server.NewUserServer(userRepository)

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, userServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("GRPC Server error %v", err)
	}

}

func setupViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
