package handler

import (
	"fmt"
	"git_server/internal/file"
	reserr "git_server/internal/server/res_err"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type (
	Repo         struct{}
	RepoDataType = map[string]interface{}
	Tree         struct {
		SHA       string        `json:"sha"`
		URL       string        `json:"url"`
		Tree      []TreeElement `json:"tree"`
		Truncated bool          `json:"truncated"`
	}
	// TreeElement.go
	TreeElement struct {
		Size *int64  `json:"size,omitempty"`
		URL  *string `json:"url"`
		Path string  `json:"path"`
		Type string  `json:"type"`
		SHA  string  `json:"sha"`
		Mode string  `json:"mode"`
	}
)

const SHA_LENGTH = 40

func (re *Repo) RepoTree(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	sha := chi.URLParam(r, "sha")
	recursive := r.URL.Query().Get("recursive")

	repoPath := filepath.Join(owner, "repo", repo)
	repoJsonPath := fmt.Sprintf("%s.json", repoPath)
	fullPath, _ := file.GetCurDir()
	fullFilePath := filepath.Join(fullPath, repoJsonPath)

	if !file.FileExist(fullFilePath) ||
		(sha != "main" && len(sha) != SHA_LENGTH) {
		reserr.ErrGitTree.Render(w, r)
		render.JSON(w, r, reserr.ErrGitTree)
		return
	}

	var data Tree
	err := file.ReadJson(repoPath, &data)
	if err != nil {
		log.Error("RepoTree#ReadJson", "err", err)
		render.Render(w, r, reserr.ErrRender(err))
		return
	}

	filteredData := make(map[string]interface{}, 0)

	if sha == "main" {
		if recursive != "" {
			render.JSON(w, r, data)
			return
		}
		treeList := make([]interface{}, 0)
		filteredData["url"] = data.URL
		filteredData["sha"] = data.SHA
		for _, val := range data.Tree {
			isRoot, _ := file.GetParentPath(val.Path)
			if isRoot {
				treeList = append(treeList, val)
			}
		}
		filteredData["tree"] = treeList
		filteredData["truncated"] = data.Truncated
		render.JSON(w, r, filteredData)
		return
	}

	fileName := ""
	for _, val := range data.Tree {
		if val.SHA == sha {
			fileName = val.Path
			filteredData["url"] = val.URL
			filteredData["sha"] = val.SHA
			filteredData["truncated"] = false
		}
	}

	treeList := make([]interface{}, 0)
	for _, val := range data.Tree {
		if strings.Contains(val.Path, fileName) && val.Path != fileName {

			newVal := val
			newVal.Path = strings.Replace(val.Path, fileName+"/", "", 1)
			if recursive != "" {
				treeList = append(treeList, newVal)
			} else {
				if file.GetFileDepth(newVal.Path) == 0 {
					treeList = append(treeList, newVal)
				}
			}
		}
	}

	filteredData["tree"] = treeList
	render.JSON(w, r, filteredData)
}
