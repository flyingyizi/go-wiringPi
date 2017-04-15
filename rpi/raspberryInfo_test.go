package rpi

import (
	"reflect"
	"testing"
)

func Test_getPostRPI2FromRevision(t *testing.T) {
	tests := []struct {
		name     string
		revision string
		wantInfo RpiInfoT
		wantErr  bool
	}{
		// TODO: Add test cases.
		{name: "pi 3", revision: "a22082", wantInfo: RpiInfoT{model: Model3B, mem: Rpi1024MB, processor: Broadcom2837,
			manufacturer: MakerEmbest, pcbRev: PcbRev1_2, overVolted: false, revision: 0xa22082}, wantErr: false},
		{name: "pi B+", revision: "900032", wantInfo: RpiInfoT{model: ModelBPlus, mem: Rpi512MB, processor: Broadcom2835,
			manufacturer: MakerSony, pcbRev: PcbRev1_2, overVolted: false, revision: 0x900032}, wantErr: false},
		{name: "pi B", revision: "0002", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfo, err := getPostRPI2FromRevision(tt.revision)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPostRPI2FromRevision() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotInfo, tt.wantInfo) {
				t.Errorf("getPostRPI2FromRevision() = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}

func Test_getPreRPI2FromRevision(t *testing.T) {
	tests := []struct {
		name     string
		revision string
		wantInfo RpiInfoT
		wantErr  bool
	}{
		// TODO: Add test cases.
		{name: "pi 3", revision: "a22082", wantErr: true},
		{name: "pi B", revision: "0002", wantInfo: RpiInfoT{model: ModelB, mem: Rpi256MB, processor: Broadcom2835,
			manufacturer: MakerEgoman, pcbRev: PcbRev1, overVolted: false, revision: 0x0002}, wantErr: false},
		{name: "pi B", revision: "0000", wantInfo: RpiInfoT{model: ModelUnknown, mem: RpiUnknownMB, processor: Broadcom2835,
			manufacturer: MakerEgoman, pcbRev: PcbRev1, overVolted: false, revision: 0x0002}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfo, err := getPreRPI2FromRevision(tt.revision)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPreRPI2FromRevision() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotInfo, tt.wantInfo) {
				t.Errorf("getPreRPI2FromRevision() = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}

func Test_get_dt_ranges(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name     string
		args     args
		wantBase int64
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBase, err := get_dt_ranges(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("get_dt_ranges() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotBase != tt.wantBase {
				t.Errorf("get_dt_ranges() = %v, want %v", gotBase, tt.wantBase)
			}
		})
	}
}
