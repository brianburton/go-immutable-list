package immutableList

import "fmt"

type Object interface{}

type Processor func(Object)
type Visitor func(int, Object)

type reporter func(message string)

type Iterator interface {
	Next() bool
	Get() Object
}

type List interface {
	Size() int
	Get(index int) Object
	GetFirst() Object
	GetLast() Object
	Append(value Object) List
	AppendList(other List) List
	Insert(indexBefore int, value Object) List
	InsertList(indexBefore int, other List) List
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

type listImpl struct {
	root node
}

var sharedEmptyListInstance List = &listImpl{sharedEmptyNode}

func Create() List {
	return sharedEmptyListInstance
}

func createListNode(root node) List {
	if root.size() == 0 {
		return sharedEmptyListInstance
	} else {
		return &listImpl{root: root}
	}
}

func (this *listImpl) FwdIterate() Iterator {
	return createIterator(this.root)
}

func (this *listImpl) Size() int {
	return this.root.size()
}

func (this *listImpl) Get(index int) Object {
	return this.root.get(index)
}

func (this *listImpl) GetFirst() Object {
	return this.root.getFirst()
}

func (this *listImpl) GetLast() Object {
	return this.root.getLast()
}

func (this *listImpl) Append(value Object) List {
	return createListNode(this.root.append(value))
}

func (this *listImpl) AppendList(other List) List {
	otherImpl := other.(*listImpl)
	return createListNode(appendNodes(this.root, otherImpl.root))
}

func (this *listImpl) Insert(indexBefore int, value Object) List {
	return createListNode(this.root.insert(indexBefore, value))
}

func (this *listImpl) InsertList(indexBefore int, other List) List {
	currentSize := this.root.size()
	if indexBefore < 0 || indexBefore > currentSize {
		panic(fmt.Sprintf("index out of bounds: size=%d index=%d", currentSize, indexBefore))
	}
	if indexBefore == 0 {
		return other.AppendList(this)
	}
	if indexBefore == currentSize {
		return this.AppendList(other)
	}
	return this.Head(indexBefore).AppendList(other).AppendList(this.Tail(indexBefore))
}

func (this *listImpl) Delete(index int) List {
	return createListNode(this.root.delete(index))
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

	var root node
	if offset == 0 {
		root = this.root.tail(limit)
	} else if limit == size {
		root = this.root.head(offset)
	} else {
		root = appendNodes(this.root.head(offset), this.root.tail(limit))
	}
	return createListNode(root)
}

func (this *listImpl) Head(length int) List {
	return createListNode(this.root.head(length))
}

func (this *listImpl) Tail(index int) List {
	return createListNode(this.root.tail(index))
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
	return createListNode(root)
}

func (this *listImpl) ForEach(proc Processor) {
	this.root.forEach(proc)
}

func (this *listImpl) Visit(offset int, limit int, visitor Visitor) {
	if offset < 0 || limit < offset || limit > this.Size() {
		panic(fmt.Sprintf("invalid offset or limit: size=%d offset=%d limit=%d", this.Size(), offset, limit))
	}
	this.root.visit(0, offset, limit, visitor)
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
	this.root.visit(0, offset, limit, func(index int, obj Object) {
		answer[index-offset] = obj
	})
	return answer
}

func (this *listImpl) Set(index int, value Object) List {
	if index == this.root.size() {
		return createListNode(this.root.append(value))
	} else {
		return createListNode(this.root.set(index, value))
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
	return createListNode(this.root.prepend(value))
}

func (this *listImpl) Pop() (Object, List) {
	switch this.Size() {
	case 0:
		panic("Pop called on empty List")
	case 1:
		value := this.root.getFirst()
		return value, sharedEmptyListInstance
	default:
		value, newRoot := this.root.pop()
		return value, createListNode(newRoot)
	}
}
