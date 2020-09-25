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
				urlString: `/v0.1/workspaces/?cursor=&limit=20`,
			},
			want:    []string{"v0.1", "workspaces"},
			want1:   url.Values{"cursor": {""}, "limit": {"20"}},
			wantErr: false,
		},
		{
			args: args{
				urlString: `/v0.1/workspaces/create/`,
			},
			want:    []string{"v0.1", "workspaces", "create"},
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				urlString: `/v0.1/workspaces/2d656454-8191-4193-94e1-4f27ed37d6b0`,
			},
			want:    []string{"v0.1", "workspaces", "2d656454-8191-4193-94e1-4f27ed37d6b0"},
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				urlString: `/v0.1/workspaces/chdcreate/`,
			},
			want:    []string{"v0.1", "workspaces", "chdcreate"},
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				urlString: `/v0.1/workspaces/topworkspaces`,
			},
			want:    []string{"v0.1", "workspaces", "topworkspaces"},
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				urlString: `/v0.1/workspaces/2d656454-8191-4193-94e1-4f27ed37d6b0/chdn`,
			},
			want:    []string{"v0.1", "workspaces", "2d656454-8191-4193-94e1-4f27ed37d6b0", "chdn"},
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				urlString: `/v0.1/workspaces/b8399c93-3ebd-48ed-8e20-067ab81067eb/getparent`,
			},
			want:    []string{"v0.1", "workspaces", "b8399c93-3ebd-48ed-8e20-067ab81067eb", "getparent"},
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				urlString: `/v0.1/channels/create/`,
			},
			want:    []string{"v0.1", "channels", "create"},
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				urlString: `/v0.1/channels/69280d16-02d5-4908-b59a-923b4e86f7bd`,
			},
			want:    []string{"v0.1", "channels", "69280d16-02d5-4908-b59a-923b4e86f7bd"},
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				urlString: `/v0.1/channels/channelbyname`,
			},
			want:    []string{"v0.1", "channels", "channelbyname"},
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				urlString: `/v0.1/messages/create`,
			},
			want:    []string{"v0.1", "messages", "create"},
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				urlString: `/v0.1/messages/3ea3c22d-2949-4afa-aae4-2bac1e15359b`,
			},
			want:    []string{"v0.1", "messages", "3ea3c22d-2949-4afa-aae4-2bac1e15359b"},
			want1:   url.Values{},
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
				s: `/v0.1/workspaces/?cursor=&limit=20`,
			},
			want:    "/v0.1/workspaces/",
			want1:   url.Values{"cursor": {""}, "limit": {"20"}},
			wantErr: false,
		},
		{
			args: args{
				s: `/v0.1/workspaces/create/`,
			},
			want:    "/v0.1/workspaces/create/",
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				s: `/v0.1/workspaces/2d656454-8191-4193-94e1-4f27ed37d6b0`,
			},
			want:    "/v0.1/workspaces/2d656454-8191-4193-94e1-4f27ed37d6b0",
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				s: `/v0.1/workspaces/chdcreate/`,
			},
			want:    "/v0.1/workspaces/chdcreate/",
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				s: `/v0.1/workspaces/topworkspaces`,
			},
			want:    "/v0.1/workspaces/topworkspaces",
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				s: `/v0.1/workspaces/2d656454-8191-4193-94e1-4f27ed37d6b0/chdn`,
			},
			want:    "/v0.1/workspaces/2d656454-8191-4193-94e1-4f27ed37d6b0/chdn",
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				s: `/v0.1/workspaces/b8399c93-3ebd-48ed-8e20-067ab81067eb/getparent`,
			},
			want:    "/v0.1/workspaces/b8399c93-3ebd-48ed-8e20-067ab81067eb/getparent",
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				s: `/v0.1/channels/create/`,
			},
			want:    "/v0.1/channels/create/",
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				s: `/v0.1/channels/69280d16-02d5-4908-b59a-923b4e86f7bd`,
			},
			want:    "/v0.1/channels/69280d16-02d5-4908-b59a-923b4e86f7bd",
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				s: `/v0.1/channels/channelbyname`,
			},
			want:    "/v0.1/channels/channelbyname",
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				s: `/v0.1/messages/create/`,
			},
			want:    "/v0.1/messages/create/",
			want1:   url.Values{},
			wantErr: false,
		},
		{
			args: args{
				s: `/v0.1/messages/3ea3c22d-2949-4afa-aae4-2bac1e15359b`,
			},
			want:    "/v0.1/messages/3ea3c22d-2949-4afa-aae4-2bac1e15359b",
			want1:   url.Values{},
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
				url: "/v0.1/workspaces/",
			},
			want: []string{"v0.1", "workspaces"},
		},
		{
			args: args{
				url: "/v0.1/workspaces/create/",
			},
			want: []string{"v0.1", "workspaces", "create"},
		},
		{
			args: args{
				url: "/v0.1/workspaces/2d656454-8191-4193-94e1-4f27ed37d6b0",
			},
			want: []string{"v0.1", "workspaces", "2d656454-8191-4193-94e1-4f27ed37d6b0"},
		},
		{
			args: args{
				url: "/v0.1/workspaces/chdcreate/",
			},
			want: []string{"v0.1", "workspaces", "chdcreate"},
		},
		{
			args: args{
				url: "/v0.1/workspaces/topworkspaces",
			},
			want: []string{"v0.1", "workspaces", "topworkspaces"},
		},
		{
			args: args{
				url: "/v0.1/workspaces/2d656454-8191-4193-94e1-4f27ed37d6b0/chdn",
			},
			want: []string{"v0.1", "workspaces", "2d656454-8191-4193-94e1-4f27ed37d6b0", "chdn"},
		},
		{
			args: args{
				url: "/v0.1/workspaces/b8399c93-3ebd-48ed-8e20-067ab81067eb/getparent",
			},
			want: []string{"v0.1", "workspaces", "b8399c93-3ebd-48ed-8e20-067ab81067eb", "getparent"},
		},
		{
			args: args{
				url: "/v0.1/channels/create/",
			},
			want: []string{"v0.1", "channels", "create"},
		},
		{
			args: args{
				url: "/v0.1/channels/69280d16-02d5-4908-b59a-923b4e86f7bd",
			},
			want: []string{"v0.1", "channels", "69280d16-02d5-4908-b59a-923b4e86f7bd"},
		},
		{
			args: args{
				url: "/v0.1/channels/channelbyname/",
			},
			want: []string{"v0.1", "channels", "channelbyname"},
		},
		{
			args: args{
				url: "/v0.1/messages/create/",
			},
			want: []string{"v0.1", "messages", "create"},
		},
		{
			args: args{
				url: "/v0.1/messages/3ea3c22d-2949-4afa-aae4-2bac1e15359b",
			},
			want: []string{"v0.1", "messages", "3ea3c22d-2949-4afa-aae4-2bac1e15359b"},
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
		{
			args: args{
				b: []byte{114, 81, 23, 127, 250, 209, 71, 164, 172, 235, 18, 212, 215, 144, 119, 79},
			},
			want:    "7251177f-fad1-47a4-aceb-12d4d790774f",
			wantErr: false,
		},
		{
			args: args{
				b: []byte{143, 91, 160, 102, 122, 10, 65, 175, 151, 222, 16, 197, 221, 193, 4, 238},
			},
			want:    "8f5ba066-7a0a-41af-97de-10c5ddc104ee",
			wantErr: false,
		},
		{
			args: args{
				b: []byte{178, 133, 22, 28, 79, 119, 71, 40, 153, 50, 8, 214, 92, 253, 190, 54},
			},
			want:    "b285161c-4f77-4728-9932-08d65cfdbe36",
			wantErr: false,
		},
		{
			args: args{
				b: []byte{153, 205, 76, 238, 199, 55, 75, 78, 139, 222, 58, 150, 6, 113, 109, 35},
			},
			want:    "99cd4cee-c737-4b4e-8bde-3a9606716d23",
			wantErr: false,
		},
		{
			args: args{
				b: []byte{101, 170, 14, 65, 29, 54, 70, 24, 169, 204, 60, 235, 171, 63, 98, 88},
			},
			want:    "65aa0e41-1d36-4618-a9cc-3cebab3f6258",
			wantErr: false,
		},
		{
			args: args{
				b: []byte{93, 12, 71, 195, 33, 215, 69, 222, 138, 29, 107, 245, 127, 125, 90, 217},
			},
			want:    "5d0c47c3-21d7-45de-8a1d-6bf57f7d5ad9",
			wantErr: false,
		},
		{
			args: args{
				b: []byte{96, 226, 159, 189, 71, 72, 64, 51, 159, 139, 202, 201, 205, 82, 133, 7},
			},
			want:    "60e29fbd-4748-4033-9f8b-cac9cd528507",
			wantErr: false,
		},
		{
			args: args{
				b: []byte{234, 23, 183, 114, 216, 139, 75, 32, 156, 164, 237, 184, 207, 207, 200, 15},
			},
			want:    "ea17b772-d88b-4b20-9ca4-edb8cfcfc80f",
			wantErr: false,
		},
		{
			args: args{
				b: []byte{232, 196, 87, 146, 115, 140, 73, 2, 130, 179, 234, 28, 35, 134, 232, 218},
			},
			want:    "e8c45792-738c-4902-82b3-ea1c2386e8da",
			wantErr: false,
		},
		{
			args: args{
				b: []byte{45, 101, 100, 84, 129, 145, 65, 147, 148, 225, 79, 39, 237, 55, 214, 176},
			},
			want:    "2d656454-8191-4193-94e1-4f27ed37d6b0",
			wantErr: false,
		},
		{
			args: args{
				b: []byte{184, 57, 156, 147, 62, 189, 72, 237, 142, 32, 6, 122, 184, 16, 103, 235},
			},
			want:    "b8399c93-3ebd-48ed-8e20-067ab81067eb",
			wantErr: false,
		},
		{
			args: args{
				b: []byte{105, 40, 13, 22, 2, 213, 73, 8, 181, 154, 146, 59, 78, 134, 247, 189},
			},
			want:    "69280d16-02d5-4908-b59a-923b4e86f7bd",
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
		{
			args: args{
				s: "7251177f-fad1-47a4-aceb-12d4d790774f",
			},
			want:    []byte{114, 81, 23, 127, 250, 209, 71, 164, 172, 235, 18, 212, 215, 144, 119, 79},
			wantErr: false,
		},
		{
			args: args{
				s: "8f5ba066-7a0a-41af-97de-10c5ddc104ee",
			},
			want:    []byte{143, 91, 160, 102, 122, 10, 65, 175, 151, 222, 16, 197, 221, 193, 4, 238},
			wantErr: false,
		},
		{
			args: args{
				s: "b285161c-4f77-4728-9932-08d65cfdbe36",
			},
			want:    []byte{178, 133, 22, 28, 79, 119, 71, 40, 153, 50, 8, 214, 92, 253, 190, 54},
			wantErr: false,
		},
		{
			args: args{
				s: "99cd4cee-c737-4b4e-8bde-3a9606716d23",
			},
			want:    []byte{153, 205, 76, 238, 199, 55, 75, 78, 139, 222, 58, 150, 6, 113, 109, 35},
			wantErr: false,
		},
		{
			args: args{
				s: "65aa0e41-1d36-4618-a9cc-3cebab3f6258",
			},
			want:    []byte{101, 170, 14, 65, 29, 54, 70, 24, 169, 204, 60, 235, 171, 63, 98, 88},
			wantErr: false,
		},
		{
			args: args{
				s: "5d0c47c3-21d7-45de-8a1d-6bf57f7d5ad9",
			},
			want:    []byte{93, 12, 71, 195, 33, 215, 69, 222, 138, 29, 107, 245, 127, 125, 90, 217},
			wantErr: false,
		},
		{
			args: args{
				s: "60e29fbd-4748-4033-9f8b-cac9cd528507",
			},
			want:    []byte{96, 226, 159, 189, 71, 72, 64, 51, 159, 139, 202, 201, 205, 82, 133, 7},
			wantErr: false,
		},
		{
			args: args{
				s: "ea17b772-d88b-4b20-9ca4-edb8cfcfc80f",
			},
			want:    []byte{234, 23, 183, 114, 216, 139, 75, 32, 156, 164, 237, 184, 207, 207, 200, 15},
			wantErr: false,
		},
		{
			args: args{
				s: "e8c45792-738c-4902-82b3-ea1c2386e8da",
			},
			want:    []byte{232, 196, 87, 146, 115, 140, 73, 2, 130, 179, 234, 28, 35, 134, 232, 218},
			wantErr: false,
		},
		{
			args: args{
				s: "2d656454-8191-4193-94e1-4f27ed37d6b0",
			},
			want:    []byte{45, 101, 100, 84, 129, 145, 65, 147, 148, 225, 79, 39, 237, 55, 214, 176},
			wantErr: false,
		},
		{
			args: args{
				s: "b8399c93-3ebd-48ed-8e20-067ab81067eb",
			},
			want:    []byte{184, 57, 156, 147, 62, 189, 72, 237, 142, 32, 6, 122, 184, 16, 103, 235},
			wantErr: false,
		},
		{
			args: args{
				s: "69280d16-02d5-4908-b59a-923b4e86f7bd",
			},
			want:    []byte{105, 40, 13, 22, 2, 213, 73, 8, 181, 154, 146, 59, 78, 134, 247, 189},
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
