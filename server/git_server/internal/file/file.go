package file

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
)

//go:embed repos.json
var reposJson []byte

func GetReposByte() []byte {
	return reposJson
}

type JsonDataType = []map[string]interface{}

func ReadJson(name string) (JsonDataType, error) {
	fileJson := fmt.Sprintf("%s.json", name)
	file, err := os.ReadFile(fileJson)
	if err != nil {
		fmt.Printf("File#ReadJson (ReadFile): %v\n", err)
		return nil, err
	}

	var jsonMap JsonDataType
	err = json.Unmarshal(file, &jsonMap)
	if err != nil {
		return nil, err
	}
	return jsonMap, nil
}

func GetReposJson() (JsonDataType, error) {
	var jsonMap JsonDataType
	err := json.Unmarshal(reposJson, &jsonMap)
	if err != nil {
		return nil, err
	}
	return jsonMap, nil
}
