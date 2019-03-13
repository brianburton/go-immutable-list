package immutableList

type branchNode struct {
	children  []node
	totalSize int
}

func createBranchNode(nodeBuffer []node, count int) node {
	totalSize := 0
	children := make([]node, count)
	for i := 0; i < count; i++ {
		children[i] = nodeBuffer[i]
		totalSize += nodeBuffer[i].size()
	}
	return branchNode{children, totalSize}
}

func (this branchNode) size() int {
	return this.totalSize
}

func (this branchNode) get(index int) Object {
	for _, child := range this.children {
		if index < child.size() {
			return child.get(index)
		}
		index -= child.size()
	}
	return nil
}

func (this branchNode) append(value Object) (node, node) {
	lastIndex := len(this.children) - 1
	replacement, extra := this.children[lastIndex].append(value)
	return replaceImpl(this.totalSize, this.children, lastIndex, replacement, extra)
}

func (this branchNode) insert(indexBefore int, value Object) (node, node) {
	childIndex, childOffset := findChildForInsert(indexBefore, this.children)
	replacement, extra := this.children[childIndex].insert(childOffset, value)
	return replaceImpl(this.totalSize, this.children, childIndex, replacement, extra)
}

func findChildForInsert(indexBefore int, children []node) (int, int) {
	childIndex := 0
	childStart := 0
	for _, child := range children {
		childEnd := childStart + child.size()
		if childStart <= indexBefore && indexBefore < childEnd {
			break
		}
		childStart = childEnd
		childIndex++
	}
	return childIndex, indexBefore - childStart
}

func replaceImpl(totalSize int, children []node, replaceIndex int, replacement node, extra node) (node, node) {
	newSize := totalSize + 1
	if extra == nil {
		newChildren := replaceNode(replaceIndex, replacement, children)
		return branchNode{children: newChildren, totalSize: newSize}, nil
	} else if len(children) < maxPerNode {
		newChildren := insertReplaceNode(replaceIndex, replacement, extra, children)
		return branchNode{children: newChildren, totalSize: newSize}, nil
	} else {
		first, second := splitInsertReplaceNode(replaceIndex, replacement, extra, children)
		child1 := branchNode{children: first, totalSize: computeNodeSize(first)}
		child2 := branchNode{children: second, totalSize: computeNodeSize(second)}
		return child1, child2
	}
}

func (this branchNode) forEach(proc Processor) {
	for _, child := range this.children {
		child.forEach(proc)
	}
}

func (this branchNode) visit(start int, limit int, visitor Visitor) {
	childStart := 0
	for _, child := range this.children {
		if childStart >= limit {
			break
		}
		childEnd := childStart + child.size()
		if childStart >= start {
			childLimit := minInt(limit, childEnd)
			child.visit(start-childStart, childLimit-childStart, func(index int, obj Object) {
				visitor(childStart+index, obj)
			})
			start = childLimit
		}
		childStart = childEnd
	}
}

func (this branchNode) height() int {
	return 1 + this.children[0].height()
}

func (this branchNode) maxCompleteHeight() int {
	if len(this.children) >= minPerNode {
		return this.height()
	} else {
		return this.children[0].maxCompleteHeight()
	}
}

func (this branchNode) visitNodesOfHeight(targetHeight int, proc nodeProcessor) {
	myHeight := this.height()
	if myHeight == targetHeight {
		proc(this)
	} else if myHeight > targetHeight {
		for _, child := range this.children {
			child.visitNodesOfHeight(targetHeight, proc)
		}
	}
}

func replaceNode(replaceIndex int, replacement node, from []node) []node {
	newNodes := make([]node, len(from))
	for i, v := range from {
		if i == replaceIndex {
			v = replacement
		}
		newNodes[i] = v
	}
	return newNodes
}

func insertReplaceNode(replaceIndex int, replacement node, extra node, from []node) []node {
	newNodes := make([]node, len(from)+1)

	index := 0
	insert := func(obj node) {
		newNodes[index] = obj
		index++
	}

	for i, v := range from {
		if i == replaceIndex {
			insert(replacement)
			insert(extra)
		} else {
			insert(v)
		}
	}

	return newNodes
}

func splitInsertReplaceNode(replaceIndex int, replacement node, extra node, from []node) ([]node, []node) {
	newLen := len(from) + 1
	secondLen := newLen / 2
	firstLen := newLen - secondLen
	first := make([]node, firstLen)
	second := make([]node, secondLen)

	index := 0
	insert := func(obj node) {
		if index < firstLen {
			first[index] = obj
		} else {
			second[index-firstLen] = obj
		}
		index++
	}

	for i, v := range from {
		if i == replaceIndex {
			insert(replacement)
			insert(extra)
		} else {
			insert(v)
		}
	}

	return first, second
}

func computeNodeSize(children []node) int {
	answer := 0
	for _, child := range children {
		answer += child.size()
	}
	return answer
}

func minInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
