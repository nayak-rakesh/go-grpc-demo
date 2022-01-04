package users

import (
	"bytes"
	"context"
	"fmt"
	"grpc-test/pb"
	"io"
	"log"
	"os"
	"sync"
)

var id = 0

type Server struct {
	users []*pb.User
	mutex sync.Mutex

	*pb.UnimplementedUserServiceServer
}

func (s *Server) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id++
	user := &pb.User{
		Id:   int32(id),
		Name: req.Name,
	}
	s.users = append(s.users, user)
	fmt.Println(user)
	return &pb.AddUserResponse{
		User: user,
	}, nil
}

func (s *Server) UploadImage(stream pb.UserService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		log.Println("can not get image id")
	}
	id := req.GetId()
	log.Println("Image id is: ", id)
	imgData := bytes.Buffer{}
	for {
		log.Println("waiting for chunk data")
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("no more data to receive")
			break
		}
		chunk := req.GetChunk()
		_, err = imgData.Write(chunk)
		if err != nil {
			log.Println("can not write chunk", err)
			return err
		}
	}
	imgPath := "./img/img.png"
	file, err := os.Create(imgPath)
	if err != nil {
		return err
	}
	_, err = imgData.WriteTo(file)
	if err != nil {
		return err
	}
	res := &pb.UploadImageResponse{
		Id: id,
	}
	err = stream.SendAndClose(res)
	if err != nil {
		log.Println("can't send response", err)
		return err
	}
	return nil
}
