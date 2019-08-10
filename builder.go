package immutableList

type Builder interface {
	Add(value Object) Builder
	Size() int
	Build() List
}

type leafBuilder struct {
	parent *branchBuilder
	count  int // only zero if Add() has never been called
	buffer [binaryArrayNodeMaxValues]Object
}

type branchBuilder struct {
	parent *branchBuilder
	left   binaryNode // never nil
	right  binaryNode // may be nil
}

func CreateBuilder() Builder {
	return &leafBuilder{}
}

func (this *leafBuilder) Add(value Object) Builder {
	if this.count == binaryArrayNodeMaxValues {
		leafNode := this.createLeafFromBuffer()
		if this.parent == nil {
			this.parent = createBranchBuilder(leafNode)
		} else {
			this.parent.addChild(leafNode)
		}
		this.buffer[0] = value
		this.count = 1
	} else {
		this.buffer[this.count] = value
		this.count += 1
	}
	return this
}

func (this *leafBuilder) Size() int {
	answer := this.count
	if this.parent != nil {
		answer += this.parent.computeSize()
	}
	return answer
}

func (this *leafBuilder) Build() List {
	var root binaryNode
	if this.count == 0 {
		root = createEmptyLeafNode()
	} else if this.parent == nil {
		root = this.createLeafFromBuffer()
	} else {
		root = this.parent.build(this.createLeafFromBuffer())
	}
	return createListNode(root)
}

func (this *leafBuilder) createLeafFromBuffer() binaryNode {
	values := make([]Object, this.count)
	copy(values[0:], this.buffer[0:this.count])
	return createMultiValueLeafNode(values)
}

func createBranchBuilder(left binaryNode) *branchBuilder {
	return &branchBuilder{left: left}
}

func (this *branchBuilder) addChild(node binaryNode) {
	if this.right == nil {
		this.right = node
	} else {
		branchNode := createBinaryBranchNode(this.left, this.right)
		if this.parent == nil {
			this.parent = createBranchBuilder(branchNode)
		} else {
			this.parent.addChild(branchNode)
		}
		this.left = node
		this.right = nil
	}
}

func (this *branchBuilder) build(extra binaryNode) binaryNode {
	var answer binaryNode
	if this.right == nil {
		answer = this.left
	} else {
		answer = createBinaryBranchNode(this.left, this.right)
	}
	if this.parent != nil {
		answer = this.parent.build(answer)
	}
	answer = answer.appendNode(extra)
	return answer
}

func (this *branchBuilder) computeSize() int {
	answer := this.left.size()
	if this.right != nil {
		answer += this.right.size()
	}
	if this.parent != nil {
		answer += this.parent.computeSize()
	}
	return answer
}
