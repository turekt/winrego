package block

import (
	"bytes"
)

type IndexRoot struct {
	HCellData
	Elements []OffsetElement
}

func (ir *IndexRoot) marshal() ([]byte, error) {
	hcData, err := ir.HCellData.marshal()
	if err != nil {
		return nil, err
	}
	return binaryWriteAll(hcData, ir.Elements, ir.Padding)
}

func (ir *IndexRoot) unmarshal(data []byte) error {
	if err := ir.HCellData.unmarshal(data[:HCellDataSize]); err != nil {
		return err
	}
	if err := ir.assertPayloadDataSize(data); err != nil {
		return err
	}
	reader := bytes.NewReader(data[HCellDataSize:])

	ir.Elements = make([]OffsetElement, ir.Metadata)
	if err := binaryBufferRead(reader, &ir.Elements); err != nil {
		return err
	}

	endPos := ir.Size() - HCellDataSize - int32(ir.Metadata*4)
	ir.Padding = make([]byte, endPos)
	return binaryBufferRead(reader, &ir.Padding)
}
