package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os/signal"
	"syscall"

	pb "github.com/burakemir/mangle-service/proto"
	"github.com/burakemir/mangle-service/service"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	server := grpc.NewServer()
	mangleServer, err := service.New()
	if err != nil {
		panic(err)
	}
	pb.RegisterMangleServer(server, mangleServer)

	basectx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		port := ":8080"
		log.Println("listen and serving on", port)

		listener, err := net.Listen("tcp", port)
		if err != nil {
			panic(err)
		}

		if err := server.Serve(listener); err != nil {
			panic(err)
		}
	}()

	<-basectx.Done()
	log.Println("bye")
	server.GracefulStop()
}
