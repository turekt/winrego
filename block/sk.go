package block

import (
	"fmt"
)

type KeySecurityData struct {
	Flink             int32
	Blink             int32
	RefCount          uint32
	SecDescriptorSize uint32
}

func (ksd *KeySecurityData) marshal() ([]byte, error) {
	return binaryWrite(ksd)
}

func (ksd *KeySecurityData) unmarshal(data []byte) error {
	return binaryRead(data, ksd)
}

type KeySecurity struct {
	HCellData
	KeySecurityData
	SecDescriptor []byte
}

func (ks *KeySecurity) marshal() ([]byte, error) {
	hcData, err := ks.HCellData.marshal()
	if err != nil {
		return nil, err
	}
	return binaryWriteAll(hcData, ks.KeySecurityData, ks.SecDescriptor, ks.Padding)
}

func (ks *KeySecurity) unmarshal(data []byte) error {
	if err := ks.HCellData.unmarshal(data[:HCellDataSize]); err != nil {
		return err
	}
	if err := ks.assertPayloadDataSize(data); err != nil {
		return err
	}
	if err := binaryRead(data[HCellDataSize:], &ks.KeySecurityData); err != nil {
		return err
	}

	const keySecurityDataEnd = HCellDataSize + KeySecurityDataSize
	dEnd := uint32(keySecurityDataEnd) + ks.SecDescriptorSize
	if dEnd > uint32(len(data)) {
		return fmt.Errorf("sec descriptor size out of bounds: size %d len %d", ks.SecDescriptorSize, len(data))
	}
	ks.SecDescriptor = data[keySecurityDataEnd:dEnd]
	ks.Padding = data[dEnd:ks.Size()]
	return nil
}
