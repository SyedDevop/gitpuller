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
	// logF, err := tea.LogToFile("debug.log", "debug")
	// if err != nil {
	// 	fmt.Println("fatal:", err)
	// 	os.Exit(1)
	// }
	// defer logF.Close()

	contents, err := f.Clint.GetCountents(nil)
	// fmt.Fprintf(logF, "FetchContent#Tea: Url = %s :: Data = %+v :: DataLen = %d \n", f.Clint.GitRepoUrl, contents, len(contents))
	if err != nil {
		return errMess{err}
	}
	f.Repo = contents
	// return errMess{errors.New("unknown error occurred will download this url: " + f.Clint.GitRepoUrl)}
	return multiSelectMsg("done")
}
