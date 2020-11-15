package usecases

import (
	"krohobor/app/adapters/archive"
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
	"reflect"
	"testing"
)

func TestNewDbDump(t *testing.T) {
	db := database.NewMemory()
    store := storage.NewFileMock("/tmp/krohobor/storage", nil)

	type args struct {
		db    database.Interface
		store storage.Interface
	}
	tests := []struct {
		name string
		args args
		want *DbDump
	}{
		 {
		 	 name: "Test",
		 	 args: args{
		 	 	 db: db,
		 	 	 store: store,
			 },
			 want: &DbDump{
				 db: db,
				 store: store,
			 },
		 },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDbDump(tt.args.db, tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDbDump() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbDump_Execute(t *testing.T) {
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
		request DbDumpRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    DbDumpResponse
		wantErr bool
	}{
		{
			name: "Test without arch - successful",
			fields: fields{
				db: db,
				store: store,
		    },
			args: args{
				request: DbDumpRequest{
					Name: "test1",
					Filename: "test",
				},
			},
			want: DbDumpResponse{},
		},
		{
			name: "Test with arch - successful",
			fields: fields{
				db: db,
				store: storeWithArch,
			},
			args: args{
				request: DbDumpRequest{
					Name: "test1",
					Filename: "test",
				},
			},
			want: DbDumpResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dd := &DbDump{
				db:    tt.fields.db,
				store: tt.fields.store,
			}
			got, err := dd.Execute(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbDump.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbDump.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
