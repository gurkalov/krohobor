package usecases

import (
	"krohobor/app/adapters/archive"
	"krohobor/app/adapters/storage"
	"krohobor/app/domain"
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
				List: []domain.Dump{
					{Name: "file1.txt"},
					{Name: "file2.txt"},
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
				List: []domain.Dump{
					{Name: "file1.txt", Size: 5},
					{Name: "file1.txt.zip"},
					{Name: "file2.txt", Size: 5},
					{Name: "file2.txt.zip"},
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
			if len(got.List) != len(tt.want.List) {
				t.Fatalf("DumpList.Execute() len = %d, want %d", len(got.List), len(tt.want.List))
			}
			for i := range got.List {
				if got.List[i].Name != tt.want.List[i].Name {
					t.Errorf("DumpList.Execute()[%d].Name = %v, want %v", i, got.List[i].Name, tt.want.List[i].Name)
				}
				if tt.want.List[i].Size > 0 && got.List[i].Size != tt.want.List[i].Size {
					t.Errorf("DumpList.Execute()[%d].Size = %v, want %v", i, got.List[i].Size, tt.want.List[i].Size)
				}
				if tt.want.List[i].Size == 0 && got.List[i].Size <= 0 {
					t.Errorf("DumpList.Execute()[%d].Size = %v, want > 0", i, got.List[i].Size)
				}
			}
		})
	}
}
