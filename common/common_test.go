package common

import (
	"net/url"
	"reflect"
	"testing"
)

func TestParseURL(t *testing.T) {
	type args struct {
		urlString string
	}
	tests := []struct {
		args    args
		want    []string
		want1   url.Values
		wantErr bool
	}{
		{
			args: args{
				urlString: `/v0.1/categories/?cursor=&limit=20`,
			},
			want:    []string{"v0.1", "categories"},
			want1:   url.Values{"cursor": {""}, "limit": {"20"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, got1, err := ParseURL(tt.args.urlString)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseURL() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ParseURL() got = %v, want %v", got, tt.want)
		}
		if !reflect.DeepEqual(got1, tt.want1) {
			t.Errorf("ParseURL() got1 = %v, want %v", got1, tt.want1)
		}
	}
}

func TestGetPathQueryString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		args    args
		want    string
		want1   url.Values
		wantErr bool
	}{
		{
			args: args{
				s: `/v0.1/categories/?cursor=&limit=20`,
			},
			want:    "/v0.1/categories/",
			want1:   url.Values{"cursor": {""}, "limit": {"20"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, got1, err := GetPathQueryString(tt.args.s)
		if (err != nil) != tt.wantErr {
			t.Errorf("GetPathQueryString() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if got != tt.want {
			t.Errorf("GetPathQueryString() got = %v, want %v", got, tt.want)
		}
		if !reflect.DeepEqual(got1, tt.want1) {
			t.Errorf("GetPathQueryString() got1 = %v, want %v", got1, tt.want1)
		}
	}
}

func TestGetPathParts(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		args args
		want []string
	}{
		{
			args: args{
				url: "/v0.1/categories/",
			},
			want: []string{"v0.1", "categories"},
		},
	}
	for _, tt := range tests {
		if got := GetPathParts(tt.args.url); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("GetPathParts() = %v, want %v", got, tt.want)
		}
	}
}

func TestUUIDBytesToStr(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		args    args
		want    string
		wantErr bool
	}{
		{
			args: args{
				b: []byte{86, 44, 59, 245, 98, 212, 79, 218, 177, 17, 162, 90, 216, 58, 114, 239},
			},
			want:    "562c3bf5-62d4-4fda-b111-a25ad83a72ef",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := UUIDBytesToStr(tt.args.b)
		if (err != nil) != tt.wantErr {
			t.Errorf("UUIDBytesToStr() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if got != tt.want {
			t.Errorf("UUIDBytesToStr() = %v, want %v", got, tt.want)
		}
	}
}

func TestUUIDStrToBytes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		args    args
		want    []byte
		wantErr bool
	}{
		{
			args: args{
				s: "562c3bf5-62d4-4fda-b111-a25ad83a72ef",
			},
			want:    []byte{86, 44, 59, 245, 98, 212, 79, 218, 177, 17, 162, 90, 216, 58, 114, 239},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, err := UUIDStrToBytes(tt.args.s)
		if (err != nil) != tt.wantErr {
			t.Errorf("UUIDStrToBytes() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("UUIDStrToBytes() = %v, want %v", got, tt.want)
		}
	}
}

func TestEncodeCursor(t *testing.T) {
	type args struct {
		cursor uint
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				cursor: 385,
			},
			want: "Mzg1",
		},
	}
	for _, tt := range tests {
		if got := EncodeCursor(tt.args.cursor); got != tt.want {
			t.Errorf("EncodeCursor() = %v, want %v", got, tt.want)
		}
	}
}

func TestDecodeCursor(t *testing.T) {
	type args struct {
		cursor string
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				cursor: "Mzg1",
			},
			want: "385",
		},
	}
	for _, tt := range tests {
		if got := DecodeCursor(tt.args.cursor); got != tt.want {
			t.Errorf("DecodeCursor() = %v, want %v", got, tt.want)
		}
	}
}
