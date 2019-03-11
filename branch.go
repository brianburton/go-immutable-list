package immutableList

type branchNode struct {
	children  []node
	totalSize int
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
	newSize := this.totalSize + 1
	oldLen := len(this.children)
	replacement, extra := this.children[oldLen-1].append(value)
	if extra == nil {
		newChildren := replaceLast(this.children, replacement)
		return branchNode{children: newChildren, totalSize: newSize}, nil
	}
	if oldLen < maxPerNode {
		newChildren := appendReplaceLast(this.children, replacement, extra)
		return branchNode{children: newChildren, totalSize: newSize}, nil
	}
	first, second := splitAppendReplaceLast(this.children, replacement, extra)
	child1 := branchNode{children: first, totalSize: computeNodeSize(first)}
	child2 := branchNode{children: second, totalSize: computeNodeSize(second)}
	return child1, child2
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

func replaceLast(from []node, replacement node) []node {
	oldLen := len(from)
	newNodes := make([]node, oldLen)
	for i, v := range from {
		newNodes[i] = v
	}
	newNodes[oldLen-1] = replacement
	return newNodes
}

func appendReplaceLast(from []node, replacement node, extra node) []node {
	oldLen := len(from)
	newNodes := make([]node, oldLen+1)
	for i, v := range from {
		newNodes[i] = v
	}
	newNodes[oldLen-1] = replacement
	newNodes[oldLen] = extra
	return newNodes
}

func splitAppendReplaceLast(from []node, replacement node, extra node) ([]node, []node) {
	oldLen := len(from)
	newLen := oldLen + 1
	firstLen := newLen - (oldLen / 2)
	secondLen := newLen - firstLen

	first := make([]node, firstLen)
	second := make([]node, secondLen)
	for i, v := range from {
		if i < firstLen {
			first[i] = v
		} else {
			second[i-firstLen] = v
		}
	}
	second[secondLen-2] = replacement
	second[secondLen-1] = extra
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
