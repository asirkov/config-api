package file

import (
	"reflect"
	"testing"
)

func TestReadFile(t *testing.T) {
	type args struct {
		name string
	}

	tests := []struct {
		name     string
		args     args
		wantData []byte
		wantErr  bool
	}{
		{
			name: "When name of file is empty",
			args: args{
				name: "",
			},
			wantData: nil,
			wantErr:  true,
		},
		{
			name: "Works ok",
			args: args{
				name: "test_data.txt",
			},
			wantData: []byte(
				`Some mock data Ushastiy plug`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := ReadFile(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("ReadFile() = '%v', want '%v'", string(gotData), string(tt.wantData))
			}
		})
	}
}
