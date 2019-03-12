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
	buffer []node
}

func CreateBuilder() Builder {
	return builderImpl{createLeafBuilder()}
}

func (this builderImpl) Add(value Object) Builder {
	this.leaves.addValue(value)
	return this
}

func (this builderImpl) Size() int {
	return this.leaves.computeSize()
}

func (this builderImpl) Build() List {
	return listImpl{this.leaves.build()}
}

func createLeafBuilder() *leafBuilder {
	var answer leafBuilder
	answer.buffer = make([]Object, maxPerNode)
	return &answer
}

func (this *leafBuilder) addValue(value Object) {
	this.buffer[this.count] = value
	this.count++
	if this.count == maxPerNode {
		if this.parent == nil {
			this.parent = createBranchBuilder()
		}
		this.parent.addNode(createLeafNode(this.buffer, minPerNode))
		for i := minPerNode; i < maxPerNode; i++ {
			this.buffer[i-minPerNode] = this.buffer[i]
		}
		this.count -= minPerNode
	}
}

func (this *leafBuilder) build() node {
	if this.count == 0 {
		return emptyNode{}
	} else if this.parent == nil {
		return createLeafNode(this.buffer, this.count)
	} else {
		return this.parent.build(createLeafNode(this.buffer, this.count))
	}
}

func (this *leafBuilder) computeSize() int {
	answer := this.count
	if this.parent != nil {
		answer += this.parent.computeSize()
	}
	return answer
}

func createBranchBuilder() *branchBuilder {
	var answer branchBuilder
	answer.buffer = make([]node, maxPerNode)
	return &answer
}

func (this *branchBuilder) addNode(node node) {
	this.buffer[this.count] = node
	this.count++
	if this.count == maxPerNode {
		if this.parent == nil {
			this.parent = createBranchBuilder()
		}
		this.parent.addNode(createBranchNode(this.buffer, minPerNode))
		for i := minPerNode; i < maxPerNode; i++ {
			this.buffer[i-minPerNode] = this.buffer[i]
		}
		this.count -= minPerNode
	}
}

func (this *branchBuilder) build(extra node) node {
	this.buffer[this.count] = extra
	var answer node
	answer = createBranchNode(this.buffer, this.count+1)
	if this.parent != nil {
		answer = this.parent.build(answer)
	}
	return answer
}

func (this *branchBuilder) buildForMerge() node {
	var answer node
	answer = createBranchNode(this.buffer, this.count)
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

func mergeLists(depth int, root1 node, root2 node) node {
	builder := createBranchBuilder()
	proc := func(n node) {
		builder.addNode(n)
	}
	root1.visitNodesOfDepth(depth, proc)
	root2.visitNodesOfDepth(depth, proc)
	return builder.buildForMerge()
}
