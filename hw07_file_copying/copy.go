package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrInvalidPath           = errors.New("invalid path")
	ErrSetOffset             = errors.New("unable to set offset")
	ErrFileStat              = errors.New("unable to get file's stat")
	ErrCopy                  = errors.New("copy failed")
)

// min return argument with minimum value.
func min(a, b int64) int64 {
	if a > b {
		return b
	}

	return a
}

// closeOrErr shortcut for defer function, not forget close files.
func closeOrErr(f io.Closer) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

// Copy limit bytes (or while not EOF if limit = 0) from file specified by fromPath to
// file specified by toPath. Coping starts at offset bytes from input.
// Return error when offset large then input file or input file not seekable.
func Copy(fromPath string, toPath string, offset int64, limit int64) error {
	var reader io.Reader

	in, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("%w: open for reading failed %s", ErrInvalidPath, err)
	}
	defer closeOrErr(in)

	inFileInfo, err := in.Stat()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFileStat, err)
	}
	inFileSize := inFileInfo.Size()

	if offset > 0 {
		if offset > inFileSize {
			return fmt.Errorf(
				"%w: given offset %d to large, file %q has len %d",
				ErrOffsetExceedsFileSize,
				offset,
				fromPath,
				inFileSize,
			)
		}
		if _, err := in.Seek(offset, io.SeekStart); err != nil {
			return fmt.Errorf("%w: set position %d in file %q failed %s", ErrSetOffset, offset, from, err)
		}
	}

	var needToCopyBytes int64
	if limit > 0 {
		reader = io.LimitReader(in, limit)
		needToCopyBytes = min(limit, inFileSize-offset)
	} else {
		reader = in
		needToCopyBytes = inFileSize
	}

	writer, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("%w: unable to open file %q for writing %s", ErrInvalidPath, toPath, err)
	}
	defer closeOrErr(writer)

	bar := pb.Full.Start64(needToCopyBytes)
	barReader := bar.NewProxyReader(reader)
	_, err = io.Copy(writer, barReader)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCopy, err)
	}
	bar.Finish()

	return nil
}
