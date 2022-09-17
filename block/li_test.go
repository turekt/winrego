package block

import (
	"reflect"
	"testing"
)

func TestIndexLeafMarshaling(t *testing.T) {
	recordIndexLeaf := []byte{
		// Size
		0x18, 0x00, 0x00, 0x00,
		// Signature
		0x6c, 0x69,
		// Element count
		0x03, 0x00,
		// Offset 1
		0x00, 0x00, 0x00, 0x00,
		// Offset 2
		0x53, 0x4c, 0x43, 0x4b,
		// Offset 3
		0x00, 0x00, 0x00, 0x00,
		// Padding
		0x00, 0x00, 0x00, 0x00,
	}
	il := &IndexLeaf{}
	if err := il.unmarshal(recordIndexLeaf); err != nil {
		t.Fatalf("failed unmarshaling li record %v", err)
	}
	got, want := il, &IndexLeaf{
		HCellData{
			BlockSize:      24,
			HCellSignature: [2]byte{0x6c, 0x69},
			Metadata:       3,
			Padding: []byte{
				0x00, 0x00, 0x00, 0x00,
			},
		},
		[]OffsetElement{0, 1262701651, 0},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("index leaf struct not equals: got|want\n%+v\n%+v", got, want)
	}

	data, err := il.marshal()
	if err != nil {
		t.Fatalf("failed marshaling li record %v", err)
	}
	if got, want := data, recordIndexLeaf; !reflect.DeepEqual(got, want) {
		t.Errorf("li records not equals: got|want\n%+v\n%+v", got, want)
	}
}
