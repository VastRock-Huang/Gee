package gee

import (
	"fmt"
	"strings"
)

//前缀树部分
//前缀树结点结构体
type node struct {
	pattern string   //待匹配的完整路由
	part string      //路由中部分前缀
	children []*node //子结点
	isWild bool      //是否动态匹配
}

//从子结点中找寻匹配当前模式字符串的一个结点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		//若子结点前缀与字符串相等, 或当前结点可动态匹配, 则返回结点
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

//从子结点中找寻匹配当前模式字符串的所有结点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		//若子结点前缀相同, 或可动态匹配
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 插入路由字符串对应的结点
//pattern完整路由, parts前缀字符串集, height高度 表示当前需访问的前缀深度
func (n *node) insert(pattern string, parts []string, depth int) {
	//前缀字符串数和深度相同, 表明为最后的叶子结点
	if len(parts) == depth {
		//if n.pattern != "" {
		//	//路由冲突待解决
		//}
		n.pattern = pattern	//记录完整路由
		return
	}
	//前缀字符串数和深度不同,非叶子结点
	part := parts[depth] 	//当前深度前缀
	child := n.matchChild(part)		//找到匹配的一个子结点
	//结点不存在则构建结点
	if child == nil {
		child = &node{
			part: part,
			//若前缀以:或*开头则为动态匹配
			isWild: part[0] == ':' || part[0] == '*',
		}
		//添加到子结点中
		n.children = append(n.children, child)
	}
	//插入下一个前缀
	child.insert(pattern, parts, depth+1)
}

//查询满足前缀字符串集的一个结点
//parts查询前缀字符串集, 查询深度
func (n *node) search(parts []string, depth int) *node {
	//若前缀字符串数与深度相等,即叶子结点;或者当前结点支持动态匹配
	if len(parts) == depth || strings.HasPrefix(n.part, "*") {
		//若当前结点不是终止结点, 则返回空
		if n.pattern == "" {
			return nil
		}
		return n	//否则返回结点
	}
	//前缀字符串数与当前深度不同,且不能动态匹配
	part := parts[depth]	//当前深度的字符串
	children := n.matchChildren(part)	//找到匹配的所有子结点
	//遍历所有匹配的结点,继续查找
	for _, child := range children {
		result := child.search(parts, depth+1)
		//若找到则返回
		if result != nil {
			return result
		}
	}
	return nil	//未找到返回空
}

//装换为字符串输出
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}",
		n.pattern, n.part, n.isWild)
}

//遍历前缀树结点的子树并添加到list中
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}