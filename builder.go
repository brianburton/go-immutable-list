package immutableList

type Builder interface {
	Add(value Object) Builder
	Size() int
	Build() List
}

type builderImpl struct {
	leaves leafBuilder
}

type leafBuilder struct {
	parent *branchBuilder
	count  int // only zero if addValue() has never been called
	buffer [binaryArrayNodeMaxValues]Object
}

type branchBuilder struct {
	parent *branchBuilder
	left   binaryNode // never nil
	right  binaryNode // may be nil
}

func CreateBuilder() Builder {
	return &builderImpl{}
}

func (this *builderImpl) Add(value Object) Builder {
	this.leaves.addValue(value)
	return this
}

func (this *builderImpl) Size() int {
	return this.leaves.computeSize()
}

func (this *builderImpl) Build() List {
	return createListNode(this.leaves.build())
}

func (this *leafBuilder) createLeafFromBuffer() binaryNode {
	values := make([]Object, this.count)
	copy(values[0:], this.buffer[0:this.count])
	return createMultiValueLeafNode(values)
}

func (this *leafBuilder) addValue(value Object) {
	if this.count == binaryArrayNodeMaxValues {
		leafNode := this.createLeafFromBuffer()
		if this.parent == nil {
			this.parent = createBranchBuilder(leafNode)
		} else {
			this.parent.addNode(leafNode)
		}
		this.buffer[0] = value
		this.count = 1
	} else {
		this.buffer[this.count] = value
		this.count += 1
	}
}

func (this *leafBuilder) build() binaryNode {
	if this.count == 0 {
		return createEmptyLeafNode()
	} else if this.parent == nil {
		return this.createLeafFromBuffer()
	} else {
		return this.parent.build(this.createLeafFromBuffer())
	}
}

func (this *leafBuilder) computeSize() int {
	answer := this.count
	if this.parent != nil {
		answer += this.parent.computeSize()
	}
	return answer
}

func createBranchBuilder(left binaryNode) *branchBuilder {
	return &branchBuilder{left: left}
}

func (this *branchBuilder) addNode(node binaryNode) {
	if this.right == nil {
		this.right = node
	} else {
		branchNode := createBinaryBranchNode(this.left, this.right)
		if this.parent == nil {
			this.parent = createBranchBuilder(branchNode)
		} else {
			this.parent.addNode(branchNode)
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
