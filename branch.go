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
	if index < 0 || index >= this.nodeSize {
		panic(fmt.Sprintf("index out of bounds: size=%d index=%d", this.nodeSize, index))
	}
	for _, child := range this.children {
		if index < child.size() {
			return child.get(index)
		}
		index -= child.size()
	}
	panic("unreachable")
}

func (this *branchNode) getFirst() Object {
	return this.children[0].getFirst()
}

func (this *branchNode) getLast() Object {
	return this.children[len(this.children)-1].getLast()
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
	return createBranchNodesFromArray(children, this.size()+other.size(), this.nodeHeight)
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
	return createBranchNodesFromArray(children, this.size()+other.size(), this.nodeHeight)
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

func (this *branchNode) head(index int) node {
	if index < 0 || index > this.nodeSize {
		panic(fmt.Sprintf("index out of bounds: size=%d index=%d", this.nodeSize, index))
	}

	if index == this.nodeSize {
		return this
	}

	childIndex, childOffset := findChildForIndex(index, this.children)
	newChild := this.children[childIndex].head(childOffset)
	if childIndex == 0 {
		return newChild
	}

	newChildren := make([]node, childIndex)
	copy(newChildren, this.children[0:childIndex])
	newBranch := createBranchNode(newChildren, computeBranchNodeSize(newChildren), this.nodeHeight)
	if newChild.size() == 0 {
		return newBranch
	}

	appended, extra := newBranch.appendNode(newChild)
	if extra != nil {
		panic("extra should never be non-nil here")
	}
	return appended
}

func (this *branchNode) tail(index int) node {
	if index < 0 || index > this.nodeSize {
		panic(fmt.Sprintf("index out of bounds: size=%d index=%d", this.nodeSize, index))
	}

	if index == this.nodeSize {
		return sharedEmptyNodeInstance
	}

	myLen := len(this.children)
	childIndex, childOffset := findChildForIndex(index, this.children)
	oldChild := this.children[childIndex]
	newChild := oldChild.tail(childOffset)
	if childIndex == (myLen - 1) {
		return newChild
	}

	newLen := myLen - (childIndex + 1)
	newChildren := make([]node, newLen)
	copy(newChildren, this.children[childIndex+1:])
	newBranch := createBranchNode(newChildren, computeBranchNodeSize(newChildren), this.nodeHeight)
	if newChild.size() == 0 {
		return newBranch
	}

	prepended, extra := newBranch.prependNode(newChild)
	if extra != nil {
		panic("extra should never be non-nil here")
	}
	return prepended
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
		return createBranchNode(newChildren, this.nodeSize-1, this.nodeHeight)
	}
}

func (this *branchNode) pop() (Object, node) {
	value, newChild := this.children[0].pop()

	var newChildren []node
	if newChild.isComplete() {
		newChildren = make([]node, len(this.children))
		newChildren[0] = newChild
		copy(newChildren[1:], this.children[1:])
	} else {
		newChildren = make([]node, len(this.children)-1)
		newChild = newChild.mergeWith(this.children[1])
		newChildren[0] = newChild
		copy(newChildren[1:], this.children[2:])
	}

	var newNode node
	if len(newChildren) == 1 {
		newNode = newChildren[0]
	} else {
		newNode = createBranchNode(newChildren, this.nodeSize-1, this.nodeHeight)
	}
	return value, newNode
}

func (this *branchNode) checkInvariants(report reporter, isRoot bool) {
	numValues := len(this.children)
	if (numValues == 0) || (numValues < minPerNode && !isRoot) {
		report(fmt.Sprintf("branchNode: too few values: numValues=%d root=%t", numValues, isRoot))
	}
	if numValues > maxPerNode {
		report(fmt.Sprintf("branchNode: too many values: %d", numValues))
	}
	if computedSize := computeBranchNodeSize(this.children); computedSize != this.nodeSize {
		report(fmt.Sprintf("branchNode: incorrect node size: actual=%d expected=%d", computedSize, this.nodeSize))
	}
	for _, child := range this.children {
		child.checkInvariants(report, false)
		if child.height() != this.nodeHeight-1 {
			report(fmt.Sprintf("branchNode: incorrect child depth: parent=%d child=%d", this.height(), child.height()))
		}
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
	return createBranchNodesFromArray(newChildren, oldSize+1, nodeHeight)
}

func createBranchNodesFromArray(newChildren []node, nodeSize int, nodeHeight int) (node, node) {
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
		return createBranchNode(first, computeBranchNodeSize(first), nodeHeight), createBranchNode(second, computeBranchNodeSize(second), nodeHeight)
	}
}

func computeBranchNodeSize(children []node) int {
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
