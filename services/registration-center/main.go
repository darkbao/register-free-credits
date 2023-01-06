package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	pb "registration-center/register"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type UserInfo struct {
	gorm.Model
	Name     string `gorm:"index;unique"`
	PhoneNum string
}

type SignInLog struct {
	gorm.Model
	Name      string `gorm:"index;unique"`
	IfMsgSent bool
}

var db *gorm.DB = nil

func init() {
	var err error
	dsn := "root:az0903AZA@tcp(localhost:3306)/register?charset=utf8&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err = db.AutoMigrate(&UserInfo{}, &SignInLog{}); err != nil {
		panic(fmt.Sprintf("db auto migrate failed: %v", err))
	}
}

type server struct {
	pb.UnimplementedRegisterServer
}

func (s *server) SignIn(_ context.Context, req *pb.SignInReq) (*pb.Empty, error) {
	log.Printf("Received req: %+v", req)
	err := db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		if err := tx.Create(&UserInfo{Name: req.Name, PhoneNum: req.PhoneNum}).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}

		if err := tx.Create(&SignInLog{Name: req.Name, IfMsgSent: false}).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRegisterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
