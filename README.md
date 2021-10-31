# New World tools

## Usage

### Pak extracter

```powershell
.\pak-extracter.exe ^
    -assets "C:\Program Files (x86)\Steam\steamapps\common\New World\assets" ^
    -output ".\extract" ^
    -filter .ext1,.ext2
```

### Datasheet converter

Supported formats are `csv` (default) and `json`

```powershell
.\datasheet-converter.exe ^
    -input ".\extract" ^
    -output ".\extract\datasheets" ^
    -format csv
```
