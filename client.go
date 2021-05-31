package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"mace_convert/convert"
	"context"
	"os"
)


func init()  {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main(){
	// 客户端 gprc

	conn, err := grpc.Dial("172.17.0.3:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Faile to conncet to gRPC server :: %v", err)
	}
	defer conn.Close( )
	client := convert.NewConvertClient(conn)
	Path, DestPath := os.Args[1], os.Args[2]

	server_data, err := client.Mace(context.Background(), &convert.Client{Path: Path, DestPath: DestPath})

	if err != nil {
		log.Fatalf("%v: ", err)
	}

	fmt.Println(server_data.GetStdout())  // 打印即可
	if server_data.GetStderr() != ""{
		log.Fatal(server_data.GetStderr())  // 打印即可
	}

}

