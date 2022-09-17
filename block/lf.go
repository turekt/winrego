package block

import (
	"bytes"
)

type FastLeaf struct {
	HCellData
	Elements []NamedElement
}

func (fl *FastLeaf) marshal() ([]byte, error) {
	hcData, err := fl.HCellData.marshal()
	if err != nil {
		return nil, err
	}
	return binaryWriteAll(hcData, fl.Elements, fl.Padding)
}

func (fl *FastLeaf) unmarshal(data []byte) error {
	if err := fl.HCellData.unmarshal(data[:HCellDataSize]); err != nil {
		return err
	}
	if err := fl.assertPayloadDataSize(data); err != nil {
		return err
	}
	reader := bytes.NewReader(data[HCellDataSize:])

	fl.Elements = make([]NamedElement, fl.Metadata)
	if err := binaryBufferRead(reader, &fl.Elements); err != nil {
		return err
	}

	endPos := fl.Size() - HCellDataSize - int32(fl.Metadata*8)
	fl.Padding = make([]byte, endPos)
	return binaryBufferRead(reader, &fl.Padding)
}
