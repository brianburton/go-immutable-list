package immutableList

const (
	minPerNode = 8
	maxPerNode = 2 * minPerNode
)

type Object interface{}

type Processor func(Object)
type Visitor func(int, Object)
type nodeProcessor func(node)
type iteratorState struct {
	next         *iteratorState
	currentNode  node
	currentIndex int
}

type Iterator interface {
	Next() bool
	Get() Object
}

type List interface {
	Size() int
	Get(index int) Object
	Append(value Object) List
	AppendList(other List) List
	Insert(indexBefore int, value Object) List
	ForEach(proc Processor)
	Visit(offset int, limit int, v Visitor)
	Select(predicate func(Object) bool) List
	Slice(offset, limit int) []Object
	Delete(index int) List
	Set(index int, value Object) List
	FwdIterate() Iterator
	height() int
}

type node interface {
	size() int
	get(index int) Object
	append(value Object) (node, node)
	appendNode(other node) (node, node)
	prependNode(other node) (node, node)
	insert(indexBefore int, value Object) (node, node)
	forEach(proc Processor)
	visit(start int, limit int, v Visitor)
	height() int
	visitNodesOfHeight(targetHeight int, proc nodeProcessor)
	isComplete() bool
	mergeWith(other node) node
	delete(index int) node
	set(index int, value Object) node
	next(state *iteratorState) (*iteratorState, Object)
}

type listImpl struct {
	root node
}

type iteratorImpl struct {
	state *iteratorState
	value Object
}

func (this *listImpl) FwdIterate() Iterator {
	var state *iteratorState
	if this.root.size() == 0 {
		state = nil
	} else {
		state = &iteratorState{currentNode: this.root}
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

func Create() List {
	return &listImpl{root: &emptyNode{}}
}

func (this *listImpl) Size() int {
	return this.root.size()
}

func (this *listImpl) height() int {
	return this.root.height()
}

func (this *listImpl) Get(index int) Object {
	return this.root.get(index)
}

func (this *listImpl) Append(value Object) List {
	replacement, extra := this.root.append(value)
	return listInsertImpl(replacement, extra)
}

func (this *listImpl) AppendList(other List) List {
	otherImpl := other.(*listImpl)
	thisSize := this.root.size()
	otherSize := otherImpl.root.size()
	if thisSize == 0 {
		return other
	} else if otherSize == 0 {
		return this
	} else if thisSize >= otherSize {
		replacement, extra := this.root.appendNode(otherImpl.root)
		return listInsertImpl(replacement, extra)
	} else {
		replacement, extra := otherImpl.root.prependNode(this.root)
		return listInsertImpl(replacement, extra)
	}
}

func (this *listImpl) Insert(indexBefore int, value Object) List {
	currentSize := this.root.size()
	if indexBefore < 0 || indexBefore > currentSize {
		return nil
	} else if indexBefore == currentSize {
		return this.Append(value)
	} else {
		replacement, extra := this.root.insert(maxInt(0, indexBefore), value)
		return listInsertImpl(replacement, extra)
	}
}

func listInsertImpl(replacement node, extra node) List {
	if extra == nil {
		return &listImpl{root: replacement}
	} else {
		children := []node{replacement, extra}
		nodeSize := replacement.size() + extra.size()
		nodeHeight := replacement.height() + 1
		return &listImpl{root: createBranchNode(children, nodeSize, nodeHeight)}
	}
}

func (this *listImpl) Delete(index int) List {
	if index < 0 || index >= this.Size() {
		return nil
	}
	newRoot := this.root.delete(index)
	return &listImpl{newRoot}
}

func (this *listImpl) ForEach(proc Processor) {
	this.root.forEach(proc)
}

func (this *listImpl) Visit(offset int, limit int, visitor Visitor) {
	offset = maxInt(0, offset)
	limit = minInt(limit, this.root.size())
	this.root.visit(offset, limit, visitor)
}

func (this *listImpl) Select(predicate func(Object) bool) List {
	answer := CreateBuilder()
	this.root.forEach(func(obj Object) {
		if predicate(obj) {
			answer.Add(obj)
		}
	})
	return answer.Build()
}

func (this *listImpl) Slice(offset, limit int) []Object {
	if offset < 0 || limit < offset || limit > this.Size() {
		return nil
	}
	if limit == offset {
		return make([]Object, 0)
	}
	answer := make([]Object, limit-offset)
	this.root.visit(offset, limit, func(index int, obj Object) {
		answer[index-offset] = obj
	})
	return answer
}

func (this *listImpl) Set(index int, value Object) List {
	size := this.Size()
	if index < 0 || index > size {
		return nil
	}
	if index == size {
		return this.Append(value)
	} else {
		newRoot := this.root.set(index, value)
		return &listImpl{newRoot}
	}
}
