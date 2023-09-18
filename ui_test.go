package swaggerui

import (
	"crypto/rand"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestSetupFs(t *testing.T) {
	ui, err := New(Spec("foo.yaml", []byte("bar")))
	assert.NoError(t, err)
	assert.NotNil(t, ui)

	uifs := ui.FileSystem()
	assert.NotNil(t, uifs)
	var hasSpecFile bool
	var hasInitializer bool
	_ = fs.WalkDir(uifs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.Type().IsRegular() {
			b, err := fs.ReadFile(uifs, d.Name())
			t.Log(d.Name())
			assert.NoErrorf(t, err, "error reading file %s: %s", d.Name(), err)
			assert.NotEmptyf(t, b, "file %s is empty", d.Name())
			if d.Name() == "foo.yaml" {
				hasSpecFile = true
				assert.Equal(t, "bar", string(b))
			}
			if d.Name() == InitializerFilename {
				hasInitializer = true
			}
		}

		return nil
	})
	assert.True(t, hasSpecFile)
	assert.True(t, hasInitializer)
}

type UiSuite struct {
	suite.Suite
}

func (suite *UiSuite) TestCustomInitializerContent() {
	var content []byte = randomBytes(1024)
	h, err := New(InitializerContent(content))
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), h)

	initjs, err := h.FileSystem().Open(InitializerFilename)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), initjs)

	b, err := fs.ReadFile(h.FileSystem(), InitializerFilename)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), b)
	assert.Equal(suite.T(), content, b)
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}
