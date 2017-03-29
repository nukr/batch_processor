test:
	# cd cmd/server && go test && cd ../..
	cd pkg/json/scanner && go test && cd ../../..
build:
	protoc --go_out=plugins=grpc:. ./pb/batch.proto
run:
	go run cmd/server/main.go
