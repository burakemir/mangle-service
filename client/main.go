package main

import (
	"context"
	"fmt"
	"io"
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
	fmt.Println("looking for answers!")
	var answer *pb.QueryAnswer
	for {
		answer, err = stream.Recv()
		if err != nil {
			break
		}
		fmt.Println(answer.GetAnswer())
	}
	if err != io.EOF {
		fmt.Println("got err %v", err)
	}
}
