package block

import (
	"fmt"
)

const (
	RegNone = iota
	RegSz
	RegExpandSz
	RegBinary
	RegDWord
	RegDWordLittleEndian
	RegDWordBigEndian
	RegLink
	RegMultiSz
	RegResourceList
	RegFullResourceDescriptor
	RegResourceRequirementsList
	RegQWord
	RegUnknown = -1
)

type KeyValueData struct {
	DataSize   int32
	DataOffset int32
	DataType   uint32
	Flags      uint16
	Spare      uint16
}

func (kvd *KeyValueData) marshal() ([]byte, error) {
	return binaryWrite(kvd)
}

func (kvd *KeyValueData) unmarshal(data []byte) error {
	return binaryRead(data, kvd)
}

type KeyValue struct {
	HCellData
	KeyValueData
	ValueName []byte
}

func (kv *KeyValue) marshal() ([]byte, error) {
	hcData, err := kv.HCellData.marshal()
	if err != nil {
		return nil, err
	}
	return binaryWriteAll(hcData, kv.KeyValueData, kv.ValueName, kv.Padding)
}

func (kv *KeyValue) unmarshal(data []byte) error {
	if err := kv.HCellData.unmarshal(data[:HCellDataSize]); err != nil {
		return err
	}
	if err := kv.assertPayloadDataSize(data); err != nil {
		return err
	}
	if err := binaryRead(data[HCellDataSize:], &kv.KeyValueData); err != nil {
		return err
	}

	const keyValueDataEnd = HCellDataSize + KeyValueDataSize
	dEnd := keyValueDataEnd + int(kv.Metadata)
	if dEnd > len(data) {
		return fmt.Errorf("name length out of bounds: name %d len %d", kv.Metadata, len(data))
	}
	kv.ValueName = data[keyValueDataEnd:dEnd]
	kv.Padding = data[dEnd:kv.Size()]
	return nil
}
