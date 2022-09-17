package block

import (
	"reflect"
	"testing"
)

func TestKeyValueMarshaling(t *testing.T) {
	recordKeyValue := []byte{
		// Size
		0x1b, 0x00, 0x00, 0x00,
		// Signature
		0x76, 0x6b,
		// Name length
		0x03, 0x00,
		// Data size
		0x00, 0x00, 0x00, 0x00,
		// Data offset
		0x00, 0x00, 0x00, 0x00,
		// Data type
		0x00, 0x00, 0x00, 0x00,
		// Flags
		0x00, 0x00,
		// Spare
		0x00, 0x00,
		// Value name string
		0x31, 0x32, 0x33,
	}
	kv := &KeyValue{}
	if err := kv.unmarshal(recordKeyValue); err != nil {
		t.Fatalf("failed unmarshaling vk record %v", err)
	}
	got, want := kv, &KeyValue{
		HCellData{
			BlockSize:      0x1b,
			HCellSignature: [2]byte{0x76, 0x6b},
			Metadata:       3,
			Padding:        []byte{},
		},
		KeyValueData{
			DataSize:   0,
			DataOffset: 0,
			DataType:   0,
			Flags:      0,
			Spare:      0,
		},
		[]byte("123"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("key value struct not equals: got|want\n%+v\n%+v", got, want)
	}

	data, err := kv.marshal()
	if err != nil {
		t.Fatalf("failed marshaling vk record %v", err)
	}
	if got, want := data, recordKeyValue; !reflect.DeepEqual(got, want) {
		t.Errorf("vk records not equals: got|want\n%+v\n%+v", got, want)
	}
}
