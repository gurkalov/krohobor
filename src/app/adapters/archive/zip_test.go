package archive

import (
	"reflect"
	"testing"
)

func TestNewZip(t *testing.T) {
	type args struct {
		dir      string
		password string
	}
	tests := []struct {
		name string
		args args
		want Zip
	}{
		{
			name: "Test",
			args: args{
				dir: "/tmp",
				password: "123456",
			},
			want: Zip{
				Dir: "/tmp",
				Password: "123456",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewZip(tt.args.dir, tt.args.password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewZip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZip_Archive(t *testing.T) {
	type args struct {
		file string
		dir  string
	}
	tests := []struct {
		name    string
		zip     Zip
		args    args
		wantErr bool
	}{
		{
			name: "Archive",
			zip: NewZipMock("/tmp", "123456", "mock"),
			args: args{
				file: "/tmp/test.zip",
				dir: "/tmp/mock",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.zip.Archive(tt.args.file, tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("Zip.Archive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZip_Unarchive(t *testing.T) {
	type fields struct {
		Dir      string
		Password string
	}
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Zip{
				Dir:      tt.fields.Dir,
				Password: tt.fields.Password,
			}
			got, err := s.Unarchive(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Zip.Unarchive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Zip.Unarchive() = %v, want %v", got, tt.want)
			}
		})
	}
}
