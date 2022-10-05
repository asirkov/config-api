package config

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "Data is empty",
			args: args{
				data: []byte(``),
			},
			want:    &Config{},
			wantErr: false,
		},
		{
			name: "With full data",
			args: args{
				data: []byte(
					`
server:
    port: 8085

database:
    user: "user"
    password: "pass"
    host: "host"
    port: 123
    dBName: "db_name"
`,
				),
			},
			want: &Config{
				Server: ServerConfig{
					Port: 8085,
				},
				Database: DatabaseConfig{
					User:     "user",
					Password: "pass",
					Host:     "host",
					Port:     123,
					DBName:   "db_name",
				},
			},
			wantErr: false,
		},
		{
			name: "Fails when data is not valid",
			args: args{
				data: []byte("Chepuha"),
			},
			want:    &Config{},
			wantErr: true,
		},
		{
			name: "Empty config when data is nil",
			args: args{
				data: nil,
			},
			want:    &Config{},
			wantErr: false,
		},
		{
			name: "When more data in file than needs",
			args: args{
				data: []byte(
					`
smth:
    someValue: 21

gateway:
    event-v1:
        cache-age: 10

server:
    port: 8085

database:
    user: "user"
    password: "pass"
    host: "host"
    port: 123
    dBName: "db_name" 

sport:
    base-url: "https://testDataSport.ls-g.net"
    timeout:
    connect: 321
    read: 987
    fake:
    size: 256
    delay: 20
`,
				),
			},
			want: &Config{
				Server: ServerConfig{
					Port: 8085,
				},
				Database: DatabaseConfig{
					User:     "user",
					Password: "pass",
					Host:     "host",
					Port:     123,
					DBName:   "db_name",
				},
			},
			wantErr: false,
		},
		{
			name: "When less data in file than needs",
			args: args{
				data: []byte(
					`
gateway:
    event-v1:
    cache-age: 10

database:
    user: "user"
    password: "pass"
    host: "host"
    port: 123
    dBName: "db_name"
`,
				),
			},
			want: &Config{
				Database: DatabaseConfig{
					User:     "user",
					Password: "pass",
					Host:     "host",
					Port:     123,
					DBName:   "db_name",
				},
			},
			wantErr: false,
		},
		{
			name: "When not valid yaml",
			args: args{
				data: []byte(
					`
gateway
event-v1
cache-age 10
`,
				),
			},
			want:    &Config{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConfig(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
