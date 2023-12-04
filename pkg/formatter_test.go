package pkg

import (
	"testing"

	"github.com/soerenschneider/fetcharr/internal/syncer"
)

func TestFormatPayload(t *testing.T) {
	type args struct {
		templ string
		stats syncer.Stats
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				templ: `{"subject":"test","message":"{{ .NumFiles }}"}`,
				stats: syncer.Stats{
					NumFiles: 5,
				},
			},
			want:    `{"subject":"test","message":"5"}`,
			wantErr: false,
		},
		{
			name: "no templating",
			args: args{
				templ: `{"subject":"test","message":"no idea"}`,
				stats: syncer.Stats{
					NumFiles: 5,
				},
			},
			want:    `{"subject":"test","message":"no idea"}`,
			wantErr: false,
		},
		{
			name: "invalid templating",
			args: args{
				templ: `{"subject":"test","message":"{{ .nonExistent }}"}`,
				stats: syncer.Stats{
					NumFiles: 5,
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Format(tt.args.templ, tt.args.stats)
			if (err != nil) != tt.wantErr {
				t.Errorf("Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Format() got = %v, want %v", got, tt.want)
			}
		})
	}
}
