package client

import (
	"os"
	"crypto/sha1"
	"crypto/rand"
	"bufio"
	"fmt"
	"io"
)

// NewID generate a random hash that can be used as random key
func NewID() string {
	data := make([]byte, 16)
	rand.Read(data)
	return SHA1(data)
}

// SHA1 is a helper to quickly compute the SHA1 hash of a []byte.
func SHA1(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// FullSHA1 is helper to compute the SHA1 of the given file path.
func FullSHA1(path string) string {
	f, _ := os.Open(path)
	defer f.Close()
	reader := bufio.NewReader(f)
	h := sha1.New()
	_, _ = io.Copy(h, reader)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// a WriteResult keeps track of the number of blobs uploaded/skipped, and basic stats.
type WriteResult struct {
	Hash string

	Size int
	SizeSkipped int
	SizeUploaded int

	BlobsCount int
	BlobsSkipped int
	BlobsUploaded int

	FilesCount int
	FilesSkipped int
	FilesUploaded int

	DirsCount int
	DirsSkipped int
	DirsUploaded int

	AlreadyExists bool
}

func (wr *WriteResult) String() string {
	return fmt.Sprintf(`Write Result:
- Size: %v (skipped:%v, uploaded:%v)
- Blobs: %d (skipped:%d, uploaded:%d)
- Files: %d (skipped:%d, uploaded:%d)
- Dirs: %d (skipped:%d, uploaded:%d)
`,
		wr.Size, wr.SizeSkipped, wr.SizeUploaded,
		wr.BlobsCount, wr.BlobsSkipped, wr.BlobsUploaded,
		wr.FilesCount, wr.FilesSkipped, wr.FilesUploaded,
		wr.DirsCount, wr.DirsSkipped, wr.DirsUploaded)
}

// Add allows two WriteResult to be added.
func (wr *WriteResult) Add(wr2 *WriteResult) {
	wr.Size += wr2.Size
	wr.SizeSkipped += wr2.SizeSkipped
	wr.SizeUploaded += wr2.SizeUploaded

	wr.BlobsCount += wr2.BlobsCount
	wr.BlobsSkipped += wr2.BlobsSkipped
	wr.BlobsUploaded += wr2.BlobsUploaded

	wr.FilesCount += wr2.FilesCount
	wr.FilesSkipped += wr2.FilesSkipped
	wr.FilesUploaded += wr2.FilesUploaded

	wr.DirsCount += wr2.DirsCount
	wr.DirsSkipped += wr2.DirsSkipped
	wr.DirsUploaded += wr2.DirsUploaded
}

// a ReadResult keeps track of the number/size of downloaded blobs.
type ReadResult struct {
	Hash string
	Size int
	BlobsCount int
//	SkippedCnt int
//	SkippedSize int
	DownloadedCnt int
	DownloadedSize int
}

// Add allow two ReadResult to be added.
func (rr *ReadResult) Add(rr2 *ReadResult) {
	rr.Size += rr2.Size
	rr.BlobsCount += rr2.BlobsCount
//	rr.SkippedCnt += rr2.SkippedCnt
//	rr.SkippedSize += rr2.SkippedSize
	rr.DownloadedCnt += rr2.DownloadedCnt
	rr.DownloadedSize += rr2.DownloadedSize
}

// MatchResult checks if a WriteResult and a ReadResult have the same size.
func MatchResult(wr *WriteResult, rr *ReadResult) bool {
//	if wr.Size == rr.Size && wr.Hash == rr.Hash &&
//	   		wr.BlobsCount == rr.BlobsCount &&
//	   		(wr.SkippedCnt + wr.UploadedCnt) == rr.DownloadedCnt &&
//	   		(wr.SkippedSize + wr.UploadedSize) == rr.DownloadedSize {
//	   	return true
//	}
	return false
}
