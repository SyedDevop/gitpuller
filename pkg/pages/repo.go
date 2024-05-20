package pages

type RepoPage struct{}

func (r *RepoPage) Title() string {
	return "Repo"
}

func (r *RepoPage) Render() string {
	return "Repo"
}
