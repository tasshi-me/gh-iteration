package flags

import (
	"github.com/spf13/cobra"
)

func Flag(flag string) Node {
	return Node{Flag: flag, Or: nil, And: nil}
}

func Or(nodes ...Node) Node {
	return Node{Flag: "", Or: nodes, And: nil}
}

func And(nodes ...Node) Node {
	return Node{Flag: "", Or: nil, And: nodes}
}

type Validator struct {
	node *Node
}

func NewValidator(node Node) Validator {
	return Validator{node: &node}
}

func (query Validator) Validate(cmd *cobra.Command) error {
	_, err := query.node.validate(cmd, true, 0)
	if err != nil {
		return err
	}
	return nil
}
