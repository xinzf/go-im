package main

import (
	// "fmt"
	"log"

	"github.com/xinzf/go-im/datacenter/user/proto"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/zookeeper"
	"golang.org/x/net/context"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())
	serverName := "im.datacenter.user"
	r := zookeeper.NewRegistry(registry.Addrs("192.168.31.152:2181"))
	service := micro.NewService(micro.Registry(r))
	client := proto.NewUserClient(serverName, service.Client())
	log.Println("request User.Detail")
	rsp, err := client.Detail(ctx, &proto.DetailRequest{Iid: "74e2e71fb619457b8c82c9bb8e544b46"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(rsp)

	// log.Println("request User.Create")
	// rsp, err = client.Create(ctx, &proto.CreateRequest{
	// 	Iid:  "",
	// 	Name: "test",
	// 	Icon: "url url",
	// 	Props: map[string]string{
	// 		"age": "324",
	// 		"sex": "fsfd",
	// 		"vip": "111",
	// 	},
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(rsp)
}
