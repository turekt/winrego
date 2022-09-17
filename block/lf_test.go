package block

import (
	"bytes"
	"reflect"
	"testing"
)

func TestFastLeafMarshaling(t *testing.T) {
	recordFastLeaf := []byte{
		// Size
		0x28, 0x00, 0x00, 0x00,
		// Signature
		0x6c, 0x66,
		// Element count
		0x03, 0x00,
		// Offset 1
		0x00, 0x00, 0x00, 0x00,
		// Name hint 1
		0x53, 0x4c, 0x43, 0x4b,
		// Offset 2
		0x53, 0x4c, 0x43, 0x4b,
		// Name hint 2
		0x00, 0x00, 0x00, 0x00,
		// Offset 3
		0x02, 0x00, 0x00, 0x00,
		// Name hint 3
		0x00, 0x00, 0x00, 0x03,
		// Padding
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}
	fl := &FastLeaf{}
	if err := fl.unmarshal(recordFastLeaf); err != nil {
		t.Fatalf("failed unmarshaling lf record %v", err)
	}
	got, want := fl, &FastLeaf{
		HCellData{
			BlockSize:      40,
			HCellSignature: [2]byte{0x6c, 0x66},
			Metadata:       3,
			Padding:        bytes.Repeat([]byte{0}, 8),
		},
		[]NamedElement{
			{0, [4]byte{0x53, 0x4c, 0x43, 0x4b}},
			{1262701651, [4]byte{0, 0, 0, 0}},
			{2, [4]byte{0, 0, 0, 3}},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("fast leaf struct not equals: got|want\n%+v\n%+v", got, want)
	}

	data, err := fl.marshal()
	if err != nil {
		t.Fatalf("failed marshaling fl record %v", err)
	}
	if got, want := data, recordFastLeaf; !reflect.DeepEqual(got, want) {
		t.Errorf("fl records not equals: got|want\n%+v\n%+v", got, want)
	}
}
