package handler

import (
	"github.com/micro/go-micro/errors"
	"github.com/micro/router-srv/db"
	proto "github.com/micro/router-srv/proto/rule"
	"github.com/micro/router-srv/rule"

	"golang.org/x/net/context"
)

type Rule struct{}

func (r *Rule) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("go.micro.srv.router.Rule.Read", "invalid id")
	}

	l, err := rule.Read(req.Id)
	if err != nil {
		if err == db.ErrNotFound {
			return errors.NotFound("go.micro.srv.router.Rule.Read", err.Error())
		}
		return errors.InternalServerError("go.micro.srv.router.Rule.Read", err.Error())
	}

	rsp.Rule = l

	return nil
}

func (r *Rule) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	if req.Rule == nil {
		return errors.BadRequest("go.micro.srv.router.Rule.Create", "invalid rule")
	}

	if len(req.Rule.Id) == 0 {
		return errors.BadRequest("go.micro.srv.router.Rule.Create", "invalid id")
	}

	if len(req.Rule.Service) == 0 {
		return errors.BadRequest("go.micro.srv.router.Rule.Create", "invalid service")
	}

	if len(req.Rule.Version) == 0 {
		return errors.BadRequest("go.micro.srv.router.Rule.Create", "invalid version")
	}

	if req.Rule.Weight < 0 || req.Rule.Weight > 100 {
		return errors.BadRequest("go.micro.srv.router.Rule.Create", "invalid weight, must be 0 to 100")
	}

	if err := rule.Create(req.Rule); err != nil {
		return errors.InternalServerError("go.micro.srv.router.Rule.Create", err.Error())
	}

	return nil
}

func (r *Rule) Update(ctx context.Context, req *proto.UpdateRequest, rsp *proto.UpdateResponse) error {
	if req.Rule == nil {
		return errors.BadRequest("go.micro.srv.router.Rule.Update", "invalid rule")
	}

	if len(req.Rule.Id) == 0 {
		return errors.BadRequest("go.micro.srv.router.Rule.Update", "invalid id")
	}

	if len(req.Rule.Service) == 0 {
		return errors.BadRequest("go.micro.srv.router.Rule.Update", "invalid service")
	}

	if len(req.Rule.Version) == 0 {
		return errors.BadRequest("go.micro.srv.router.Rule.Update", "invalid version")
	}

	if req.Rule.Weight < 0 || req.Rule.Weight > 100 {
		return errors.BadRequest("go.micro.srv.router.Rule.Update", "invalid weight, must be 0 to 100")
	}

	if err := rule.Update(req.Rule); err != nil {
		return errors.InternalServerError("go.micro.srv.router.Rule.Update", err.Error())
	}

	return nil
}

func (r *Rule) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("go.micro.srv.router.Rule.Delete", "invalid id")
	}

	if err := rule.Delete(req.Id); err != nil {
		if err == db.ErrNotFound {
			return nil
		}
		return errors.InternalServerError("go.micro.srv.router.Rule.Delete", err.Error())
	}

	return nil
}

func (r *Rule) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.SearchResponse) error {
	if req.Limit <= 0 {
		req.Limit = 10
	}

	if req.Offset <= 0 {
		req.Offset = 0
	}

	rules, err := rule.Search(req.Service, req.Version, req.Limit, req.Offset)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.router.Rule.Search", err.Error())
	}

	rsp.Rules = rules

	return nil
}
