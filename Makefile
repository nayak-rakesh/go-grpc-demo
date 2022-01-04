gen:
	protoc --proto_path=proto proto/*.proto --go_opt=module=grpc-test --go-grpc_opt=module=grpc-test --go_out=. --go-grpc_out=.

clean:
	rm pb/*.go

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go