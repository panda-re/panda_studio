package util

import (
	"hash"
	"io"
)

func NewWriterHasher(writer io.Writer, hasher hash.Hash) io.Writer {
	return io.MultiWriter(writer, hasher)
}

func NewReaderHasher(reader io.Reader, hasher hash.Hash) io.Reader {
	rdPipe, wrPipe := io.Pipe()
	mWriter := NewWriterHasher(wrPipe, hasher)

	go func() {
		// todo: make this cancellable
		defer wrPipe.Close()
		_, err := io.Copy(mWriter, reader)
		if err != nil {
			wrPipe.CloseWithError(err)
		}
	}()

	return rdPipe
}