package multiSelect

import (
	"sync"

	"github.com/SyedDevop/gitpuller/cmd/api"
)

type Node struct {
	SelectedRepo map[int]struct{}
	Repo         []api.TreeElement
}

// TODO: Rename the FolderRepo to SelectedFolders and SelectedRepo to SelectedFiles
type ContentTree struct {
	Tree         map[string]*Node
	CurPath      string
	RootPath     string
	SelectedRepo map[string][]api.TreeElement
	FolderRepo   []api.TreeElement
	Mu           sync.Mutex
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
// - []api.Repo: Slice of "dir" type selected repositories.
func (c *ContentTree) AppendSelected(name string) {
	for _, repos := range c.Tree {
		for selectRepo := range repos.SelectedRepo {
			if repos.Repo[selectRepo].Type == "tree" {
				c.FolderRepo = append(c.FolderRepo, repos.Repo[selectRepo])
			} else {
				c.SelectedRepo[name] = append(c.SelectedRepo[name], repos.Repo[selectRepo])
			}
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
