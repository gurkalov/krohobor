// +build aws

package storage

import (
	"krohobor/app/adapters/archive"
	"reflect"
	"testing"
)

const (
	bucket = "krohobor-backup"
)

func TestNewAwsS3(t *testing.T) {
	arch := archive.NewZipMock(storageDir, "")

	type args struct {
		bucket string
		arch   archive.Interface
	}
	tests := []struct {
		name string
		args args
		want AwsS3
	}{
		{
			name: "Test",
			args: args{
				bucket: bucket,
				arch:   arch,
			},
			want: AwsS3{
				bucket:  bucket,
				archive: arch,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAwsS3Test(tt.args.bucket, tt.args.arch)
			tt.want.client = got.client

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAwsS3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAwsS3_Check(t *testing.T) {
	arch := archive.NewZipMock(storageDir, "")

	type fields struct {
		bucket  string
		archive archive.Interface
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Check - successful",
			fields: fields{
				bucket:  bucket,
				archive: arch,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAwsS3Test(tt.fields.bucket, tt.fields.archive)
			if err := s.Check(); (err != nil) != tt.wantErr {
				t.Errorf("AwsS3.Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAwsS3_Filename(t *testing.T) {
	arch := archive.NewZipMock(storageDir, "")

	type fields struct {
		bucket  string
		archive archive.Interface
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Filename with arch - successful",
			fields: fields{
				bucket:  bucket,
				archive: arch,
			},
			args: args{
				filename: "test.txt",
			},
			want: "/tmp/test.txt",
		},
		{
			name: "Filename without arch - successful",
			fields: fields{
				bucket:  bucket,
				archive: nil,
			},
			args: args{
				filename: "test.txt",
			},
			want: "/tmp/test.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAwsS3Test(tt.fields.bucket, tt.fields.archive)
			if got := s.Filename(tt.args.filename); got != tt.want {
				t.Errorf("AwsS3.Filename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAwsS3_Read(t *testing.T) {
	arch := archive.NewZipMock(storageDir, "")

	type fields struct {
		bucket  string
		archive archive.Interface
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Read with arch first - successful",
			fields: fields{
				bucket:  bucket,
				archive: arch,
			},
			args: args{
				filename: "file1",
			},
			want:    "/tmp/krohobor/storage/file1",
			wantErr: false,
		},
		{
			name: "Read with arch not exists - error",
			fields: fields{
				bucket:  bucket,
				archive: arch,
			},
			args: args{
				filename: "file404",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Read without arch first - successful",
			fields: fields{
				bucket:  bucket,
				archive: nil,
			},
			args: args{
				filename: "file1",
			},
			want:    "/tmp/file1",
			wantErr: false,
		},
		{
			name: "Read without arch not exists - error",
			fields: fields{
				bucket:  bucket,
				archive: nil,
			},
			args: args{
				filename: "file404",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAwsS3Test(tt.fields.bucket, tt.fields.archive)
			got, err := s.Read(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("AwsS3.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AwsS3.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAwsS3_Write(t *testing.T) {
	arch := archive.NewZipMock(storageDir, "")

	type fields struct {
		bucket  string
		archive archive.Interface
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Write with arch first - successful",
			fields: fields{
				bucket:  bucket,
				archive: arch,
			},
			args: args{
				filename: "krohobor/storage/mock/file1",
			},
			wantErr: false,
		},
		{
			name: "Write with arch second - successful",
			fields: fields{
				bucket:  bucket,
				archive: arch,
			},
			args: args{
				filename: "krohobor/storage/mock/file2",
			},
			wantErr: false,
		},
		{
			name: "Write without arch first - successful",
			fields: fields{
				bucket:  bucket,
				archive: nil,
			},
			args: args{
				filename: "file3",
			},
			wantErr: false,
		},
		{
			name: "Write without arch second - successful",
			fields: fields{
				bucket:  bucket,
				archive: nil,
			},
			args: args{
				filename: "file4",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAwsS3Test(tt.fields.bucket, tt.fields.archive)
			if err := s.Write(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("AwsS3.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAwsS3_Delete(t *testing.T) {
	arch := archive.NewZipMock(storageDir, "")

	type fields struct {
		bucket  string
		archive archive.Interface
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Delete with arch first - successful",
			fields: fields{
				bucket:  bucket,
				archive: arch,
			},
			args: args{
				filename: "file1",
			},
			wantErr: false,
		},
		{
			name: "Delete with arch not exists - successful",
			fields: fields{
				bucket:  bucket,
				archive: arch,
			},
			args: args{
				filename: "file404",
			},
			wantErr: false,
		},
		{
			name: "Delete without arch first - successful",
			fields: fields{
				bucket:  bucket,
				archive: nil,
			},
			args: args{
				filename: "file1",
			},
			wantErr: false,
		},
		{
			name: "Delete without arch not exists - successful",
			fields: fields{
				bucket:  bucket,
				archive: nil,
			},
			args: args{
				filename: "file404",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAwsS3Test(tt.fields.bucket, tt.fields.archive)
			if err := s.Delete(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("AwsS3.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAwsS3_List(t *testing.T) {
	arch := archive.NewZipMock(storageDir, "")

	type fields struct {
		bucket  string
		archive archive.Interface
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name: "List with arch - successful",
			fields: fields{
				bucket:  bucket,
				archive: arch,
			},
			want: []string{
				"file1",
				"file2",
			},
			wantErr: false,
		},
		{
			name: "List without arch - successful",
			fields: fields{
				bucket:  bucket,
				archive: nil,
			},
			want: []string{
				"file1",
				"file2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAwsS3Test(tt.fields.bucket, tt.fields.archive)
			got, err := s.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("AwsS3.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AwsS3.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
