package archive

import (
	"reflect"
	"testing"
)

const (
	tmpDir = "/tmp/krohobor/archive"
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
				dir:      "/tmp",
				password: "123456",
			},
			want: Zip{
				Dir:      "/tmp",
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

func TestZip_Check(t *testing.T) {
	type fields struct {
		Dir      string
		Password string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Check",
			fields: fields{
				Dir: tmpDir,
				Password: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Zip{
				Dir:      tt.fields.Dir,
				Password: tt.fields.Password,
			}
			if err := s.Check(); (err != nil) != tt.wantErr {
				t.Errorf("Zip.Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZip_Archive(t *testing.T) {
	NewZipMock(tmpDir, "123456")

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
			name: "Archive with password - successful",
			zip:  NewZip(tmpDir, "123456"),
			args: args{
				file: tmpDir + "/test-with-password.zip",
				dir:  tmpDir + "/mock",
			},
			wantErr: false,
		},
		{
			name: "Archive without password - successful",
			zip:  NewZip(tmpDir, ""),
			args: args{
				file: tmpDir + "/test-without-password.zip",
				dir:  tmpDir + "/mock",
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
	NewZipMock(tmpDir, "123456")

	type args struct {
		file string
	}
	tests := []struct {
		name    string
		zip     Zip
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Unarchive with password - successful",
			zip:  NewZip(tmpDir, "123456"),
			args: args{
				file: tmpDir + "/test-mock-with-password.zip",
			},
			want:    tmpDir + "/file2",
			wantErr: false,
		},
		{
			name: "Unarchive with password - empty password - error",
			zip:  NewZip(tmpDir, ""),
			args: args{
				file: tmpDir + "/test-mock-with-password.zip",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Unarchive with password - wrong password - error",
			zip:  NewZip(tmpDir, "1234"),
			args: args{
				file: tmpDir + "/test-mock-with-password.zip",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Unarchive without password - any password - successful",
			zip:  NewZip(tmpDir, "111111"),
			args: args{
				file: tmpDir + "/test-mock-without-password.zip",
			},
			want:    tmpDir + "/file2",
			wantErr: false,
		},
		{
			name: "Unarchive without password - empty password - successful",
			zip:  NewZip(tmpDir, ""),
			args: args{
				file: tmpDir + "/test-mock-without-password.zip",
			},
			want:    tmpDir + "/file2",
			wantErr: false,
		},
		{
			name: "Unarchive empty folder - error",
			zip:  NewZip(tmpDir, "123456"),
			args: args{
				file: tmpDir + "/test-mock-empty-folder.zip",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.zip.Unarchive(tt.args.file)
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

