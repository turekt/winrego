package block

import (
	"reflect"
	"testing"
)

func TestHashLeafMarshaling(t *testing.T) {
	recordHashLeaf := []byte{
		// Size
		0x28, 0x00, 0x00, 0x00,
		// Signature
		0x6c, 0x68,
		// Element count
		0x03, 0x00,
		// Offset 1
		0x00, 0x00, 0x00, 0x00,
		// Name hash 1
		0x00, 0x00, 0x00, 0x00,
		// Offset 2
		0x00, 0x00, 0x00, 0x00,
		// Name hash 2
		0x00, 0x00, 0x00, 0x00,
		// Offset 3
		0x00, 0x00, 0x00, 0x00,
		// Name hash 3
		0x00, 0x00, 0x00, 0x00,
		// Padding
		0x53, 0x4c, 0x43, 0x4b,
		0x00, 0x00, 0x00, 0x00,
	}
	hl := &HashLeaf{}
	if err := hl.unmarshal(recordHashLeaf); err != nil {
		t.Fatalf("failed unmarshaling lh record %v", err)
	}
	got, want := hl, &HashLeaf{
		HCellData{
			BlockSize:      40,
			HCellSignature: [2]byte{0x6c, 0x68},
			Metadata:       3,
			Padding: []byte{
				0x53, 0x4c, 0x43, 0x4b,
				0x00, 0x00, 0x00, 0x00,
			},
		},
		[]NamedElement{
			{0, [4]byte{0, 0, 0, 0}},
			{0, [4]byte{0, 0, 0, 0}},
			{0, [4]byte{0, 0, 0, 0}},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("hash leaf struct not equals: got|want\n%+v\n%+v", got, want)
	}

	data, err := hl.marshal()
	if err != nil {
		t.Fatalf("failed marshaling hl record %v", err)
	}
	if got, want := data, recordHashLeaf; !reflect.DeepEqual(got, want) {
		t.Errorf("hl records not equals: got|want\n%+v\n%+v", got, want)
	}
}
