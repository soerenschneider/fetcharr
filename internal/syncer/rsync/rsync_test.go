package rsync

import (
	"reflect"
	"testing"

	"github.com/soerenschneider/fetcharr/internal/syncer"
)

func Test_parseRsyncOutput(t *testing.T) {
	type args struct {
		output string
	}
	tests := []struct {
		name string
		args args
		want syncer.Stats
	}{
		{
			name: "linux, regular",
			args: args{
				output: `
Number of files: 4 (reg: 1, dir: 3)
Number of created files: 1 (reg: 1)
Number of deleted files: 0
Number of regular files transferred: 1
Total file size: 104,857,600 bytes
Total transferred file size: 104,857,600 bytes
Literal data: 104,857,600 bytes
Matched data: 0 bytes
File list size: 112
File list generation time: 0.001 seconds
File list transfer time: 0.000 seconds
Total bytes sent: 67
Total bytes received: 104,883,364

sent 67 bytes  received 104,883,364 bytes  9,988,898.19 bytes/sec
total size is 104,857,600  speedup is 1.00`,
			},
			want: syncer.Stats{
				NumFiles:                    4,
				NumCreatedFiles:             1,
				NumDeletedFiles:             0,
				NumTransferredFiles:         1,
				TotalFileSize:               104857600,
				TotalBytesSent:              67,
				TotalBytesReceived:          104883364,
				TotalFileSizeHumanized:      "100.00 MB",
				TotalBytesSentHumanized:     "67 B",
				TotalBytesReceivedHumanized: "100.02 MB",
			},
		},
		{
			name: "linux - empty",
			args: args{
				output: `
Number of files: 3 (dir: 3)
Number of created files: 0
Number of deleted files: 0
Number of regular files transferred: 0
Total file size: 0 bytes
Total transferred file size: 0 bytes
Literal data: 0 bytes
Matched data: 0 bytes
File list size: 81
File list generation time: 0.001 seconds
File list transfer time: 0.000 seconds
Total bytes sent: 36
Total bytes received: 94

sent 36 bytes  received 94 bytes  86.67 bytes/sec
total size is 0  speedup is 0.00`,
			},
			want: syncer.Stats{
				NumFiles:                    3,
				NumCreatedFiles:             0,
				NumDeletedFiles:             0,
				NumTransferredFiles:         0,
				TotalFileSize:               0,
				TotalBytesSent:              36,
				TotalBytesReceived:          94,
				TotalFileSizeHumanized:      "0 B",
				TotalBytesSentHumanized:     "36 B",
				TotalBytesReceivedHumanized: "94 B",
			},
		},
		{
			name: "linux",
			args: args{
				output: `Number of files: 5 (reg: 2, dir: 3)
Number of created files: 2 (reg: 2)
Number of deleted files: 4
Number of regular files transferred: 2
Total file size: 280,527 bytes
Total transferred file size: 280,527 bytes
Literal data: 280,527 bytes
Matched data: 0 bytes
File list size: 140
File list generation time: 0.001 seconds
File list transfer time: 0.000 seconds
Total bytes sent: 94
Total bytes received: 280,822

sent 94 bytes  received 280,822 bytes  112,366.40 bytes/sec
total size is 280,527  speedup is 1.00`,
			},
			want: syncer.Stats{
				NumFiles:                    5,
				NumCreatedFiles:             2,
				NumDeletedFiles:             4,
				NumTransferredFiles:         2,
				TotalFileSize:               280527,
				TotalBytesSent:              94,
				TotalBytesReceived:          280822,
				TotalFileSizeHumanized:      "273.95 KB",
				TotalBytesSentHumanized:     "94 B",
				TotalBytesReceivedHumanized: "274.24 KB",
			},
		},
		{
			name: "mac osx",
			args: args{output: `downloaded=/path/to/file1.txt
downloaded=file2.txt
Number of files: 4
Number of files transferred: 2
Total file size: 410032 bytes
Total transferred file size: 410032 bytes
Literal data: 410032 bytes
Matched data: 0 bytes
File list size: 81
File list generation time: 0.001 seconds
File list transfer time: 0.000 seconds
Total bytes sent: 66
Total bytes received: 410323

sent 66 bytes  received 410323 bytes  273592.67 bytes/sec
total size is 410032  speedup is 1.00`},
			want: syncer.Stats{
				Files: []string{
					"file1.txt",
					"file2.txt",
				},
				NumFiles:                    4,
				NumCreatedFiles:             0,
				NumDeletedFiles:             0,
				NumTransferredFiles:         2,
				TotalFileSize:               410032,
				TotalBytesSent:              66,
				TotalBytesReceived:          410323,
				TotalFileSizeHumanized:      "400.42 KB",
				TotalBytesSentHumanized:     "66 B",
				TotalBytesReceivedHumanized: "400.71 KB",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseRsyncOutput(tt.args.output); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseRsyncOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}
