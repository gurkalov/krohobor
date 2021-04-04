package usecases

import (
	"krohobor/app/adapters/archive"
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
	"reflect"
	"testing"
)

func TestNewDumpCreate(t *testing.T) {
	db := database.NewMemory()
	store := storage.NewFileMock("/tmp/krohobor/storage", nil)

	type args struct {
		db    database.Interface
		store storage.Interface
	}
	tests := []struct {
		name string
		args args
		want *DumpCreate
	}{
		{
			name: "Test",
			args: args{
				db:    db,
				store: store,
			},
			want: &DumpCreate{
				db:    db,
				store: store,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDumpCreate(tt.args.db, tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDumpCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDumpCreate_Execute(t *testing.T) {
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
		request DumpCreateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    DumpCreateResponse
		wantErr bool
	}{
		{
			name: "Create dump for one database without arch - successful",
			fields: fields{
				db:    db,
				store: store,
			},
			args: args{
				request: DumpCreateRequest{
					DbNames:  []string{"test1"},
					Filename: "test",
				},
			},
			want: DumpCreateResponse{},
		},
		{
			name: "Create dump for one database with arch - successful",
			fields: fields{
				db:    db,
				store: storeWithArch,
			},
			args: args{
				request: DumpCreateRequest{
					DbNames:  []string{"test1"},
					Filename: "test",
				},
			},
			want: DumpCreateResponse{},
		},
		{
			name: "Create dump for all databases without arch - successful",
			fields: fields{
				db:    db,
				store: store,
			},
			args: args{
				request: DumpCreateRequest{
					DbNames:  []string{},
					Filename: "test-all",
				},
			},
			want: DumpCreateResponse{},
		},
		{
			name: "Create dump for all databases with arch - successful",
			fields: fields{
				db:    db,
				store: storeWithArch,
			},
			args: args{
				request: DumpCreateRequest{
					DbNames:  []string{},
					Filename: "test-all",
				},
			},
			want: DumpCreateResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dd := &DumpCreate{
				db:    tt.fields.db,
				store: tt.fields.store,
			}
			got, err := dd.Execute(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DumpCreate.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DumpCreate.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
