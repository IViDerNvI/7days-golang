package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
// 返回一个字符串数组，数组中的元素是 pattern 的各个部分，遇到 * 时停止
// 如 /p/:lang/*path/www 返回 [p, :lang, *path]
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 添加路由
// @param method 请求方法
// @param pattern 路由
// @param handler 处理方法
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	// 构成 GET POST ... 为第一层的 key
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	// 插入到 trie 树中，设置处理方法
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// 获取路由
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	// 从根节点开始查找满足 searchParts 的节点
	n := root.search(searchParts, 0)

	// 如果找到了，将参数填充到 params 中
	if n != nil {
		// 获取 pattern 的各个部分
		parts := parsePattern(n.pattern)

		// 将 searchParts 中的值填充到 params 中
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

// 获取所有路由
func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

// 处理请求
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
