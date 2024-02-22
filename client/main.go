package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pb "github.com/burakemir/mangle-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewMangleClient(conn)
	var query = "reachable(/a, X)"
	if len(os.Args) > 1 {
		query = os.Args[1]
	}
	stream, err := client.Query(context.Background(), &pb.QueryRequest{Query: query})
	if err != nil {
		panic(err)
	}
	log.Printf("[query %q]", query)
	var answer *pb.QueryAnswer
	n := 0
	for {
		answer, err = stream.Recv()
		if err != nil {
			break
		}
		n++
		fmt.Println(answer.GetAnswer())
	}
	if err != io.EOF {
		log.Printf("got err %v", err)
	}
	log.Printf("[got %d answers]", n)
}
