package block

import (
	"reflect"
	"testing"
)

func TestHBinMarshaling(t *testing.T) {
	recordHBin := []byte{
		// hbin
		0x68, 0x62, 0x69, 0x6e,
		// HBin data offset
		0x00, 0x00, 0x00, 0x00,
		// HBin size
		0x98, 0x00, 0x00, 0x00,
		// Reserved
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// Timestamp
		0x46, 0x10, 0xec, 0xc7, 0xdd, 0xfc, 0xd2, 0x01,
		// Spare
		0x00, 0x00, 0x00, 0x00,
		// HCell size
		0x88, 0xff, 0xff, 0xff,
		// Signature "nk"
		0x6e, 0x6b,
		// Flags
		0x2c, 0x00,
		// Timestamp
		0xbc, 0xfc, 0xf8, 0xab, 0x2b, 0xfc, 0xd2, 0x01,
		// Access bits
		0x00, 0x00, 0x00, 0x00,
		// Parent
		0xa8, 0x07, 0x00, 0x00,
		// Subkeys count
		0x02, 0x00, 0x00, 0x00,
		// VSubkeys count
		0x00, 0x00, 0x00, 0x00,
		// Subkeys list offset
		0xf8, 0x03, 0x00, 0x00,
		// VSubkeys list offset
		0xff, 0xff, 0xff, 0xff,
		// Key values count
		0x00, 0x00, 0x00, 0x00,
		// Key values list offset
		0xff, 0xff, 0xff, 0xff,
		// Key security offset
		0x98, 0x00, 0x00, 0x00,
		// Class name offset
		0xff, 0xff, 0xff, 0xff,
		// Largest subkey name length
		0x20, 0x00, 0x00, 0x00,
		// Largest subkey class name length
		0x00, 0x00, 0x00, 0x00,
		// Larest value name length
		0x00, 0x00, 0x00, 0x00,
		// Largest value data size
		0x00, 0x00, 0x00, 0x00,
		// WorkVar
		0x00, 0x00, 0x00, 0x00,
		// Key name length
		0x26, 0x00,
		// Class name length
		0x00, 0x00,
		// Key name
		0x7b, 0x34, 0x39, 0x65, 0x64, 0x65, 0x37, 0x37,
		0x66, 0x2d, 0x34, 0x62, 0x32, 0x66, 0x2d, 0x34,
		0x35, 0x62, 0x38, 0x2d, 0x62, 0x31, 0x66, 0x38,
		0x2d, 0x35, 0x62, 0x63, 0x37, 0x34, 0x30, 0x31,
		0x38, 0x32, 0x62, 0x64, 0x66, 0x7d, 0x00, 0x00,
	}
	hb := &HBin{}
	if err := hb.unmarshal(recordHBin); err != nil {
		t.Errorf("failed to unmarshal record: %v", err)
	}

	expect := &HBin{
		HBinHeader{
			HBinSignature:  0x6e696268,
			HBinDataOffset: 0,
			HBinSize:       152,
			Reserved1:      0,
			Timestamp:      0x01d2fcddc7ec1046,
			Spare:          0,
		},
		[]HCell{
			&KeyNode{
				HCellData{
					BlockSize:        -120,
					HCellSignature:   [2]byte{0x6e, 0x6b},
					Metadata:         0x2c,
					Padding:          []byte{0, 0},
					ParentHBin:       hb,
					ParentHBinOffset: hb.HBinDataOffset + HBinHeaderSize,
				},
				KeyNodeData{
					LastWTimestamp:         0x01d2fc2babf8fcbc,
					AccessBits:             0,
					Parent:                 0x07a8,
					SubkeysCount:           2,
					VSubkeysCount:          0,
					SubkeysListOffset:      0x3f8,
					VSubkeysListOffset:     -1,
					KeyValuesCount:         0,
					KeyValuesListOffset:    -1,
					KeySecurityOffset:      0x98,
					ClassNameOffset:        -1,
					LSubkeyNameLength:      0x20,
					LSubkeyClassNameLength: 0,
					LValueNameLength:       0,
					LValueDataSize:         0,
					WorkVar:                0,
					KeyNameLength:          0x26,
					ClassNameLength:        0,
				},
				[]byte("{49ede77f-4b2f-45b8-b1f8-5bc740182bdf}"),
			},
		},
		nil,
	}
	if got, want := hb.HBinHeader, expect.HBinHeader; !reflect.DeepEqual(got, want) {
		t.Errorf("hbin headers not equal (got|want):\n%+v\n%+v", got, want)
	}
	if got, want := len(hb.Cells), len(expect.Cells); got != want {
		t.Fatalf("hbin cell number not equal: got %d, want %d", got, want)
	}
	if got, want := hb.Cells[0], expect.Cells[0]; !reflect.DeepEqual(got, want) {
		t.Errorf("hbin cell structs not equal (got|want):\n%+v\n%+v", got, want)
	}

	data, err := hb.marshal()
	if err != nil {
		t.Errorf("failed to marshal record: %v", err)
	}

	if got, want := data, recordHBin; !reflect.DeepEqual(got, want) {
		t.Errorf("hbin data not equal (got|want):\n%+v\n%+v", got, want)
	}
}
