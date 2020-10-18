package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
	"io"
	"os"
	"runtime/debug"
	"time"
)

type AwsS3 struct {
	Bucket string
}

func (s AwsS3) Read(filename string) error {
	fileParam := "/tmp/download_backup.zip"
	file, err := os.Create(fileParam)
	if err != nil {
		return err
	}

	defer file.Close()

	config, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return err
	}

	downloader := s3manager.NewDownloader(config)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(filename),
		})
	if err != nil {
		return err
	}

	return nil
}

func (s AwsS3) Write(filename string) error {
	config, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(config)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	readLogger := NewReadLogger(file, config.Logger)

	dt := time.Now()
	nowDate := dt.Format("2006-01-02_15-04")
	key := "backup_" + nowDate + ".zip"

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: &s.Bucket,
		Key:    &key,
		Body:   readLogger,
	}, func(u *s3manager.Uploader) {
		u.Concurrency = 1
		u.RequestOptions = append(u.RequestOptions, func(r *aws.Request) {
		})
	})

	return err
}

func (s AwsS3) Delete(filename string) error {
	config, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return err
	}

	svc := s3.New(config)
	svc.DeleteObjectRequest(&s3.DeleteObjectInput{
		Bucket: &s.Bucket,
		Key:    aws.String(filename),
	})

	return nil
}

func (s AwsS3) List() ([]string, error) {
	config, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return []string{}, err
	}

	var result []string
	svc := s3.New(config)
	req := svc.ListObjectsRequest(&s3.ListObjectsInput{Bucket: &s.Bucket})
	p := s3.NewListObjectsPaginator(req)
	for p.Next(context.TODO()) {
		page := p.CurrentPage()
		for _, obj := range page.Contents {
			result = append(result, *obj.Key)
		}
	}

	return result, nil
}

// Logger is a logger use for logging the readers usage.
type Logger interface {
	Log(args ...interface{})
}

// ReadSeeker interface provides the interface for a Reader, Seeker, and ReadAt.
type ReadSeeker interface {
	io.ReadSeeker
	io.ReaderAt
}

// ReadLogger wraps an reader with logging for access.
type ReadLogger struct {
	reader ReadSeeker
	logger Logger
}

// NewReadLogger a ReadLogger that wraps the passed in ReadSeeker (Reader,
// Seeker, ReadAt) with a logger.
func NewReadLogger(r ReadSeeker, logger Logger) *ReadLogger {
	return &ReadLogger{
		reader: r,
		logger: logger,
	}
}

// Seek offsets the reader's current position for the next read.
func (s *ReadLogger) Seek(offset int64, mode int) (int64, error) {
	newOffset, err := s.reader.Seek(offset, mode)
	msg := fmt.Sprintf(
		"ReadLogger.Seek(offset:%d, mode:%d) (newOffset:%d, err:%v)",
		offset, mode, newOffset, err)
	if err != nil {
		msg += fmt.Sprintf("\n\tStack:\n%s", string(debug.Stack()))
	}

	//s.logger.Log(msg)
	return newOffset, err
}

// Read attempts to read from the reader, returning the bytes read, or error.
func (s *ReadLogger) Read(b []byte) (int, error) {
	n, err := s.reader.Read(b)
	msg := fmt.Sprintf(
		"ReadLogger.Read(len(bytes):%d) (read:%d, err:%v)",
		len(b), n, err)
	if err != nil {
		msg += fmt.Sprintf("\n\tStack:\n%s", string(debug.Stack()))
	}

	s.logger.Log(msg)
	return n, err
}

// ReadAt will read the underlying reader starting at the offset.
func (s *ReadLogger) ReadAt(b []byte, offset int64) (int, error) {
	n, err := s.reader.ReadAt(b, offset)
	msg := fmt.Sprintf(
		"ReadLogger.ReadAt(len(bytes):%d, offset:%d) (read:%d, err:%v)",
		len(b), offset, n, err)
	if err != nil {
		msg += fmt.Sprintf("\n\tStack:\n%s", string(debug.Stack()))
	}

	s.logger.Log(msg)
	return n, err
}
