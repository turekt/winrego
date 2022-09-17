package block

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

const (
	BaseBlockSize = 4096
)

var (
	ErrInvalidBlock = errors.New("block data is empty or not properly aligned")
)

func Marshal(rb RegistryBlock) ([]byte, error) {
	return rb.marshal()
}

func Unmarshal(rb RegistryBlock, data []byte) error {
	return rb.unmarshal(data)
}

func UnmarshalBlock(data []byte) (RegistryBlock, error) {
	if len(data) == 0 || len(data)%8 != 0 {
		return nil, ErrInvalidBlock
	}

	var rb RegistryBlock
	switch string(data[0:4]) {
	case "regf":
		rb = &BaseBlock{}
	case "hbin":
		rb = &HBin{}
	default:
		return UnmarshalHCell(data)
	}

	err := Unmarshal(rb, data)
	return rb, err
}

func UnmarshalHCell(data []byte) (HCell, error) {
	if len(data) == 0 || len(data)%8 != 0 {
		return nil, ErrInvalidBlock
	}
	hcType := string(data[4:6])

	var hc HCell
	switch hcType {
	case "li":
		hc = &IndexLeaf{}
	case "lf":
		hc = &FastLeaf{}
	case "lh":
		hc = &HashLeaf{}
	case "ri":
		hc = &IndexRoot{}
	case "nk":
		hc = &KeyNode{}
	case "vk":
		hc = &KeyValue{}
	case "sk":
		hc = &KeySecurity{}
	case "db":
		hc = &BigData{}
	default:
		hc = &DataRecord{}
	}

	err := hc.unmarshal(data)
	return hc, err
}

type RegistryBlock interface {
	marshal() ([]byte, error)
	unmarshal(data []byte) error
	Size() int32
	Offset() int32
	Signature() string
}

type BaseBlock struct {
	RegfHeader       uint32
	Sequence1        uint32
	Sequence2        uint32
	LastWTimestamp   uint64 // LastWrittenTimestamp
	Major            uint32
	Minor            uint32
	FileType         uint32
	FileFormat       uint32
	RootCellOffset   uint32
	HBinSize         uint32
	ClusteringFactor uint32
	FileName         [64]byte
	RmId             [16]byte
	LogId            [16]byte
	Flags            uint32
	TmId             [16]byte
	GUIDSignature    uint32
	LastRTimestamp   uint64 // LastReorganizedTimestamp
	Reserved1        [332]byte
	Checksum         uint32
	Reserved2        [3528]byte
	ThawTmId         [16]byte
	ThawRmId         [16]byte
	ThawLogId        [16]byte
	BootType         uint32
	BootRecover      uint32
}

func (b *BaseBlock) marshal() ([]byte, error) {
	return binaryWrite(b)
}

func (b *BaseBlock) unmarshal(data []byte) error {
	if len(data) < BaseBlockSize {
		return fmt.Errorf("base block data size is %d, expected at least %d", len(data), BaseBlockSize)
	}
	return binaryRead(data, b)
}

func (b *BaseBlock) Size() int32 {
	return BaseBlockSize
}

func (b *BaseBlock) Offset() int32 {
	return 0
}

func (b *BaseBlock) Signature() string {
	return string(uint32toba(b.RegfHeader))
}

func ParseFiletime(ft uint64) time.Time {
	// From https://github.com/Velocidex/regparser/blob/8e74df808b0a4609952bbf8643cbb3bab6b6a438/helpers.go#L17-L19
	return time.Unix(int64(((ft - 11644473600000*10000) / 10000000)), 0)
}

func binaryRead(data []byte, s any) error {
	reader := bytes.NewReader(data)
	return binaryBufferRead(reader, s)
}

func binaryReadAll(data []byte, s ...any) error {
	reader := bytes.NewReader(data)
	for _, target := range s {
		if err := binaryBufferRead(reader, target); err != nil {
			return err
		}
	}
	return nil
}

func binaryBufferRead(reader *bytes.Reader, s any) error {
	return binary.Read(reader, binary.LittleEndian, s)
}

func binaryWrite(s any) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binaryBufferWrite(buf, s)
	return buf.Bytes(), err
}

func binaryWriteAll(s ...any) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, target := range s {
		if err := binaryBufferWrite(buf, target); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func binaryBufferWrite(buf *bytes.Buffer, s any) error {
	return binary.Write(buf, binary.LittleEndian, s)
}

func uint32toba(u uint32) []byte {
	return []byte{uint8(u >> 24), uint8(u >> 16), uint8(u >> 8), uint8(u & 0xff)}
}
