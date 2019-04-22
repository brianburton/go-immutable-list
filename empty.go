package immutableList

import "fmt"

type emptyNode struct {
}

var sharedEmptyNodeInstance node = &emptyNode{}

func (_ *emptyNode) height() int {
	return 1
}

func (_ *emptyNode) size() int {
	return 0
}

func (_ *emptyNode) get(index int) Object {
	panic(fmt.Sprintf("get called on emptyNode: size=0 index=%d", index))
}

func (_ *emptyNode) getFirst() Object {
	panic("getFirst called on emptyNode")
}

func (_ *emptyNode) pop() (Object, node) {
	panic("pop called on emptyNode")
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

func (_ *emptyNode) head(index int) node {
	if index != 0 {
		panic(fmt.Sprintf("index out of bounds: size=0 index=%d", index))
	}
	return sharedEmptyNodeInstance
}

func (_ *emptyNode) tail(index int) node {
	if index != 0 {
		panic(fmt.Sprintf("index out of bounds: size=0 index=%d", index))
	}
	return sharedEmptyNodeInstance
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

func (this *emptyNode) checkInvariants(report reporter, isRoot bool) {
	if this != sharedEmptyNodeInstance {
		report("emptyNode is not sharedEmptyNodeInstance")
	}
}
