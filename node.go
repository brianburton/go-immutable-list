package immutableList

import "fmt"

type node interface {
	size() int
	get(index int) Object
	getFirst() Object
	getLast() Object
	append(value Object) node
	prepend(value Object) node
	appendNode(n node) node
	prependNode(n node) node
	insert(index int, value Object) node
	delete(index int) node
	set(index int, value Object) node
	head(index int) node
	tail(index int) node
	pop() (Object, node)
	depth() int
	forEach(proc Processor)
	visit(base int, start int, limit int, v Visitor)
	checkInvariants(report reporter, isRoot bool)
	rotateLeft(parentLeft node) node
	rotateRight(parentRight node) node
	next(state *iteratorState) (*iteratorState, Object)
	left() node
	right() node
}

type iteratorState struct {
	next         *iteratorState
	currentNode  node
	currentIndex int
}

type iteratorImpl struct {
	state *iteratorState
	value Object
}

func createIterator(n node) Iterator {
	var state *iteratorState
	if n.size() == 0 {
		state = nil
	} else {
		state = &iteratorState{currentNode: n}
	}
	return &iteratorImpl{state: state}
}

func (this *iteratorImpl) Next() bool {
	if this.state == nil {
		return false
	}
	this.state, this.value = this.state.currentNode.next(this.state)
	return true
}

func (this *iteratorImpl) Get() Object {
	return this.value
}

const (
	maxValuesPerLeaf = 32
)

type leafNode struct {
	values []Object
}

func createSingleValueLeafNode(value Object) node {
	values := make([]Object, 1)
	values[0] = value
	return createMultiValueLeafNode(values)
}

func createMultiValueLeafNode(values []Object) node {
	return &leafNode{values: values}
}

func (a *leafNode) get(index int) Object {
	return a.values[index]
}

func (a *leafNode) getFirst() Object {
	return a.values[0]
}

func (a *leafNode) getLast() Object {
	return a.values[len(a.values)-1]
}

func (a *leafNode) pop() (Object, node) {
	return a.values[0], a.delete(0)
}

func (a *leafNode) set(index int, value Object) node {
	currentSize := len(a.values)
	if index < 0 || index >= currentSize {
		panic(fmt.Sprintf("invalid index for leaf node: %d", index))
	}
	newValues := make([]Object, currentSize)
	copy(newValues, a.values)
	newValues[index] = value
	return createMultiValueLeafNode(newValues)
}

func (a *leafNode) insert(index int, value Object) node {
	currentSize := len(a.values)
	if index < 0 || index > currentSize {
		panic(fmt.Sprintf("invalid index for leaf node: %d", index))
	}

	if index == 0 {
		return a.prepend(value)
	} else if index == currentSize {
		return a.append(value)
	} else if currentSize < maxValuesPerLeaf {
		values := make([]Object, currentSize+1)
		copy(values[0:], a.values[0:index])
		values[index] = value
		copy(values[(index+1):], a.values[index:])
		return createMultiValueLeafNode(values)
	} else {
		left := make([]Object, index)
		copy(left[0:], a.values[0:index])

		right := make([]Object, currentSize+1-index)
		right[0] = value
		copy(right[1:], a.values[index:])
		return createBranchNode(createMultiValueLeafNode(left), createMultiValueLeafNode(right))
	}
}

func (a *leafNode) delete(index int) node {
	currentSize := len(a.values)
	if index < 0 || index >= currentSize {
		panic(fmt.Sprintf("invalid index for leaf node: %d", index))
	}
	if len(a.values) == 1 {
		return createEmptyLeafNode()
	}
	values := make([]Object, currentSize-1)
	if index == 0 {
		copy(values[0:], a.values[1:])
	} else if index == currentSize-1 {
		copy(values[0:], a.values[0:(currentSize-1)])
	} else {
		copy(values[0:], a.values[0:index])
		copy(values[index:], a.values[(index+1):])
	}
	return createMultiValueLeafNode(values)
}

func (a *leafNode) append(value Object) node {
	currentSize := len(a.values)
	if currentSize < maxValuesPerLeaf {
		values := make([]Object, currentSize+1)
		copy(values[0:], a.values[0:])
		values[currentSize] = value
		return createMultiValueLeafNode(values)
	} else {
		values := make([]Object, 1)
		values[0] = value
		return createBranchNode(a, createMultiValueLeafNode(values))
	}
}

func (a *leafNode) prepend(value Object) node {
	currentSize := len(a.values)
	if currentSize < maxValuesPerLeaf {
		values := make([]Object, currentSize+1)
		values[0] = value
		copy(values[1:], a.values[0:])
		return createMultiValueLeafNode(values)
	} else {
		values := make([]Object, 1)
		values[0] = value
		return createBranchNode(createMultiValueLeafNode(values), a)
	}
}

func (a *leafNode) forEach(proc Processor) {
	for _, value := range a.values {
		proc(value)
	}
}

func (a *leafNode) visit(base int, start int, limit int, v Visitor) {
	size := len(a.values)
	if limit > size {
		limit = size
	}
	for i := start; i < limit; i++ {
		v(base+i, a.values[i])
	}
}

func (a *leafNode) head(index int) node {
	currentSize := len(a.values)
	if index < 0 || index > currentSize {
		panic(fmt.Sprintf("invalid index for leaf node: %d", index))
	}
	if index == 0 {
		return createEmptyLeafNode()
	} else if index == currentSize {
		return a
	} else {
		values := make([]Object, index)
		copy(values[0:], a.values[0:index])
		return createMultiValueLeafNode(values)
	}
}

func (a *leafNode) tail(index int) node {
	currentSize := len(a.values)
	if index < 0 || index > currentSize {
		panic(fmt.Sprintf("invalid index for leaf node: %d", index))
	}
	if index == 0 {
		return a
	} else if index == currentSize {
		return createEmptyLeafNode()
	} else {
		values := make([]Object, currentSize-index)
		copy(values[0:], a.values[index:])
		return createMultiValueLeafNode(values)
	}
}

func (a *leafNode) left() node {
	panic("not implemented for leaf nodes")
}

func (a *leafNode) right() node {
	panic("not implemented for leaf nodes")
}

func (a *leafNode) depth() int {
	return 0
}

func (a *leafNode) size() int {
	return len(a.values)
}

func (a *leafNode) appendNode(n node) node {
	if n.size() == 0 {
		return a
	}
	if o, matches := n.(*leafNode); matches {
		combinedSize := a.size() + o.size()
		if combinedSize <= maxValuesPerLeaf {
			return appendLeafNodeValues(combinedSize, a, o)
		}
	}
	return createBranchNode(a, n)
}

func (a *leafNode) prependNode(n node) node {
	if n.size() == 0 {
		return a
	}
	if o, matches := n.(*leafNode); matches {
		combinedSize := o.size() + a.size()
		if combinedSize <= maxValuesPerLeaf {
			return appendLeafNodeValues(combinedSize, o, a)
		}
	}
	return createBranchNode(n, a)
}

func (a *leafNode) next(state *iteratorState) (*iteratorState, Object) {
	if state == nil || state.currentNode != a {
		state = &iteratorState{currentNode: a, next: state}
	}
	value := a.values[state.currentIndex]
	state.currentIndex++
	if state.currentIndex == len(a.values) {
		return state.next, value
	} else {
		return state, value
	}
}

func appendLeafNodeValues(combinedSize int, a *leafNode, b *leafNode) node {
	values := make([]Object, combinedSize)
	copy(values[0:], a.values)
	copy(values[a.size():], b.values)
	return createMultiValueLeafNode(values)
}

func (a *leafNode) checkInvariants(report reporter, isRoot bool) {
	currentSize := len(a.values)
	if currentSize < 1 || currentSize > maxValuesPerLeaf {
		report(fmt.Sprintf("incorrect size: currentSize=%d", currentSize))
	}
}

func (a *leafNode) rotateLeft(parentLeft node) node {
	panic("not implemented for leaf node")
}

func (a *leafNode) rotateRight(parentRight node) node {
	panic("not implemented for leaf node")
}

type emptyNode struct {
}

var sharedEmptyNode node = &emptyNode{}

func createEmptyLeafNode() node {
	return sharedEmptyNode
}

func (e *emptyNode) get(index int) Object {
	panic("not implemented for empty nodes")
}

func (b *emptyNode) getFirst() Object {
	panic("not implemented for empty nodes")
}

func (b *emptyNode) getLast() Object {
	panic("not implemented for empty nodes")
}

func (b *emptyNode) pop() (Object, node) {
	panic("not implemented for empty nodes")
}

func (b *emptyNode) set(index int, value Object) node {
	panic("not implemented for empty nodes")
}

func (e *emptyNode) insert(index int, value Object) node {
	if index == 0 {
		return createSingleValueLeafNode(value)
	} else {
		panic(fmt.Sprintf("invalid index for empty node: %d", index))
	}
}

func (b *emptyNode) delete(index int) node {
	panic("not implemented for empty nodes")
}

func (b *emptyNode) head(index int) node {
	if index == 0 {
		return b
	} else {
		panic(fmt.Sprintf("invalid index for empty node: %d", index))
	}
}

func (b *emptyNode) tail(index int) node {
	if index == 0 {
		return b
	} else {
		panic(fmt.Sprintf("invalid index for empty node: %d", index))
	}
}

func (e *emptyNode) append(value Object) node {
	return createSingleValueLeafNode(value)
}

func (e *emptyNode) prepend(value Object) node {
	return createSingleValueLeafNode(value)
}

func (e *emptyNode) forEach(proc Processor) {
}

func (e *emptyNode) visit(base int, start int, limit int, v Visitor) {
}

func (e *emptyNode) left() node {
	panic("not implemented for empty nodes")
}

func (e *emptyNode) right() node {
	panic("not implemented for empty nodes")
}

func (e *emptyNode) depth() int {
	return 0
}

func (e *emptyNode) size() int {
	return 0
}

func (e *emptyNode) checkInvariants(report reporter, isRoot bool) {
	if !isRoot {
		report("emptyNode: should not exist below root")
	}
}

func (e *emptyNode) rotateLeft(parentLeft node) node {
	panic("not implemented for leaf nodes")
}

func (e *emptyNode) rotateRight(parentRight node) node {
	panic("not implemented for leaf nodes")
}

func (b *emptyNode) appendNode(n node) node {
	if n.depth() != 0 {
		panic("appending branch to leaf")
	}
	return n
}

func (b *emptyNode) prependNode(n node) node {
	if n.depth() != 0 {
		panic("prepending branch to leaf")
	}
	return n
}

func (e *emptyNode) next(state *iteratorState) (*iteratorState, Object) {
	return nil, nil
}

type branchNode struct {
	leftChild  node
	rightChild node
	mySize     int
	myDepth    int
}

func createBranchNode(leftChild node, rightChild node) node {
	return &branchNode{
		leftChild:  leftChild,
		rightChild: rightChild,
		mySize:     leftChild.size() + rightChild.size(),
		myDepth:    1 + maxDepth(leftChild, rightChild),
	}
}

func createBalancedBranchNode(left node, right node) node {
	diff := left.depth() - right.depth()
	if diff > 1 {
		return left.rotateRight(right)
	} else if diff < -1 {
		return right.rotateLeft(left)
	} else {
		return createBranchNode(left, right)
	}
}

func (b *branchNode) append(value Object) node {
	return createBalancedBranchNode(b.leftChild, b.rightChild.append(value))
}

func (b *branchNode) prepend(value Object) node {
	return createBalancedBranchNode(b.leftChild.prepend(value), b.rightChild)
}

func (b *branchNode) forEach(proc Processor) {
	b.leftChild.forEach(proc)
	b.rightChild.forEach(proc)
}

func (b *branchNode) visit(base int, start int, limit int, v Visitor) {
	visitNode(b.leftChild, 0, base, start, limit, v)
	visitNode(b.rightChild, b.leftChild.size(), base, start, limit, v)
}

func visitNode(node node, offset int, base int, start int, limit int, v Visitor) {
	base += offset
	start -= offset
	limit -= offset
	if start < 0 {
		start = 0
	}
	if limit > node.size() {
		limit = node.size()
	}
	if limit > start {
		node.visit(base, start, limit, v)
	}
}

func maxDepth(leftChild node, rightChild node) int {
	leftDepth, rightDepth := leftChild.depth(), rightChild.depth()
	if leftDepth > rightDepth {
		return leftDepth
	} else {
		return rightDepth
	}
}

func depthDiff(leftChild node, rightChild node) int {
	leftDepth, rightDepth := leftChild.depth(), rightChild.depth()
	if leftDepth > rightDepth {
		return leftDepth - rightDepth
	} else {
		return rightDepth - leftDepth
	}
}

func (b *branchNode) get(index int) Object {
	leftSize := b.leftChild.size()
	if index < leftSize {
		return b.leftChild.get(index)
	} else {
		return b.rightChild.get(index - leftSize)
	}
}

func (b *branchNode) getFirst() Object {
	return b.leftChild.getFirst()
}

func (b *branchNode) getLast() Object {
	return b.rightChild.getLast()
}

func (b *branchNode) pop() (Object, node) {
	value, newLeft := b.leftChild.pop()
	if newLeft.size() == 0 {
		return value, b.rightChild
	} else {
		return value, createBalancedBranchNode(newLeft, b.rightChild)
	}
}

func (b *branchNode) set(index int, value Object) node {
	leftSize := b.leftChild.size()
	if index < leftSize {
		return createBranchNode(b.leftChild.set(index, value), b.rightChild)
	} else {
		return createBranchNode(b.leftChild, b.rightChild.set(index-leftSize, value))
	}
}

func (b *branchNode) rotateLeft(parentLeft node) node {
	if b.leftChild.depth() > b.rightChild.depth() {
		return createBranchNode(createBranchNode(parentLeft, b.leftChild.left()), createBranchNode(b.leftChild.right(), b.rightChild))
	} else {
		return createBranchNode(createBranchNode(parentLeft, b.leftChild), b.rightChild)
	}
}

func (b *branchNode) rotateRight(parentRight node) node {
	if b.leftChild.depth() >= b.rightChild.depth() {
		return createBranchNode(b.leftChild, createBranchNode(b.rightChild, parentRight))
	} else {
		return createBranchNode(createBranchNode(b.leftChild, b.rightChild.left()), createBranchNode(b.rightChild.right(), parentRight))
	}
}

func (b *branchNode) insert(index int, value Object) node {
	var newLeft node
	var newRight node
	leftSize := b.leftChild.size()
	if index < leftSize {
		newLeft = b.leftChild.insert(index, value)
		newRight = b.rightChild
	} else {
		newLeft = b.leftChild
		newRight = b.rightChild.insert(index-leftSize, value)
	}
	return createBalancedBranchNode(newLeft, newRight)
}

func (b *branchNode) delete(index int) node {
	var newLeft, newRight node
	leftSize := b.leftChild.size()
	if index < leftSize {
		newLeft = b.leftChild.delete(index)
		newRight = b.rightChild
		if newLeft.size() == 0 {
			return newRight
		}
	} else {
		newLeft = b.leftChild
		newRight = b.rightChild.delete(index - leftSize)
		if newRight.size() == 0 {
			return newLeft
		}
	}
	return createBalancedBranchNode(newLeft, newRight)
}

func (b *branchNode) head(index int) node {
	leftSize := b.leftChild.size()
	if index < leftSize {
		return b.leftChild.head(index)
	} else {
		newRight := b.rightChild.head(index - leftSize)
		return appendNodes(b.leftChild, newRight)
	}
}

func (b *branchNode) tail(index int) node {
	leftSize := b.leftChild.size()
	if index < leftSize {
		newLeft := b.leftChild.tail(index)
		return appendNodes(newLeft, b.rightChild)
	} else {
		return b.rightChild.tail(index - leftSize)
	}
}

func (b *branchNode) left() node {
	return b.leftChild
}

func (b *branchNode) right() node {
	return b.rightChild
}

func (b *branchNode) depth() int {
	return b.myDepth
}

func (b *branchNode) size() int {
	return b.mySize
}

func (b *branchNode) appendNode(n node) node {
	if n.depth() > b.depth() {
		panic("appending larger node to smaller node")
	}
	if depthDiff(n, b) <= 1 {
		return createBranchNode(b, n)
	} else {
		return createBalancedBranchNode(b.leftChild, b.rightChild.appendNode(n))
	}
}

func (b *branchNode) prependNode(n node) node {
	if n.depth() > b.depth() {
		panic("prepending larger node to smaller node")
	}
	if depthDiff(n, b) <= 1 {
		return createBranchNode(n, b)
	} else {
		return createBalancedBranchNode(b.leftChild.prependNode(n), b.rightChild)
	}
}

func appendNodes(a node, b node) node {
	if a.size() == 0 {
		return b
	} else if b.size() == 0 {
		return a
	} else if a.depth() < b.depth() {
		return b.prependNode(a)
	} else {
		return a.appendNode(b)
	}
}

func (b *branchNode) checkInvariants(report reporter, isRoot bool) {
	if b.depth() != maxDepth(b.leftChild, b.rightChild)+1 {
		report(fmt.Sprintf("incorrect depth: depth=%d leftDepth=%d rightDepth=%d", b.depth(), b.leftChild.depth(), b.rightChild.depth()))
	}
	if depthDiff(b.leftChild, b.rightChild) > 1 {
		report(fmt.Sprintf("invalid child depths: leftDepth=%d rightDepth=%d", b.leftChild.depth(), b.rightChild.depth()))
	}
	if b.size() != b.leftChild.size()+b.rightChild.size() {
		report(fmt.Sprintf("incorrect size: size=%d leftSize=%d rightSize=%d", b.size(), b.leftChild.size(), b.rightChild.size()))
	}
	b.leftChild.checkInvariants(report, false)
	b.rightChild.checkInvariants(report, false)
}

func (b *branchNode) next(state *iteratorState) (*iteratorState, Object) {
	if state == nil || state.currentNode != b {
		state = &iteratorState{currentNode: b, next: state}
	}
	switch state.currentIndex {
	case 0:
		state.currentIndex = 1
		return b.leftChild.next(state)
	case 1:
		state.currentIndex = 2
		return b.rightChild.next(state.next)
	default:
		panic("invalid index in iterator state")
	}
}
