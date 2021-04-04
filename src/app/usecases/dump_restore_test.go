package usecases

import (
	"krohobor/app/adapters/archive"
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
	"reflect"
	"testing"
)

func TestNewDumpRestore(t *testing.T) {
	db := database.NewMemory()
	store := storage.NewFileMock("/tmp/krohobor/storage", nil)

	type args struct {
		db    database.Interface
		store storage.Interface
	}
	tests := []struct {
		name string
		args args
		want *DumpRestore
	}{
		{
			name: "Test",
			args: args{
				db:    db,
				store: store,
			},
			want: &DumpRestore{
				db:    db,
				store: store,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDumpRestore(tt.args.db, tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDumpRestore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDumpRestore_Execute(t *testing.T) {
	dir := "/tmp/krohobor/storage"

	db := database.NewMemory()
	store := storage.NewFileMock(dir, nil)

	arch := archive.NewZipMock(dir, "")
	storeWithArch := storage.NewFileMock(dir, arch)

	type fields struct {
		db    database.Interface
		store storage.Interface
	}
	type args struct {
		request DumpRestoreRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    DumpRestoreResponse
		wantErr bool
	}{
		{
			name: "Test without arch - successful",
			fields: fields{
				db:    db,
				store: store,
			},
			args: args{
				request: DumpRestoreRequest{
					Name:     "test1",
					Filename: "file1.txt",
				},
			},
			want: DumpRestoreResponse{},
		},
		{
			name: "Test with arch - successful",
			fields: fields{
				db:    db,
				store: storeWithArch,
			},
			args: args{
				request: DumpRestoreRequest{
					Name:     "test1",
					Filename: "file1.txt",
				},
			},
			want: DumpRestoreResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dl := &DumpRestore{
				db:    tt.fields.db,
				store: tt.fields.store,
			}
			got, err := dl.Execute(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DumpRestore.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DumpRestore.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
