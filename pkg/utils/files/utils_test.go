package files

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestFileExists_Not(t *testing.T) {
	exists, err := FileExists("/bla/bli")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestFileExists_IsDirectory(t *testing.T) {
	_, err := FileExists("/")
	assert.Error(t, err)
}

func TestFileExists_Does(t *testing.T) {

	testDir, done, err := createTestDir()
	if !assert.NoError(t, err) {
		return
	}
	defer done()

	testFile, err := createTestFile(testDir)
	if !assert.NoError(t, err) {
		return
	}

	exists, err := FileExists(testFile)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestDirectoryExists_Does(t *testing.T) {
	testDir, done, err := createTestDir()
	if !assert.NoError(t, err) {
		return
	}
	defer done()

	exists, err := DirectoryExists(testDir)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestDirectoryExists_IsFile(t *testing.T) {
	testDir, done, err := createTestDir()
	if !assert.NoError(t, err) {
		return
	}
	defer done()

	testFile, err := createTestFile(testDir)
	if !assert.NoError(t, err) {
		return
	}

	_, err = DirectoryExists(testFile)
	assert.Error(t, err)
}

func TestCreaterDirectoryIfNotExists_IsFile(t *testing.T) {
	testDir, done, err := createTestDir()
	if !assert.NoError(t, err) {
		return
	}
	defer done()

	testFile, err := createTestFile(testDir)
	if !assert.NoError(t, err) {
		return
	}

	_, err = CreateDirectoryIfNotExists(testFile)
	assert.Error(t, err)
}

func createTestFile(dir string) (string, error) {
	tmpFile := path.Join(dir, uuid.NewV4().String())
	if err := os.WriteFile(tmpFile, []byte("hello"), os.ModePerm); err != nil {
		return "", err
	}
	return tmpFile, nil
}

func createTestDir() (string, func(), error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", nil, err
	}
	tmp := path.Join(wd, ".utilstests")
	dirExists, err := CreateDirectoryIfNotExists(tmp)
	if err != nil {
		return "", nil, err
	}
	return tmp, func() {
		if !dirExists {
			defer os.RemoveAll(tmp)
			defer os.Remove(tmp)
		}
	}, nil

}
