package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
	"io"
	"io/ioutil"
	"krohobor/app/adapters/archive"
	"krohobor/app/adapters/config"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

type AwsS3 struct {
	bucket  string
	archive archive.Interface
	client  *s3.Client
}

func NewAwsS3(cfg config.AwsS3Config, arch archive.Interface) AwsS3 {
	conf, err := external.LoadDefaultAWSConfig(
		external.WithRegion(cfg.Region),
		external.WithCredentialsValue(aws.Credentials{
			AccessKeyID:     cfg.KeyId,
			SecretAccessKey: cfg.AccessKey,
		}),
	)
	if err != nil {
		return AwsS3{}
	}

	return AwsS3{cfg.Catalog, arch, s3.New(conf)}
}

func NewAwsS3Test(bucket string, arch archive.Interface) AwsS3 {
	conf, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return AwsS3{}
	}

	client := s3.New(conf)

	iter := s3manager.NewDeleteListIterator(client, &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})
	if err := s3manager.NewBatchDeleteWithClient(client).Delete(context.TODO(), iter); err != nil {
		fmt.Println(err)
	}

	reqDelete := client.DeleteBucketRequest(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	_, err = reqDelete.Send(context.TODO())
	if err != nil {
		fmt.Println(err)
	}

	reqCreate := client.CreateBucketRequest(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	_, err = reqCreate.Send(context.TODO())
	if err != nil {
		fmt.Println(err)
	}

	awsS3 := AwsS3{bucket, arch, client}

	d1 := []byte("hello")
	if err := ioutil.WriteFile("/tmp/file1", d1, 0644); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile("/tmp/file2", d1, 0644); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile("/tmp/file3", d1, 0644); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile("/tmp/file4", d1, 0644); err != nil {
		panic(err)
	}

	_ = awsS3.Write("file1")
	_ = awsS3.Write("file2")

	return awsS3
}

func (s AwsS3) Check() error {
	if s.archive != nil {
		if err := s.archive.Check(); err != nil {
			return err
		}
	}

	conf, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return err
	}

	svc := s3.New(conf)
	req := svc.GetBucketVersioningRequest(&s3.GetBucketVersioningInput{
		Bucket: &s.bucket,
	})

	_, err = req.Send(context.TODO())

	return err
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

	downloader := s3manager.NewDownloader(s.client.Config)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
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

	readLogger := NewReadLogger(file, s.client.Config.Logger)

	key := filepath.Base(fname)

	uploader := s3manager.NewUploader(s.client.Config)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: &s.bucket,
		Key:    &key,
		Body:   readLogger,
	}, func(u *s3manager.Uploader) {
		u.Concurrency = 1
		u.RequestOptions = append(u.RequestOptions, func(r *aws.Request) {
		})
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

	req := s.client.DeleteObjectRequest(&s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    aws.String(filename),
	})

	_, err := req.Send(context.TODO())

	return err
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
	var result []string

	req := s.client.ListObjectsRequest(&s3.ListObjectsInput{Bucket: &s.bucket})
	p := s3.NewListObjectsPaginator(req)
	for p.Next(context.TODO()) {
		page := p.CurrentPage()
		for _, obj := range page.Contents {
			if s.archive != nil && s.archive.Ext() != "" {
				if s.archive.Ext() == filepath.Ext(*obj.Key) {
					result = append(result, strings.TrimSuffix(*obj.Key, s.archive.Ext()))
				}
			} else {
				result = append(result, *obj.Key)
			}
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
