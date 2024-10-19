package json

import (
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	type args struct {
		src []byte
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name:    "string",
			args:    args{src: []byte(`"hello"`)},
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "wrong string, extra",
			args:    args{src: []byte(`"hello""hello"`)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "number",
			args:    args{src: []byte("1.23")},
			want:    1.23,
			wantErr: false,
		},
		{
			name:    "null",
			args:    args{src: []byte("null")},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "true",
			args:    args{src: []byte("true")},
			want:    true,
			wantErr: false,
		},
		{
			name:    "false",
			args:    args{src: []byte("false")},
			want:    false,
			wantErr: false,
		},
		{
			name:    "array",
			args:    args{src: []byte(`["hello", 1.23, null, true, false, []]`)},
			want:    []any{"hello", 1.23, nil, true, false, []any{}},
			wantErr: false,
		},
		{
			name:    "empty array",
			args:    args{src: []byte(`[]`)},
			want:    []any{},
			wantErr: false,
		},
		{
			name:    "wrong array, many [",
			args:    args{src: []byte(`[[]`)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "wrong array, many ]",
			args:    args{src: []byte(`[]]`)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "wrong array, trailing comma",
			args:    args{src: []byte(`[null,]`)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "wrong array, many `,`",
			args:    args{src: []byte(`[null,,null]`)},
			want:    nil,
			wantErr: true,
		},
		{
			name: "object",
			args: args{src: []byte(`{
				"str": "hello",
				"num": 1.23,
				"true": true,
				"false": false,
				"array": [[null], [true]],
				"object": {"1": {"str": "hello"}, "2": {"num": 1.23}}
			}`)},
			want: map[string]any{
				"str":   "hello",
				"num":   1.23,
				"true":  true,
				"false": false,
				"array": []any{[]any{nil}, []any{true}},
				"object": map[string]any{
					"1": map[string]any{
						"str": "hello",
					},
					"2": map[string]any{
						"num": 1.23,
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "empty object",
			args:    args{src: []byte(`{}`)},
			want:    map[string]any{},
			wantErr: false,
		},
		{
			name:    "wrong object, no string key",
			args:    args{src: []byte(`{1: "value",}`)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "wrong object, no colon",
			args:    args{src: []byte(`{"key" "value",}`)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "wrong object, many {",
			args:    args{src: []byte(`{{}`)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "wrong object, many }",
			args:    args{src: []byte(`{}}`)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "wrong object, triling comma",
			args:    args{src: []byte(`{"key": "value",}`)},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}
