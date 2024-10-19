package json

import (
	"reflect"
	"testing"
)

func Test_tokenize(t *testing.T) {
	type args struct {
		src []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []token
		wantErr bool
	}{
		{
			name:    "string",
			args:    args{src: []byte(`"hello"`)},
			want:    []token{newStringToken("hello")},
			wantErr: false,
		},
		{
			name:    "blank string",
			args:    args{src: []byte(`""`)},
			want:    []token{newStringToken("")},
			wantErr: false,
		},
		{
			name:    "multi-bytes string",
			args:    args{src: []byte(`"こんにちは👋"`)},
			want:    []token{newStringToken("こんにちは👋")},
			wantErr: false,
		},
		{
			name:    "number",
			args:    args{src: []byte(`1.23`)},
			want:    []token{newNumberToken(1.23)},
			wantErr: false,
		},
		{
			name:    "float number starting with 0",
			args:    args{src: []byte(`0.12`)},
			want:    []token{newNumberToken(0.12)},
			wantErr: false,
		},
		{
			name:    "float number starting with double zero",
			args:    args{src: []byte(`00.1`)},
			want:    []token{newNumberToken(0.1)},
			wantErr: false,
		},
		{
			name:    "int number",
			args:    args{src: []byte(`123`)},
			want:    []token{newNumberToken(123)},
			wantErr: false,
		},
		{
			name:    "wrong number, finishing with dot",
			args:    args{src: []byte(`0.`)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "wrong string, missing last double quote",
			args:    args{src: []byte(`"hello`)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "null",
			args:    args{src: []byte("null")},
			want:    []token{newNullToken()},
			wantErr: false,
		},
		{
			name:    "wrong null, starting from n but wrong",
			args:    args{src: []byte("neko")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "wrong null, starting from n but short",
			args:    args{src: []byte("nul")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "true",
			args:    args{src: []byte("true")},
			want:    []token{newTrueToken()},
			wantErr: false,
		},
		{
			name:    "wrong true, starting from t, but wrong",
			args:    args{src: []byte("trua")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "wrong true, starting from t, but short",
			args:    args{src: []byte("tru")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "false",
			args:    args{src: []byte("false")},
			want:    []token{newFalseToken()},
			wantErr: false,
		},
		{
			name:    "wrong false, starting from f, but wrong",
			args:    args{src: []byte("falso")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "wrong false, starting from f, but short",
			args:    args{src: []byte("fals")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "[",
			args:    args{src: []byte("[")},
			want:    []token{newLBracketToken()},
			wantErr: false,
		},
		{
			name:    "]",
			args:    args{src: []byte("]")},
			want:    []token{newRBracketToken()},
			wantErr: false,
		},
		{
			name:    "{",
			args:    args{src: []byte("{")},
			want:    []token{newLBraceToken()},
			wantErr: false,
		},
		{
			name:    "}",
			args:    args{src: []byte("}")},
			want:    []token{newRBraceToken()},
			wantErr: false,
		},
		{
			name:    ":",
			args:    args{src: []byte(":")},
			want:    []token{newColonToken()},
			wantErr: false,
		},
		{
			name:    ",",
			args:    args{src: []byte(",")},
			want:    []token{newCommaToken()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokenize(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}
