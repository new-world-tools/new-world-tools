# New World tools

## Downloads

Download the compiled binaries in [Releases](https://github.com/new-world-tools/new-world-tools/releases)

## Usage

### Pak extracter

Powershell:
```powershell
.\pak-extracter.exe `
    -input "C:\Program Files (x86)\Steam\steamapps\common\New World\assets" `
    -output ".\extract"
```
or
```powershell
.\pak-extracter.exe `
    -input "C:\Program Files (x86)\Steam\steamapps\common\New World\assets\server\server.pak" `
    -output ".\extract"
```

Optional:
```powershell
    -threads 3 `
    -decompress-azcs `
    -fix-luac `
    -hash ".\extract\files.sha1"
```

Optional regexp filters:
```powershell
    -include "\.datasheet$" `
    -exclude "^coatgen|^slices" `
    -include-priority
```

### Datasheet converter

Supported formats are `csv` (default), `json` and `yaml`

Powershell:
```powershell
.\datasheet-converter.exe `
    -input ".\extract\sharedassets\springboardentitites\datatables" `
    -output ".\extract\datasheets" `
    -format csv
```

Optional:
```powershell
    -threads 3 `
    -localization ".\extract\localization\en-us" ^
    -keep-structure
```

### Object stream converter

Converts object streams (slices, timelines, various .*db) to json. Supports `Amazon Compressed Stream` (AZCS) by default.

Powershell:
```powershell
.\object-stream-converter.exe `
    -input ".\extract\slices" `
    -output ".\extract\objects-streams"
```

Optional:
```powershell
    -with-indents ^
    -indents-size 2 ^
    -threads 3
```
