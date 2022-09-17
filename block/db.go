package block

type BigData struct {
	HCellData
	DataOffset int32
}

func (bd *BigData) marshal() ([]byte, error) {
	hcData, err := bd.HCellData.marshal()
	if err != nil {
		return nil, err
	}
	return binaryWriteAll(hcData, bd.DataOffset, bd.Padding)
}

func (bd *BigData) unmarshal(data []byte) error {
	if err := bd.HCellData.unmarshal(data); err != nil {
		return err
	}
	if err := bd.assertPayloadDataSize(data); err != nil {
		return err
	}
	if err := binaryReadAll(data[HCellDataSize:], &bd.DataOffset); err != nil {
		return err
	}

	dEnd := HCellDataSize + 4
	bd.Padding = data[dEnd:bd.Size()]
	return nil
}
