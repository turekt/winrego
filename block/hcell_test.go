package block

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDataRecordMarshaling(t *testing.T) {
	recordData := [][]byte{
		{
			// Size
			0x10, 0x00, 0x00, 0x00,
			// Data
			0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
			0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
		},
		{
			// Size
			0x08, 0x00, 0x00, 0x00,
			// Data
			0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42,
		},
		{
			// Size
			0x08, 0x00, 0x00, 0x00,
			// Data
			0x43, 0x43, 0x43, 0x43,
		},
	}
	expect := []*DataRecord{
		&DataRecord{
			HCellData: HCellData{BlockSize: 16},
			Data:      []byte("AAAAAAAAAAAAAAAA"),
		},
		&DataRecord{
			HCellData: HCellData{BlockSize: 8},
			Data:      []byte("BBBBBBBB"),
		},
		&DataRecord{
			HCellData: HCellData{BlockSize: 8},
			Data:      []byte("CCCC"),
		},
	}

	for i := 0; i < len(recordData); i++ {
		dr := &DataRecord{}
		if err := dr.unmarshal(recordData[i]); err != nil {
			t.Fatalf("failed unmarshaling data record %d %v", i, err)
		}
		if got, want := dr, expect[i]; !reflect.DeepEqual(got, want) {
			t.Errorf("data record struct %d not equals: got|want\n%+v\n%+v", i, got, want)
		}

		data, err := dr.marshal()
		if err != nil {
			t.Fatalf("failed marshaling data record %d %v", i, err)
		}
		if got, want := data, recordData[i]; !reflect.DeepEqual(got, want) {
			t.Errorf("data record struct %d not equals: got|want\n%+v\n%+v", i, got, want)
		}
	}
}

func TestInsufficientHCellDataProvided(t *testing.T) {
	recordHCellCorrupt := [][]byte{
		{
			0x20, 0x00, 0x00, 0x00,
			0x6c, 0x66,
			0x03, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x53, 0x4c, 0x43, 0x4b,
			0x53, 0x4c, 0x43, 0x4b,
			0x00, 0x00, 0x00, 0x00,
			0x02, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x03,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
		},
		{
			0x20, 0x00, 0x00, 0x00,
			0x72, 0x69,
			0x03, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x53, 0x4c, 0x43, 0x4b,
			0x00, 0x00, 0x00, 0x00,
		},
		{
			0x2e, 0x00, 0x00, 0x00,
			0x73, 0x6b,
			0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x16, 0x00, 0x00, 0x00,
		},
	}
	errMsgFunc := func(dataLen int, expectLen byte) error {
		return fmt.Errorf("provided data of size %d, expected at least %d", dataLen, expectLen)
	}
	expect := []error{
		nil,
		errMsgFunc(len(recordHCellCorrupt[1]), recordHCellCorrupt[1][0]),
		errMsgFunc(len(recordHCellCorrupt[2]), recordHCellCorrupt[2][0]),
	}

	hc := &HCellData{}
	for i := 0; i < len(expect); i++ {
		if err := hc.unmarshal(recordHCellCorrupt[i]); err != nil {
			t.Errorf("failed to marshal record %d: %v", i, err)
		}
		got, want := hc.assertPayloadDataSize(recordHCellCorrupt[i]), expect[i]
		switch {
		case got == nil && want != nil:
			fallthrough
		case got != nil && want == nil:
			fallthrough
		case fmt.Sprint(got) != fmt.Sprint(want):
			t.Errorf("error mismatch: got %v, want %v", got, want)
		}
	}
}
