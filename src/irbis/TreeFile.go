package irbis

type TreeNode struct {
	Value    string
	Children []TreeNode
	level    int
}

func (node *TreeNode) Add(value string) *TreeNode {
	child := TreeNode{Value: value}
	node.Children = append(node.Children, child)
	return node
}

func (node *TreeNode) String() string {
	return node.Value
}

type TreeFile struct {
	Roots []TreeNode
}

func arrange1(list []TreeNode, level int) {
	count := len(list)
	index := 0
	for index < count {
		next := arrange2(list, level, index, count)
		index = next
	}
}

func arrange2(list []TreeNode, level, index, count int) int {
	next := index + 1
	level2 := level + 1
	parent := list[index]
	for next < count {
		child := list[next]
		if child.level < level {
			break
		}
		if child.level == level2 {
			parent.Children = append(parent.Children, child)
		}
		next++
	}
	return next
}

func countIndent(text string) (result int) {
	length := len(text)
	for i := 0; i < length; i++ {
		if text[i] == '\t' {
			result++
		} else {
			break
		}
	}
	return
}

func (tree *TreeFile) AddRoot(value string) *TreeNode {
	result := new(TreeNode)
	result.Value = value
	tree.Roots = append(tree.Roots, *result)
	result = &tree.Roots[len(tree.Roots)-1]
	return result
}

func (tree *TreeFile) Parse(lines []string) {
	// TODO implement properly

	if len(lines) == 0 {
		return
	}

	currentLevel := 0
	line := lines[0]
	if countIndent(line) != 0 {
		panic("Wrong indent")
	}
	list := []TreeNode{{Value: line}}

	maxLevel := 0
	for _, item := range list {
		if item.level > maxLevel {
			maxLevel = item.level
		}
	}

	for _, line := range lines[1:] {
		if len(line) == 0 {
			continue
		}
		level := countIndent(line)
		if level > (currentLevel + 1) {
			panic("Wrong level")
		}
		currentLevel = level
		line = line[currentLevel:]
		node := TreeNode{Value: line, level: currentLevel}
		list = append(list, node)
	}

	for level := 0; level < maxLevel; level++ {
		arrange1(list, level)
	}

	for i := range list {
		item := list[i]
		if item.level == 0 {
			tree.Roots = append(tree.Roots, item)
		}
	}
}
