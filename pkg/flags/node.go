package flags

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mshrtsr/gh-iteration/pkg/log"
	"github.com/spf13/cobra"
)

type Node struct {
	Flag string
	Or   []Node
	And  []Node
}

func (node Node) printRelation() string {
	switch {
	case len(node.Flag) > 0:
		return "--" + node.Flag
	case node.Or != nil:
		arr := make([]string, 0, len(node.Or))
		for _, n := range node.Or {
			arr = append(arr, n.printRelation())
		}
		str := "(" + strings.Join(arr, " | ") + ")"
		return str
	case node.And != nil:
		arr := make([]string, 0, len(node.Or))
		for _, n := range node.And {
			arr = append(arr, n.printRelation())
		}
		str := "(" + strings.Join(arr, " & ") + ")"
		return str
	}
	return "invalid state"
}

func (node Node) flagList() []string {
	switch {
	case len(node.Flag) > 0:
		return []string{"--" + node.Flag}
	case node.Or != nil:
		arr := make([]string, 0, len(node.Or))
		for _, n := range node.Or {
			arr = append(arr, n.flagList()...)
		}
		return arr
	case node.And != nil:
		arr := make([]string, 0, len(node.Or))
		for _, n := range node.And {
			arr = append(arr, n.flagList()...)
		}
		return arr
	}
	return nil
}

func (node Node) validate(cmd *cobra.Command, required bool, depth int) (bool, error) {
	switch {
	case len(node.Flag) > 0:
		return node.validateAsFlag(cmd, depth)
	case node.Or != nil:
		return node.validateAsOr(cmd, required, depth)
	case node.And != nil:
		return node.validateAsAnd(cmd, required, depth)
	}
	return false, errors.New("invalid state") //nolint:goerr113
}

func (node Node) validateAsFlag(cmd *cobra.Command, depth int) (bool, error) {
	filler := strings.Repeat("  ", depth)
	log.Trace(filler + "Validate Flag Node: " + node.printRelation())
	log.Trace(filler + fmt.Sprintf("Flag Node %s changed: %v", node.Flag, cmd.Flag(node.Flag).Changed))
	return cmd.Flag(node.Flag).Changed, nil
}

func (node Node) validateAsOr(cmd *cobra.Command, required bool, depth int) (bool, error) {
	filler := strings.Repeat("  ", depth)

	log.Trace(filler + "Validate OR Node: " + node.printRelation())
	changed := false
	var changedNode Node
	var invalidNodes []Node
	var err error
	for _, childNode := range node.Or {
		nodeChanged, e := childNode.validate(cmd, false, depth+1)
		err = e
		if nodeChanged {
			if changed {
				invalidNodes = append(invalidNodes, childNode)
			} else {
				changedNode = childNode
			}
		}
		changed = changed || nodeChanged
	}

	log.Trace(filler + "changedNode" + fmt.Sprint(changedNode))
	log.Trace(filler + "invalidNodes" + fmt.Sprint(invalidNodes))

	// Error on current node
	if changed {
		if len(invalidNodes) > 0 {
			var arr []string
			for _, invalidNode := range invalidNodes {
				arr = append(arr, invalidNode.flagList()...)
			}
			return true, fmt.Errorf("when you set %s, you cannot set %s", changedNode.flagList(), arr) //nolint:goerr113
		}
	} else {
		if required {
			return false, fmt.Errorf("you must set one of %s", node.flagList()) //nolint:goerr113
		}
	}

	// Error on children node
	if err != nil {
		return changed, err
	}
	return changed, nil
}

//nolint:cyclop
func (node Node) validateAsAnd(cmd *cobra.Command, required bool, depth int) (bool, error) {
	filler := strings.Repeat("  ", depth)

	log.Trace(filler + "Validate AND Node: " + node.printRelation())
	changed := false
	var changedNodes []Node
	var unchangedNodes []Node
	var err error
	for _, childNode := range node.And {
		nodeChanged, e := childNode.validate(cmd, required, depth+1)
		err = e
		if nodeChanged {
			changedNodes = append(changedNodes, childNode)
		} else {
			unchangedNodes = append(unchangedNodes, childNode)
		}
		changed = changed || nodeChanged
	}
	if required {
		if len(unchangedNodes) > 0 {
			return changed, fmt.Errorf("you must set all of %s", node.printRelation()) //nolint:goerr113
		}
	}

	if changed && len(unchangedNodes) > 0 {
		var changedFlags []string
		for _, changedNode := range changedNodes {
			changedFlags = append(changedFlags, changedNode.flagList()...)
		}
		var unchangedFlags []string
		for _, unchangedNode := range unchangedNodes {
			unchangedFlags = append(unchangedFlags, unchangedNode.flagList()...)
		}

		return changed, fmt.Errorf("when you set %s, you must set %s", changedFlags, unchangedFlags) //nolint:goerr113
	}

	if required && !changed {
		return changed, fmt.Errorf("you must set all of %s", node.printRelation()) //nolint:goerr113
	}

	if err != nil {
		return changed, err
	}
	// All the children were changed/not changed
	return changed, nil
}
