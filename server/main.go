package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	pb "github.com/burakemir/mangle-service/proto"
	"github.com/burakemir/mangle-service/service"
	"google.golang.org/grpc"
)

var source = flag.String("source", "", "path to source to evaluate")

func readSource() io.Reader {
	sourceBytes, err := os.ReadFile(*source)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("read source from %q", *source)
	return strings.NewReader(string(sourceBytes))
}

func main() {
	flag.Parse()

	mangleService := service.New()
	if *source != "" {
		if err := mangleService.UpdateFromSource(readSource()); err != nil {
			log.Fatal(err)
		}
		log.Printf("updated db.")
	} else {
		log.Println("no --source given.")
	}
	server := grpc.NewServer()
	pb.RegisterMangleServer(server, mangleService)

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
