package block

import (
	"fmt"
)

type KeyNodeFlag uint16

const (
	KeyVolatile KeyNodeFlag = iota << 1
	KeyHiveExit
	KeyHiveEntry
	KeyNoDelete
	KeySymLink
	KeyCompName
	KeyPredefHandle
	KeyVirtMirrored
	KeyVirtTarget
	KeyVirtualStore
)

type KeyNodeData struct {
	LastWTimestamp         uint64
	AccessBits             uint32
	Parent                 int32
	SubkeysCount           int32
	VSubkeysCount          int32 // Volatile Subkey
	SubkeysListOffset      int32
	VSubkeysListOffset     int32 // Volatile Subkey
	KeyValuesCount         int32
	KeyValuesListOffset    int32
	KeySecurityOffset      int32
	ClassNameOffset        int32
	LSubkeyNameLength      int32 // Largest Subkey
	LSubkeyClassNameLength int32 // Largest Subkey
	LValueNameLength       int32 // Largest Value
	LValueDataSize         int32 // Largest Value
	WorkVar                int32
	KeyNameLength          int16
	ClassNameLength        int16
}

func (knd *KeyNodeData) marshal() ([]byte, error) {
	return binaryWrite(knd)
}

func (knd *KeyNodeData) unmarshal(data []byte) error {
	return binaryRead(data, knd)
}

type KeyNode struct {
	HCellData
	KeyNodeData
	KeyName []byte
}

func (kn *KeyNode) marshal() ([]byte, error) {
	hcData, err := kn.HCellData.marshal()
	if err != nil {
		return nil, err
	}
	return binaryWriteAll(hcData, kn.KeyNodeData, kn.KeyName, kn.Padding)
}

func (kn *KeyNode) unmarshal(data []byte) error {
	if err := kn.HCellData.unmarshal(data[:HCellDataSize]); err != nil {
		return err
	}
	if err := kn.assertPayloadDataSize(data); err != nil {
		return err
	}
	if err := binaryRead(data[HCellDataSize:], &kn.KeyNodeData); err != nil {
		return err
	}

	const keyNodeDataEnd = HCellDataSize + KeyNodeDataSize
	dEnd := keyNodeDataEnd + int(kn.KeyNameLength)
	if dEnd > len(data) {
		return fmt.Errorf("key name length out of bounds: name %d len %d", kn.KeyNameLength, len(data))
	}
	kn.KeyName = data[keyNodeDataEnd:dEnd]
	kn.Padding = data[dEnd:kn.Size()]
	return nil
}
