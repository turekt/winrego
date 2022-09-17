package block

import (
	"reflect"
	"testing"
)

func TestKeyNodeMarshaling(t *testing.T) {
	recordKeyNode := []byte{
		// Size
		0x5c, 0x00, 0x00, 0x00,
		// Signature
		0x6e, 0x6b,
		// Flags
		0x20, 0x00,
		// Last written timestamp
		0x03, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
		// Access bits
		0x04, 0x00, 0x00, 0x00,
		// Parent
		0x00, 0x00, 0x00, 0x00,
		// Subkeys count
		0x00, 0x00, 0x00, 0x00,
		// Volatile subkeys count
		0x00, 0x00, 0x00, 0x00,
		// Subkeys list offset
		0x00, 0x00, 0x00, 0x00,
		// Volatile subkeys count
		0x00, 0x00, 0x00, 0x00,
		// Key values count
		0x00, 0x00, 0x00, 0x00,
		// Key values list offset
		0x00, 0x00, 0x00, 0x00,
		// Key security offset
		0x00, 0x00, 0x00, 0x00,
		// Class name offset
		0x00, 0x00, 0x00, 0x00,
		// Largest subkey name length
		0x00, 0x00, 0x00, 0x00,
		// Largest subkey class name length
		0x00, 0x00, 0x00, 0x00,
		// Largest value name length
		0x00, 0x00, 0x00, 0x00,
		// Largest value data size
		0x00, 0x00, 0x00, 0x00,
		// WorkVar
		0x0c, 0x00, 0x00, 0x00,
		// Key name length
		0x0c, 0x00,
		// Class name length
		0x00, 0x00,
		// Key name string
		0x4f, 0x70, 0x65, 0x6e,
		0x57, 0x69, 0x74, 0x68,
		0x4c, 0x69, 0x73, 0x74,
	}
	kn := &KeyNode{}
	if err := kn.unmarshal(recordKeyNode); err != nil {
		t.Fatalf("failed unmarshaling nk record %v", err)
	}
	got, want := kn, &KeyNode{
		HCellData{
			BlockSize:      0x5c,
			HCellSignature: [2]byte{0x6e, 0x6b},
			Metadata:       32,
			Padding:        []byte{},
		},
		KeyNodeData{
			LastWTimestamp:         4294967299,
			AccessBits:             4,
			Parent:                 0,
			SubkeysCount:           0,
			SubkeysListOffset:      0,
			VSubkeysListOffset:     0,
			KeyValuesCount:         0,
			KeyValuesListOffset:    0,
			KeySecurityOffset:      0,
			ClassNameOffset:        0,
			LSubkeyNameLength:      0,
			LSubkeyClassNameLength: 0,
			LValueNameLength:       0,
			LValueDataSize:         0,
			WorkVar:                12,
			KeyNameLength:          12,
			ClassNameLength:        0,
		},
		[]byte("OpenWithList"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("key node struct not equals: got|want\n%+v\n%+v", got, want)
	}

	data, err := kn.marshal()
	if err != nil {
		t.Fatalf("failed marshaling nk record %v", err)
	}
	if got, want := data, recordKeyNode; !reflect.DeepEqual(got, want) {
		t.Errorf("nk records not equals: got|want\n%+v\n%+v", got, want)
	}
}
