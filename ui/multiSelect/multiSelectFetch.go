package multiSelect

import (
	"github.com/SyedDevop/gitpuller/util"
	tea "github.com/charmbracelet/bubbletea"
)

func (f *Fetch) fetchContent() tea.Msg {
	contents, err := f.Clint.GetCountents()
	if err != nil {
		return errMess{err}
	}
	f.Repo = util.GetRepoFromContent(*contents)
	// return errMess{errors.New("unknown error occurred will download this url: " + f.Clint.GitRepoUrl)}
	return multiSelectMsg("done")
}
