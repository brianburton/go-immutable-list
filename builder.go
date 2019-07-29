package immutableList

const (
	minBranchBuilderHeight = 2
	maxPerBuilderNode      = minPerNode + maxPerNode
)

type Builder interface {
	Add(value Object) Builder
	Size() int
	Build() List
}

type builderImpl struct {
	leaves *leafBuilder
}

type leafBuilder struct {
	parent *branchBuilder
	count  int
	buffer []Object
}

type branchBuilder struct {
	parent *branchBuilder
	count  int
	height int
	buffer []node
}

func CreateBuilder() Builder {
	return &builderImpl{createLeafBuilder()}
}

func (this *builderImpl) Add(value Object) Builder {
	this.leaves.addValue(value)
	return this
}

func (this *builderImpl) Size() int {
	return this.leaves.computeSize()
}

func (this *builderImpl) Build() List {
	return createListForBuilder(this.leaves.build())
}

func createLeafBuilder() *leafBuilder {
	var answer leafBuilder
	answer.buffer = make([]Object, maxPerBuilderNode)
	return &answer
}

func (this *leafBuilder) createLeafNodeFromBuffer(start int, count int) node {
	contents := make([]Object, count)
	copy(contents, this.buffer[start:])
	return createLeafNode(contents)
}

func (this *leafBuilder) addValue(value Object) {
	this.buffer[this.count] = value
	this.count += 1
	if this.count == maxPerBuilderNode {
		if this.parent == nil {
			this.parent = createBranchBuilder(minBranchBuilderHeight)
		}
		this.parent.addNode(this.createLeafNodeFromBuffer(0, maxPerNode))
		copy(this.buffer[0:], this.buffer[maxPerNode:])
		this.count = minPerNode
	}
}

func (this *leafBuilder) build() node {
	if this.count <= maxPerNode {
		if this.parent == nil {
			return this.createLeafNodeFromBuffer(0, this.count)
		} else {
			return this.parent.build(this.createLeafNodeFromBuffer(0, this.count), nil)
		}
	}

	parent := this.parent
	if parent == nil {
		parent = createBranchBuilder(minBranchBuilderHeight)
	}
	split := this.count - minPerNode
	return parent.build(this.createLeafNodeFromBuffer(0, split), this.createLeafNodeFromBuffer(split, minPerNode))
}

func (this *leafBuilder) computeSize() int {
	answer := this.count
	if this.parent != nil {
		answer += this.parent.computeSize()
	}
	return answer
}

func createBranchBuilder(height int) *branchBuilder {
	return &branchBuilder{buffer: make([]node, maxPerBuilderNode+1), height: height}
}

func (this *branchBuilder) addNode(node node) {
	this.buffer[this.count] = node
	this.count += 1
	if this.count == maxPerBuilderNode {
		if this.parent == nil {
			this.parent = createBranchBuilder(this.height + 1)
		}
		this.parent.addNode(this.createBranchNodeFromBuffer(0, maxPerNode))
		copy(this.buffer[0:], this.buffer[maxPerNode:])
		this.count = minPerNode
	}
}

func (this *branchBuilder) build(extra1 node, extra2 node) node {
	this.buffer[this.count] = extra1
	count := this.count + 1
	if extra2 != nil {
		this.buffer[this.count+1] = extra2
		count += 1
	}

	var answer node
	if count <= maxPerNode {
		answer = this.createBranchNodeFromBuffer(0, count)
		if this.parent != nil {
			answer = this.parent.build(answer, nil)
		}
	} else {
		split := count - minPerNode
		if split > maxPerNode {
			split = maxPerNode
		}
		node1 := this.createBranchNodeFromBuffer(0, split)
		node2 := this.createBranchNodeFromBuffer(split, count-split)
		parent := this.parent
		if parent == nil {
			parent = createBranchBuilder(this.height + 1)
		}
		answer = parent.build(node1, node2)
	}

	return answer
}

func (this *branchBuilder) computeSize() int {
	answer := 0
	for i := 0; i < this.count; i++ {
		answer += this.buffer[i].size()
	}
	if this.parent != nil {
		answer += this.parent.computeSize()
	}
	return answer
}

func (this *branchBuilder) createBranchNodeFromBuffer(startOffset int, count int) node {
	children := make([]node, count)
	copy(children, this.buffer[startOffset:])
	nodeSize := computeBranchNodeSize(children)
	nodeHeight := this.height
	return createBranchNode(children, nodeSize, nodeHeight)
}
