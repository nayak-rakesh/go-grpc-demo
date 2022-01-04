package main

import (
	"bufio"
	"context"
	"fmt"
	"grpc-test/pb"
	"io"
	"log"

	"net/http"
	// "os"
	"strconv"

	"google.golang.org/grpc"
)

var userService pb.UserServiceClient

func main() {
	config()
	http.HandleFunc("/name", addUser)
	http.HandleFunc("/upload", uploadImage)
	log.Println("client server stared at 9040")
	log.Fatal(http.ListenAndServe(":9040", nil))

	// test()
}

func config() {
	conn, err := grpc.Dial("0.0.0.0:8080", grpc.WithInsecure())
	// defer conn.Close()
	if err != nil {
		log.Fatal(err)

	}
	userService = pb.NewUserServiceClient(conn)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	req := &pb.AddUserRequest{
		Name: "John Doe",
	}
	res, err := userService.AddUser(context.TODO(), req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Response", res)
	fmt.Fprint(w, res)
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	// file, err := os.Open("../temp/male.png")
	// if err != nil {
	// 	log.Println("could not open", err)
	// }
	// defer file.Close()

	file, _, err := r.FormFile("profile")
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err)
	}
	defer file.Close()

	id := 1
	stream, err := userService.UploadImage(context.TODO())
	if err != nil {
		log.Println(err)
	}
	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Id{
			Id: strconv.Itoa(id),
		},
	}
	err = stream.Send(req)
	if err != nil {
		log.Println("couldn't send id", err)
	}
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
		}
		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_Chunk{
				Chunk: buffer[:n],
			},
		}
		err = stream.Send(req)
		if err != nil {
			log.Println(err)
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Println(err)
	}
	log.Println("successfully uploaded", res.Id)
	fmt.Fprintf(w, "upload success %s", res.Id)

}

func test() {
	conn, err := grpc.Dial("0.0.0.0:8080", grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		log.Fatal(err)

	}
	req := &pb.AddUserRequest{
		Name: "Rakesh",
	}
	userService := pb.NewUserServiceClient(conn)
	res, err := userService.AddUser(context.TODO(), req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Response", res)
}
