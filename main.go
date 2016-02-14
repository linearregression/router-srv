package main

import (
	log "github.com/golang/glog"
	"github.com/micro/cli"
	"github.com/micro/go-micro"

	"github.com/micro/router-srv/db"
	"github.com/micro/router-srv/db/mysql"
	"github.com/micro/router-srv/handler"
	"github.com/micro/router-srv/router"

	label "github.com/micro/router-srv/proto/label"
	proto "github.com/micro/router-srv/proto/router"
	rule "github.com/micro/router-srv/proto/rule"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.router"),
		micro.Version("latest"),

		micro.Flags(
			cli.StringFlag{
				Name:   "database_url",
				EnvVar: "DATABASE_URL",
				Usage:  "The database URL e.g root@tcp(127.0.0.1:3306)/router",
			},
		),

		micro.Action(func(c *cli.Context) {
			if len(c.String("database_url")) > 0 {
				mysql.Url = c.String("database_url")
			}
		}),
	)

	service.Init(
		micro.BeforeStart(func() error {
			router.Init(service)
			return nil
		}),
	)

	proto.RegisterRouterHandler(service.Server(), new(handler.Router))
	label.RegisterLabelHandler(service.Server(), new(handler.Label))
	rule.RegisterRuleHandler(service.Server(), new(handler.Rule))

	// subcriber to stats
	service.Server().Subscribe(service.Server().NewSubscriber(router.StatsTopic, router.ProcessStats))

	// initialise database
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
