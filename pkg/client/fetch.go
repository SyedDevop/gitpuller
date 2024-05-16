package client

import (
	"github.com/SyedDevop/gitpuller/pkg/git"

	tea "github.com/charmbracelet/bubbletea"
)

type Fetch struct {
	Err       error
	Clint     *Clint
	FetchMess string
	Repo      []git.TreeElement
	FethDone  bool
}

type (
	FetchDoneMess string
	FetchErrMess  struct{ error }
)

func (e FetchErrMess) Error() string { return e.error.Error() }

func (f *Fetch) FetchContent() tea.Msg {
	// logF, err := tea.LogToFile("debug.log", "debug")
	// if err != nil {
	// 	fmt.Println("fatal:", err)
	// 	os.Exit(1)
	// }
	// defer logF.Close()

	contents, err := f.Clint.GetCountents(nil)
	// fmt.Fprintf(logF, "FetchContent#Tea: Url = %s :: Data = %+v :: DataLen = %d \n", f.Clint.GitRepoUrl, contents, len(contents))
	if err != nil {
		return FetchErrMess{err}
	}
	f.Repo = contents
	// return errMess{errors.New("unknown error occurred will download this url: " + f.Clint.GitRepoUrl)}
	return FetchDoneMess("done")
}
