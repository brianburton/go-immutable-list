package immutableList

import "fmt"

type leafNode struct {
	contents []Object
}

func createLeafNode(contents []Object) node {
	return &leafNode{contents}
}

func (this *leafNode) size() int {
	return len(this.contents)
}

func (this *leafNode) get(index int) Object {
	if index < 0 || index >= len(this.contents) {
		panic(fmt.Sprintf("index out of bounds: size=%d index=%d", len(this.contents), index))
	}
	return this.contents[index]
}

func (this *leafNode) append(value Object) (node, node) {
	return this.insert(len(this.contents), value)
}

func (this *leafNode) appendNode(other node) (node, node) {
	return appendLeafNodeImpl(this, other.(*leafNode))
}

func (this *leafNode) prependNode(other node) (node, node) {
	return appendLeafNodeImpl(other.(*leafNode), this)
}

func (this *leafNode) insert(indexBefore int, value Object) (node, node) {
	newContents := make([]Object, len(this.contents)+1)
	copy(newContents, this.contents[0:indexBefore])
	newContents[indexBefore] = value
	copy(newContents[indexBefore+1:], this.contents[indexBefore:])
	return createLeafNodesFromArray(newContents)
}

func (this *leafNode) set(index int, value Object) node {
	newContents := make([]Object, len(this.contents))
	copy(newContents, this.contents)
	newContents[index] = value
	return createLeafNode(newContents)
}

func (this *leafNode) head(index int) node {
	myLen := len(this.contents)
	if index < 0 || index > myLen {
		panic(fmt.Sprintf("index out of bounds: size=%d index=%d", myLen, index))
	}
	if index == 0 {
		return sharedEmptyNodeInstance
	}
	if index == myLen {
		return this
	}
	newContents := make([]Object, index)
	copy(newContents, this.contents[0:index])
	return createLeafNode(newContents)
}

func (this *leafNode) forEach(proc Processor) {
	for _, value := range this.contents {
		proc(value)
	}
}

func (this *leafNode) visit(start int, limit int, v Visitor) {
	for i := start; i < limit; i++ {
		v(i, this.contents[i])
	}
}

func (this *leafNode) height() int {
	return 1
}

func (this *leafNode) isComplete() bool {
	return len(this.contents) >= minPerNode
}

func (this *leafNode) mergeWith(other node) node {
	otherLeaf := other.(*leafNode)
	myLen := len(this.contents)
	newContents := make([]Object, myLen+len(otherLeaf.contents))
	copy(newContents[0:], this.contents)
	copy(newContents[myLen:], otherLeaf.contents)
	return createLeafNode(newContents)
}

func (this *leafNode) delete(index int) node {
	oldLen := len(this.contents)
	if oldLen == 1 {
		return sharedEmptyNodeInstance
	} else {
		newContents := make([]Object, oldLen-1)
		copy(newContents[0:], this.contents[0:index])
		copy(newContents[index:], this.contents[index+1:])
		return createLeafNode(newContents)
	}
}

func (this *leafNode) next(state *iteratorState) (*iteratorState, Object) {
	if state == nil || state.currentNode != this {
		state = &iteratorState{currentNode: this, next: state}
	}
	value := this.contents[state.currentIndex]
	state.currentIndex++
	if state.currentIndex == len(this.contents) {
		return state.next, value
	} else {
		return state, value
	}
}

func (this *leafNode) checkInvariants(report reporter, isRoot bool) {
	numValues := len(this.contents)
	if numValues > maxPerNode {
		report(fmt.Sprintf("leafNode: too many values: %d", numValues))
	}
	if numValues < minPerNode && !isRoot {
		report(fmt.Sprintf("leafNode: too few values: %d", numValues))
	}
}

func appendLeafNodeImpl(a *leafNode, b *leafNode) (node, node) {
	myLen := len(a.contents)
	newContents := make([]Object, myLen+len(b.contents))
	copy(newContents[0:], a.contents)
	copy(newContents[myLen:], b.contents)
	return createLeafNodesFromArray(newContents)
}

func createLeafNodesFromArray(newContents []Object) (node, node) {
	newLen := len(newContents)
	if newLen <= maxPerNode {
		return createLeafNode(newContents), nil
	} else {
		firstLen := newLen / 2
		secondLen := newLen - firstLen
		first := make([]Object, firstLen)
		second := make([]Object, secondLen)
		copy(first, newContents[0:])
		copy(second, newContents[firstLen:])
		return createLeafNode(first), createLeafNode(second)
	}
}
