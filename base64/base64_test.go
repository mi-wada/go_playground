package base64

import (
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	type args struct {
		src []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "1 byte",
			args: args{
				src: []byte{0b100000_00},
			},
			want: []byte{'g', 'A', '=', '='},
		},
		{
			name: "2 bytes",
			args: args{
				src: []byte{0b000000_00, 0b1000_0000},
			},
			want: []byte{'A', 'I', 'A', '='},
		},
		{
			name: "3 bytes",
			args: args{
				src: []byte{0b000000_00, 0b0000_0000, 0b00_000000},
			},
			want: []byte{'A', 'A', 'A', 'A'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Encode(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
