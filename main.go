package main

import (
	log "github.com/golang/glog"
	"github.com/micro/go-micro"

	"github.com/micro/router-srv/handler"
	"github.com/micro/router-srv/router"

	proto "github.com/micro/router-srv/proto/router"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.router"),
		micro.Version("latest"),
	)

	service.Init()

	proto.RegisterRouterHandler(service.Server(), new(handler.Router))

	// subcriber to stats
	service.Server().Subscribe(service.Server().NewSubscriber(router.StatsTopic, router.ProcessStats))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
