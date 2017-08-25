package user

import (
	"fmt"
	"log"

	"github.com/xinzf/go-im/datacenter/user/proto"
	"github.com/xinzf/go-utils"

	rds "github.com/garyburd/redigo/redis"
	"github.com/pquerna/ffjson/ffjson"
	"golang.org/x/net/context"
)

const (
	REDIS_KEY_FORMAT = "im:data:user:%s"
)

type UserHandle struct {
}

func (this *UserHandle) setCache(u UserModel) error {
	props, _ := ffjson.Marshal(u.Props)
	// key := fmt.Sprintf(REDIS_KEY_FORMAT, u.Iid)
	_, err := redisPool.Get().Do("HMSET",
		fmt.Sprintf(REDIS_KEY_FORMAT, u.Iid),
		"iid",
		u.Iid,
		"name",
		u.Name,
		"icon",
		u.Icon,
		"create_at",
		u.CreateAt,
		"props",
		string(props),
	)
	return err
}

func (this *UserHandle) getCache(iid string) (u UserModel, err error) {
	client := redisPool.Get()
	defer client.Close()

	replay, err := client.Do("HGETALL", fmt.Sprintf(REDIS_KEY_FORMAT, iid))
	data, err := rds.StringMap(replay, err)

	if err != nil {
		log.Println(err)
		return
	}

	u.Iid = data["iid"]
	u.Name = data["name"]
	u.Icon = data["icon"]
	u.CreateAt = data["create_at"]
	u.Props = data["props"]

	return
}

func (this *UserHandle) Detail(ctx context.Context, req *proto.DetailRequest, rsp *proto.DetailResponse) error {
	log.Println("Received User.Detail request")

	var props map[string]string
	var u UserModel
	var err error

	u, err = this.getCache(req.Iid)
	if err != nil || len(u.Iid) == 0 {
		if err != nil {
			log.Println(err)
		}

		err = db.Where("iid=?", req.Iid).First(&u).Error
		if err != nil {
			log.Println(err)
			return err
		}

		err = this.setCache(u)
		if err != nil {
			log.Println(err)
		}
	}

	ffjson.Unmarshal([]byte(u.Props), &props)

	rsp.Iid = u.Iid
	rsp.Name = u.Name
	rsp.Icon = u.Icon
	rsp.Props = props
	rsp.CreateAt = u.CreateAt

	return nil
}

func (this *UserHandle) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.DetailResponse) error {
	props, _ := ffjson.Marshal(req.Props)
	u := UserModel{
		Name:     req.Name,
		Icon:     req.Icon,
		Props:    string(props),
		CreateAt: utils.NowTime(),
	}
	err := db.Create(&u).Error
	if err != nil {
		return err
	}

	this.setCache(u)

	rsp.Iid = u.Iid
	rsp.Name = u.Name
	rsp.Icon = u.Icon
	rsp.CreateAt = u.CreateAt
	rsp.Props = req.Props

	return nil
}
