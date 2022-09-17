package block

import (
	"bytes"
)

type HashLeaf struct {
	HCellData
	Elements []NamedElement
}

func (hl *HashLeaf) marshal() ([]byte, error) {
	hcData, err := hl.HCellData.marshal()
	if err != nil {
		return nil, err
	}
	return binaryWriteAll(hcData, hl.Elements, hl.Padding)
}

func (hl *HashLeaf) unmarshal(data []byte) error {
	if err := hl.HCellData.unmarshal(data[:HCellDataSize]); err != nil {
		return err
	}
	if err := hl.assertPayloadDataSize(data); err != nil {
		return err
	}
	reader := bytes.NewReader(data[HCellDataSize:])

	hl.Elements = make([]NamedElement, hl.Metadata)
	if err := binaryBufferRead(reader, &hl.Elements); err != nil {
		return err
	}

	endPos := hl.Size() - HCellDataSize - int32(hl.Metadata*8)
	hl.Padding = make([]byte, endPos)
	return binaryBufferRead(reader, &hl.Padding)
}
