package block

import (
	"fmt"
)

const (
	HCellSizeLength     = 4
	HCellDataSize       = 8
	NamedElementSize    = 8
	KeyValueDataSize    = 16
	KeySecurityDataSize = 16
	KeyNodeDataSize     = 72
)

type HCell interface {
	RegistryBlock
	setParentHBin(hbin *HBin)
	setOffset(offset int32)
}

type HCellData struct {
	// Size of this HCell, including this field
	BlockSize int32
	// One of ["li", "lh", "ri", "vk", "nk", "db", "sk"]
	HCellSignature [2]byte
	// All HCell types have these two bytes,  but there
	// is a different meaning depending on type
	// This is either:
	//   - Metadata - number of elements
	//   - NameLength - selfexplanatory
	//   - Reserved - in Key Security
	//   - Flags - in Key Node
	Metadata uint16
	// Points to the start of HBin where this cell is located
	ParentHBin *HBin
	// Offset to this cell from parent HBin
	ParentHBinOffset int32
	// Optional padding that this cell contains
	// Needed for precise marshaling
	Padding []byte
}

func (hcd *HCellData) unmarshal(data []byte) error {
	if len(data) < HCellDataSize {
		return fmt.Errorf("hc data size %d, expected at least %d", len(data), HCellDataSize)
	}
	return binaryReadAll(data, &hcd.BlockSize, &hcd.HCellSignature, &hcd.Metadata)
}

func (hcd *HCellData) marshal() ([]byte, error) {
	return binaryWriteAll(hcd.BlockSize, hcd.HCellSignature, hcd.Metadata)
}

func (hcd *HCellData) setParentHBin(hbin *HBin) {
	hcd.ParentHBin = hbin
}

func (hcd *HCellData) setOffset(offset int32) {
	hcd.ParentHBinOffset = offset
}

func (hcd *HCellData) assertPayloadDataSize(data []byte) error {
	if int32(len(data)) < hcd.Size() {
		return fmt.Errorf("provided data of size %d, expected at least %d", len(data), hcd.Size())
	}
	return nil
}

func (hcd *HCellData) Size() int32 {
	size := hcd.BlockSize
	// if size is >0 then the cell is unallocated
	if hcd.BlockSize < 0 {
		size *= -1
	}
	return size
}

func (hcd *HCellData) Offset() int32 {
	return hcd.ParentHBinOffset
}

func (hcd *HCellData) Signature() string {
	return fmt.Sprintf("%s", hcd.HCellSignature)
}

type OffsetElement int32

type NamedElement struct {
	Offset int32
	Name   [4]byte
}

type DataRecord struct {
	HCellData
	Data []byte
}

func (dr *DataRecord) marshal() ([]byte, error) {
	return binaryWriteAll(dr.BlockSize, dr.Data)
}

func (dr *DataRecord) unmarshal(data []byte) error {
	if err := binaryRead(data, &dr.BlockSize); err != nil {
		return err
	}
	if err := dr.assertPayloadDataSize(data); err != nil {
		return err
	}
	dr.Data = data[HCellSizeLength:dr.Size()]
	return nil
}
