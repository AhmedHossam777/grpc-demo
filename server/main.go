package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	grpc_demo "github.com/AhmedHossam777/grpc-demo/pb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type taskService struct {
	grpc_demo.UnimplementedTaskServiceServer

	mu     sync.Mutex
	tasks  []*grpc_demo.Task
	nextID int
}

func (s *taskService) CreateTask(
	ctx context.Context, req *grpc_demo.CreateTaskRequest,
) (*grpc_demo.CreateTaskResponse, error) {
	if req.GetTitle() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "title is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	task := &grpc_demo.Task{
		Id:          fmt.Sprintf("task-%d", s.nextID),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Completed:   false,
	}

	s.tasks = append(s.tasks, task)
	log.Printf("Created task: %s - %s", task.Id, task.Title)

	return &grpc_demo.CreateTaskResponse{Task: task}, nil
}

func (s *taskService) ListTasks(
	ctx context.Context,
	req *grpc_demo.ListTasksRequest,
) (*grpc_demo.ListTaskResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("Listing %d tasks", len(s.tasks))

	return &grpc_demo.ListTaskResponse{Task: s.tasks}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	grpc_demo.RegisterTaskServiceServer(grpcServer, &taskService{})
	log.Println("gRPC server listening on :50051")

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
