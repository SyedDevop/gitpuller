package repo

import (
	"sync"

	"github.com/SyedDevop/gitpuller/pkg/git"
)

type Node struct {
	SelectedRepo map[int]struct{}
	Repo         []git.TreeElement
}

// TODO: Rename the FolderRepo to SelectedFolders and SelectedRepo to SelectedFiles
type StateTree struct {
	Tree         map[string]*Node
	CurPath      string
	RootPath     string
	SelectedRepo map[string][]git.TreeElement
	FolderRepo   []git.TreeElement
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

func (s *StateTree) SelectedRepoLen() int {
	tempTen := 0
	for _, v := range s.SelectedRepo {
		tempTen += len(v)
	}
	return tempTen
}

// AppendSelected compiles all selected repositories from the ContentTree's Tree map into the SelectedRepo slice.
// It filters and returns a slice of repositories with type "dir".
//
// Returns:
// - []api.Repo: Slice of "dir" type selected repositories.
func (s *StateTree) AppendSelected() {
	for key, repos := range s.Tree {
		for selectRepo := range repos.SelectedRepo {
			if repos.Repo[selectRepo].Type == "tree" {
				s.FolderRepo = append(s.FolderRepo, repos.Repo[selectRepo])
			} else {
				s.SelectedRepo[key] = append(s.SelectedRepo[key], repos.Repo[selectRepo])
			}
		}
	}
}

// UpdateTreesSelected updates the selection status of a repository at the current path (CurPath) identified by index.
func (s *StateTree) UpdateTreesSelected(index int) {
	path := s.CurPath
	if treeData, ok := s.Tree[path]; ok {
		treeData.UpdateSelectedRepo(index)
	}
}

// SelectAllCurTreeRepo selects all repositories at the current path within the tree.
func (s *StateTree) SelectAllCurTreeRepo() {
	path := s.CurPath
	if treeData, ok := s.Tree[path]; ok {
		treeData.SelecteAllRepo()
	}
}

// RemoveAllCurTreeRepo deselects (removes) all selected repositories at the current path within the tree.
func (s *StateTree) RemoveAllCurTreeRepo() {
	path := s.CurPath
	if treeData, ok := s.Tree[path]; ok {
		treeData.RemoveAllRepo()
	}
}
