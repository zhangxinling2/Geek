package web

import "strings"

type router struct {
	trees map[string]*node
}

func NewRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

func (r *router) AddRouter(method string, path string, handleFunc HandlerFunc) {
	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		children := root.childrenOrCreate(seg)
		root = children
	}
	root.handler = handleFunc
}

type node struct {
	path string

	children map[string]*node

	handler HandlerFunc
}

func (n *node) childrenOrCreate(seg string) *node {
	if n.children == nil {
		n.children = map[string]*node{}
	}
	res, ok := n.children[seg]
	if !ok {
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}
