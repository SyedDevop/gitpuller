package file

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed repos.json
var reposJson []byte

func GetReposByte() []byte {
	return reposJson
}

type JsonDataType = []map[string]interface{}

func GetCurDir() (string, bool) {
	_, filename, _, ok := runtime.Caller(0)

	return filepath.Dir(filename), ok
}

func ReadJson(name string) (JsonDataType, error) {
	fileJson := fmt.Sprintf("%s.json", name)

	curDir, ok := GetCurDir()
	if !ok {
		fmt.Println("File#ReadJson (curDir): unable to determine current file path")
		return nil, fmt.Errorf("unable to determine current file path")
	}

	file, err := os.ReadFile(filepath.Join(curDir, fileJson))
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
