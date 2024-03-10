package multiSelect

import (
	"github.com/SyedDevop/gitpuller/cmd/api"

	tea "github.com/charmbracelet/bubbletea"
)

type Fetch struct {
	Err       error
	Clint     *api.Clint
	FetchMess string
	Repo      []api.TreeElement
	FethDone  bool
}

func (f *Fetch) fetchContent() tea.Msg {
	contents, err := f.Clint.GetCountents(nil)
	if err != nil {
		return errMess{err}
	}
	f.Repo = contents
	// return errMess{errors.New("unknown error occurred will download this url: " + f.Clint.GitRepoUrl)}
	return multiSelectMsg("done")
}
