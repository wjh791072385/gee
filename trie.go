package goWebGee

import "strings"

type node struct {
	pattern  string
	part     string
	isWild   bool // 是否精确匹配，part 含有 : 或 * 时为true
	children []*node
}

// 返回第一个匹配成功的结点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

//	返回所有匹配成功的结点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// pattern是整个路径，parts是pattern按照/分割开的切片
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		//只有最后一个结点才会设置pattern
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	// insert
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	//查找出所有匹配的结点，继续向下匹配
	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil

}
