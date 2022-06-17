package main

import (
	"context"
	//"bigSystem/all_packaged_library/logtool"
	"bigSystem/svc/user/pb"
	//usergrpc "bigSystem/svc/user/Transport/grpc"
	//"go.uber.org/zap"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {

	conn, err := grpc.Dial("127.0.0.1:8078", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	svc := pb.NewUserClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := svc.RpcUserLogin(ctx, &pb.Login{
		Account:  "gfj",
		Password: "123456",
	})
	if err != nil {
		log.Fatalf("could not put: %v", err)
	}

	log.Println(r)

}
