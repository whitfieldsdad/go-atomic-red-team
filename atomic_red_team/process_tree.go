package atomic_red_team

import (
	"time"

	"github.com/mitchellh/go-ps"
	"github.com/pkg/errors"
)

type LightweightProcessTree struct {
	CreateTime time.Time
	UpdateTime time.Time
	pidToPpid  map[int]int // Map of PID -> PPID
}

// GetDepth returns to the depth of the process tree.
func (tree *LightweightProcessTree) GetDepth() int {
	maxDepth := 0
	for pid := range tree.pidToPpid {
		depth := tree.GetProcessDepth(pid)
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	return maxDepth
}

// GetProcessDepth returns the depth of the process in the process tree.
func (tree *LightweightProcessTree) GetProcessDepth(pid int) int {
	ancestors, _ := tree.GetAncestors(pid)
	depth := len(ancestors)
	return depth
}

// GetFamilyMembers returns the PID of the process and all of its ancestors, descendants, siblings, and siblings' descendants.
func (tree *LightweightProcessTree) GetFamilyMembers(pid int) ([]int, error) {
	if !tree.ContainsPID(pid) {
		return nil, errors.Errorf("PID %d not found", pid)
	}
	ancestors, _ := tree.GetAncestors(pid)
	descendants, _ := tree.GetDescendants(pid)
	siblings, _ := tree.GetSiblings(pid)
	siblingDescendants, _ := tree.GetSiblingDescendants(pid)

	var family []int
	family = append(family, ancestors...)
	family = append(family, pid)
	family = append(family, descendants...)
	family = append(family, siblings...)
	family = append(family, siblingDescendants...)
	return family, nil
}

// HasPID returns true if the process tree contains the specified PID.
func (tree *LightweightProcessTree) ContainsPID(pid int) bool {
	ok := tree.ContainsPPID(pid)
	if ok {
		return true
	}
	for _, ppid := range tree.pidToPpid {
		if ppid == pid {
			return true
		}
	}
	return false
}

// HasPPID returns true if the process tree contains the specified PPID.
func (tree *LightweightProcessTree) ContainsPPID(pid int) bool {
	_, ok := tree.pidToPpid[pid]
	return ok
}

// GetFamily returns the PID of the process and all of its ancestors and descendants.
func (tree *LightweightProcessTree) GetAncestors(pid int) ([]int, error) {
	if !tree.ContainsPID(pid) {
		return nil, errors.Errorf("PID %d not found", pid)
	}
	var ancestors []int
	cursor := pid
	for {
		ppid, ok := tree.pidToPpid[cursor]
		if !ok {
			break
		}
		ancestors = append(ancestors, ppid)
		cursor = ppid
	}
	return ancestors, nil
}

func (tree *LightweightProcessTree) GetPPID(pid int) (int, error) {
	ppid, ok := tree.pidToPpid[pid]
	if !ok {
		return -1, errors.Errorf("PID %d not found", pid)
	}
	return ppid, nil
}

func (tree *LightweightProcessTree) GetChildren(pid int) ([]int, error) {
	ppid := pid
	found := false
	var children []int
	for childPid, parentPid := range tree.pidToPpid {
		if parentPid == ppid {
			children = append(children, childPid)
			found = true
		}
	}
	if !found {
		return nil, errors.Errorf("PID %d not found", pid)
	}
	return children, nil
}

func (tree *LightweightProcessTree) GetSiblings(pid int) ([]int, error) {
	var siblings []int
	ppid, ok := tree.pidToPpid[pid]
	if !ok {
		return nil, errors.Errorf("PID %d not found", pid)
	}
	for child, parent := range tree.pidToPpid {
		if parent == ppid && child != pid {
			siblings = append(siblings, child)
		}
	}
	return siblings, nil
}

func (tree *LightweightProcessTree) GetSiblingDescendants(pid int) ([]int, error) {
	siblings, err := tree.GetSiblings(pid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get siblings")
	}
	var descendants []int
	for _, sibling := range siblings {
		descendant, err := tree.GetDescendants(sibling)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get descendants of sibling %d", sibling)
		}
		descendants = append(descendants, descendant...)
	}
	return descendants, nil
}

func (tree *LightweightProcessTree) GetDescendants(pid int) ([]int, error) {
	var descendants []int
	for {
		children, err := tree.GetChildren(pid)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get children of %d", pid)
		}
		if len(children) == 0 {
			break
		}
		descendants = append(descendants, children...)
		for _, child := range children {
			childDescendants, err := tree.GetDescendants(child)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get descendants of %d", child)
			}
			descendants = append(descendants, childDescendants...)
		}
	}
	return descendants, nil
}

func NewLightweightProcessTree() *LightweightProcessTree {
	now := time.Now()
	return &LightweightProcessTree{
		CreateTime: now,
		UpdateTime: now,
		pidToPpid:  make(map[int]int),
	}
}

// GetLightweightProcessTree returns a directed graph of processes that only includes the PID and PPID of each process.
func GetLightweightProcessTree() (*LightweightProcessTree, error) {
	processes, err := ps.Processes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list processes")
	}
	tree := NewLightweightProcessTree()
	for _, process := range processes {
		pid := process.Pid()
		ppid := process.PPid()
		tree.pidToPpid[pid] = ppid
	}
	return tree, nil
}
