package block

import (
	"bytes"
	"fmt"
)

const (
	HBinHeaderSize = 32
)

type HBinHeader struct {
	HBinSignature  uint32
	HBinDataOffset int32
	HBinSize       int32
	Reserved1      uint64
	Timestamp      uint64
	Spare          uint32
}

func (hbHeader *HBinHeader) marshal() ([]byte, error) {
	return binaryWrite(hbHeader)
}

func (hbHeader *HBinHeader) unmarshal(data []byte) error {
	if len(data) < HBinHeaderSize {
		return fmt.Errorf("hb header of size %d, expected at least %d", len(data), HBinHeaderSize)
	}
	return binaryRead(data, hbHeader)
}

func (hbHeader *HBinHeader) Size() int32 {
	return HBinHeaderSize
}

func (hbHeader *HBinHeader) Offset() int32 {
	return hbHeader.HBinDataOffset
}

func (hbHeader *HBinHeader) Signature() string {
	return string(uint32toba(hbHeader.HBinSignature))
}

type HBinData []HBin

func (hbins *HBinData) marshal() ([]byte, error) {
	var buf bytes.Buffer
	for _, hbin := range *hbins {
		b, err := hbin.marshal()
		if err != nil {
			return buf.Bytes(), err
		}
		buf.Write(b)
	}
	return buf.Bytes(), nil
}

func (hbins *HBinData) unmarshal(data []byte) error {
	hb := &HBin{}
	*hbins = make(HBinData, 0)
	for start := int32(0); start < int32(len(data)); start += hb.HBinHeader.HBinSize {
		if err := hb.unmarshal(data[start : start+hb.HBinHeader.HBinSize]); err != nil {
			return err
		}
		hb.HBinDataPtr = hbins
		*hbins = append(*hbins, *hb)
	}
	return nil
}

func (hbins *HBinData) Size() int32 {
	return int32(len(*hbins))
}

func (hbins *HBinData) Offset() int32 {
	return BaseBlockSize
}

func (hbins *HBinData) Signature() string {
	return "hbins"
}

type HBin struct {
	HBinHeader
	// List of cells that this HBin contains
	// The cells are ordered sequentially as they
	// are stored on the file system
	Cells []HCell
	// Pointer to the start of the HBin data area
	HBinDataPtr *HBinData
}

func (hb *HBin) marshal() ([]byte, error) {
	hbHeader, err := hb.HBinHeader.marshal()
	if err != nil {
		return hbHeader, err
	}

	buf := new(bytes.Buffer)
	if err := binaryBufferWrite(buf, hbHeader); err != nil {
		return buf.Bytes(), err
	}

	for _, hc := range hb.Cells {
		data, err := hc.marshal()
		if err != nil {
			return buf.Bytes(), err
		}

		if err = binaryBufferWrite(buf, data); err != nil {
			return buf.Bytes(), err
		}
	}

	return buf.Bytes(), nil
}

func (hb *HBin) unmarshal(data []byte) error {
	if err := hb.HBinHeader.unmarshal(data[:HBinHeaderSize]); err != nil {
		return err
	}
	hb.Cells = make([]HCell, 0)

	var cellSize int32
	for i := int32(HBinHeaderSize); i < hb.HBinHeader.HBinSize; i += cellSize {
		if err := binaryRead(data[i:i+HCellSizeLength], &cellSize); err != nil {
			return err
		}

		if cellSize < 0 {
			cellSize *= -1
		}

		cell, err := UnmarshalHCell(data[i : i+cellSize])
		if err != nil {
			return err
		}

		cell.setParentHBin(hb)
		cell.setOffset(i)
		hb.Cells = append(hb.Cells, cell)
	}

	return nil
}

func (hb *HBin) Size() int32 {
	return hb.HBinHeader.HBinSize
}
