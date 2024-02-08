package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test Get Parent Path
func TestGetParentPath(t *testing.T) {
	pathsList := []string{
		"/Home/username/projects/helloworld",
		"Home/username/projects/helloworld/",
		"/home",
		"Home/",
	}

	pathSteps1 := []string{
		"/Home/username/projects",
		"/Home/username",
		"/Home",
	}

	// pathSteps2 := []string{
	// 	"/Home/username/projects",
	// 	"/Home/username",
	// 	"Home",
	// }

	workingPath := pathsList[0]
	for i := 1; i <= 3; i++ {
		fmt.Println(i)
		isParent, parentPath := GetParentPath(workingPath)
		if i == 3 {
			assert.Equal(t, parentPath, pathSteps1[2], "The path should be a parent from :last")
			assert.Equal(t, isParent, true, "Is parent should be true from :last")
		}
		assert.Equal(t, parentPath, pathSteps1[i-1], "The path should be a parent :else")
		assert.Equal(t, isParent, false, "Is parent should be False :else")
		workingPath = parentPath
	}

	// workingPath2 := pathsList[1]
	// for i := 1; i <= 3; i++ {
	// 	isParent, parentPath := GetParentPath(workingPath2)
	// 	if i == 3 {
	// 		assert.Equal(t, parentPath, pathSteps2[2], "The path should be a parent")
	// 		assert.Equal(t, isParent, true, "Is parent should be true")
	// 	}
	// 	assert.Equal(t, parentPath, pathSteps2[i-1], "The path should be a parent")
	// 	assert.Equal(t, isParent, false, "Is parent should be False")
	// }

	isParent, parentPath := GetParentPath(pathsList[2])
	assert.Equal(t, parentPath, "/home", "The parent should be /home")
	assert.Equal(t, isParent, true, "Is parent should be true for /home")

	isParent1, parentPath1 := GetParentPath(pathsList[3])
	assert.Equal(t, parentPath1, "Home", "The path should be a parent")
	assert.Equal(t, isParent1, true, "Is parent should be False")
}
