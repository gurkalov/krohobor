package usecases

import (
	"krohobor/app/adapters/database"
	"krohobor/app/domain"
	"reflect"
	"testing"
)

func TestNewDbList(t *testing.T) {
	db := database.NewMemory()

	type args struct {
		db database.Interface
	}
	tests := []struct {
		name string
		args args
		want *DbList
	}{
		{
			name: "Test",
			args: args{
				db: db,
			},
			want: &DbList{
				db: db,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDbList(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDbList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbList_Execute(t *testing.T) {
	db := database.NewMemory()

	type fields struct {
		db database.Interface
	}
	type args struct {
		request DbListRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    DbListResponse
		wantErr bool
	}{
		{
			name: "List - successful",
			fields: fields{
				db: db,
			},
			args: args{
				request: DbListRequest{},
			},
			want: DbListResponse{
				List: []domain.Database{
					{
						Name: "test1",
						Size: 100000,
					},
					{
						Name: "test2",
						Size: 100000,
					},
					{
						Name: "test3",
						Size: 100000,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dl := &DbList{
				db: tt.fields.db,
			}
			got, err := dl.Execute(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbList.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbList.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
