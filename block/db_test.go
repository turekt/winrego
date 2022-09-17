package block

import (
	"reflect"
	"testing"
)

func TestBigDataMarshaling(t *testing.T) {
	recordBigData := []byte{
		// Size
		0x10, 0x00, 0x00, 0x00,
		// Signature
		0x64, 0x62,
		// Segment count
		0x02, 0x00,
		// Segment list offset
		0x09, 0x00, 0x00, 0x00,
		// Padding
		0x00, 0x00, 0x00, 0x00,
	}
	bd := &BigData{}
	if err := bd.unmarshal(recordBigData); err != nil {
		t.Fatalf("failed unmarshaling db record %v", err)
	}
	got, want := bd, &BigData{
		HCellData{
			BlockSize:      16,
			HCellSignature: [2]byte{0x64, 0x62},
			Metadata:       2,
			Padding:        []byte{0, 0, 0, 0},
		},
		int32(9),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("big data struct not equals: got|want\n%+v\n%+v", got, want)
	}

	data, err := bd.marshal()
	if err != nil {
		t.Fatalf("failed marshaling db record %v", err)
	}
	if got, want := data, recordBigData; !reflect.DeepEqual(got, want) {
		t.Errorf("db records not equals: got|want\n%+v\n%+v", got, want)
	}
}
