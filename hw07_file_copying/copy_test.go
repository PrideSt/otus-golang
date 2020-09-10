package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// fileGetData returns content of file specifies by filePath
func fileGetData(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer closeOrErr(file)

	expData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	return expData
}

// filesEqual compare content of two files
func filesEqual(t *testing.T, expectedPath, actualPath string) {
	expDt := fileGetData(expectedPath)
	actDt := fileGetData(actualPath)

	require.Equal(t, expDt, actDt)
}

// fileSizeEqual compare whether file specified by path has size equal expectedSize
func fileSizeEqual(t *testing.T, path string, expectedSize int64) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer closeOrErr(f)

	inFileInfo, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	require.Equal(t, expectedSize, inFileInfo.Size())
}

func TestCopy(t *testing.T) {
	for _, tt := range [...]struct {
		name         string
		inputPath    string
		offset       int64
		limit        int64
		expectedPath string
		err          error
		errMsg       string
	}{
		{
			name:         `copy full file`,
			inputPath:    `testdata/input.txt`,
			expectedPath: `testdata/out_offset0_limit0.txt`,
		},
		{
			name:         `copy with limit`,
			inputPath:    `testdata/input.txt`,
			limit:        10,
			expectedPath: `testdata/out_offset0_limit10.txt`,
		},
		{
			name:         `copy with large limit`,
			inputPath:    `testdata/input.txt`,
			limit:        1000,
			expectedPath: `testdata/out_offset0_limit1000.txt`,
		},
		{
			name:         `copy with limit out of the input file`,
			inputPath:    `testdata/input.txt`,
			limit:        10000,
			expectedPath: `testdata/out_offset0_limit10000.txt`,
		},
		{
			name:         `copy with limit and offset`,
			inputPath:    `testdata/input.txt`,
			offset:       100,
			limit:        1000,
			expectedPath: `testdata/out_offset100_limit1000.txt`,
		},
		{
			name:         `copy with limit and large offset, offset with limit is out of input`,
			inputPath:    `testdata/input.txt`,
			offset:       6000,
			limit:        1000,
			expectedPath: `testdata/out_offset6000_limit1000.txt`,
		},
		{
			name:      `offset is out of input`,
			inputPath: `testdata/input.txt`,
			offset:    99999,
			err:       ErrOffsetExceedsFileSize,
			errMsg:    `offset exceeds file size: given offset 99999 to large, file "testdata/input.txt" has len 6617`,
		},
		{
			name:      `input file not exists`,
			inputPath: `testdata/no-input-file.txt`,
			err:       ErrInvalidPath,
			errMsg:    `invalid path: open for reading failed open testdata/no-input-file.txt: no such file or directory`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var tmpFile *os.File
			var err error

			if len(tt.expectedPath) > 0 {
				baseDir := "/tmp"
				tmpFile, err = ioutil.TempFile(baseDir, "test-copy")
				if err != nil {
					log.Fatal(err)
				}

				closeOrErr(tmpFile)

				defer func() {
					err := os.Remove(tmpFile.Name())
					if err != nil {
						log.Fatal(err)
					}
				}()
			}

			if tmpFile == nil {
				err = Copy(tt.inputPath, "", tt.offset, tt.limit)
			} else {
				err = Copy(tt.inputPath, tmpFile.Name(), tt.offset, tt.limit)
			}

			if tt.err == nil {
				require.NoError(t, err)
				filesEqual(t, tt.expectedPath, tmpFile.Name())
			} else {
				require.EqualError(t, err, tt.errMsg)
				require.True(t, errors.Is(err, tt.err))
			}
		})
	}
}

func TestDevRandom(t *testing.T) {
	for _, tt := range [...]struct {
		name   string
		offset int64
		limit  int64
		err    error
		errMsg string
	}{
		{
			name:  `copy with limit`,
			limit: 1000,
		},
		{
			name:   `copy with limit and offset`,
			offset: 100,
			limit:  1000,
			err:    ErrOffsetExceedsFileSize,
			errMsg: `offset exceeds file size: given offset 100 to large, file "/dev/urandom" has len 0`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var tmpFile *os.File
			var err error

			inputPath := "/dev/urandom"
			baseDir := "/tmp"
			tmpFile, err = ioutil.TempFile(baseDir, "test-dev-random")
			if err != nil {
				log.Fatal(err)
			}

			closeOrErr(tmpFile)

			defer func() {
				err := os.Remove(tmpFile.Name())
				if err != nil {
					log.Fatal(err)
				}
			}()

			err = Copy(inputPath, tmpFile.Name(), tt.offset, tt.limit)

			if tt.err == nil {
				require.NoError(t, err)
				fileSizeEqual(t, tmpFile.Name(), tt.limit)
			} else {
				require.EqualError(t, err, tt.errMsg)
				require.True(t, errors.Is(err, tt.err))
			}
		})
	}
}

func TestEmptyToPath(t *testing.T) {
	err := Copy(`testdata/input.txt`, "", 0, 0)

	require.EqualError(t, err, `invalid path: unable to open file "" for writing open : no such file or directory`)
	require.True(t, errors.Is(err, ErrInvalidPath))
}

func TestInputIsADir(t *testing.T) {
	var tmpFile *os.File
	var err error

	inputPath := "testdata/its_a_dir"

	baseDir := "/tmp"
	tmpFile, err = ioutil.TempFile(baseDir, "test-from-dir")
	if err != nil {
		log.Fatal(err)
	}

	closeOrErr(tmpFile)

	defer func() {
		err := os.Remove(tmpFile.Name())
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = Copy(inputPath, tmpFile.Name(), 0, 0)

	require.EqualError(t, err, "copy failed: read testdata/its_a_dir: is a directory")
	require.True(t, errors.Is(err, ErrCopy))
}
