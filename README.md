# New World tools

## Support me

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/zelenin)

## Downloads

Download the compiled binaries in [Releases](https://github.com/new-world-tools/new-world-tools/releases)

### Dependencies

[Microsoft Visual C++ 2015 Redistributable](https://www.microsoft.com/en-us/download/details.aspx?id=52685)

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

### Asset catalog parser

Powershell:
```powershell
.\asset-catalog-parser.exe `
    -input ".\extract\assetcatalog.catalog" `
    -asset-info-output "./asset-info.csv"
```
