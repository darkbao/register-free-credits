package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "registration-center/register"
)

var (
	addr  = flag.String("addr", "localhost:50051", "the address to connect to")
	name  = flag.String("name", "nobody", "Name to sign in")
	phone = flag.String("phone", "11011011011", "Phone number")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() { _ = conn.Close() }()
	c := pb.NewRegisterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, err = c.SignIn(ctx, &pb.SignInReq{
		Name:     *name,
		PhoneNum: *phone,
	}); err != nil {
		log.Fatalf("could not sign in: %v", err)
	}
	log.Println("Registered successfully")
}
