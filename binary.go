package immutableList

import "fmt"

type binaryNode interface {
	get(index int) Object
	getFirst() Object
	getLast() Object
	pop() (Object, binaryNode)
	set(index int, value Object) binaryNode
	insert(index int, value Object) binaryNode
	delete(index int) binaryNode
	head(index int) binaryNode
	tail(index int) binaryNode
	left() binaryNode
	right() binaryNode
	depth() int
	size() int
	appendNode(n binaryNode) binaryNode
	prependNode(n binaryNode) binaryNode
	checkInvariants(report reporter, isRoot bool)
	rotateLeft(parentLeft binaryNode) binaryNode
	rotateRight(parentRight binaryNode) binaryNode
}

type binaryLeafNode struct {
	leftValue  Object
	rightValue Object
}

func createBinaryLeafNode(leftValue Object, rightValue Object) binaryNode {
	return &binaryLeafNode{
		leftValue:  leftValue,
		rightValue: rightValue,
	}
}

func (b *binaryLeafNode) get(index int) Object {
	switch index {
	case 0:
		return b.leftValue
	case 1:
		return b.rightValue
	default:
		panic(fmt.Sprintf("invalid index for binaryLeafNode: %d", index))
	}
}

func (b *binaryLeafNode) getFirst() Object {
	return b.leftValue
}

func (b *binaryLeafNode) getLast() Object {
	return b.rightValue
}

func (b *binaryLeafNode) pop() (Object, binaryNode) {
	return b.leftValue, createSingleLeafNode(b.rightValue)
}

func (b *binaryLeafNode) set(index int, value Object) binaryNode {
	switch index {
	case 0:
		return createBinaryLeafNode(value, b.rightValue)
	case 1:
		return createBinaryLeafNode(b.leftValue, value)
	default:
		panic(fmt.Sprintf("invalid index for binaryLeftNode: %d", index))
	}
}

func (b *binaryLeafNode) insert(index int, value Object) binaryNode {
	switch index {
	case 0:
		return createBinaryBranchNode(createSingleLeafNode(value), b)
	case 1:
		return createBinaryBranchNode(createSingleLeafNode(b.leftValue), createBinaryLeafNode(value, b.rightValue))
	case 2:
		return createBinaryBranchNode(b, createSingleLeafNode(value))
	default:
		panic(fmt.Sprintf("invalid index for binaryLeftNode: %d", index))
	}
}

func (b *binaryLeafNode) delete(index int) binaryNode {
	switch index {
	case 0:
		return createSingleLeafNode(b.rightValue)
	case 1:
		return createSingleLeafNode(b.leftValue)
	case 2:
		return createEmptyBinaryNode()
	default:
		panic(fmt.Sprintf("invalid index for binaryLeftNode: %d", index))
	}
}

func (b *binaryLeafNode) head(index int) binaryNode {
	switch index {
	case 0:
		return createEmptyBinaryNode()
	case 1:
		return createSingleLeafNode(b.leftValue)
	case 2:
		return b
	default:
		panic(fmt.Sprintf("invalid index for binaryLeftNode: %d", index))
	}
}

func (b *binaryLeafNode) tail(index int) binaryNode {
	switch index {
	case 0:
		return b
	case 1:
		return createSingleLeafNode(b.rightValue)
	case 2:
		return createEmptyBinaryNode()
	default:
		panic(fmt.Sprintf("invalid index for binaryLeftNode: %d", index))
	}
}

func (b *binaryLeafNode) left() binaryNode {
	panic("not implemented for leaf nodes")
}

func (b *binaryLeafNode) right() binaryNode {
	panic("not implemented for leaf nodes")
}

func (b *binaryLeafNode) depth() int {
	return 0
}

func (b *binaryLeafNode) size() int {
	return 2
}

func (b *binaryLeafNode) checkInvariants(report reporter, isRoot bool) {
}

func (b *binaryLeafNode) rotateLeft(parentLeft binaryNode) binaryNode {
	panic("not implemented for leaf nodes")
}

func (b *binaryLeafNode) rotateRight(parentRight binaryNode) binaryNode {
	panic("not implemented for leaf nodes")
}

func (b *binaryLeafNode) appendNode(n binaryNode) binaryNode {
	if n.depth() != 0 {
		panic("appending branch to leaf")
	}
	return createBinaryBranchNode(b, n)
}

func (b *binaryLeafNode) prependNode(n binaryNode) binaryNode {
	if n.depth() != 0 {
		panic("appending branch to leaf")
	}
	return createBinaryBranchNode(n, b)
}

type singleLeafNode struct {
	value Object
}

func createSingleLeafNode(value Object) binaryNode {
	return &singleLeafNode{value: value}
}

func (s *singleLeafNode) get(index int) Object {
	if index == 0 {
		return s.value
	} else {
		panic(fmt.Sprintf("invalid index for singleLeafNode: %d", index))
	}
}

func (b *singleLeafNode) getFirst() Object {
	return b.value
}

func (b *singleLeafNode) getLast() Object {
	return b.value
}

func (b *singleLeafNode) pop() (Object, binaryNode) {
	return b.value, createEmptyBinaryNode()
}

func (b *singleLeafNode) set(index int, value Object) binaryNode {
	switch index {
	case 0:
		return createSingleLeafNode(value)
	default:
		panic(fmt.Sprintf("invalid index for binaryLeftNode: %d", index))
	}
}

func (s *singleLeafNode) insert(index int, value Object) binaryNode {
	switch index {
	case 0:
		return createBinaryLeafNode(value, s.value)
	case 1:
		return createBinaryLeafNode(s.value, value)
	default:
		panic(fmt.Sprintf("invalid index for singleLeafNode: %d", index))
	}
}

func (b *singleLeafNode) delete(index int) binaryNode {
	if index == 0 {
		return createEmptyBinaryNode()
	} else {
		panic(fmt.Sprintf("invalid index for binaryLeftNode: %d", index))
	}
}

func (b *singleLeafNode) head(index int) binaryNode {
	switch index {
	case 0:
		return createEmptyBinaryNode()
	case 1:
		return b
	default:
		panic(fmt.Sprintf("invalid index for binaryLeftNode: %d", index))
	}
}

func (b *singleLeafNode) tail(index int) binaryNode {
	switch index {
	case 0:
		return b
	case 1:
		return createEmptyBinaryNode()
	default:
		panic(fmt.Sprintf("invalid index for binaryLeftNode: %d", index))
	}
}

func (s *singleLeafNode) left() binaryNode {
	panic("not implemented for leaf nodes")
}

func (s *singleLeafNode) right() binaryNode {
	panic("not implemented for leaf nodes")
}

func (s *singleLeafNode) depth() int {
	return 0
}

func (s *singleLeafNode) size() int {
	return 1
}

func (s *singleLeafNode) checkInvariants(report reporter, isRoot bool) {
}

func (s *singleLeafNode) rotateLeft(parentLeft binaryNode) binaryNode {
	panic("not implemented for leaf nodes")
}

func (s *singleLeafNode) rotateRight(parentRight binaryNode) binaryNode {
	panic("not implemented for leaf nodes")
}

func (b *singleLeafNode) appendNode(n binaryNode) binaryNode {
	if n.depth() != 0 {
		panic("appending branch to leaf")
	}
	return createBinaryBranchNode(b, n)
}

func (b *singleLeafNode) prependNode(n binaryNode) binaryNode {
	if n.depth() != 0 {
		panic("appending branch to leaf")
	}
	return createBinaryBranchNode(n, b)
}

type emptyLeafNode struct {
}

var sharedEmptyBinaryNode binaryNode = &emptyLeafNode{}

func createEmptyBinaryNode() binaryNode {
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
		return createSingleLeafNode(value)
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

type binaryBranchNode struct {
	leftChild  binaryNode
	rightChild binaryNode
	mySize     int
	myDepth    int
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
