package main

import (
	"fmt"
	"log"
	"net"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"

	batchService "github.com/nukr/batch_processor/pb"
	"github.com/nukr/batch_processor/pkg/db"
	"github.com/nukr/batch_processor/pkg/db/rethinkdb"
)

type server struct {
	db db.DB
}

func (s *server) Create(
	ctx context.Context,
	q *batchService.Query,
) (*batchService.Response, error) {
	return &batchService.Response{
		Msg: "create",
	}, nil
}

func (s *server) Update(
	ctx context.Context,
	q *batchService.Query,
) (*batchService.Response, error) {
	affected, err := s.db.Update(q.GetSelector(), q.GetDocument())
	return handleResponse(affected, err)
}

func handleResponse(n int32, err error) (*batchService.Response, error) {
	return &batchService.Response{
		Affected: n,
		Msg:      fmt.Sprintf("updated %d row(s)", n),
	}, err
}

func (s *server) Delete(
	ctx context.Context,
	q *batchService.Query,
) (*batchService.Response, error) {
	return &batchService.Response{
		Msg: "delete",
	}, nil
}

func (s *server) Get(
	q *batchService.Query,
	stream batchService.BatchService_GetServer,
) error {
	for i := 0; i < 1000; i++ {
		stream.Send(&batchService.Response{
			Msg: "get",
		})
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":33333")
	if err != nil {
		log.Fatal("listen err => ", err)
	}
	grpcServer := grpc.NewServer()
	db := &rethinkdb.DB{}
	s := server{db: db}
	batchService.RegisterBatchServiceServer(grpcServer, &s)
	grpcServer.Serve(lis)
}
