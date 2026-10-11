package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"krohobor/app/adapters/archive"
	cfgApp "krohobor/app/adapters/config"
	"krohobor/app/domain"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type AwsS3 struct {
	bucket   string
	archive  archive.Interface
	client   *s3.Client
	progress func(done, total int)
}

func NewAwsS3(cfg cfgApp.AwsS3Config, arch archive.Interface) *AwsS3 {
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
		return &AwsS3{}
	}

	client := s3.NewFromConfig(conf, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
		o.UsePathStyle = true // bucket in the URL path (required for MinIO/S3-compatible endpoints)
	})

	return &AwsS3{bucket: cfg.Catalog, archive: arch, client: client}
}

// SetProgress registers a callback invoked as a dump is uploaded to S3.
func (s *AwsS3) SetProgress(p func(done, total int)) {
	s.progress = p
}

func (s AwsS3) Check() error {
	if s.archive != nil {
		if err := s.archive.Check(); err != nil {
			return err
		}
	}

	if err := s.ensureBucket(context.TODO()); err != nil {
		return err
	}

	_, err := s.client.GetBucketVersioning(context.TODO(), &s3.GetBucketVersioningInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return err
	}

	return nil
}

// ensureBucket creates the bucket if it does not exist yet.
// This makes the storage self-sufficient for local testing (fresh MinIO/librefs
// instance) without a separate bucket bootstrap step.
func (s AwsS3) ensureBucket(ctx context.Context) error {
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err == nil {
		return nil
	}

	var notFound *types.NotFound
	if !errors.As(err, &notFound) {
		return err
	}

	_, createErr := s.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if createErr == nil {
		return nil
	}

	var alreadyOwned *types.BucketAlreadyOwnedByYou
	var alreadyExists *types.BucketAlreadyExists
	if errors.As(createErr, &alreadyOwned) || errors.As(createErr, &alreadyExists) {
		return nil
	}
	return createErr
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

	obj, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return "", err
	}
	defer obj.Body.Close()

	if _, err = io.Copy(file, obj.Body); err != nil {
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

	if err := s.ensureBucket(context.TODO()); err != nil {
		return err
	}

	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	key := filepath.Base(fname)

	body := io.Reader(file)
	if s.progress != nil {
		var total int
		if info, statErr := file.Stat(); statErr == nil {
			total = int(info.Size())
		}
		body = &progressReader{r: file, total: total, cb: s.progress}
	}

	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   body,
	}, func(o *s3.Options) {
		// The counting progress reader is not seekable, and the SDK would
		// otherwise pre-read the whole body to compute a checksum header.
		// MinIO/librefs do not require checksums, so skip the calculation.
		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
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

func (s AwsS3) List() ([]domain.Dump, error) {
	var out []domain.Dump

	result, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return out, err
	}

	for _, item := range result.Contents {
		out = append(out, domain.Dump{Name: *item.Key, Size: int(*item.Size)})
	}

	return out, nil
}

// progressReader wraps an io.Reader and reports bytes read to a callback.
type progressReader struct {
	r     io.ReadSeeker
	done  int
	total int
	cb    func(done, total int)
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.r.Read(b)
	p.done += n
	if p.done > p.total {
		p.done = p.total // the SDK may rewind and re-read the body
	}
	if p.cb != nil {
		p.cb(p.done, p.total)
	}
	return n, err
}

func (p *progressReader) Seek(offset int64, whence int) (int64, error) {
	return p.r.Seek(offset, whence)
}
