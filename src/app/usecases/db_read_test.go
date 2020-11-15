package usecases

import (
	"krohobor/app/adapters/database"
	"krohobor/app/domain"
	"reflect"
	"testing"
)

func TestNewDbRead(t *testing.T) {
	db := database.NewMemory()

	type args struct {
		db database.Interface
	}
	tests := []struct {
		name string
		args args
		want *DbRead
	}{
		{
			name: "Test",
			args: args{
				db: db,
			},
			want: &DbRead{
				db: db,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDbRead(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDbRead() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbRead_Execute(t *testing.T) {
	db := database.NewMemory()

	type fields struct {
		db database.Interface
	}
	type args struct {
		request DbReadRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    DbReadResponse
		wantErr bool
	}{
		{
			name: "List - successful",
			fields: fields{
				db: db,
			},
			args: args{
				request: DbReadRequest{
					Name: "test1",
				},
			},
			want: DbReadResponse{
				List: []domain.Table{
					{
						Name: "table1",
						Size: 10000,
						Count: 0,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dl := &DbRead{
				db: tt.fields.db,
			}
			got, err := dl.Execute(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbRead.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbRead.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
