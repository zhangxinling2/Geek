package web

import (
	"fmt"
	"regexp"
	"strings"
)

type router struct {
	// trees 是按照 HTTP 方法来组织的
	// 如 GET => *node
	trees map[string]*node
}

func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

// addRoute 注册路由。
// method 是 HTTP 方法
// - 已经注册了的路由，无法被覆盖。例如 /user/home 注册两次，会冲突
// - path 必须以 / 开始并且结尾不能有 /，中间也不允许有连续的 /
// - 不能在同一个位置注册不同的参数路由，例如 /user/:id 和 /user/:name 冲突
// - 不能在同一个位置同时注册通配符路由和参数路由，例如 /user/:id 和 /user/* 冲突
// - 同名路径参数，在路由匹配的时候，值会被覆盖。例如 /user/:id/abc/:id，那么 /user/123/abc/456 最终 id = 456
func (r *router) addRoute(method string, path string, handler HandleFunc) {
	//判断路径为空
	if path == "" {
		panic("web: 路由是空字符串")
	}

	//获取方法树，如果没有则创建
	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	//对根节点特殊处理
	if path == "/" {
		if root.handler != nil {
			panic("web: 路由冲突[/]")
		}
		root.handler = handler
		return
	}
	//检测/的三种情况
	if path[0] != '/' {
		panic("web: 路由必须以 / 开头")
	}

	if path != "/" && strings.HasSuffix(path, "/") {
		panic("web: 路由不能以 / 结尾")
	}
	if strings.Contains(path, "//") {
		panic("web: 非法路由。不允许使用 //a/b, /a//b 之类的路由, [" + path + "]")
	}
	//没有问题则分割路径，之后添加节点
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		child := root.childOrCreate(seg)
		root = child
	}
	if root.handler != nil {
		panic("web: 路由冲突[" + path + "]")
	}
	//为节点添加方法
	root.handler = handler
}

// findRoute 查找对应的节点
// 注意，返回的 node 内部 HandleFunc 不为 nil 才算是注册了路由
func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	//先检查路径是否为空
	if path == "" {
		return nil, false
	}
	//找到路径所代表的方法树
	root, ok := r.trees[method]
	if !ok {
		panic("此方法还未注册路由")
	}
	//对根节点做特殊处理
	if path == "/" {
		if root.handler == nil {
			panic("web:此路由未注册")
		}
		return &matchInfo{n: root}, true
	}
	//将path切割并一层一层的向下搜索
	segs := strings.Split(path[1:], "/")
	info := &matchInfo{}
	for _, seg := range segs {
		child, ok := root.childOf(seg)
		if !ok {
			//还有一种通配符在末尾可以匹配后续所有路径
			if root.typ != nodeTypeAny {
				return nil, false
			} else {
				return &matchInfo{n: root}, true
			}
			return nil, false
		}
		if child.paramName != "" {
			info.addValue(child.paramName, seg)
		}
		root = child
	}
	info.n = root
	if root.handler == nil {
		return &matchInfo{n: root}, true
	}
	// expected: map[string]string{"id":"123"} 		actual  : map[string]string(nil)		debug发现:id的handle没有注册上,在注册路由时没有判断是否已存在路径路由
	if root.typ == nodeTypeParam || root.typ == nodeTypeReg {
		info := &matchInfo{n: root}
		info.addValue(root.paramName, segs[len(segs)-1])
		return info, true
	}
	return info, true
}

type nodeType int

const (
	// 静态路由
	nodeTypeStatic = iota
	// 正则路由
	nodeTypeReg
	// 路径参数路由
	nodeTypeParam
	// 通配符路由
	nodeTypeAny
)

// node 代表路由树的节点
// 路由树的匹配顺序是：
// 1. 静态完全匹配
// 2. 正则匹配，形式 :param_name(reg_expr)
// 3. 路径参数匹配：形式 :param_name
// 4. 通配符匹配：*
// 这是不回溯匹配
type node struct {
	typ nodeType

	path string
	// children 子节点
	// 子节点的 path => node
	children map[string]*node
	// handler 命中路由之后执行的逻辑
	handler HandleFunc

	// 通配符 * 表达的节点，任意匹配
	starChild *node

	paramChild *node
	// 正则路由和参数路由都会使用这个字段
	paramName string

	// 正则表达式
	regChild *node
	regExpr  *regexp.Regexp
}

// child 返回子节点
// 第一个返回值 *node 是命中的节点
// 第二个返回值 bool 代表是否命中
func (n *node) childOf(path string) (*node, bool) {
	//按优先级依次匹配
	child, ok := n.children[path]
	//如果静态路由中没有此节点，继续按优先级向下匹配：静态，正则，路径，通配符
	if !ok {
		//正则匹配,正则匹配需判断是否符合正则表达式,否则看是否有路径匹配
		if n.regChild != nil {
			if n.regChild.regExpr.MatchString(path) {
				return n.regChild, true
			}
		}
		//路径匹配
		if n.paramChild != nil {
			return n.paramChild, true
		}
		//通配符匹配
		if n.starChild != nil {
			return n.starChild, true
		}
		return nil, false
	}
	return child, true
}

// childOrCreate 查找子节点，
// 首先会判断 path 是不是通配符路径
// 其次判断 path 是不是参数路径，即以 : 开头的路径
// 最后会从 children 里面查找，
// 如果没有找到，那么会创建一个新的节点，并且保存在 node 里面
func (n *node) childOrCreate(path string) *node {

	//如果是通配符，判断此节点是否已经注册了其他路由,或重复注册，如果都没有则注册
	if path == "*" {
		if n.paramChild != nil {
			panic("web: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [" + path + "]")
		}
		if n.regChild != nil {
			panic("web: 非法路由，已有正则路由。不允许同时注册通配符路由和正则路由 [" + path + "]")
		}
		if n.starChild != nil {
			return n.starChild
		}
		n.starChild = &node{
			path: "*",
			typ:  nodeTypeAny,
		}
		return n.starChild
	}
	//如果是开头:，再判断是否注册了通配符之后需要判断是路径或是正则表达式，在进行相应的操作。
	if path[0] == ':' {
		isReg := checkPathReg(path)
		//如果是正则表达式，需要解析正则表达式,得到paramName和表达式Regexp,不是则注册路径路由
		if isReg {
			if n.starChild != nil {
				panic("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和正则路由 [" + path + "]")
			}
			if n.paramChild != nil {
				panic("web: 非法路由，已有路径参数路由。不允许同时注册正则路由和参数路由 [" + path + "]")
			}
			name, reg := parseReg(path)
			if n.regChild != nil {
				if n.regChild.regExpr != reg || n.paramName != name {
					panic(fmt.Sprintf("web: 路由冲突，正则路由冲突，已有 %s，新注册 %s", n.regChild.path, path))
				}
				return n.regChild
			}
			n.regChild = &node{
				path:      path,
				typ:       nodeTypeReg,
				regExpr:   reg,
				paramName: name,
			}
			return n.regChild
		} else {
			if n.starChild != nil {
				panic("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [" + path + "]")
			}
			if n.regChild != nil {
				panic("web: 非法路由，已有正则路由。不允许同时注册正则路由和参数路由 [" + path + "]")
			}
			if n.paramChild != nil {
				if n.paramChild.path != path {
					panic("web: 路由冲突，参数路由冲突，已有 " + n.paramChild.path + "，新注册 " + path)
				}
				return n.paramChild
			}
			n.paramChild = &node{
				path:      path,
				typ:       nodeTypeParam,
				paramName: path[1:],
			}
			return n.paramChild
		}
	}
	//对map为空做特殊处理  解决assignment to entry in nil map
	if n.children == nil {
		n.children = map[string]*node{}
	}
	//最后查找静态路由
	child, ok := n.children[path]
	if !ok {
		child = &node{
			path: path,
			typ:  nodeTypeStatic,
		}
		n.children[path] = child
		return child
	}
	return child
}

//parseReg 解析正则表达式
//第一个返回paramName
//第二个返回正则表达式解析
func parseReg(path string) (string, *regexp.Regexp) {
	segs := strings.Split(path, "(")
	name := segs[0]
	seg := segs[1]
	reg, err := regexp.Compile(seg[:len(seg)-1])
	if err != nil {
		panic("正则表达式解析错误")
	}
	return name[1:], reg
}

//checkPathType 判断是路径还是通配符
//返回是否是正则表达式
func checkPathReg(path string) bool {
	segs := strings.Split(path, "(")
	if len(segs) == 2 {
		if strings.HasSuffix(segs[1], ")") {
			return true
		}
	}
	return false
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func (m *matchInfo) addValue(key string, value string) {
	if m.pathParams == nil {
		// 大多数情况，参数路径只会有一段
		m.pathParams = map[string]string{key: value}
	}
	m.pathParams[key] = value
}
