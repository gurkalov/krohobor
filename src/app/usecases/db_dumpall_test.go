package usecases

import (
	"krohobor/app/adapters/archive"
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
	"reflect"
	"testing"
)

func TestNewDbDumpAll(t *testing.T) {
	db := database.NewMemory()
	store := storage.NewFileMock("/tmp/krohobor/storage", nil)

	type args struct {
		db    database.Interface
		store storage.Interface
	}
	tests := []struct {
		name string
		args args
		want *DbDumpAll
	}{
		{
			name: "Test",
			args: args{
				db: db,
				store: store,
			},
			want: &DbDumpAll{
				db: db,
				store: store,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDbDumpAll(tt.args.db, tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDbDumpAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbDumpAll_Execute(t *testing.T) {
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
		request DbDumpAllRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    DbDumpAllResponse
		wantErr bool
	}{
		{
			name: "Test without arch - successful",
			fields: fields{
				db: db,
				store: store,
			},
			args: args{
				request: DbDumpAllRequest{
					Filename: "test-all",
				},
			},
			want: DbDumpAllResponse{},
		},
		{
			name: "Test with arch - successful",
			fields: fields{
				db: db,
				store: storeWithArch,
			},
			args: args{
				request: DbDumpAllRequest{
					Filename: "test-all-arch",
				},
			},
			want: DbDumpAllResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ddl := &DbDumpAll{
				db:    tt.fields.db,
				store: tt.fields.store,
			}
			got, err := ddl.Execute(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbDumpAll.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbDumpAll.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
