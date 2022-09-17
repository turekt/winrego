package block

import (
	"reflect"
	"testing"
)

func TestIndexRootMarshaling(t *testing.T) {
	recordIndexRoot := []byte{
		// Size
		0x18, 0x00, 0x00, 0x00,
		// Signature
		0x72, 0x69,
		// Element count
		0x03, 0x00,
		// Offset 1
		0x00, 0x00, 0x00, 0x00,
		// Offset 2
		0x00, 0x00, 0x00, 0x00,
		// Offset 3
		0x53, 0x4c, 0x43, 0x4b,
		// Padding
		0x00, 0x00, 0x00, 0x00,
	}
	ir := &IndexRoot{}
	if err := ir.unmarshal(recordIndexRoot); err != nil {
		t.Fatalf("failed unmarshaling ri record %v", err)
	}
	got, want := ir, &IndexRoot{
		HCellData{
			BlockSize:      24,
			HCellSignature: [2]byte{0x72, 0x69},
			Metadata:       3,
			Padding:        []byte{0x00, 0x00, 0x00, 0x00},
		},
		[]OffsetElement{0, 0, 1262701651},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("index root struct not equals: got|want\n%+v\n%+v", got, want)
	}

	data, err := ir.marshal()
	if err != nil {
		t.Fatalf("failed marshaling ri record %v", err)
	}
	if got, want := data, recordIndexRoot; !reflect.DeepEqual(got, want) {
		t.Errorf("ri records not equals: got|want\n%+v\n%+v", got, want)
	}
}
