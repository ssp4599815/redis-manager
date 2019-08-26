package models

type Node struct {

	NodeId    string
	Addr      string
	Flags     string
	MasterId  string
	Connected string
}

// 插入
func (node *Node) Create()  {
}
