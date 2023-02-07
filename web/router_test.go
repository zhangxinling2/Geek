package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_AddRouter(t *testing.T) {
	//构造路由树
	testRoute := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
	}
	var mockHandler HandlerFunc = func(ctx Context) {

	}
	r := NewRouter()
	for _, route := range testRoute {
		r.AddRouter(route.method, route.path, mockHandler)
	}
	//判断两者相等
	wantRouteTree := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path: "/",
				children: map[string]*node{
					"user": &node{
						path: "user",
						children: map[string]*node{
							"home": &node{
								path:    "home",
								handler: mockHandler,
							},
						},
					},
				},
			},
		},
	}
	msg, ok := wantRouteTree.equal(r)
	assert.True(t, ok, msg)
}

func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("没有相同的HTTP方法"), false
		}
		//v ,dst 要相等
		msg, equal := v.equal(dst)
		if !equal {
			return msg, false
		}
	}
	return "", true
}
func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("没有相同的路径"), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("子节点路径不匹配"), false
	}
	//比较 handler
	nHandler := reflect.ValueOf(n.handler)
	yHandler := reflect.ValueOf(y.handler)
	if nHandler != yHandler {
		return fmt.Sprintf("handler 不相等"), false
	}
	for path, c := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点%s不存在", path), false
		}
		msg, ok := c.equal(dst)
		if !ok {
			return msg, false
		}
	}
	return "", true
}
