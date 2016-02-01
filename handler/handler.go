package handler

import (
	"github.com/micro/go-micro/errors"
	proto "github.com/micro/router-srv/proto/router"
	"github.com/micro/router-srv/router"

	"golang.org/x/net/context"
)

type Router struct{}

func (r *Router) Stats(ctx context.Context, req *proto.StatsRequest, rsp *proto.StatsResponse) error {
	if len(req.Service) == 0 {
		return errors.BadRequest("go.micro.srv.router.Stats", "invalid service name")
	}

	stats, err := router.Stats(req.Service, req.NodeId)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.router.Stats", err.Error())
	}

	rsp.Stats = stats

	return nil
}

func (r *Router) Select(ctx context.Context, req *proto.SelectRequest, rsp *proto.SelectResponse) error {
	// TODO: process filters

	if len(req.Service) == 0 {
		return errors.BadRequest("go.micro.srv.router.Select", "invalid service name")
	}

	services, err := router.Select(req.Service)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.router.Select", err.Error())
	}

	rsp.Services = services
	rsp.Expires = int64(router.DefaultExpiry)

	return nil
}

func (r *Router) Mark(ctx context.Context, req *proto.MarkRequest, rsp *proto.MarkResponse) error {
	return nil
}
