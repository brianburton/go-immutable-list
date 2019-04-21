package immutableList

import "fmt"

type branchNode struct {
	children   []node
	nodeSize   int
	nodeHeight int
}

func (this *branchNode) next(state *iteratorState) (*iteratorState, Object) {
	if state == nil || state.currentNode != this {
		state = &iteratorState{currentNode: this, next: state}
	}
	child := this.children[state.currentIndex]
	state.currentIndex++
	if state.currentIndex == len(this.children) {
		return child.next(state.next)
	} else {
		return child.next(state)
	}
}

func createBranchNode(children []node, nodeSize int, nodeHeight int) node {
	return &branchNode{children, nodeSize, nodeHeight}
}

func (this *branchNode) size() int {
	return this.nodeSize
}

func (this *branchNode) get(index int) Object {
	for _, child := range this.children {
		if index < child.size() {
			return child.get(index)
		}
		index -= child.size()
	}
	return nil
}

func (this *branchNode) appendNode(other node) (node, node) {
	var children []node
	thisHeight := this.height()
	thatHeight := other.height()
	thisLen := len(this.children)
	if thatHeight > thisHeight {
		panic(fmt.Sprintf("appendNode called with larger node as argument: thisHeight=%d thatHeight=%d", thisHeight, thatHeight))
	} else if thatHeight == thisHeight {
		that := other.(*branchNode)
		thatLen := len(that.children)
		children = make([]node, thisLen+thatLen)
		copy(children[0:], this.children)
		copy(children[thisLen:], that.children)
	} else {
		replacement, extra := this.children[thisLen-1].appendNode(other)
		if extra == nil {
			children = make([]node, thisLen)
			copy(children, this.children[0:thisLen-1])
			children[thisLen-1] = replacement
		} else {
			newLen := thisLen + 1
			children = make([]node, newLen)
			copy(children, this.children[0:thisLen-1])
			children[newLen-2] = replacement
			children[newLen-1] = extra
		}
	}
	return splitIfNecessary(children, this.size()+other.size(), this.nodeHeight)
}

func (this *branchNode) prependNode(other node) (node, node) {
	var children []node
	thisHeight := this.height()
	thatHeight := other.height()
	myLen := len(this.children)
	if thatHeight > thisHeight {
		panic(fmt.Sprintf("prependNode called with larger node as argument: thisHeight=%d thatHeight=%d", thisHeight, thatHeight))
	} else if thatHeight == thisHeight {
		thatBranch := other.(*branchNode)
		thatLen := len(thatBranch.children)
		children = make([]node, myLen+thatLen)
		copy(children[0:], thatBranch.children)
		copy(children[thatLen:], this.children)
	} else {
		replacement, extra := this.children[0].prependNode(other)
		if extra == nil {
			children = make([]node, myLen)
			children[0] = replacement
			copy(children[1:], this.children[1:])
		} else {
			children = make([]node, myLen+1)
			children[0] = replacement
			children[1] = extra
			copy(children[2:], this.children[1:])
		}
	}
	return splitIfNecessary(children, this.size()+other.size(), this.nodeHeight)
}

func (this *branchNode) append(value Object) (node, node) {
	lastIndex := len(this.children) - 1
	replacement, extra := this.children[lastIndex].append(value)
	return replaceImpl(this.nodeSize, this.nodeHeight, this.children, lastIndex, replacement, extra)
}

func (this *branchNode) insert(indexBefore int, value Object) (node, node) {
	childIndex, childOffset := findChildForIndex(indexBefore, this.children)
	replacement, extra := this.children[childIndex].insert(childOffset, value)
	return replaceImpl(this.nodeSize, this.nodeHeight, this.children, childIndex, replacement, extra)
}

func (this *branchNode) set(index int, value Object) node {
	childIndex, childOffset := findChildForIndex(index, this.children)
	newNode := this.children[childIndex].set(childOffset, value)
	newChildren := make([]node, len(this.children))
	copy(newChildren, this.children)
	newChildren[childIndex] = newNode
	return createBranchNode(newChildren, this.nodeSize, this.nodeHeight)
}

func (this *branchNode) forEach(proc Processor) {
	for _, child := range this.children {
		child.forEach(proc)
	}
}

func (this *branchNode) visit(start int, limit int, visitor Visitor) {
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

func (this *branchNode) height() int {
	return this.nodeHeight
}

func (this *branchNode) visitNodesOfHeight(targetHeight int, proc nodeProcessor) {
	myHeight := this.height()
	if myHeight == targetHeight {
		proc(this)
	} else if myHeight > targetHeight {
		for _, child := range this.children {
			child.visitNodesOfHeight(targetHeight, proc)
		}
	}
}

func (this *branchNode) isComplete() bool {
	return len(this.children) >= minPerNode
}

func (this *branchNode) mergeWith(other node) node {
	otherBranch := other.(*branchNode)
	myLen := len(this.children)
	otherLen := len(otherBranch.children)
	newChildren := make([]node, myLen+otherLen)
	copy(newChildren[0:], this.children)
	copy(newChildren[myLen:], otherBranch.children)
	newTotalSize := this.nodeSize + otherBranch.nodeSize
	return createBranchNode(newChildren, newTotalSize, this.nodeHeight)
}

func (this *branchNode) delete(index int) node {
	childIndex, childOffset := findChildForIndex(index, this.children)
	oldLen := len(this.children)
	var newChildren []node
	newChild := this.children[childIndex].delete(childOffset)
	if newChild.isComplete() {
		newChildren = make([]node, oldLen)
		copy(newChildren, this.children)
		newChildren[childIndex] = newChild
	} else {
		newChildren = make([]node, oldLen-1)
		if childIndex == 0 {
			newChild = newChild.mergeWith(this.children[1])
			newChildren[0] = newChild
			copy(newChildren[1:], this.children[2:])
		} else {
			newChild = this.children[childIndex-1].mergeWith(newChild)
			copy(newChildren[0:], this.children[0:childIndex-1])
			newChildren[childIndex-1] = newChild
			copy(newChildren[childIndex:], this.children[childIndex+1:])
		}
	}
	if len(newChildren) == 1 {
		return newChildren[0]
	} else {
		return createBranchNode(newChildren, computeNodeSize(newChildren), this.nodeHeight)
	}
}

func findChildForIndex(indexBefore int, children []node) (childIndex int, childOffset int) {
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

func replaceImpl(oldSize int, nodeHeight int, oldChildren []node, replaceIndex int, replacement node, extra node) (node, node) {
	var newChildren []node
	oldLen := len(oldChildren)
	if extra == nil {
		newChildren = make([]node, oldLen)
		copy(newChildren, oldChildren)
		newChildren[replaceIndex] = replacement
	} else {
		newChildren = make([]node, oldLen+1)
		copy(newChildren, oldChildren[0:replaceIndex])
		newChildren[replaceIndex] = replacement
		newChildren[replaceIndex+1] = extra
		copy(newChildren[replaceIndex+2:], oldChildren[replaceIndex+1:])
	}
	return splitIfNecessary(newChildren, oldSize+1, nodeHeight)
}

func splitIfNecessary(newChildren []node, nodeSize int, nodeHeight int) (node, node) {
	newLen := len(newChildren)
	if newLen <= maxPerNode {
		return createBranchNode(newChildren, nodeSize, nodeHeight), nil
	} else {
		firstLen := newLen / 2
		secondLen := newLen - firstLen
		first := make([]node, firstLen)
		second := make([]node, secondLen)
		copy(first, newChildren[0:])
		copy(second, newChildren[firstLen:])
		return createBranchNode(first, computeNodeSize(first), nodeHeight), createBranchNode(second, computeNodeSize(second), nodeHeight)
	}
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
