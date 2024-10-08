// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/goodsinfo.proto

package goodsinfo

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Goodsinfo service

type GoodsinfoService interface {
	// AddGoods: 定义增加商品的微服务, 这里的写法和gRPC中的写法一致
	AddGoods(ctx context.Context, in *AddRequest, opts ...client.CallOption) (*AddResponse, error)
}

type goodsinfoService struct {
	c    client.Client
	name string
}

func NewGoodsinfoService(name string, c client.Client) GoodsinfoService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "goodsinfo"
	}
	return &goodsinfoService{
		c:    c,
		name: name,
	}
}

func (c *goodsinfoService) AddGoods(ctx context.Context, in *AddRequest, opts ...client.CallOption) (*AddResponse, error) {
	req := c.c.NewRequest(c.name, "Goodsinfo.AddGoods", in)
	out := new(AddResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Goodsinfo service

type GoodsinfoHandler interface {
	// AddGoods: 定义增加商品的微服务, 这里的写法和gRPC中的写法一致
	AddGoods(context.Context, *AddRequest, *AddResponse) error
}

func RegisterGoodsinfoHandler(s server.Server, hdlr GoodsinfoHandler, opts ...server.HandlerOption) error {
	type goodsinfo interface {
		AddGoods(ctx context.Context, in *AddRequest, out *AddResponse) error
	}
	type Goodsinfo struct {
		goodsinfo
	}
	h := &goodsinfoHandler{hdlr}
	return s.Handle(s.NewHandler(&Goodsinfo{h}, opts...))
}

type goodsinfoHandler struct {
	GoodsinfoHandler
}

func (h *goodsinfoHandler) AddGoods(ctx context.Context, in *AddRequest, out *AddResponse) error {
	return h.GoodsinfoHandler.AddGoods(ctx, in, out)
}
