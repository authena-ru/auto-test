package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/authena-ru/auto-test/internal/server"
	"github.com/authena-ru/auto-test/internal/service"
)

func main() {
	lis, err := net.Listen("tcp", ":"+os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", lis)
	}

	servReg := grpc.NewServer()
	serv := server.New(service.CheckAttemptToPassTestingTask)
	server.RegisterAutoTestServiceServer(servReg, serv)

	log.Printf("server listening at %v", lis.Addr())

	if err := servReg.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
