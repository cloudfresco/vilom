package common

import (
	"reflect"
	"testing"
	"time"
)

func TestGetSelectorForPasswdRecoveryToken(t *testing.T) {
	type args struct {
		token     string
		requestID string
	}
	tests := []struct {
		args    args
		want    [64]byte
		want1   string
		wantErr bool
	}{
		{
			args: args{
				token:     "Lmh-sldIxvy5jGV-i0AY4dHJZR46Hw3nhpCKLJg6tk-kyvAkVR4epn-MAaCpJ15kX1aW9N_Kv8NrEjhC1c0MnQ==",
				requestID: "bk4tsqg91jatm09q91i0",
			},
			want:    [64]byte{195, 72, 156, 85, 156, 44, 188, 137, 51, 208, 1, 211, 70, 208, 231, 128, 175, 49, 173, 227, 150, 7, 232, 206, 218, 156, 165, 133, 58, 23, 115, 86, 117, 14, 85, 44, 3, 228, 39, 69, 253, 37, 110, 66, 40, 69, 211, 118, 154, 22, 35, 151, 176, 204, 179, 179, 42, 209, 166, 165, 223, 171, 8, 242},
			want1:   "h6+V8h0qeLIPNqpBgY9VKm6ARbd05Q6NSfCZoXooxf/mPx1mTS2Fsruok00OofvyIw3xirnZILiZiRDHqK//Gw==",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got, got1, err := GetSelectorForPasswdRecoveryToken(tt.args.token, tt.args.requestID)
		if (err != nil) != tt.wantErr {
			t.Errorf("GetSelectorForPasswdRecoveryToken() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("GetSelectorForPasswdRecoveryToken() got = %v, want %v", got, tt.want)
		}
		if got1 != tt.want1 {
			t.Errorf("GetSelectorForPasswdRecoveryToken() got1 = %v, want %v", got1, tt.want1)
		}
	}
}

func TestValidatePasswdRecoveryToken(t *testing.T) {
	type args struct {
		verifierBytes [64]byte
		verifier      string
		tokenExpiry   time.Time
		requestID     string
	}
	tests := []struct {
		args    args
		wantErr bool
	}{
		{
			args: args{
				verifierBytes: [64]byte{195, 72, 156, 85, 156, 44, 188, 137, 51, 208, 1, 211, 70, 208, 231, 128, 175, 49, 173, 227, 150, 7, 232, 206, 218, 156, 165, 133, 58, 23, 115, 86, 117, 14, 85, 44, 3, 228, 39, 69, 253, 37, 110, 66, 40, 69, 211, 118, 154, 22, 35, 151, 176, 204, 179, 179, 42, 209, 166, 165, 223, 171, 8, 242},
				verifier:      "w0icVZwsvIkz0AHTRtDngK8xreOWB+jO2pylhToXc1Z1DlUsA+QnRf0lbkIoRdN2mhYjl7DMs7Mq0aal36sI8g==",
				tokenExpiry:   time.Now().UTC().Truncate(time.Second),
				requestID:     "bk4tsqg91jatm09q91i0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if err := ValidatePasswdRecoveryToken(tt.args.verifierBytes, tt.args.verifier, tt.args.tokenExpiry, tt.args.requestID); (err != nil) != tt.wantErr {
			t.Errorf("ValidatePasswdRecoveryToken() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}
