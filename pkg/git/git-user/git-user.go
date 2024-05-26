package git

type (
	ReposLink struct {
		NextLink, LastLink, PrevLink *string
		FirstLink                    string
		CurrentPage                  int
		PageCount                    int
	}
	ReposLinkIterator interface {
		Next() ReposLink
	}
	GitUser struct {
		ReposLink *ReposLink
		Name      string
	}
)

func NetwReposLink(firstUrl string) *ReposLinkIterator {
	repoLinlk := &ReposLink{
		FirstLink:   firstUrl,
		CurrentPage: 1,
		PageCount:   0,
	}
	return repoLinlk
}

func (r *ReposLink) Next() *ReposLink {
	return r
}
