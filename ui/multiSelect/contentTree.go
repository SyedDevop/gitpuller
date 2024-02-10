package multiSelect

import types "github.com/SyedDevop/gitpuller/mytypes"

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
func (t *Node) UpdateSelectedRepo(key int) {
	if _, ok := t.SelectedRepo[key]; ok {
		delete(t.SelectedRepo, key)
	} else {
		t.SelectedRepo[key] = struct{}{}
	}
}

// SelecteAllRepo selects all repositories within this node.
// It does so by adding all indices to SelectedRepo if the number of repositories is greater than the number of selected repositories.
func (t *Node) SelecteAllRepo() {
	if len(t.Repo) > len(t.SelectedRepo) {
		for i := 0; i < len(t.Repo); i++ {
			t.SelectedRepo[i] = struct{}{}
		}
	}
}

// RemoveAllRepo deselects (removes) all selected repositories within this node.
func (t *Node) RemoveAllRepo() {
	if len(t.SelectedRepo) > 0 {
		for i := 0; i <= len(t.Repo); i++ {
			delete(t.SelectedRepo, i)
		}
	}
}

// AppendSelected aggregates all selected repositories across all nodes in the tree into the SelectedRepo slice of the ContentTree structure.
func (c *ContentTree) AppendSelected() {
	for _, repos := range c.Tree {
		for selectRepo := range repos.SelectedRepo {
			c.SelectedRepo = append(c.SelectedRepo, repos.Repo[selectRepo])
		}
	}
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
