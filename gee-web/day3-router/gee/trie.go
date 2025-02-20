package gee

import (
	"fmt"
	"strings"
)

// trie 节点
type node struct {
	// pattern 是待匹配路由，例如 /p/:lang
	pattern string

	// part 是路由中的key部分，例如 :lang
	part string

	// 该节点的子节点
	children []*node

	// 是否是精确匹配，part 含有 : 或 * 时为false
	isWild bool
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

func (n *node) insert(pattern string, parts []string, height int) {
	// 递归终止条件
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// 获取当前层级的部分
	part := parts[height]

	// 查找当前层级是否已经存在该部分
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	// 递归调用
	child.insert(pattern, parts, height+1)
}

// search 查找 parts 中的路由
func (n *node) search(parts []string, height int) *node {
	// 递归终止条件
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]             // 获取当前层级的部分
	children := n.matchChildren(part) // 查找当前层级是否存在该部分

	// 遍历当前层级的子节点
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

func (n *node) travel(list *([]*node)) {
	// 递归终止条件
	if n.pattern != "" {
		*list = append(*list, n)
	}

	// 遍历当前层级的子节点
	for _, child := range n.children {
		// 递归调用
		child.travel(list)
	}
}

// 第一次匹配，用于插入
func (n *node) matchChild(part string) *node {
	// 遍历当前层级的子节点
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 全部匹配，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
