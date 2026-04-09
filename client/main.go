package main

import (
	"context"
	"log"
	"time"
	
	grpc_demo "github.com/AhmedHossam777/grpc-demo/pb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	
	client := grpc_demo.NewTaskServiceClient(conn)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	tasks := []struct {
		title string
		desc  string
	}{
		{
			"Learn gRPC basics",
			"Understand proto files, code generation, and unary RPCs",
		},
		{
			"Add streaming RPCs",
			"Implement server-streaming and bidirectional streaming",
		},
		{"Add interceptors", "gRPC equivalent of middleware — logging, auth, etc."},
	}
	
	for _, t := range tasks {
		resp, err := client.CreateTask(
			ctx, &grpc_demo.CreateTaskRequest{
				Title:       t.title,
				Description: t.desc,
			},
		)
		if err != nil {
			log.Fatalf("CreateTask failed: %v", err)
		}
		log.Printf("Created: %s — %s", resp.Task.Id, resp.Task.Title)
	}
	
	listResp, err := client.ListTasks(ctx, &grpc_demo.ListTasksRequest{})
	if err != nil {
		log.Fatalf("ListTasks failed: %v", err)
	}
	
	log.Println("\n--- All Tasks ---")
	for _, task := range listResp.Task {
		log.Printf("  [%s] %s (completed: %v)", task.Id, task.Title, task.Completed)
	}
}
