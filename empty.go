package immutableList

type emptyNode struct {
}

func (_ emptyNode) size() int {
	return 0
}

func (_ emptyNode) get(index int) Object {
	return nil
}

func (_ emptyNode) append(value Object) (node, node) {
	return leafNode{contents: []Object{value}}, nil
}

func (_ emptyNode) forEach(proc Processor) {}
