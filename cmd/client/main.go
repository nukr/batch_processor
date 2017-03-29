package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	batch "github.com/nukr/batch_processor/pb"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":33333", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var document []byte
	{
		file, err := os.Open("cmd/client/document.json")
		if err != nil {
			log.Fatal(err)
		}
		document, err = ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
	}

	var selector []byte
	{
		file, err := os.Open("cmd/client/selector.json")
		if err != nil {
			log.Fatal(err)
		}
		selector, err = ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
	}

	client := batch.NewBatchServiceClient(conn)
	ctx := context.Background()
	client.Update(ctx, &batch.Query{
		Selector: selector,
		Document: document,
	})
}
