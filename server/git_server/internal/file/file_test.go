package file_test

import (
	"fmt"
	"git_server/internal/file"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	data, err := file.ReadJson("thePrimeagen")
	if err != nil {
		fmt.Print(err)
	}
	assert.Nil(t, err, "Got a error reading file ")
	assert.NotNil(t, data, "Data should not be empty")
}
