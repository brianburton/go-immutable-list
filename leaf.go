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
	if index >= 0 && index < len(this.contents) {
		return this.contents[index]
	} else {
		panic(fmt.Sprintf("index out of bounds: size=%d index=%d", len(this.contents), index))
	}
}

func (this *leafNode) append(value Object) (node, node) {
	return this.insert(len(this.contents), value)
}

func (this *leafNode) appendNode(other node) (node, node) {
	return appendNodeImpl(this, other.(*leafNode))
}

func (this *leafNode) prependNode(other node) (node, node) {
	return appendNodeImpl(other.(*leafNode), this)
}

func (this *leafNode) insert(indexBefore int, value Object) (node, node) {
	currentSize := len(this.contents)
	if currentSize < maxPerNode {
		return createLeafNode(insertObject(indexBefore, value, this.contents)), nil
	} else {
		first, second := splitInsertObject(indexBefore, value, this.contents)
		return createLeafNode(first), createLeafNode(second)
	}
}

func (this *leafNode) set(index int, value Object) node {
	newContents := make([]Object, len(this.contents))
	copy(newContents, this.contents)
	newContents[index] = value
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

func insertObject(insertIndex int, extra Object, from []Object) []Object {
	newObjects := make([]Object, len(from)+1)
	copy(newObjects[0:], from[0:insertIndex])
	newObjects[insertIndex] = extra
	copy(newObjects[insertIndex+1:], from[insertIndex:])
	return newObjects
}

func splitInsertObject(insertIndex int, extra Object, from []Object) ([]Object, []Object) {
	newContents := insertObject(insertIndex, extra, from)
	newLen := len(newContents)
	secondLen := newLen / 2
	firstLen := newLen - secondLen
	first := make([]Object, firstLen)
	copy(first, newContents[0:firstLen])
	second := make([]Object, secondLen)
	copy(second, newContents[firstLen:])
	return first, second
}

func (this *leafNode) isComplete() bool {
	return len(this.contents) >= minPerNode
}

func (this *leafNode) mergeWith(other node) node {
	otherLeaf := other.(*leafNode)
	myLen := len(this.contents)
	otherLen := len(otherLeaf.contents)
	newLen := myLen + otherLen
	newContents := make([]Object, newLen)
	copy(newContents[0:], this.contents)
	copy(newContents[myLen:], otherLeaf.contents)
	return createLeafNode(newContents)
}

func (this *leafNode) delete(index int) node {
	oldLen := len(this.contents)
	if oldLen == 1 {
		return &emptyNode{}
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

func appendNodeImpl(a *leafNode, b *leafNode) (node, node) {
	myLen := len(a.contents)
	otherLen := len(b.contents)
	newLen := myLen + otherLen
	if newLen <= maxPerNode {
		newContents := make([]Object, newLen)
		copy(newContents[0:], a.contents)
		copy(newContents[myLen:], b.contents)
		return createLeafNode(newContents), nil
	} else {
		newContents := make([]Object, newLen)
		copy(newContents[0:], a.contents)
		copy(newContents[myLen:], b.contents)
		first := make([]Object, minPerNode)
		second := make([]Object, newLen-minPerNode)
		copy(first[0:], newContents[0:minPerNode])
		copy(second[0:], newContents[minPerNode:])
		return createLeafNode(first), createLeafNode(second)
	}
}
