package i2c

import "testing"

func Test_getEndian(t *testing.T) {
	tests := []struct {
		name    string
		wantRet bool
	}{
		// TODO: Add test cases.
		{name: "x86", wantRet: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := getEndian(); gotRet != tt.wantRet {
				t.Errorf("getEndian() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}
