package gee

//前缀树部分
import (
	"fmt"
	"strings"
)

//前缀树结点结构体
type node struct {
	pattern   string  //待匹配的完整模式字符串
	part      string  //路由中该结点对应的前缀字符串
	children  []*node //静态路由子结点
	wildChild *node   //动态路由子结点
}

//转换字符串输出
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s}",
		n.pattern, n.part)
}

// 插入路由字符串对应的结点
//pattern完整路由, parts前缀字符串集, depth深度 表示当前需访问的前缀深度
func (n *node) insert(pattern string, parts []string, depth int) {
	//前缀字符串数和深度相同, 表明为最后的叶子结点
	if len(parts) == depth {
		//若当前结点已经存有路由,则证明为重复插入相同路由,引发错误
		if n.pattern != "" {
			panic("a handle is already registered for path '" +
				n.pattern + "'")
		}
		//添加完整路径
		n.pattern = pattern
		return
	}
	//前缀字符串数和深度不同,非叶子结点
	part := parts[depth] //当前结点深度对于的前缀字符串
	var child *node      //用于记录匹配的子结点
	//若为动态路由匹配
	if part[0] == ':' || part[0] == '*' {
		//缺少动态路由前缀名称, 引发错误
		if len(part) == 1 {
			panic("dynamic route missing name")
		}
		//若当前无动态路由子结点则构建
		if n.wildChild == nil {
			n.wildChild = &node{
				part: part,
			}
		}
		child = n.wildChild
		//若子结点动态路由前缀和当前前缀字符串不同,
		//则为相同路径下有多个动态路由, 引发错误
		if child.part != part {
			panic("'" + part + "' in new path '" + pattern +
				"' conflicts with existing wildcard '" +
				child.part + "'")
		}
	} else { //静态路由结点
		//遍历已存在的结点
		for _, ch := range n.children {
			if ch.part == part {
				child = ch
				break
			}
		}
		//不存在则构建
		if child == nil {
			child = &node{
				part: part,
			}
			//添加到静态子结点集中
			n.children = append(n.children, child)
		}
	}
	//递归插入下一个前缀
	child.insert(pattern, parts, depth+1)
}

//查询满足前缀字符串集的一个结点
func (n *node) search(parts []string, depth int) *node {
	//若前缀字符串数与深度相等,即叶子结点;或者当前结点支持动态匹配
	if len(parts) == depth || strings.HasPrefix(n.part, "*") {
		//若当前结点不是终止结点, 则返回空
		if n.pattern == "" {
			return nil
		}
		return n //否则返回结点
	}
	part := parts[depth] //当前深度的前缀字符串
	//遍历所有静态前缀子结点,看是否是静态路由前缀
	for _, child := range n.children {
		if child.part == part {
			return child.search(parts, depth+1)
		}
	}
	//若不是静态路由前缀且动态路由结点未创建
	//则不匹配返回空
	if n.wildChild == nil {
		return nil
	}
	//动态路由子结点继续查找
	return n.wildChild.search(parts, depth+1)
}

//遍历前缀树结点的子树并添加到list中
func (n *node) travel(list *[]*node) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	//添加静态子结点的子树
	for _, child := range n.children {
		child.travel(list)
	}
	//添加动态子结点的子树
	if n.wildChild != nil {
		n.wildChild.travel(list)
	}
}
