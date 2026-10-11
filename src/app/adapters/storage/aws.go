package storage

import (
	"context"
	"fmt"
	"io"
	"krohobor/app/adapters/archive"
	cfgApp "krohobor/app/adapters/config"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsS3 struct {
	bucket  string
	archive archive.Interface
	client  *s3.Client
}

func NewAwsS3(cfg cfgApp.AwsS3Config, arch archive.Interface) AwsS3 {
	creds := credentials.NewStaticCredentialsProvider(
		cfg.KeyId,
		cfg.AccessKey,
		"", // Session token (empty for most S3-compatible services)
	)

	conf, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		return AwsS3{}
	}

	client := s3.NewFromConfig(conf, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
	})

	return AwsS3{cfg.Catalog, arch, client}
}

func (s AwsS3) Check() error {
	if s.archive != nil {
		if err := s.archive.Check(); err != nil {
			return err
		}
	}

	_, err := s.client.GetBucketVersioning(context.TODO(), &s3.GetBucketVersioningInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s AwsS3) Filename(filename string) string {
	return "/tmp/" + filename
}

func (s AwsS3) Read(filename string) (string, error) {
	if s.archive != nil {
		filename = filename + s.archive.Ext()
	}

	fileParam := s.Filename(filepath.Base(filename))

	file, err := os.Create(fileParam)
	if err != nil {
		return "", err
	}

	defer file.Close()

	downloader := s3manager.NewDownloader(s.client)
	_, err = downloader.Download(context.TODO(), file, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return "", err
	}

	if s.archive != nil {
		if file, err := s.archive.Unarchive(fileParam); err == nil {
			return file, err
		} else {
			return "", err
		}
	}

	return fileParam, nil
}

func (s AwsS3) Write(filename string) error {
	filename = s.Filename(filename)
	fname := filename
	if s.archive != nil {
		archFile := fname + s.archive.Ext()
		if err := s.archive.Archive(archFile, fname); err != nil {
			return err
		}
		fname = archFile
	}

	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	key := filepath.Base(fname)
	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		return err
	}

	if err := os.Remove(filename); err != nil {
		return err
	}
	if s.archive != nil {
		if err := os.Remove(fname); err != nil {
			return err
		}
	}

	return nil
}

func (s AwsS3) Delete(filename string) error {
	if s.archive != nil {
		filename = filename + s.archive.Ext()
	}

	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	})

	if err != nil {
		return err
	}

	return nil
}

func (s AwsS3) Clean(filename string) error {
	if s.archive != nil {
		if err := os.Remove(filename + s.archive.Ext()); err != nil {
			fmt.Println(err)
		}
	}
	return os.Remove(filename)
}

func (s AwsS3) List() ([]string, error) {
	var out []string

	result, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return out, err
	}

	for _, item := range result.Contents {
		out = append(out, *item.Key)
	}

	return out, nil
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

	s.logger.Log(msg)
	return s.reader.Seek(offset, mode)
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
