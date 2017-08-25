package main

import (
	// "fmt"
	"log"

	"github.com/xinzf/go-im/datacenter/device/proto"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/zookeeper"
	"golang.org/x/net/context"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())
	serverName := "im.datacenter.device"
	r := zookeeper.NewRegistry(registry.Addrs("192.168.31.152:2181"))
	service := micro.NewService(micro.Registry(r))
	client := proto.NewDeviceClient(serverName, service.Client())
	// log.Println("request Device.Online")
	rsp, err := client.Online(ctx, &proto.OnlineRequest{
		Iid:  "74e2e71fb619457b8c82c9bb8e544b46",
		Node: "192.168.31.124",
		Props: map[string]string{
			"type": "mac",
			"ip":   "192.168.31.124",
			"tag":  "test",
		},
	})

	// log.Println("request Device.Offline")
	// rsp, err := client.Offline(ctx, &proto.OfflineRequest{
	// 	Iid:      "74e2e71fb619457b8c82c9bb8e544b46",
	// 	DeviceId: "f895f146d3584fd1a5eec437c0dc1413",
	// })

	// rsp, err := client.List(ctx, &proto.ListRequest{
	// 	Iid: "74e2e71fb619457b8c82c9bb8e544b46",
	// })

	// rsp, err := client.Detail(ctx, &proto.DetailRequest{
	// 	Iid:      "74e2e71fb619457b8c82c9bb8e544b46",
	// 	DeviceId: "f895f146d3584fd1a5eec437c0dc1413",
	// })

	if err != nil {
		log.Fatal(err)
	}
	log.Println(rsp)
}
