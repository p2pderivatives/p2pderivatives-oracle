package utils_test

import (
	"p2pderivatives-oracle/internal/utils"
	"p2pderivatives-oracle/test"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var PasswordDirectory = filepath.Join(test.VectorsDirectoryPath, "password")

var PasswordFiles = map[string]string{
	"pass_0.txt":             "qfVC6IFiCUjjJNI7cX+oZdHcvCQbmJXWFvtXWf7Oq0M=",
	"pass_1.txt":             "uggt9zMYhZX777fPqr7gS5gRNNAKfzHTsB9JC7uUDj4=",
	"pass_with_new_line.txt": "Cp2XZwgUhenKFUWP0jlVVnWy9ifUhYz+DXSy56tQ8yI=",
}

func TestReadPasswordFile_Success(t *testing.T) {
	for path, expectedPass := range PasswordFiles {
		pass, err := utils.ReadFirstLineFromFile(filepath.Join(PasswordDirectory, path))
		assert.NoError(t, err)
		assert.Equal(t, expectedPass, pass)
	}
}
