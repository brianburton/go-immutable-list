package immutableList

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
	return &listImpl{this.leaves.build()}
}

func createLeafBuilder() *leafBuilder {
	var answer leafBuilder
	answer.buffer = make([]Object, maxPerNode)
	return &answer
}

func (this *leafBuilder) createLeafNodeOfLength(count int) node {
	contents := make([]Object, count)
	copy(contents, this.buffer)
	return createLeafNode(contents)
}

func (this *leafBuilder) addValue(value Object) {
	this.buffer[this.count] = value
	this.count++
	if this.count == maxPerNode {
		if this.parent == nil {
			this.parent = createBranchBuilder(2)
		}
		this.parent.addNode(this.createLeafNodeOfLength(minPerNode))
		copy(this.buffer[0:], this.buffer[minPerNode:this.count])
		this.count -= minPerNode
	}
}

func (this *leafBuilder) build() node {
	if this.count == 0 {
		return sharedEmptyInstance
	} else if this.parent == nil {
		return this.createLeafNodeOfLength(this.count)
	} else {
		return this.parent.build(this.createLeafNodeOfLength(this.count))
	}
}

func (this *leafBuilder) computeSize() int {
	answer := this.count
	if this.parent != nil {
		answer += this.parent.computeSize()
	}
	return answer
}

func createBranchBuilder(height int) *branchBuilder {
	return &branchBuilder{buffer: make([]node, maxPerNode), height: height}
}

func (this *branchBuilder) addNode(node node) {
	this.buffer[this.count] = node
	this.count++
	if this.count == maxPerNode {
		if this.parent == nil {
			this.parent = createBranchBuilder(this.height + 1)
		}
		this.parent.addNode(this.createBranchNodeOfLength(minPerNode))
		copy(this.buffer[0:], this.buffer[minPerNode:this.count])
		this.count -= minPerNode
	}
}

func (this *branchBuilder) build(extra node) node {
	this.buffer[this.count] = extra
	var answer node
	answer = this.createBranchNodeOfLength(this.count + 1)
	if this.parent != nil {
		answer = this.parent.build(answer)
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

func (this *branchBuilder) createBranchNodeOfLength(count int) node {
	children := make([]node, count)
	copy(children, this.buffer)
	nodeSize := computeBranchNodeSize(children)
	nodeHeight := this.height
	return createBranchNode(children, nodeSize, nodeHeight)
}
