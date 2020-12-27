package usecases

import (
	"krohobor/app/adapters/database"
	"reflect"
	"testing"
)

func TestNewDbDelete(t *testing.T) {
	db := database.NewMemory()

	type args struct {
		db database.Interface
	}
	tests := []struct {
		name string
		args args
		want *DbDelete
	}{
		{
			name: "Test",
			args: args{
				db: db,
			},
			want: &DbDelete{
				db: db,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDbDelete(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDbDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbDelete_Execute(t *testing.T) {
	db := database.NewMemory()

	type fields struct {
		db database.Interface
	}
	type args struct {
		request DbDeleteRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    DbDeleteResponse
		wantErr bool
	}{
		{
			name: "Delete - successful",
			fields: fields{
				db: db,
			},
			args: args{
				request: DbDeleteRequest{
					Name: "test1",
				},
			},
			want: DbDeleteResponse{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dl := &DbDelete{
				db: tt.fields.db,
			}
			got, err := dl.Execute(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DbDelete.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbDelete.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
