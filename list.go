package immutableList

const (
	minPerNode = 8
	maxPerNode = 2 * minPerNode
)

type Object interface{}

type Processor func(Object)
type Visitor func(int, Object)
type nodeProcessor func(node)

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
}

type node interface {
	size() int
	get(index int) Object
	append(value Object) (node, node)
	insert(indexBefore int, value Object) (node, node)
	forEach(proc Processor)
	visit(start int, limit int, v Visitor)
	height() int
	maxCompleteHeight() int
	visitNodesOfHeight(targetHeight int, proc nodeProcessor)
	isComplete() bool
	mergeWith(other node) node
	delete(index int) node
	set(index int, value Object) node
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

func (this listImpl) AppendList(other List) List {
	otherImpl := other.(listImpl)

	var answer List
	thisHeight := this.root.maxCompleteHeight()
	otherHeight := otherImpl.root.maxCompleteHeight()
	if otherHeight < 1 {
		answer = this
		other.ForEach(func(object Object) {
			answer = answer.Append(object)
		})
	} else if thisHeight < 1 {
		answer = other
		index := 0
		this.ForEach(func(object Object) {
			answer = answer.Insert(index, object)
			index++
		})
	} else {
		commonHeight := minInt(thisHeight, otherHeight)
		newRoot := mergeLists(commonHeight, this.root, otherImpl.root)
		answer = listImpl{newRoot}
	}
	return answer
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

func (this listImpl) Delete(index int) List {
	if index < 0 || index >= this.Size() {
		return this
	}
	newRoot := this.root.delete(index)
	return listImpl{newRoot}
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

func (this listImpl) Set(index int, value Object) List {
	size := this.Size()
	if index < 0 || index > size {
		return nil
	}
	if index == size {
		return this.Append(value)
	} else {
		newRoot := this.root.set(index, value)
		return listImpl{newRoot}
	}
}
