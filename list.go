package immutableList

const (
	minPerNode = 9
	maxPerNode = 2 * minPerNode
)

type Object interface{}

type Processor func(Object)
type Visitor func(int, Object)

type List interface {
	Size() int
	Get(index int) Object
	Append(value Object) List
	Insert(indexBefore int, value Object) List
	ForEach(proc Processor)
	Visit(offset int, limit int, v Visitor)
	Select(predicate func(Object) bool) List
	Slice(offset, limit int) []Object
}

type node interface {
	size() int
	get(index int) Object
	append(value Object) (node, node)
	insert(indexBefore int, value Object) (node, node)
	forEach(proc Processor)
	visit(start int, limit int, v Visitor)
}

type listImpl struct {
	root node
}

func Create() List {
	return listImpl{root: emptyNode{}}
}

func (this listImpl) Size() int {
	return this.root.size()
}

func (this listImpl) Get(index int) Object {
	return this.root.get(index)
}

func (this listImpl) Append(value Object) List {
	replacement, extra := this.root.append(value)
	return listInsertImpl(replacement, extra)
}

func (this listImpl) Insert(indexBefore int, value Object) List {
	currentSize := this.root.size()
	if indexBefore >= currentSize {
		return this.Append(value)
	} else {
		replacement, extra := this.root.insert(maxInt(0, indexBefore), value)
		return listInsertImpl(replacement, extra)
	}
}

func listInsertImpl(replacement node, extra node) List {
	if extra == nil {
		return listImpl{root: replacement}
	} else {
		children := []node{replacement, extra}
		totalSize := replacement.size() + extra.size()
		return listImpl{root: branchNode{children: children, totalSize: totalSize}}
	}
}

func (this listImpl) ForEach(proc Processor) {
	this.root.forEach(proc)
}

func (this listImpl) Visit(offset int, limit int, visitor Visitor) {
	offset = maxInt(0, offset)
	limit = minInt(limit, this.root.size())
	this.root.visit(offset, limit, visitor)
}

func (this listImpl) Select(predicate func(Object) bool) List {
	answer := Create()
	this.root.forEach(func(obj Object) {
		if predicate(obj) {
			answer = answer.Append(obj)
		}
	})
	return answer
}

func (this listImpl) Slice(offset, limit int) []Object {
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
