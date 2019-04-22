package immutableList

import "fmt"

const (
	minPerNode = 8
	maxPerNode = 2 * minPerNode
)

type Object interface{}

type Processor func(Object)
type Visitor func(int, Object)
type iteratorState struct {
	next         *iteratorState
	currentNode  node
	currentIndex int
}

type reporter func(message string)

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
	Head(length int) List
	Tail(index int) List
	SubList(offset int, limit int) List
	ForEach(proc Processor)
	Visit(offset int, limit int, v Visitor)
	Select(predicate func(Object) bool) List
	Slice(offset, limit int) []Object
	Delete(index int) List
	DeleteRange(offset int, limit int) List
	Set(index int, value Object) List
	FwdIterate() Iterator
	checkInvariants(r reporter)

	IsEmpty() bool
	Push(value Object) List
	Pop() (Object, List)
}

type node interface {
	size() int
	get(index int) Object
	append(value Object) (node, node)
	appendNode(other node) (node, node)
	prependNode(other node) (node, node)
	insert(indexBefore int, value Object) (node, node)
	head(index int) node
	tail(index int) node
	forEach(proc Processor)
	visit(start int, limit int, v Visitor)
	height() int
	isComplete() bool
	mergeWith(other node) node
	delete(index int) node
	set(index int, value Object) node
	next(state *iteratorState) (*iteratorState, Object)
	checkInvariants(r reporter, isRoot bool)
}

type listImpl struct {
	root node
}

type iteratorImpl struct {
	state *iteratorState
	value Object
}

var sharedEmptyListInstance List = &listImpl{sharedEmptyNodeInstance}

func Create() List {
	return sharedEmptyListInstance
}

func createListForBuilder(root node) List {
	if root.size() == 0 {
		return sharedEmptyListInstance
	} else {
		return &listImpl{root}
	}
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

func (this *listImpl) Size() int {
	return this.root.size()
}

func (this *listImpl) Get(index int) Object {
	return this.root.get(index)
}

func (this *listImpl) Append(value Object) List {
	replacement, extra := this.root.append(value)
	return createListNode(replacement, extra)
}

func (this *listImpl) AppendList(other List) List {
	otherImpl := other.(*listImpl)
	thisSize := this.root.size()
	otherSize := otherImpl.root.size()
	if thisSize == 0 {
		return other
	}
	if otherSize == 0 {
		return this
	}

	var replacement, extra node
	if thisSize >= otherSize {
		replacement, extra = this.root.appendNode(otherImpl.root)
	} else {
		replacement, extra = otherImpl.root.prependNode(this.root)
	}
	return createListNode(replacement, extra)
}

func (this *listImpl) Insert(indexBefore int, value Object) List {
	currentSize := this.root.size()
	if indexBefore < 0 || indexBefore > currentSize {
		panic(fmt.Sprintf("index out of bounds: size=%d index=%d", currentSize, indexBefore))
	}
	if indexBefore == currentSize {
		return this.Append(value)
	}
	replacement, extra := this.root.insert(indexBefore, value)
	return createListNode(replacement, extra)
}

func (this *listImpl) Delete(index int) List {
	size := this.Size()
	if index < 0 || index >= size {
		panic(fmt.Sprintf("index out of bounds: size=%d index=%d", size, index))
	}
	if size == 1 {
		return sharedEmptyListInstance
	} else {
		newRoot := this.root.delete(index)
		return &listImpl{newRoot}
	}
}

func (this *listImpl) DeleteRange(offset int, limit int) List {
	size := this.Size()
	if offset < 0 || limit < offset || limit > size {
		panic(fmt.Sprintf("invalid offset or limit: size=%d offset=%d limit=%d", size, offset, limit))
	}
	if offset == 0 && limit == size {
		return sharedEmptyListInstance
	}
	if offset == limit {
		return this
	}

	var root1, root2 node
	if offset == 0 {
		root1 = this.root.tail(limit)
	} else if limit == size {
		root1 = this.root.head(offset)
	} else {
		prefix := this.root.head(offset)
		suffix := this.root.tail(limit)
		if prefix.size() >= suffix.size() {
			root1, root2 = prefix.appendNode(suffix)
		} else {
			root1, root2 = suffix.prependNode(prefix)
		}
	}
	return createListNode(root1, root2)
}

func (this *listImpl) Head(length int) List {
	size := this.Size()
	if length < 0 || length > size {
		panic(fmt.Sprintf("length out of bounds: size=%d length=%d", size, length))
	}
	if length == 0 {
		return sharedEmptyListInstance
	}

	root := this.root.head(length)
	return &listImpl{root}
}

func (this *listImpl) Tail(index int) List {
	size := this.Size()
	if index < 0 || index > size {
		panic(fmt.Sprintf("index out of bounds: size=%d index=%d", size, index))
	}
	if index == size {
		return sharedEmptyListInstance
	}

	root := this.root.tail(index)
	return &listImpl{root}
}

func (this *listImpl) SubList(offset int, limit int) List {
	size := this.Size()
	if offset < 0 || limit < offset || limit > size {
		panic(fmt.Sprintf("invalid offset or limit: size=%d offset=%d limit=%d", size, offset, limit))
	}
	if offset == 0 && limit == size {
		return this
	}
	if offset == limit {
		return sharedEmptyListInstance
	}

	var root node
	if offset == 0 {
		root = this.root.head(limit)
	} else if limit == size {
		root = this.root.tail(offset)
	} else {
		root = this.root.head(limit).tail(offset)
	}
	return &listImpl{root}
}

func (this *listImpl) ForEach(proc Processor) {
	this.root.forEach(proc)
}

func (this *listImpl) Visit(offset int, limit int, visitor Visitor) {
	if offset < 0 || limit < offset || limit > this.Size() {
		panic(fmt.Sprintf("invalid offset or limit: size=%d offset=%d limit=%d", this.Size(), offset, limit))
	}
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
		panic(fmt.Sprintf("invalid offset or limit: size=%d offset=%d limit=%d", this.Size(), offset, limit))
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
		panic(fmt.Sprintf("index out of bounds: size=%d index=%d", size, index))
	}
	if index == size {
		return this.Append(value)
	} else {
		newRoot := this.root.set(index, value)
		return &listImpl{newRoot}
	}
}

func (this *listImpl) checkInvariants(report reporter) {
	if this.Size() == 0 && this != sharedEmptyListInstance {
		report("empty list is not the sharedEmptyListInstance")
	}
	this.root.checkInvariants(report, true)
}

func (this *listImpl) IsEmpty() bool {
	return this.root.size() == 0
}

func (this *listImpl) Push(value Object) List {
	return this.Insert(0, value)
}

func (this *listImpl) Pop() (Object, List) {
	if this.Size() == 0 {
		panic("Pop called on empty List")
	}
	value := this.Get(0)
	popped := this.Delete(0)
	return value, popped
}

func createListNode(replacement node, extra node) List {
	if extra == nil {
		return &listImpl{root: replacement}
	} else {
		children := []node{replacement, extra}
		nodeSize := replacement.size() + extra.size()
		nodeHeight := replacement.height() + 1
		return &listImpl{root: createBranchNode(children, nodeSize, nodeHeight)}
	}
}
