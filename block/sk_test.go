package block

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestKeySecurityMarshaling(t *testing.T) {
	recordKeySecurity := []byte{
		// Size
		0x2e, 0x00, 0x00, 0x00,
		// Signature
		0x73, 0x6b,
		// Reserved
		0x00, 0x00,
		// Flink
		0x00, 0x00, 0x00, 0x00,
		// Blink
		0x00, 0x00, 0x00, 0x00,
		// Reference count
		0x00, 0x00, 0x00, 0x00,
		// Security descriptor size
		0x16, 0x00, 0x00, 0x00,
		// Security descriptor value
		0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa,
		0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa,
	}
	ks := &KeySecurity{}
	if err := ks.unmarshal(recordKeySecurity); err != nil {
		t.Fatalf("failed unmarshaling sk record %v", err)
	}
	got, want := ks, &KeySecurity{
		HCellData{
			BlockSize:      46,
			HCellSignature: [2]byte{0x73, 0x6b},
			Metadata:       0,
			Padding:        []byte{},
		},
		KeySecurityData{
			Flink:             0,
			Blink:             0,
			RefCount:          0,
			SecDescriptorSize: 22,
		},
		bytes.Repeat([]byte{0xaa}, 22),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("key security struct not equals: got|want\n%+v\n%+v", got, want)
	}

	data, err := ks.marshal()
	if err != nil {
		t.Fatalf("failed marshaling sk record %v", err)
	}
	if got, want := data, recordKeySecurity; !reflect.DeepEqual(got, want) {
		t.Errorf("sk records not equals: got|want\n%+v\n%+v", got, want)
	}
}

func TestCorruptKeySecurityMarshaling(t *testing.T) {
	recordKeySecurityCorrupt := [][]byte{
		{
			0x20, 0x00, 0x00, 0x00,
			0x73, 0x6b,
			0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x36, 0x00, 0x00, 0x00,
			0xaa, 0xaa, 0xaa, 0xaa,
			0xaa, 0xaa, 0xaa, 0xaa,
		},
		{
			0x20, 0x00, 0x00, 0x00,
			0x73, 0x6b,
			0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x04, 0x00, 0x00, 0x00,
			0xaa, 0xaa, 0xaa, 0xaa,
			0xbb, 0xbb, 0xbb, 0xbb,
		},
		{
			0x10, 0x00, 0x00, 0x00,
			0x73, 0x6b,
			0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
		},
	}
	expect := []struct {
		Success bool
		Result  *KeySecurity
		Error   func(err error) bool
	}{
		{
			Success: false,
			Result:  nil,
			Error:   isOutOfBoundsErr,
		},
		{
			Success: true,
			Result: &KeySecurity{
				HCellData{
					BlockSize:      32,
					HCellSignature: [2]byte{0x73, 0x6b},
					Metadata:       0,
					Padding:        bytes.Repeat([]byte{0xbb}, 4),
				},
				KeySecurityData{
					Flink:             0,
					Blink:             0,
					RefCount:          0,
					SecDescriptorSize: 4,
				},
				bytes.Repeat([]byte{0xaa}, 4),
			},
			Error: isOutOfBoundsErr,
		},
		{
			Success: false,
			Result:  nil,
			Error:   nil,
		},
	}

	for i := 0; i < len(recordKeySecurityCorrupt); i++ {
		e := expect[i]
		ks := &KeySecurity{}
		if err := ks.unmarshal(recordKeySecurityCorrupt[i]); err != nil {
			if e.Success {
				t.Fatalf("record %d should have been a success, got %v", i, err)
			}
			if e.Error != nil && !e.Error(err) {
				t.Fatalf("expected out of bounds err on record %d, got %v", i, err)
			}
			continue
		}
		if !e.Success {
			t.Fatalf("record %d should have been an error, got %+v", i, ks)
		}
		if got, want := ks, e.Result; !reflect.DeepEqual(got, want) {
			t.Fatalf("result mismatch, got|want:\n%+v\n%+v", got, want)
		}
	}
}

func isOutOfBoundsErr(err error) bool {
	return err != nil && strings.Contains(fmt.Sprint(err), "out of bounds")
}
