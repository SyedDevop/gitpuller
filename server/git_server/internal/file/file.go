package file

import (
	_ "embed"
	"encoding/json"
)

//go:embed repos.json
var reposJson []byte

func GetReposByte() []byte {
	return reposJson
}

func GetReposJson() ([]map[string]interface{}, error) {
	var jsonMap []map[string]interface{}
	err := json.Unmarshal(reposJson, &jsonMap)
	if err != nil {
		return nil, err
	}
	return jsonMap, nil
}
