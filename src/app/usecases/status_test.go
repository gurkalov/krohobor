package usecases

import (
	"krohobor/app/adapters/config"
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
	"reflect"
	"testing"
)

func TestNewStatus(t *testing.T) {
	dir := "/tmp/krohobor/storage"

	cfg := config.LoadMock()
	db := database.NewMemory()
	store := storage.NewFileMock(dir, nil)

	type args struct {
		cfg     config.Config
		db      database.Interface
		storage storage.Interface
	}
	tests := []struct {
		name string
		args args
		want *Status
	}{
		{
			name: "Test",
			args: args{
				cfg: cfg,
				db: db,
				storage: store,
			},
			want: &Status{
				cfg: cfg,
				db: db,
				storage: store,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStatus(tt.args.cfg, tt.args.db, tt.args.storage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatus_Execute(t *testing.T) {
	dir := "/tmp/krohobor/storage"

	cfg := config.LoadMock()
	db := database.NewMemory()
	store := storage.NewFileMock(dir, nil)

	type fields struct {
		cfg     config.Config
		db      database.Interface
		storage storage.Interface
	}
	type args struct {
		request StatusRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    StatusResponse
		wantErr bool
	}{
		{
			name: "Status - successful",
			fields: fields{
				cfg: cfg,
				db: db,
				storage: store,
			},
			args: args{
				request: StatusRequest{
					Target: "",
				},
			},
			want: StatusResponse{
				Db: struct {
					Check bool
					Host  string
				}{Check: true, Host: ":"},
				Storage: struct {
					Check bool
					Host  string
				}{Check: true, Host: ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Status{
				cfg:     tt.fields.cfg,
				db:      tt.fields.db,
				storage: tt.fields.storage,
			}
			got, err := s.Execute(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Status.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Status.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
