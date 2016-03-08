package scheduler

type Node struct {
	NodeId    string `json:"nodeId"` // uniq
	Used      bool   `json:"used"`
	Reachable bool   `json:"reachable"`

	Meta *NodeMeta `json:"meta"`
}

type NodeMeta struct {
	Name       string            `json:"name"`
	IP         string            `json:"ip"`
	CPU        int               `json:"cpu"`
	MemoryMB   int               `json:"memoryMB"`
	DiskMB     int               `json:"diskMB"`
	Attributes map[string]string `json:"attributes"` //other
}
