package usecases

import (
	"krohobor/app/adapters/archive"
	"krohobor/app/adapters/storage"
	"reflect"
	"testing"
)

func TestNewDumpList(t *testing.T) {
	store := storage.NewFileMock("/tmp/krohobor/storage", nil)

	type args struct {
		store storage.Interface
	}
	tests := []struct {
		name string
		args args
		want *DumpList
	}{
		{
			name: "Test",
			args: args{
				store: store,
			},
			want: &DumpList{
				store: store,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDumpList(tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDumpList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDumpList_Execute(t *testing.T) {
	dir := "/tmp/krohobor/storage"

	store := storage.NewFileMock(dir, nil)

	arch := archive.NewZipMock(dir, "")
	storeWithArch := storage.NewFileMock(dir, arch)

	type fields struct {
		store storage.Interface
	}
	type args struct {
		request DumpListRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    DumpListResponse
		wantErr bool
	}{
		{
			name: "List - successful",
			fields: fields{
				store: storeWithArch,
			},
			args: args{
				request: DumpListRequest{},
			},
			want: DumpListResponse{
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
				request: DumpListRequest{},
			},
			want: DumpListResponse{
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
			dl := &DumpList{
				store: tt.fields.store,
			}
			got, err := dl.Execute(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DumpList.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DumpList.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
