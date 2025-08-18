package config

import (
	"reflect"
	"testing"

	"github.com/MilindGour/jellyfin-media-renamer/testdata"
)

func TestJmrConfig_ParseFromBytes(t *testing.T) {
	type args struct {
		config []byte
	}
	tests := []struct {
		name string
		j    JmrConfig
		args args
		want *Config
	}{
		{
			name: "Parse mock config",
			j:    JmrConfig{},
			args: args{config: testdata.ConfigJsonMock},
			want: &Config{
				"1.2.3",
				"7749",
				AllowedExtensions{[]string{".srt"}, []string{".mp4"}},
				[]DirConfig{{"name1", "path1"}},
			},
		},
		{
			name: "Nil config",
			j:    JmrConfig{},
			args: args{config: []byte{}},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := JmrConfig{}
			if got := j.ParseFromBytes(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JmrConfig.ParseFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
