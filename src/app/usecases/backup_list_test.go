package usecases

import (
	"krohobor/app/adapters/archive"
	"krohobor/app/adapters/storage"
	"reflect"
	"testing"
)

func TestNewBackupList(t *testing.T) {
	store := storage.NewFileMock("/tmp/krohobor/storage", nil)

	type args struct {
		store storage.Interface
	}
	tests := []struct {
		name string
		args args
		want *BackupList
	}{
		{
			name: "Test",
			args: args{
				store: store,
			},
			want: &BackupList{
				store: store,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBackupList(tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBackupList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackupList_Execute(t *testing.T) {
	dir := "/tmp/krohobor/storage"

	store := storage.NewFileMock(dir, nil)

	arch := archive.NewZipMock(dir, "")
	storeWithArch := storage.NewFileMock(dir, arch)

	type fields struct {
		store storage.Interface
	}
	type args struct {
		request BackupListRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    BackupListResponse
		wantErr bool
	}{
		{
			name: "List - successful",
			fields: fields{
				store: storeWithArch,
			},
			args: args{
				request: BackupListRequest{},
			},
			want: BackupListResponse{
				List: []string{
					"file1.txt",
					"file2.txt",
				},
			},
		},
		{
			name: "List with archive - successful",
			fields: fields{
				store: store,
			},
			args: args{
				request: BackupListRequest{},
			},
			want: BackupListResponse{
				List: []string{
					"file1.txt",
					"file1.txt.zip",
					"file2.txt",
					"file2.txt.zip",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dl := &BackupList{
				store: tt.fields.store,
			}
			got, err := dl.Execute(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("BackupList.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BackupList.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
