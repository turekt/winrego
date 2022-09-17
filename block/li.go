package block

import (
	"bytes"
)

type IndexLeaf struct {
	HCellData
	Elements []OffsetElement
}

func (il *IndexLeaf) marshal() ([]byte, error) {
	hcData, err := il.HCellData.marshal()
	if err != nil {
		return nil, err
	}
	return binaryWriteAll(hcData, il.Elements, il.Padding)
}

func (il *IndexLeaf) unmarshal(data []byte) error {
	if err := il.HCellData.unmarshal(data[:HCellDataSize]); err != nil {
		return err
	}
	if err := il.assertPayloadDataSize(data); err != nil {
		return err
	}
	reader := bytes.NewReader(data[HCellDataSize:])

	il.Elements = make([]OffsetElement, il.Metadata)
	if err := binaryBufferRead(reader, &il.Elements); err != nil {
		return err
	}

	endPos := il.Size() - HCellDataSize - int32(il.Metadata*4)
	il.Padding = make([]byte, endPos)
	return binaryBufferRead(reader, &il.Padding)
}
