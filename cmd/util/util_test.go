package util

import (
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

	pathSteps2 := []string{
		"Home/username/projects",
		"Home/username",
		"Home",
	}

	workingPath := pathsList[0]
	for i := 1; i <= 10; i++ {
		isParent, parentPath := GetParentPath(workingPath)
		workingPath = parentPath
		if i >= 4 {
			assert.Equal(t, pathSteps1[2], parentPath, "The path should be a parent from :last")
			assert.Equal(t, true, isParent, "Is parent should be true from :last")
			continue
		}
		assert.Equal(t, pathSteps1[i-1], parentPath, "The path should be a parent :else")
		assert.Equal(t, false, isParent, "Is parent should be False :else")
	}

	workingPath2 := pathsList[1]
	for i := 1; i <= 10; i++ {
		isParent, parentPath := GetParentPath(workingPath2)
		workingPath2 = parentPath
		if i >= 4 {
			assert.Equal(t, pathSteps2[2], parentPath, "The path should be a parent")
			assert.Equal(t, true, isParent, "Is parent should be true")
			continue
		}
		assert.Equal(t, pathSteps2[i-1], parentPath, "The path should be a parent")
		assert.Equal(t, false, isParent, "Is parent should be False")
	}

	isParent, parentPath := GetParentPath(pathsList[2])
	assert.Equal(t, "/home", parentPath, "The parent should be /home")
	assert.Equal(t, true, isParent, "Is parent should be true for /home")

	isParent1, parentPath1 := GetParentPath(pathsList[3])
	assert.Equal(t, "Home", parentPath1, "The path should be a parent")
	assert.Equal(t, true, isParent1, "Is parent should be False")

	isRoot, parentPath2 := GetParentPath("/")
	assert.Equal(t, "/", parentPath2, "Passed '/' and should get '/'")
	assert.Equal(t, true, isRoot, "This should be true")

	isRoot, file := GetParentPath("test.json")
	assert.Equal(t, "test.json", file, "Passed 'test.json' and should get 'test.json'")
	assert.Equal(t, true, isRoot, "This should be true")

	isRoot, file = GetParentPath("cmd")
	assert.Equal(t, "cmd", file, "Passed 'cmd' and should get 'cmd'")
	assert.Equal(t, true, isRoot, "This should be true")
}
