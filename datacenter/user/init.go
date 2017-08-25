package user

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/xinzf/go-im/datacenter/user/proto"

	rds "github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/jinzhu/configor"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/zookeeper"
)

var db *gorm.DB
var redisPool *rds.Pool
var ctx context.Context
var Shutdown context.CancelFunc

func initDB() {
	url := fmt.Sprintf(
		"%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local",
		config.Db.User,
		config.Db.Pswd,
		config.Db.Host,
		config.Db.Name,
	)

	var err error
	db, err = gorm.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}

	db.DB().SetMaxOpenConns(config.Db.MaxOpenConns)
	db.DB().SetMaxIdleConns(config.Db.MaxIdleConns)
}

func initRedis() {
	redisPool = &rds.Pool{
		MaxIdle:     config.Redis.MaxIdleConns,
		MaxActive:   config.Redis.MaxOpenConns,
		IdleTimeout: 240 * time.Second,
		Dial: func() (rds.Conn, error) {
			c, err := rds.Dial("tcp", config.Redis.Host)
			if err != nil {
				return nil, err
			}

			if len(config.Redis.Password) > 0 {
				if _, err := c.Do("AUTH", config.Redis.Password); err != nil {
					c.Close()
					return nil, err
				}
			}

			c.Do("SELECT", config.Redis.Db)
			return c, err
		},
		TestOnBorrow: func(c rds.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func Run() {
	service := micro.NewService(
		micro.Flags(
			cli.StringFlag{
				Name:  "config",
				Usage: "config file path",
			},
		),
	)

	service.Init(
		micro.Action(func(c *cli.Context) {
			configor.Load(&config, c.String("config"))
		}),
	)

	initDB()
	initRedis()

	ctx, Shutdown = context.WithCancel(context.Background())
	r := zookeeper.NewRegistry(registry.Addrs(config.Discovery.Addrs))
	service.Init(
		micro.Registry(r),
		micro.RegisterTTL(time.Second*time.Duration(config.Discovery.Ttl)),
		micro.RegisterInterval(time.Second*time.Duration(config.Discovery.Interval)),
		micro.Name(config.Name),
		micro.Version(config.Version),
		micro.Context(ctx),
	)

	proto.RegisterUserHandler(service.Server(), new(UserHandle))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
