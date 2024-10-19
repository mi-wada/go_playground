package base64

import (
	"reflect"
	"testing"
)

func TestNewEncoding(t *testing.T) {
	tests := []struct {
		name        string
		encoder     string
		shouldPanic bool
	}{
		{
			name:        "Valid 64-byte encoder",
			encoder:     "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/",
			shouldPanic: false,
		},
		{
			name:        "Invalid encoder with less than 64 bytes",
			encoder:     "ABC",
			shouldPanic: true,
		},
		{
			name:        "Invalid encoder with more than 64 bytes",
			encoder:     "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.shouldPanic {
						t.Errorf("NewEncoding() panicked unexpectedly")
					}
				} else {
					if tt.shouldPanic {
						t.Errorf("NewEncoding() did not panic as expected")
					}
				}
			}()
			NewEncoding(tt.encoder)
		})
	}
}

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
			encoding := StdEncoding

			if got := encoding.Encode(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
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
				src: []byte{'g', 'A', '=', '='},
			},
			want: []byte{0b100000_00},
		},
		{
			name: "2 byte",
			args: args{
				src: []byte{'A', 'I', 'A', '='},
			},
			want: []byte{0b000000_00, 0b1000_0000},
		},
		{
			name: "3 byte",
			args: args{
				src: []byte{'A', 'A', 'A', 'A'},
			},
			want: []byte{0b000000_00, 0b0000_0000, 0b00_000000},
		},
	}
	for _, tt := range tests {
		encoding := StdEncoding
		t.Run(tt.name, func(t *testing.T) {
			if got := encoding.Decode(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}
