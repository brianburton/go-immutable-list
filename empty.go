package immutableList

type emptyNode struct {
}

func (_ *emptyNode) visitNodesOfHeight(targetHeight int, proc nodeProcessor) {
}

func (_ *emptyNode) height() int {
	return 1
}

func (_ *emptyNode) size() int {
	return 0
}

func (_ *emptyNode) get(index int) Object {
	return nil
}

func (_ *emptyNode) append(value Object) (node, node) {
	return createLeafNode([]Object{value}), nil
}

func (_ *emptyNode) appendNode(other node) (node, node) {
	return other, nil
}

func (_ *emptyNode) prependNode(other node) (node, node) {
	return other, nil
}

func (_ *emptyNode) insert(_ int, value Object) (node, node) {
	return createLeafNode([]Object{value}), nil
}

func (this *emptyNode) set(index int, value Object) node {
	return this
}

func (_ *emptyNode) forEach(proc Processor) {}

func (_ *emptyNode) visit(start int, limit int, v Visitor) {}

func (_ *emptyNode) isComplete() bool {
	return false
}

func (_ *emptyNode) mergeWith(other node) node {
	return other
}

func (this *emptyNode) delete(index int) node {
	return this
}

func (_ *emptyNode) next(state *iteratorState) (*iteratorState, Object) {
	return nil, nil
}
