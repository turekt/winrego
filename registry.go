package winrego

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"

	"github.com/turekt/winrego/block"
)

type RegRModeFlag int

const (
	// Reads file pointer to Registry struct
	ReadFP RegRModeFlag = 1 << iota
	// Reads base block header data to Registry struct
	ReadHeader
	// Reads HBins data (unmarshaled) to Registry struct
	ReadHBins
	// Reads RawHiveData to Registry struct
	ReadHBinsRaw
	// Reads padding data to Registry struct
	ReadRemData
	// Reads and unmarshals all bytes to Registry struct fields
	ReadAllUnmarshal = ReadHeader | ReadHBins | ReadRemData
	// Reads base block header and sets the remaining bytes to raw fields
	ReadAllRaw = ReadHeader | ReadHBinsRaw | ReadRemData
	// Reads absolutely everything
	ReadAll = ReadFP | ReadHeader | ReadHBins | ReadHBinsRaw | ReadRemData
)

type RegWModeFlag int

const (
	_ RegWModeFlag = 1 << iota
	// Writes base block header to bytes
	WriteHeader
	// Writes HBins data (unmarshaled) to bytes
	WriteHBins
	// Writes RawHiveData to bytes
	WriteHBinsRaw
	// Writes padding data to bytes
	WriteRemData
	// Writes the unmarshaled data to bytes, with padding
	WriteAllMarshal = WriteHeader | WriteHBins | WriteRemData
	// Writes base block header and raw data including padding
	WriteAllRaw = WriteHeader | WriteHBinsRaw | WriteRemData
)

type Registry struct {
	// Unmarshaled base block data
	block.BaseBlock
	// Unmarshaled HBin data
	HBins block.HBinData
	// Padding data
	RemnantData []byte
	// File pointer
	File *os.File
	// Raw hive data, base block not included
	RawHiveData []byte
}

func OpenRegistry(filepath string, mode RegRModeFlag) (*Registry, error) {
	r := &Registry{}
	if (mode & ReadFP) != 0 {
		fp, err := os.Open(filepath)
		if err != nil {
			return nil, err
		}
		r.File = fp
	}

	if mode != ReadFP {
		data, err := os.ReadFile(filepath)
		if err != nil {
			return r, err
		}

		if err := r.Load(data, mode); err != nil {
			return r, err
		}
	}

	return r, nil
}

func (r *Registry) Save(filepath string, mode RegWModeFlag) error {
	if filepath == "" {
		if r.File == nil {
			return errors.New("both filepath and reg file pointer are not specified")
		}
		filepath = r.File.Name()
	}

	data, err := r.Bytes(mode)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}

func (r *Registry) Load(data []byte, mode RegRModeFlag) error {
	if len(data) < block.BaseBlockSize {
		return errors.New("not enough data supplied to load reg file header")
	}

	hbSize := binary.LittleEndian.Uint32(data[40:44])
	if (mode & ReadHeader) != 0 {
		if err := block.Unmarshal(&r.BaseBlock, data[:block.BaseBlockSize]); err != nil {
			return err
		}
		hbSize = r.HBinSize
	}

	hbEnd := block.BaseBlockSize + hbSize
	if uint32(len(data)) < hbEnd {
		return errors.New("not enough data supplied to load reg file content")
	}

	if (mode & ReadHBins) != 0 {
		hbins := new(block.HBinData)
		err := block.Unmarshal(hbins, data[block.BaseBlockSize:hbEnd])
		if err != nil {
			return err
		}
		r.HBins = *hbins
	}

	if (mode & ReadHBinsRaw) != 0 {
		r.RawHiveData = data[block.BaseBlockSize:hbEnd]
	}

	if (mode & ReadRemData) != 0 {
		r.RemnantData = data[hbEnd:]
	}

	return nil
}

func (r *Registry) Bytes(mode RegWModeFlag) ([]byte, error) {
	var buf bytes.Buffer

	if (mode & WriteHeader) != 0 {
		b, err := block.Marshal(&r.BaseBlock)
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}

	if (mode & WriteHBins) != 0 {
		b, err := block.Marshal(&r.HBins)
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}

	if (mode & WriteHBinsRaw) != 0 {
		buf.Write(r.RawHiveData)
	}

	if (mode & WriteRemData) != 0 {
		buf.Write(r.RemnantData)
	}

	return buf.Bytes(), nil
}
