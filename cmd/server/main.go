package main

import (
	"grpc-test/pb"
	"grpc-test/users"
	"io/ioutil"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer() //
	pb.RegisterUserServiceServer(srv, &users.Server{})
	log.Infoln("serving grpc on: ", listener.Addr().String())
	if err = srv.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
