# WinReGo

A low level library for offline handling of Windows Registry files.

## Motivation

There are a lot of great Windows Registry parsers out there, e.g:
- https://github.com/msuhanov/yarp
- https://github.com/Velocidex/regparser
- https://github.com/williballenthin/python-registry
- https://github.com/libyal/libregf

All of these parsers are doing a great job in providing read access to Windows Registry files with a few differences in implementation and use case coverage.

WinReGo tries to be different by:
- providing ability to load fields to memory in addition to capability of reading data via offsets directly from file
- provide writing ability for all memory block types
- provide with a simple and low level API to build more complex components or libraries
- provide the functionality in pure Go

The library is currently in early stages so breaking changes are possible and additions are expected.

## Resources

This library is being built from documentation and resources provided by others:
- code reference:
  - https://github.com/msuhanov/yarp
  - https://github.com/Velocidex/regparser
  - https://github.com/williballenthin/python-registry
  - https://github.com/libyal/libregf
- documentation:
  - https://github.com/msuhanov/regf
  - https://github.com/libyal/libregf/tree/main/documentation

Testing data and files are available in the referenced resource repositories:
- https://github.com/libyal/libregf/tree/main/tests/data
- https://github.com/msuhanov/yarp/tree/master/hives_for_manual_tests
- https://github.com/msuhanov/yarp/tree/master/hives_for_tests
- https://github.com/msuhanov/yarp/tree/master/records_for_tests
- https://github.com/Velocidex/regparser/tree/master/testdata
- https://github.com/williballenthin/python-registry/tree/master/tests/reg_samples
- https://github.com/williballenthin/python-registry/tree/master/testing/reg_samples

Setting up test data used in tests:
```sh
declare -a urls=(
	"https://raw.githubusercontent.com/msuhanov/yarp/master/hives_for_manual_tests/FuseHive"
	"https://raw.githubusercontent.com/msuhanov/yarp/master/hives_for_manual_tests/FuseHive2"
	"https://raw.githubusercontent.com/msuhanov/yarp/master/hives_for_manual_tests/FuseHive3"
	"https://raw.githubusercontent.com/msuhanov/yarp/master/hives_for_manual_tests/FuseHive4"
	"https://raw.githubusercontent.com/libyal/winreg-kb/main/test_data/SAM"
	"https://raw.githubusercontent.com/williballenthin/python-registry/master/testing/reg_samples/new_log_1/SYSTEM"
)
dir="_testdata"
mkdir -p "${dir}"
for url in "${urls[@]}"; do
	wget "${url}" -P "${dir}"
done
```
