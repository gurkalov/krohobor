package storage

import (
	"krohobor/app/adapters/archive"
	"reflect"
	"testing"
)

const (
	storageDir = "/tmp/krohobor/storage"
)

func TestNewFile(t *testing.T) {
	arch := archive.NewZipMock(storageDir, "")

	type args struct {
		dir  string
		arch archive.Interface
	}
	tests := []struct {
		name string
		args args
		want File
	}{
		{
			name: "Test",
			args: args{
				dir:  storageDir,
				arch: arch,
			},
			want: File{
				dir:     storageDir,
				archive: arch,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFile(tt.args.dir, tt.args.arch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_Check(t *testing.T) {
	arch := archive.NewZipMock(storageDir, "")
	NewFileMock(storageDir, arch)

	tests := []struct {
		name    string
		file    File
		wantErr bool
	}{
		{
			name:    "Check with archive - successful",
			file:    NewFile(storageDir, arch),
			wantErr: false,
		},
		{
			name:    "Check witour archive - successful",
			file:    NewFile(storageDir, nil),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.file.Check(); (err != nil) != tt.wantErr {
				t.Errorf("File.Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFile_Read(t *testing.T) {
	arch := archive.NewZipMock(storageDir, "")
	NewFileMock(storageDir, arch)

	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		file    File
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Read exists with arch - successful",
			file: NewFile(storageDir, arch),
			args: args{
				filename: "file1.txt",
			},
			want:    storageDir + "/file1.txt",
			wantErr: false,
		},
		{
			name: "Read exists without arch - successful",
			file: NewFile(storageDir, nil),
			args: args{
				filename: "file1.txt",
			},
			want:    storageDir + "/file1.txt",
			wantErr: false,
		},
		{
			name: "Read not exists with arch - error",
			file: NewFile(storageDir, arch),
			args: args{
				filename: "file404.txt",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Read not exists without arch - error",
			file: NewFile(storageDir, nil),
			args: args{
				filename: "file404.txt",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.file.Read(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("File.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("File.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_Write(t *testing.T) {
	arch := archive.NewZipMock(storageDir, "")
	NewFileMock(storageDir, arch)

	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		file    File
		args    args
		wantErr bool
	}{
		{
			name: "Write with arch - successful",
			file: NewFile(storageDir, arch),
			args: args{
				filename: "file1.txt",
			},
			wantErr: false,
		},
		{
			name: "Write without arch - successful",
			file: NewFile(storageDir, nil),
			args: args{
				filename: "file1.txt",
			},
			wantErr: false,
		},
		{
			name: "Write with arch - error",
			file: NewFile(storageDir, arch),
			args: args{
				filename: "file404.txt",
			},
			wantErr: true,
		},
		{
			name: "Write without arch - error",
			file: NewFile(storageDir, nil),
			args: args{
				filename: "file404.txt",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.file.Write(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("File.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFile_Delete(t *testing.T) {
	arch := archive.NewZipMock(storageDir, "")
	NewFileMock(storageDir, arch)

	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		file    File
		args    args
		wantErr bool
	}{
		{
			name: "Delete exists - successful",
			file: NewFile(storageDir, arch),
			args: args{
				"file1.txt",
			},
			wantErr: false,
		},
		{
			name: "Delete not exists - successful",
			file: NewFile(storageDir, arch),
			args: args{
				"file404.txt",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.file.Delete(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("File.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFile_List(t *testing.T) {
	arch := archive.NewZipMock(storageDir, "")
	NewFileMock(storageDir, arch)

	tests := []struct {
		name    string
		file    File
		want    []string
		wantErr bool
	}{
		{
			name: "List with archive - successful",
			file: NewFile(storageDir, arch),
			want: []string{
				"file1.txt",
				"file2.txt",
			},
			wantErr: false,
		},
		{
			name: "List without archive - successful",
			file: NewFile(storageDir, nil),
			want: []string{
				"file1.txt",
				"file1.txt.zip",
				"file2.txt",
				"file2.txt.zip",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.file.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("File.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("File.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
