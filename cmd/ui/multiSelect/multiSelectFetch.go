package multiSelect

import (
	"github.com/SyedDevop/gitpuller/cmd/api"
	"github.com/SyedDevop/gitpuller/cmd/util"

	tea "github.com/charmbracelet/bubbletea"
)

type Fetch struct {
	Err       error
	Clint     *api.Clint
	FetchMess string
	Repo      []api.Repo
	FethDone  bool
}

func (f *Fetch) fetchContent() tea.Msg {
	contents, err := f.Clint.GetCountents(nil)
	if err != nil {
		return errMess{err}
	}
	f.Repo = util.GetRepoFromContent(*contents)
	// return errMess{errors.New("unknown error occurred will download this url: " + f.Clint.GitRepoUrl)}
	return multiSelectMsg("done")
}
