package node

type DefaultNode struct {
	StatisticNode
}

func NewDefaultNode() *DefaultNode {
	return &DefaultNode{StatisticNode: *NewStatisticNode()}
}
