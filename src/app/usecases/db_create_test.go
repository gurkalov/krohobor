package usecases

import (
	"krohobor/app/adapters/database"
	"reflect"
	"testing"
)

func TestNewDbCreate(t *testing.T) {
	db := database.NewMemory()

	type args struct {
		db database.Interface
	}
	tests := []struct {
		name string
		args args
		want *DbCreate
	}{
		{
			name: "Test",
			args: args{
				db: db,
			},
			want: &DbCreate{
				db: db,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDbCreate(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDbCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbCreate_Execute(t *testing.T) {
	db := database.NewMemory()

	type fields struct {
		db database.Interface
	}
	type args struct {
		request DbCreateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    DbCreateResponse
		wantErr bool
	}{
		{
			name: "Create - successful",
			fields: fields{
				db: db,
			},
			args: args{
				request: DbCreateRequest{
					Name: "test5",
				},
			},
			want: DbCreateResponse{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dl := &DbCreate{
				db: tt.fields.db,
			}
			got, err := dl.Execute(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbCreate.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbCreate.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
