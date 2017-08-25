package device

import (
	"fmt"
	// "log"

	"github.com/xinzf/go-im/datacenter/device/proto"
	"github.com/xinzf/go-utils"

	rds "github.com/garyburd/redigo/redis"
	"github.com/pquerna/ffjson/ffjson"
	"golang.org/x/net/context"
)

const (
	REDIS_ONLINE_KEY_FORMAT = "im:online:%s"
)

type DeviceHandle struct {
}

func (this *DeviceHandle) Online(ctx context.Context, req *proto.OnlineRequest, rsp *proto.DetailResponse) error {
	var err error
	client := redisPool.Get()
	defer client.Close()

	deviceId := utils.UUID()
	onlineTime := utils.NowTime()

	data, _ := ffjson.Marshal(map[string]interface{}{
		"deviceId": deviceId,
		"iid":      req.Iid,
		"time":     onlineTime,
		"props":    req.Props,
		"node":     req.Node,
	})

	_, err = client.Do(
		"HSET",
		fmt.Sprintf(REDIS_ONLINE_KEY_FORMAT, req.Iid),
		deviceId,
		data,
	)

	if err != nil {
		return err
	}

	rsp.DeviceId = deviceId
	rsp.Iid = req.Iid
	rsp.OnlineTime = onlineTime
	rsp.Node = req.Node
	rsp.Props = req.Props

	return nil
}

func (this *DeviceHandle) Offline(ctx context.Context, req *proto.OfflineRequest, rsp *proto.OperationResponse) error {
	var err error
	client := redisPool.Get()
	defer client.Close()

	_, err = client.Do(
		"HDEL",
		fmt.Sprintf(REDIS_ONLINE_KEY_FORMAT, req.Iid),
		req.DeviceId,
	)
	if err != nil {
		return err
	}

	rsp.Code = 0
	rsp.Ret = true
	rsp.Msg = ""
	return nil
}

func (this *DeviceHandle) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	client := redisPool.Get()
	defer client.Close()

	replay, err := client.Do(
		"HGETALL",
		fmt.Sprintf(REDIS_ONLINE_KEY_FORMAT, req.Iid),
	)

	if err != nil {
		return err
	}

	data, err := rds.StringMap(replay, err)
	if err != nil {
		return err
	}

	var p map[string]interface{}
	for _, v := range data {
		ffjson.Unmarshal([]byte(v), &p)
		detail := &proto.DetailResponse{
			DeviceId:   p["deviceId"].(string),
			Iid:        p["iid"].(string),
			OnlineTime: p["time"].(string),
			Node:       p["node"].(string),
			Props:      make(map[string]string),
		}

		var nprops = p["props"].(map[string]interface{})
		for pk, pv := range nprops {
			detail.Props[pk] = pv.(string)
		}

		rsp.Devices = append(rsp.Devices, detail)
	}

	return nil
}

func (this *DeviceHandle) Detail(ctx context.Context, req *proto.DetailRequest, rsp *proto.DetailResponse) error {
	client := redisPool.Get()
	defer client.Close()

	d, err := rds.String(client.Do(
		"HGET",
		fmt.Sprintf(REDIS_ONLINE_KEY_FORMAT, req.Iid),
		req.DeviceId,
	))

	if err != nil {
		return err
	}

	var data map[string]interface{}
	ffjson.Unmarshal([]byte(d), &data)

	rsp.DeviceId = data["deviceId"].(string)
	rsp.Iid = data["iid"].(string)
	rsp.OnlineTime = data["time"].(string)
	rsp.Node = data["node"].(string)
	rsp.Props = make(map[string]string)

	p := data["props"].(map[string]interface{})
	for pk, pv := range p {
		rsp.Props[pk] = pv.(string)
	}

	return nil
}
