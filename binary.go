package immutableList

import "fmt"

type binaryNode interface {
	size() int
	get(index int) Object
	getFirst() Object
	getLast() Object
	append(value Object) binaryNode
	prepend(value Object) binaryNode
	appendNode(n binaryNode) binaryNode
	prependNode(n binaryNode) binaryNode
	insert(index int, value Object) binaryNode
	delete(index int) binaryNode
	set(index int, value Object) binaryNode
	head(index int) binaryNode
	tail(index int) binaryNode
	pop() (Object, binaryNode)
	depth() int
	forEach(proc Processor)
	visit(start int, limit int, v Visitor)
	checkInvariants(report reporter, isRoot bool)
	rotateLeft(parentLeft binaryNode) binaryNode
	rotateRight(parentRight binaryNode) binaryNode
	next(state *binaryIteratorState) (*binaryIteratorState, Object)
	left() binaryNode
	right() binaryNode
}

type binaryIteratorState struct {
	next         *binaryIteratorState
	currentNode  binaryNode
	currentIndex int
}

type binaryIteratorImpl struct {
	state *binaryIteratorState
	value Object
}

func createBinaryIterator(n binaryNode) Iterator {
	var state *binaryIteratorState
	if n.size() == 0 {
		state = nil
	} else {
		state = &binaryIteratorState{currentNode: n}
	}
	return &binaryIteratorImpl{state: state}
}

func (this *binaryIteratorImpl) Next() bool {
	if this.state == nil {
		return false
	}
	this.state, this.value = this.state.currentNode.next(this.state)
	return true
}

func (this *binaryIteratorImpl) Get() Object {
	return this.value
}

const (
	binaryArrayNodeMaxValues = 32
)

type arrayLeafNode struct {
	values []Object
}

func createSingleValueLeafNode(value Object) binaryNode {
	values := make([]Object, 1)
	values[0] = value
	return createMultiValueLeafNode(values)
}

func createMultiValueLeafNode(values []Object) binaryNode {
	return &arrayLeafNode{values: values}
}

func (a *arrayLeafNode) get(index int) Object {
	return a.values[index]
}

func (a *arrayLeafNode) getFirst() Object {
	return a.values[0]
}

func (a *arrayLeafNode) getLast() Object {
	return a.values[len(a.values)-1]
}

func (a *arrayLeafNode) pop() (Object, binaryNode) {
	return a.values[0], a.delete(0)
}

func (a *arrayLeafNode) set(index int, value Object) binaryNode {
	currentSize := len(a.values)
	if index < 0 || index >= currentSize {
		panic(fmt.Sprintf("invalid index for arrayLeafNode: %d", index))
	}
	newValues := make([]Object, currentSize)
	copy(newValues, a.values)
	newValues[index] = value
	return createMultiValueLeafNode(newValues)
}

func (a *arrayLeafNode) insert(index int, value Object) binaryNode {
	currentSize := len(a.values)
	if index < 0 || index > currentSize {
		panic(fmt.Sprintf("invalid index for arrayLeafNode: %d", index))
	}

	if index == 0 {
		return a.prepend(value)
	} else if index == currentSize {
		return a.append(value)
	} else if currentSize < binaryArrayNodeMaxValues {
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
		return createBinaryBranchNode(createMultiValueLeafNode(left), createMultiValueLeafNode(right))
	}
}

func (a *arrayLeafNode) delete(index int) binaryNode {
	currentSize := len(a.values)
	if index < 0 || index >= currentSize {
		panic(fmt.Sprintf("invalid index for arrayLeafNode: %d", index))
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

func (a *arrayLeafNode) append(value Object) binaryNode {
	currentSize := len(a.values)
	if currentSize < binaryArrayNodeMaxValues {
		values := make([]Object, currentSize+1)
		copy(values[0:], a.values[0:])
		values[currentSize] = value
		return createMultiValueLeafNode(values)
	} else {
		values := make([]Object, 1)
		values[0] = value
		return createBinaryBranchNode(a, createMultiValueLeafNode(values))
	}
}

func (a *arrayLeafNode) prepend(value Object) binaryNode {
	currentSize := len(a.values)
	if currentSize < binaryArrayNodeMaxValues {
		values := make([]Object, currentSize+1)
		values[0] = value
		copy(values[1:], a.values[0:])
		return createMultiValueLeafNode(values)
	} else {
		values := make([]Object, 1)
		values[0] = value
		return createBinaryBranchNode(createMultiValueLeafNode(values), a)
	}
}

func (a *arrayLeafNode) forEach(proc Processor) {
	for _, value := range a.values {
		proc(value)
	}
}

func (a *arrayLeafNode) visit(start int, limit int, v Visitor) {
	currentSize := len(a.values)
	if start < 0 || start > currentSize {
		panic(fmt.Sprintf("invalid index for arrayLeafNode: %d", start))
	}
	for i := start; i < limit; i++ {
		v(i, a.values[i])
	}
}

func (a *arrayLeafNode) head(index int) binaryNode {
	currentSize := len(a.values)
	if index < 0 || index > currentSize {
		panic(fmt.Sprintf("invalid index for arrayLeafNode: %d", index))
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

func (a *arrayLeafNode) tail(index int) binaryNode {
	currentSize := len(a.values)
	if index < 0 || index > currentSize {
		panic(fmt.Sprintf("invalid index for arrayLeafNode: %d", index))
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

func (a *arrayLeafNode) left() binaryNode {
	panic("not implemented for leaf node")
}

func (a *arrayLeafNode) right() binaryNode {
	panic("not implemented for leaf node")
}

func (a *arrayLeafNode) depth() int {
	return 0
}

func (a *arrayLeafNode) size() int {
	return len(a.values)
}

func (a *arrayLeafNode) appendNode(n binaryNode) binaryNode {
	if n.size() == 0 {
		return a
	}
	if o, matches := n.(*arrayLeafNode); matches {
		combinedSize := a.size() + o.size()
		if combinedSize <= binaryArrayNodeMaxValues {
			return appendLeafNodeValues(combinedSize, a, o)
		}
	}
	return createBinaryBranchNode(a, n)
}

func (a *arrayLeafNode) prependNode(n binaryNode) binaryNode {
	if n.size() == 0 {
		return a
	}
	if o, matches := n.(*arrayLeafNode); matches {
		combinedSize := o.size() + a.size()
		if combinedSize <= binaryArrayNodeMaxValues {
			return appendLeafNodeValues(combinedSize, o, a)
		}
	}
	return createBinaryBranchNode(n, a)
}

func (a *arrayLeafNode) next(state *binaryIteratorState) (*binaryIteratorState, Object) {
	if state == nil || state.currentNode != a {
		state = &binaryIteratorState{currentNode: a, next: state}
	}
	value := a.values[state.currentIndex]
	state.currentIndex++
	if state.currentIndex == len(a.values) {
		return state.next, value
	} else {
		return state, value
	}
}

func appendLeafNodeValues(combinedSize int, a *arrayLeafNode, b *arrayLeafNode) binaryNode {
	values := make([]Object, combinedSize)
	copy(values[0:], a.values)
	copy(values[a.size():], b.values)
	return createMultiValueLeafNode(values)
}

func (a *arrayLeafNode) checkInvariants(report reporter, isRoot bool) {
	currentSize := len(a.values)
	if currentSize < 1 || currentSize > binaryArrayNodeMaxValues {
		report(fmt.Sprintf("incorrect size: currentSize=%d", currentSize))
	}
}

func (a *arrayLeafNode) rotateLeft(parentLeft binaryNode) binaryNode {
	panic("not implemented for leaf node")
}

func (a *arrayLeafNode) rotateRight(parentRight binaryNode) binaryNode {
	panic("not implemented for leaf node")
}

type emptyLeafNode struct {
}

var sharedEmptyBinaryNode binaryNode = &emptyLeafNode{}

func createEmptyLeafNode() binaryNode {
	return sharedEmptyBinaryNode
}

func (e *emptyLeafNode) get(index int) Object {
	panic("not implemented for emptyLeafNodes")
}

func (b *emptyLeafNode) getFirst() Object {
	panic("not implemented for emptyLeafNodes")
}

func (b *emptyLeafNode) getLast() Object {
	panic("not implemented for emptyLeafNodes")
}

func (b *emptyLeafNode) pop() (Object, binaryNode) {
	panic("not implemented for emptyLeafNodes")
}

func (b *emptyLeafNode) set(index int, value Object) binaryNode {
	panic("not implemented for emptyLeafNodes")
}

func (e *emptyLeafNode) insert(index int, value Object) binaryNode {
	if index == 0 {
		return createSingleValueLeafNode(value)
	} else {
		panic(fmt.Sprintf("invalid index for emptyLeafNode: %d", index))
	}
}

func (b *emptyLeafNode) delete(index int) binaryNode {
	panic("not implemented for emptyLeafNodes")
}

func (b *emptyLeafNode) head(index int) binaryNode {
	if index == 0 {
		return b
	} else {
		panic(fmt.Sprintf("invalid index for emptyLeafNode: %d", index))
	}
}

func (b *emptyLeafNode) tail(index int) binaryNode {
	if index == 0 {
		return b
	} else {
		panic(fmt.Sprintf("invalid index for emptyLeafNode: %d", index))
	}
}

func (e *emptyLeafNode) append(value Object) binaryNode {
	return createSingleValueLeafNode(value)
}

func (e *emptyLeafNode) prepend(value Object) binaryNode {
	return createSingleValueLeafNode(value)
}

func (e *emptyLeafNode) forEach(proc Processor) {
}

func (e *emptyLeafNode) visit(start int, limit int, v Visitor) {
	if start != 0 || limit != 0 {
		panic(fmt.Sprintf("invalid index for emptyLeafNode: start=%d limit=%d", start, limit))
	}
}

func (e *emptyLeafNode) left() binaryNode {
	panic("not implemented for emptyLeafNodes")
}

func (e *emptyLeafNode) right() binaryNode {
	panic("not implemented for emptyLeafNodes")
}

func (e *emptyLeafNode) depth() int {
	return 0
}

func (e *emptyLeafNode) size() int {
	return 0
}

func (e *emptyLeafNode) checkInvariants(report reporter, isRoot bool) {
	if !isRoot {
		report("emptyLeafNode: should not exist below root")
	}
}

func (e *emptyLeafNode) rotateLeft(parentLeft binaryNode) binaryNode {
	panic("not implemented for leaf nodes")
}

func (e *emptyLeafNode) rotateRight(parentRight binaryNode) binaryNode {
	panic("not implemented for leaf nodes")
}

func (b *emptyLeafNode) appendNode(n binaryNode) binaryNode {
	if n.depth() != 0 {
		panic("appending branch to leaf")
	}
	return n
}

func (b *emptyLeafNode) prependNode(n binaryNode) binaryNode {
	if n.depth() != 0 {
		panic("appending branch to leaf")
	}
	return n
}

func (e *emptyLeafNode) next(state *binaryIteratorState) (*binaryIteratorState, Object) {
	return nil, nil
}

type binaryBranchNode struct {
	leftChild  binaryNode
	rightChild binaryNode
	mySize     int
	myDepth    int
}

func (b *binaryBranchNode) append(value Object) binaryNode {
	return createBalancedBinaryBranchNode(b.leftChild, b.rightChild.append(value))
}

func (b *binaryBranchNode) prepend(value Object) binaryNode {
	return createBalancedBinaryBranchNode(b.leftChild.prepend(value), b.rightChild)
}

func (b *binaryBranchNode) forEach(proc Processor) {
	b.leftChild.forEach(proc)
	b.rightChild.forEach(proc)
}

func (b *binaryBranchNode) visit(start int, limit int, v Visitor) {
	leftSize := b.leftChild.size()
	if start < leftSize {
		b.leftChild.visit(start, limit, v)
	}
	if limit > leftSize {
		b.rightChild.visit(start-leftSize, limit-leftSize, v)
	}
}

func maxDepth(leftChild binaryNode, rightChild binaryNode) int {
	leftDepth, rightDepth := leftChild.depth(), rightChild.depth()
	if leftDepth > rightDepth {
		return leftDepth
	} else {
		return rightDepth
	}
}

func depthDiff(leftChild binaryNode, rightChild binaryNode) int {
	leftDepth, rightDepth := leftChild.depth(), rightChild.depth()
	if leftDepth > rightDepth {
		return leftDepth - rightDepth
	} else {
		return rightDepth - leftDepth
	}
}

func createBinaryBranchNode(leftChild binaryNode, rightChild binaryNode) binaryNode {
	return &binaryBranchNode{
		leftChild:  leftChild,
		rightChild: rightChild,
		mySize:     leftChild.size() + rightChild.size(),
		myDepth:    1 + maxDepth(leftChild, rightChild),
	}
}

func createBalancedBinaryBranchNode(left binaryNode, right binaryNode) binaryNode {
	diff := left.depth() - right.depth()
	if diff > 1 {
		return left.rotateRight(right)
	} else if diff < -1 {
		return right.rotateLeft(left)
	} else {
		return createBinaryBranchNode(left, right)
	}
}

func (b *binaryBranchNode) get(index int) Object {
	leftSize := b.leftChild.size()
	if index < leftSize {
		return b.leftChild.get(index)
	} else {
		return b.rightChild.get(index - leftSize)
	}
}

func (b *binaryBranchNode) getFirst() Object {
	return b.leftChild.getFirst()
}

func (b *binaryBranchNode) getLast() Object {
	return b.rightChild.getLast()
}

func (b *binaryBranchNode) pop() (Object, binaryNode) {
	value, newLeft := b.leftChild.pop()
	if newLeft.size() == 0 {
		return value, b.rightChild
	} else {
		return value, createBalancedBinaryBranchNode(newLeft, b.rightChild)
	}
}

func (b *binaryBranchNode) set(index int, value Object) binaryNode {
	leftSize := b.leftChild.size()
	if index < leftSize {
		return createBinaryBranchNode(b.leftChild.set(index, value), b.rightChild)
	} else {
		return createBinaryBranchNode(b.leftChild, b.rightChild.set(index-leftSize, value))
	}
}

func (b *binaryBranchNode) rotateLeft(parentLeft binaryNode) binaryNode {
	if b.leftChild.depth() > b.rightChild.depth() {
		return createBinaryBranchNode(createBinaryBranchNode(parentLeft, b.leftChild.left()), createBinaryBranchNode(b.leftChild.right(), b.rightChild))
	} else {
		return createBinaryBranchNode(createBinaryBranchNode(parentLeft, b.leftChild), b.rightChild)
	}
}

func (b *binaryBranchNode) rotateRight(parentRight binaryNode) binaryNode {
	if b.leftChild.depth() >= b.rightChild.depth() {
		return createBinaryBranchNode(b.leftChild, createBinaryBranchNode(b.rightChild, parentRight))
	} else {
		return createBinaryBranchNode(createBinaryBranchNode(b.leftChild, b.rightChild.left()), createBinaryBranchNode(b.rightChild.right(), parentRight))
	}
}

func (b *binaryBranchNode) insert(index int, value Object) binaryNode {
	var newLeft binaryNode
	var newRight binaryNode
	leftSize := b.leftChild.size()
	if index < leftSize {
		newLeft = b.leftChild.insert(index, value)
		newRight = b.rightChild
	} else {
		newLeft = b.leftChild
		newRight = b.rightChild.insert(index-leftSize, value)
	}
	return createBalancedBinaryBranchNode(newLeft, newRight)
}

func (b *binaryBranchNode) delete(index int) binaryNode {
	var newLeft, newRight binaryNode
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
	return createBalancedBinaryBranchNode(newLeft, newRight)
}

func (b *binaryBranchNode) head(index int) binaryNode {
	leftSize := b.leftChild.size()
	if index < leftSize {
		return b.leftChild.head(index)
	} else {
		newRight := b.rightChild.head(index - leftSize)
		return appendBinaryNodes(b.leftChild, newRight)
	}
}

func (b *binaryBranchNode) tail(index int) binaryNode {
	leftSize := b.leftChild.size()
	if index < leftSize {
		newLeft := b.leftChild.tail(index)
		return appendBinaryNodes(newLeft, b.rightChild)
	} else {
		return b.rightChild.tail(index - leftSize)
	}
}

func (b *binaryBranchNode) left() binaryNode {
	return b.leftChild
}

func (b *binaryBranchNode) right() binaryNode {
	return b.rightChild
}

func (b *binaryBranchNode) depth() int {
	return b.myDepth
}

func (b *binaryBranchNode) size() int {
	return b.mySize
}

func (b *binaryBranchNode) appendNode(n binaryNode) binaryNode {
	if n.depth() > b.depth() {
		panic("appending larger node to smaller node")
	}
	if depthDiff(n, b) <= 1 {
		return createBinaryBranchNode(b, n)
	} else {
		return createBalancedBinaryBranchNode(b.leftChild, b.rightChild.appendNode(n))
	}
}

func (b *binaryBranchNode) prependNode(n binaryNode) binaryNode {
	if n.depth() > b.depth() {
		panic("prepending larger node to smaller node")
	}
	if depthDiff(n, b) <= 1 {
		return createBinaryBranchNode(n, b)
	} else {
		return createBalancedBinaryBranchNode(b.leftChild.prependNode(n), b.rightChild)
	}
}

func appendBinaryNodes(a binaryNode, b binaryNode) binaryNode {
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

func (b *binaryBranchNode) checkInvariants(report reporter, isRoot bool) {
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

func (b *binaryBranchNode) next(state *binaryIteratorState) (*binaryIteratorState, Object) {
	if state == nil || state.currentNode != b {
		state = &binaryIteratorState{currentNode: b, next: state}
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
