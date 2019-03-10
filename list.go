package immutableList

const (
	minPerNode = 9
	maxPerNode = 2 * minPerNode
)

type Object interface{}

type Processor func(Object)

type List interface {
	Size() int
	Get(index int) Object
	Append(value Object) List
	ForEach(proc Processor)
	Select(predicate func(Object)bool) List
}

type node interface {
	size() int
	get(index int) Object
	append(value Object) (node, node)
	forEach(proc Processor)
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

func (this listImpl) Select(predicate func(Object)bool) List {
	list := Create()
	answer := &list
	this.root.forEach(func(obj Object) {
		if predicate(obj) {
			*answer = (*answer).Append(obj)
		}
	})
	return *answer
}