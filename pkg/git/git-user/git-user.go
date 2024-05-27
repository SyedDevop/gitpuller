package gituser

type (
	ReposLink struct {
		NextLink, LastLink, PrevLink *string
		FirstLink                    string
		CurrentPage                  int
		PageCount                    int
	}
	Link struct {
		Url string
		Rel string
	}
	// ReposLinkIterator interface {
	// 	Next() ReposLink
	// }
	GitUser struct {
		ReposLink *ReposLink
		Name      string
	}
)

func NetwReposLink(firstUrl string) *ReposLink {
	repoLinlk := &ReposLink{
		FirstLink:   firstUrl,
		CurrentPage: 1,
		PageCount:   0,
	}
	return repoLinlk
}

func (r *ReposLink) Next() *ReposLink {
	panic("Implement me the next iterator for ReposLink")
}
