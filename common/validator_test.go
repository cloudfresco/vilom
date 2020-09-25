package common

import (
	"testing"
	"time"
)

func TestValidator_IsGreaterThan(t *testing.T) {
	validator := NewValidator()
	type args struct {
		fieldName  string
		fieldValue int
		high       int
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `Price`,
				fieldValue: int(20),
				high:       int(25),
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Price`,
				fieldValue: int(30),
				high:       int(25),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		if got := tt.v.IsGreaterThan(tt.args.fieldName, tt.args.fieldValue, tt.args.high); got != tt.want {
			t.Errorf("Validator.IsGreaterThan() = %v, want %v", got, tt.want)
		}
	}
}

func TestValidator_IsTimeBefore(t *testing.T) {
	validator := NewValidator()
	mtime := time.Now().UTC()
	myDate1 := mtime.AddDate(0, 0, -14)
	myDate2 := mtime.AddDate(0, 0, 14)
	type args struct {
		fieldName  string
		fieldValue time.Time
		max        time.Time
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `Destination Time`,
				fieldValue: myDate1,
				max:        mtime,
			},
			want: true,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Destination Time`,
				fieldValue: myDate2,
				max:        mtime,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		if got := tt.v.IsTimeBefore(tt.args.fieldName, tt.args.fieldValue, tt.args.max); got != tt.want {
			t.Errorf("Validator.IsTimeBefore() = %v, want %v", got, tt.want)
		}
	}
}

func TestValidator_IsStrNotEmpty(t *testing.T) {
	validator := NewValidator()
	type args struct {
		fieldName  string
		fieldValue string
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `Workspace Name`,
				fieldValue: "",
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Workspace Name`,
				fieldValue: "Output Direct Component",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		if got := tt.v.IsStrNotEmpty(tt.args.fieldName, tt.args.fieldValue); got != tt.want {
			t.Errorf("Validator.IsStrNotEmpty() = %v, want %v", got, tt.want)
		}
	}
}

func TestValidator_IsInt64Negative(t *testing.T) {
	validator := NewValidator()
	type args struct {
		fieldName  string
		fieldValue int64
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `Price`,
				fieldValue: int64(204),
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Price`,
				fieldValue: int64(-204),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		if got := tt.v.IsInt64Negative(tt.args.fieldName, tt.args.fieldValue); got != tt.want {
			t.Errorf("Validator.IsInt64Negative() = %v, want %v", got, tt.want)
		}
	}
}

func TestValidator_IsInt64NonNegative(t *testing.T) {
	validator := NewValidator()
	type args struct {
		fieldName  string
		fieldValue int64
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `Price`,
				fieldValue: int64(-204),
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Price`,
				fieldValue: int64(204),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		if got := tt.v.IsInt64NonNegative(tt.args.fieldName, tt.args.fieldValue); got != tt.want {
			t.Errorf("Validator.IsInt64NonNegative() = %v, want %v", got, tt.want)
		}
	}
}

func TestValidator_IsInt64Positive(t *testing.T) {
	validator := NewValidator()
	type args struct {
		fieldName  string
		fieldValue int64
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `Price`,
				fieldValue: int64(-204),
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Price`,
				fieldValue: int64(204),
			},
			want: true,
		},
	}
	for _, tt := range tests {

		if got := tt.v.IsInt64Positive(tt.args.fieldName, tt.args.fieldValue); got != tt.want {
			t.Errorf("Validator.IsInt64Positive() = %v, want %v", got, tt.want)
		}

	}
}

func TestValidator_IsInt64NonPositive(t *testing.T) {
	validator := NewValidator()
	type args struct {
		fieldName  string
		fieldValue int64
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `Price`,
				fieldValue: int64(204),
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Price`,
				fieldValue: int64(-204),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		if got := tt.v.IsInt64NonPositive(tt.args.fieldName, tt.args.fieldValue); got != tt.want {
			t.Errorf("Validator.IsInt64NonPositive() = %v, want %v", got, tt.want)
		}
	}
}

func TestValidator_IsEmail(t *testing.T) {
	validator := NewValidator()
	type args struct {
		fieldName  string
		fieldValue string
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `Email`,
				fieldValue: "efev$?sdecd§/az@gmail.com",
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Email`,
				fieldValue: "abcd@gmail.com",
			},
			want: true,
		},
	}
	for _, tt := range tests {

		if got := tt.v.IsEmail(tt.args.fieldName, tt.args.fieldValue); got != tt.want {
			t.Errorf("Validator.IsEmail() = %v, want %v", got, tt.want)
		}

	}
}

func TestValidator_IsDateFormat(t *testing.T) {
	validator := NewValidator()
	type args struct {
		fieldName  string
		fieldValue string
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `Start Date`,
				fieldValue: "12/13/2019",
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Start Date`,
				fieldValue: "13/12/2019",
			},
			want: true,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Start Date`,
				fieldValue: "13/12/20df",
			},
			want: false,
		},
	}
	for _, tt := range tests {

		if got := tt.v.IsDateFormat(tt.args.fieldName, tt.args.fieldValue); got != tt.want {
			t.Errorf("Validator.IsDateFormat() = %v, want %v", got, tt.want)
		}

	}
}

func TestValidator_IsPhoneNumber(t *testing.T) {
	validator := NewValidator()
	type args struct {
		fieldName  string
		fieldValue string
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `Phone Number`,
				fieldValue: "90191919908ee",
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Phone Number`,
				fieldValue: "1 (234) 568-9871",
			},
			want: true,
		},
	}
	for _, tt := range tests {

		if got := tt.v.IsPhoneNumber(tt.args.fieldName, tt.args.fieldValue); got != tt.want {
			t.Errorf("Validator.IsPhoneNumber() = %v, want %v", got, tt.want)
		}
	}
}

func TestValidator_IsUUID4(t *testing.T) {
	validator := NewValidator()
	type args struct {
		fieldName  string
		fieldValue string
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `UUID`,
				fieldValue: "b987fbc9-4bed-5078-af07-9141ba07c9f3",
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `UUID`,
				fieldValue: "527f67e3-89f5-70v7-73y2-b78bc31tybbf",
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `UUID`,
				fieldValue: "57b73598-8764-4ad0-a76a-679bb6640eb1",
			},
			want: true,
		},
	}
	for _, tt := range tests {

		if got := tt.v.IsUUID4(tt.args.fieldName, tt.args.fieldValue); got != tt.want {
			t.Errorf("Validator.IsUUID4() = %v, want %v", got, tt.want)
		}

	}
}

func TestValidator_IsAlpha(t *testing.T) {
	validator := NewValidator()
	type args struct {
		fieldName  string
		fieldValue string
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `Name`,
				fieldValue: `\u`,
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Name`,
				fieldValue: "yoOgth",
			},
			want: true,
		},
	}
	for _, tt := range tests {

		if got := tt.v.IsAlpha(tt.args.fieldName, tt.args.fieldValue); got != tt.want {
			t.Errorf("Validator.IsAlpha() = %v, want %v", got, tt.want)
		}
	}
}

func TestValidator_IsAlphaNumeric(t *testing.T) {
	validator := NewValidator()
	type args struct {
		fieldName  string
		fieldValue string
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `Name`,
				fieldValue: `\ufff0`,
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Name`,
				fieldValue: "YUI900",
			},
			want: true,
		},
	}
	for _, tt := range tests {

		if got := tt.v.IsAlphaNumeric(tt.args.fieldName, tt.args.fieldValue); got != tt.want {
			t.Errorf("Validator.IsAlphaNumeric() = %v, want %v", got, tt.want)
		}
	}
}

func TestValidator_IsDigits(t *testing.T) {
	validator := NewValidator()
	type args struct {
		fieldName  string
		fieldValue string
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `Number`,
				fieldValue: `TY456`,
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Number`,
				fieldValue: "89362",
			},
			want: true,
		},
	}
	for _, tt := range tests {

		if got := tt.v.IsDigits(tt.args.fieldName, tt.args.fieldValue); got != tt.want {
			t.Errorf("Validator.IsDigits() = %v, want %v", got, tt.want)
		}
	}
}

func TestValidator_IsCreditCard(t *testing.T) {
	validator := NewValidator()
	type args struct {
		fieldName  string
		fieldValue string
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `Card Number`,
				fieldValue: `39444691789989`,
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Card Number`,
				fieldValue: "4716461583322103",
			},
			want: true,
		},
	}
	for _, tt := range tests {

		if got := tt.v.IsCreditCard(tt.args.fieldName, tt.args.fieldValue); got != tt.want {
			t.Errorf("Validator.IsCreditCard() = %v, want %v", got, tt.want)
		}

	}
}

func TestValidator_IsStrLenGtMin(t *testing.T) {
	validator := NewValidator()
	type args struct {
		fieldName  string
		fieldValue string
		min        int
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `Workspace Name`,
				fieldValue: `Output`,
				min:        int(8),
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Workspace Name`,
				fieldValue: "Output Direct Component",
				min:        int(8),
			},
			want: true,
		},
	}
	for _, tt := range tests {

		if got := tt.v.IsStrLenGtMin(tt.args.fieldName, tt.args.fieldValue, tt.args.min); got != tt.want {
			t.Errorf("Validator.IsStrLenGtMin() = %v, want %v", got, tt.want)
		}
	}
}

func TestValidator_IsStrLenLtMax(t *testing.T) {
	validator := NewValidator()
	type args struct {
		fieldName  string
		fieldValue string
		max        int
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `Workspace Name`,
				fieldValue: `Output Direct Component`,
				max:        int(8),
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Workspace Name`,
				fieldValue: "Output",
				max:        int(8),
			},
			want: true,
		},
	}
	for _, tt := range tests {

		if got := tt.v.IsStrLenLtMax(tt.args.fieldName, tt.args.fieldValue, tt.args.max); got != tt.want {
			t.Errorf("Validator.IsStrLenLtMax() = %v, want %v", got, tt.want)
		}
	}
}

func TestValidator_IsStrLenBetMinMax(t *testing.T) {
	validator := NewValidator()
	type args struct {
		fieldName  string
		fieldValue string
		min        int
		max        int
	}
	tests := []struct {
		v    *Validator
		args args
		want bool
	}{
		{
			v: validator,
			args: args{
				fieldName:  `Workspace Name`,
				fieldValue: "",
				min:        int(1),
				max:        int(50),
			},
			want: false,
		},
		{
			v: validator,
			args: args{
				fieldName:  `Workspace Name`,
				fieldValue: "Output Direct Component",
				min:        int(1),
				max:        int(50),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		if got := tt.v.IsStrLenBetMinMax(tt.args.fieldName, tt.args.fieldValue, tt.args.min, tt.args.max); got != tt.want {
			t.Errorf("Validator.IsStrLenBetMinMax() = %v, want %v", got, tt.want)
		}
	}
}
