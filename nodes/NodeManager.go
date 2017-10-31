package nodes

import (
	nomad "github.com/hashicorp/nomad/api"
)

type NodeManager struct {
	nomadClient *nomad.Client
}

func NewManager(nomadClient *nomad.Client) *NodeManager {
	return &NodeManager{
		nomadClient: nomadClient,
	}
}

func (m *NodeManager) List() ([]*NodeSummary, error) {
	nodes, _, err := m.nomadClient.Nodes().List(&nomad.QueryOptions{})

	if err != nil {
		return nil, err
	}

	res := make([]*NodeSummary, len(nodes))
	for i, v := range nodes {
		res[i] = &NodeSummary{
			Status: v.Status,
			Id: v.ID,
			Name: v.Name,
		}
	}



	return res, nil
}

func (m *NodeManager) Get(nodeId string) (*Node, error) {
	node, _, err := m.nomadClient.Nodes().Info(nodeId, &nomad.QueryOptions{})

	if err != nil {
		return nil, err
	}

	return &Node{
		Id: node.ID,
		Name: node.Name,
		Attributes: node.Attributes,
		Resources: &NodeResources{
			CPU: node.Resources.CPU,
			Disk: node.Resources.DiskMB,
			IOPS: node.Resources.IOPS,
			Memory: node.Resources.MemoryMB,
		},
	}, nil
}
