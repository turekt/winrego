package winrego

import (
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/turekt/winrego/block"
)

var (
	testResDir = flag.String("testResDir", "_testdata", "directory containing test registry files")
	system     = filepath.Join(*testResDir, "SYSTEM")
	sam        = filepath.Join(*testResDir, "SAM")
	fuseHive1  = filepath.Join(*testResDir, "FuseHive")
	fuseHive2  = filepath.Join(*testResDir, "FuseHive2")
	fuseHive3  = filepath.Join(*testResDir, "FuseHive3")
	fuseHive4  = filepath.Join(*testResDir, "FuseHive4")
)

func testFilesList(t *testing.T) []string {
	testFiles := []string{
		fuseHive1,
		fuseHive2,
		fuseHive3,
		fuseHive4,
		system,
		sam,
	}

	var foundTestFiles []string
	for _, testFile := range testFiles {
		if _, err := os.Stat(testFile); err == nil {
			foundTestFiles = append(foundTestFiles, testFile)
		}
	}
	if len(foundTestFiles) == 0 {
		t.Skip("none of the test files were found")
	}

	return foundTestFiles
}

func TestBaseBlockMarshalCycle(t *testing.T) {
	testFiles := testFilesList(t)

	for _, testFile := range testFiles {
		data, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("failed to open %s: %v", fuseHive1, err)
		}

		bb := &block.BaseBlock{}
		if err := block.Unmarshal(bb, data); err != nil {
			t.Fatalf("failed parsing registry header: %v", err)
		}

		mData, err := block.Marshal(bb)
		if err != nil {
			t.Fatalf("failed marshaling registry header: %v", err)
		}

		if got, want := mData, data[:block.BaseBlockSize]; !reflect.DeepEqual(got, want) {
			t.Fatalf("fs base block != inmem base block; got|want:\n%v\n%v", got, want)
		}
	}
}

func TestReadFlags(t *testing.T) {
	testCases := []struct {
		Mode      RegRModeFlag
		CheckFunc func(r *Registry)
	}{
		{
			ReadFP,
			func(r *Registry) {
				if r.BaseBlock.RegfHeader != 0 {
					t.Errorf("regf header in ReadFP: %v", r.BaseBlock.RegfHeader)
				}
				if r.HBins.Size() != 0 {
					t.Errorf("hbins parsed in ReadFP: %d", r.HBins.Size())
				}
				if r.File == nil {
					t.Errorf("file pointer not loaded in ReadFP")
				}
				if len(r.RawHiveData) != 0 {
					t.Errorf("raw hive data loaded in ReadFP: %d", len(r.RawHiveData))
				}
			},
		},
		{
			ReadHeader,
			func(r *Registry) {
				if r.BaseBlock.RegfHeader != 0x66676572 {
					t.Errorf("regf header in ReadHeader: %v", r.BaseBlock.RegfHeader)
				}
				if r.HBins.Size() != 0 {
					t.Errorf("hbins parsed in ReadHeader: %d", r.HBins.Size())
				}
				if r.File != nil {
					t.Errorf("file pointer loaded in ReadHeader")
				}
				if len(r.RawHiveData) != 0 {
					t.Errorf("raw hive data loaded in ReadHeader: %d", len(r.RawHiveData))
				}
			},
		},
		{
			ReadHBins,
			func(r *Registry) {
				if r.BaseBlock.RegfHeader != 0 {
					t.Errorf("regf header in ReadHBins: %v", r.BaseBlock.RegfHeader)
				}
				if r.HBins.Size() == 0 {
					t.Errorf("hbins not parsed in ReadHBins")
				}
				if r.File != nil {
					t.Errorf("file pointer loaded in ReadHBins")
				}
				if len(r.RawHiveData) != 0 {
					t.Errorf("raw hive data loaded in ReadHBins: %d", len(r.RawHiveData))
				}
			},
		},
		{
			ReadHBinsRaw,
			func(r *Registry) {
				if r.BaseBlock.RegfHeader != 0 {
					t.Errorf("regf header in ReadHBinsRaw: %v", r.BaseBlock.RegfHeader)
				}
				if r.HBins.Size() != 0 {
					t.Errorf("hbins parsed in ReadHBinsRaw: %d", r.HBins.Size())
				}
				if r.File != nil {
					t.Errorf("file pointer loaded in ReadHBinsRaw")
				}
				if len(r.RawHiveData) == 0 {
					t.Errorf("raw hive data not loaded in ReadHBinsRaw")
				}
			},
		},
		{
			ReadRemData,
			func(r *Registry) {
				if r.BaseBlock.RegfHeader != 0 {
					t.Errorf("regf header in ReadRemData: %v", r.BaseBlock.RegfHeader)
				}
				if r.HBins.Size() != 0 {
					t.Errorf("hbins parsed in ReadRemData: %d", r.HBins.Size())
				}
				if r.File != nil {
					t.Errorf("file pointer loaded in ReadRemData")
				}
				if len(r.RawHiveData) != 0 {
					t.Errorf("raw hive data loaded in ReadRemData: %d", len(r.RawHiveData))
				}
			},
		},
		{
			ReadAllUnmarshal,
			func(r *Registry) {
				if r.BaseBlock.RegfHeader != 0x66676572 {
					t.Errorf("regf header in ReadAllUnmarshal: %v", r.BaseBlock.RegfHeader)
				}
				if r.HBins.Size() == 0 {
					t.Errorf("hbins not parsed in ReadAllUnmarshal")
				}
				if r.File != nil {
					t.Errorf("file pointer loaded in ReadAllUnmarshal")
				}
				if len(r.RawHiveData) != 0 {
					t.Errorf("raw hive data loaded in ReadAllUnmarshal: %d", len(r.RawHiveData))
				}
			},
		},
		{
			ReadAllRaw,
			func(r *Registry) {
				if r.BaseBlock.RegfHeader != 0x66676572 {
					t.Errorf("regf header in ReadAllRaw: %v", r.BaseBlock.RegfHeader)
				}
				if r.HBins.Size() != 0 {
					t.Errorf("hbins parsed in ReadAllRaw: %d", r.HBins.Size())
				}
				if r.File != nil {
					t.Errorf("file pointer loaded in ReadAllRaw")
				}
				if len(r.RawHiveData) == 0 {
					t.Errorf("raw hive data not loaded in ReadAllRaw")
				}
			},
		},
		{
			ReadAll,
			func(r *Registry) {
				if r.BaseBlock.RegfHeader != 0x66676572 {
					t.Errorf("regf header in ReadAllRaw: %v", r.BaseBlock.RegfHeader)
				}
				if r.HBins.Size() == 0 {
					t.Errorf("hbins not parsed in ReadAllUnmarshal")
				}
				if r.File == nil {
					t.Errorf("file pointer not loaded in ReadFP")
				}
				if len(r.RawHiveData) == 0 {
					t.Errorf("raw hive data not loaded in ReadAllRaw")
				}
			},
		},
	}

	testFile := testFilesList(t)[0]
	for _, tc := range testCases {
		r, err := OpenRegistry(testFile, tc.Mode)
		if err != nil {
			t.Errorf("failed to open registry: %v", err)
		}

		tc.CheckFunc(r)
	}
}

func TestWriteFlags(t *testing.T) {
	r := &Registry{
		BaseBlock: block.BaseBlock{
			RegfHeader: 0x66676572,
		},
		HBins: []block.HBin{
			block.HBin{
				HBinHeader: block.HBinHeader{
					HBinSignature: 0x6e696268,
				},
				Cells: []block.HCell{
					&block.DataRecord{
						HCellData: block.HCellData{
							BlockSize: 4,
						},
						Data: []byte{0x61, 0x62, 0x63, 0x64},
					},
				},
				HBinDataPtr: nil,
			},
		},
		RemnantData: []byte{
			0x64, 0x63, 0x62, 0x61,
		},
		RawHiveData: []byte{
			0x74, 0x73, 0x65, 0x74,
		},
	}
	testCases := []struct {
		Mode      RegWModeFlag
		CheckFunc func(data []byte)
	}{
		{
			WriteHeader,
			func(data []byte) {
				if len(data) != block.BaseBlockSize {
					t.Errorf("len(data) (%d) != 1024 in WriteHeader", len(data))
				}
				if got, want := string(data[0:4]), "regf"; got != want {
					t.Errorf("regf header not written in WriteHeader mode: %s", got)
				}
			},
		},
		{
			WriteHBins,
			func(data []byte) {
				if len(data) != 40 {
					t.Errorf("len(data) (%d) != 40 in WriteHBins", len(data))
				}
				if got, want := string(data[0:4]), "hbin"; got != want {
					t.Errorf("hbin header not written in WriteHBins mode: %s", got)
				}
				if data[32] != 4 {
					t.Errorf("size %d != 4", data[4])
				}
				if got, want := string(data[36:]), "abcd"; got != want {
					t.Errorf("cell data %s != %s", got, want)
				}
			},
		},
		{
			WriteHBinsRaw,
			func(data []byte) {
				if len(data) != 4 {
					t.Errorf("len(data) (%d) != 4 in WriteHBinsRaw", len(data))
				}
				if got, want := string(data), "tset"; got != want {
					t.Errorf("data %s != %s", got, want)
				}
			},
		},
		{
			WriteRemData,
			func(data []byte) {
				if len(data) != 4 {
					t.Errorf("rem data length %d != 4", len(data))
				}
				if got, want := string(data), "dcba"; got != want {
					t.Errorf("data content not equals, got %s, want %s", got, want)
				}
			},
		},
		{
			WriteAllMarshal,
			func(data []byte) {
				if got, want := len(data), 4140; got != want {
					t.Errorf("WriteAllMarshal data size not equals: got %d, want %d", got, want)
				}
				got, want := data[:4], []byte{
					0x72, 0x65, 0x67, 0x66,
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("regf header not equals: got %s, want %s", got, want)
				}
				got, want = data[block.BaseBlockSize:block.BaseBlockSize+4], []byte{
					0x68, 0x62, 0x69, 0x6e,
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("WriteAllMarshal content not equal (got|want):\n%v\n%v", got, want)
				}
				got, want = data[block.BaseBlockSize+32:], []byte{
					0x04, 0x00, 0x00, 0x00,
					0x61, 0x62, 0x63, 0x64,
					0x64, 0x63, 0x62, 0x61,
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("WriteAllMarshal cell content not equal (got|want):\n%v\n%v", got, want)
				}
			},
		},
		{
			WriteAllRaw,
			func(data []byte) {
				if got, want := len(data), 4104; got != want {
					t.Errorf("WriteAllRaw data size not equals: got %d, want %d", got, want)
				}
				got, want := data[:4], []byte{
					0x72, 0x65, 0x67, 0x66,
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("regf header not equals: got %s, want %s", got, want)
				}
				got, want = data[block.BaseBlockSize:block.BaseBlockSize+8], []byte{
					0x74, 0x73, 0x65, 0x74,
					0x64, 0x63, 0x62, 0x61,
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("WriteAllMarshal content not equal (got|want):\n%v\n%v", got, want)
				}
			},
		},
	}

	for _, tc := range testCases {
		data, err := r.Bytes(tc.Mode)
		if err != nil {
			t.Errorf("error marshaling registry: %v", err)
		}

		tc.CheckFunc(data)
	}
}

func TestHBinHeaderMarshalCycle(t *testing.T) {
	testFiles := testFilesList(t)

	for _, testFile := range testFiles {
		data, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("failed to open %s: %v", fuseHive1, err)
		}
		data = data[block.BaseBlockSize:]

		hb := &block.HBin{}
		if err := block.Unmarshal(hb, data); err != nil {
			t.Fatalf("failed parsing hbin: %v", err)
		}

		mData, err := block.Marshal(hb)
		if err != nil {
			t.Fatalf("failed marshaling hbin: %v", err)
		}

		got, want := mData[:block.HBinHeaderSize], data[:block.HBinHeaderSize]
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("fs hbin header != inmem hbin header; got|want:\n%v\n%v", got, want)
		}
	}
}

func TestRegistryMarshalCycle(t *testing.T) {
	testFiles := testFilesList(t)

	for _, testFile := range testFiles {
		data, err := os.ReadFile(testFile)
		if err != nil {
			t.Errorf("failed to read file %s: %v", testFile, err)
			continue
		}

		r := &Registry{}
		if err := r.Load(data, ReadAllUnmarshal); err != nil {
			t.Errorf("failed parsing registry %s: %v", testFile, err)
			continue
		}

		b, err := r.Bytes(WriteAllMarshal)
		if err != nil {
			t.Fatalf("failed marshaling registry: %v", err)
		}
		if got, want := b, data; !reflect.DeepEqual(got, want) {
			f, err := os.CreateTemp("", "winrego-regtest-*")
			defer f.Close()
			if err != nil {
				t.Fatalf("failed to create temp file for comparison")
			}

			if _, err := f.Write(got); err != nil {
				t.Fatalf("failed to write registry bytes to temp file")
			}

			if err := f.Close(); err != nil {
				t.Fatalf("failed to close temp file")
			}
			t.Errorf("registry bytes not equal for %s, written to %s for comparison", testFile, f.Name())
		}
	}
}

func TestRegistryMarshalHBinContent(t *testing.T) {
	if _, err := os.Stat(fuseHive1); os.IsNotExist(err) {
		t.Skip("FuseHive file was not found")
	}

	r, err := OpenRegistry(fuseHive1, ReadAllUnmarshal)
	if err != nil {
		t.Errorf("failed to read registry file: %v", err)
	}

	expectHBinOffsets := []int32{
		0x00,
		0x1000,
		0x3000,
		0x7000,
		0xb000,
		0xf000,
		0x13000,
		0x17000,
		0x1b000,
		0x1f000,
	}
	for i := 0; i < len(r.HBins); i++ {
		if got, want := r.HBins[i].HBinDataOffset, expectHBinOffsets[i]; got != want {
			t.Errorf("hbin offset mismatch: got %x, want %x", got, want)
		}
	}
}
