package usecases

import (
	"krohobor/app/adapters/archive"
	"krohobor/app/adapters/storage"
	"reflect"
	"testing"
)

func TestNewDumpDelete(t *testing.T) {
	store := storage.NewFileMock("/tmp/krohobor/storage", nil)

	type args struct {
		store storage.Interface
	}
	tests := []struct {
		name string
		args args
		want *DumpDelete
	}{
		{
			name: "Test",
			args: args{
				store: store,
			},
			want: &DumpDelete{
				store: store,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDumpDelete(tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDumpDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDumpDelete_Execute(t *testing.T) {
	dir := "/tmp/krohobor/storage"

	store := storage.NewFileMock(dir, nil)

	arch := archive.NewZipMock(dir, "")
	storeWithArch := storage.NewFileMock(dir, arch)

	type fields struct {
		store storage.Interface
	}
	type args struct {
		request DumpDeleteRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    DumpDeleteResponse
		wantErr bool
	}{
		{
			name: "Delete dump without arch - successful",
			fields: fields{
				store: store,
			},
			args: args{
				request: DumpDeleteRequest{
					Name: "file1.txt",
				},
			},
			want: DumpDeleteResponse{},
		},
		{
			name: "Delete dump arch - successful",
			fields: fields{
				store: storeWithArch,
			},
			args: args{
				request: DumpDeleteRequest{
					Name: "file1.txt",
				},
			},
			want: DumpDeleteResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dl := &DumpDelete{
				store: tt.fields.store,
			}
			got, err := dl.Execute(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DumpDelete.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DumpDelete.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
