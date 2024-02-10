package multiSelect

import (
	types "github.com/SyedDevop/gitpuller/mytypes"
)

type Node struct {
	SelectedRepo map[int]struct{}
	Repo         []types.Repo
}

type ContentTree struct {
	Tree         map[string]*Node
	CurPath      string
	RootPath     string
	SelectedRepo []types.Repo
}

// UpdateSelectedRepo toggles the selection status of a repository identified by key.
// If the repository is already selected, it is deselected (removed from SelectedRepo);
// if it is not selected, it is added
func (n *Node) UpdateSelectedRepo(key int) {
	if _, ok := n.SelectedRepo[key]; ok {
		delete(n.SelectedRepo, key)
	} else {
		n.SelectedRepo[key] = struct{}{}
	}
}

// SelecteAllRepo selects all repositories within this node.
// It does so by adding all indices to SelectedRepo if the number of repositories is greater than the number of selected repositories.
func (n *Node) SelecteAllRepo() {
	if len(n.Repo) > len(n.SelectedRepo) {
		for i := 0; i < len(n.Repo); i++ {
			n.SelectedRepo[i] = struct{}{}
		}
	}
}

// RemoveAllRepo deselects (removes) all selected repositories within this node.
func (n *Node) RemoveAllRepo() {
	if len(n.SelectedRepo) > 0 {
		for i := 0; i <= len(n.Repo); i++ {
			delete(n.SelectedRepo, i)
		}
	}
}

// AppendSelected compiles all selected repositories from the ContentTree's Tree map into the SelectedRepo slice.
// It filters and returns a slice of repositories with type "dir".
//
// Returns:
// - []types.Repo: Slice of "dir" type selected repositories.
func (c *ContentTree) AppendSelected() []types.Repo {
	dirRepo := make([]types.Repo, 0)
	for _, repos := range c.Tree {
		for selectRepo := range repos.SelectedRepo {
			if repos.Repo[selectRepo].Type == "dir" {
				dirRepo = append(dirRepo, repos.Repo[selectRepo])
			}
			c.SelectedRepo = append(c.SelectedRepo, repos.Repo[selectRepo])
		}
	}
	return dirRepo
}

// UpdateTreesSelected updates the selection status of a repository at the current path (CurPath) identified by index.
func (c *ContentTree) UpdateTreesSelected(index int) {
	path := c.CurPath
	if treeData, ok := c.Tree[path]; ok {
		treeData.UpdateSelectedRepo(index)
	}
}

// SelectAllCurTreeRepo selects all repositories at the current path within the tree.
func (c *ContentTree) SelectAllCurTreeRepo() {
	path := c.CurPath
	if treeData, ok := c.Tree[path]; ok {
		treeData.SelecteAllRepo()
	}
}

// RemoveAllCurTreeRepo deselects (removes) all selected repositories at the current path within the tree.
func (c *ContentTree) RemoveAllCurTreeRepo() {
	path := c.CurPath
	if treeData, ok := c.Tree[path]; ok {
		treeData.RemoveAllRepo()
	}
}
