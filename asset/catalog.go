package asset

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/new-world-tools/new-world-tools/reader"
	"io"
	"strings"
)

type AssetCatalog struct {
	Signature                            []byte
	Version                              uint32
	FileSize1                            uint32
	Field4                               uint32
	GuidBlockOffset                      uint32
	TypeBlockOffset                      uint32
	DirBlockOffset                       uint32
	FileNameBlockOffset                  uint32
	FileSize2                            uint32
	AssetIdToInfoNumEntries              uint32
	AssetIdToInfo                        []*AssetIdToInfoRef
	Unknown1                             uint32 // AssetDependenciesNumEntries?
	AssetPathToIdNumEntries              uint32
	AssetPathToId                        []*AssetPathToIdRef
	LegacyAssetIdToRealAssetIdNumEntries uint32
	LegacyAssetIdToRealAssetId           []*LegacyAssetIdToRealAssetIdRef
}

type AssetId struct {
	Guid  string
	SubId uint32
}

func (assetId *AssetId) String() string {
	return fmt.Sprintf("{%s}:%x", strings.ToUpper(assetId.Guid), assetId.SubId)
}

type AssetInfo struct {
	AssetId      *AssetId
	RelativePath string
	SizeBytes    uint32
	AssetType    string
}

type AssetIdToInfoRef struct {
	Guid1Index     uint32
	SubId1         uint32
	Guid2Index     uint32
	SubId2         uint32
	TypeIndex      uint32
	Field6         uint32
	FileSize       uint32
	Field8         uint32
	DirOffset      uint32
	FileNameOffset uint32
}

func (ref *AssetIdToInfoRef) Load(rs io.ReadSeeker, cat *AssetCatalog) (*AssetInfo, error) {
	if ref.Guid1Index != ref.Guid2Index {
		return nil, fmt.Errorf("indexes do not match: %d and %d", ref.Guid1Index, ref.Guid2Index)
	}
	if ref.SubId1 != ref.SubId2 {
		return nil, fmt.Errorf("subIds do not match: %d and %d", ref.SubId1, ref.SubId2)
	}

	rs.Seek(int64(cat.GuidBlockOffset+16*ref.Guid2Index), io.SeekStart)
	data, err := reader.ReadBytes(rs, 16)
	if err != nil {
		return nil, err
	}
	u, err := uuid.FromBytes(data)
	if err != nil {
		return nil, err
	}
	guid := u.String()

	rs.Seek(int64(cat.TypeBlockOffset+16*ref.TypeIndex), io.SeekStart)
	data, err = reader.ReadBytes(rs, 16)
	if err != nil {
		return nil, err
	}
	u, err = uuid.FromBytes(data)
	if err != nil {
		return nil, err
	}
	typ := u.String()

	rs.Seek(int64(cat.DirBlockOffset+ref.DirOffset), io.SeekStart)
	str, err := reader.ReadNullTerminatedString(rs)
	if err != nil {
		return nil, err
	}
	dir := str

	rs.Seek(int64(cat.FileNameBlockOffset+ref.FileNameOffset), io.SeekStart)
	str, err = reader.ReadNullTerminatedString(rs)
	if err != nil {
		return nil, err
	}
	fileName := str

	return &AssetInfo{
		AssetId: &AssetId{
			Guid:  guid,
			SubId: ref.SubId2,
		},
		RelativePath: dir + fileName,
		SizeBytes:    ref.FileSize,
		AssetType:    typ,
	}, nil
}

type AssetPathToId struct {
	Uuid    string
	AssetId *AssetId
}

type AssetPathToIdRef struct {
	UuidIndex uint32
	GuidIndex uint32
	SubId     uint32
}

func (ref *AssetPathToIdRef) Load(rs io.ReadSeeker, cat *AssetCatalog) (*AssetPathToId, error) {
	rs.Seek(int64(cat.GuidBlockOffset+16*ref.UuidIndex), io.SeekStart)
	data, err := reader.ReadBytes(rs, 16)
	if err != nil {
		return nil, err
	}
	u, err := uuid.FromBytes(data)
	if err != nil {
		return nil, err
	}
	uid := u.String()

	rs.Seek(int64(cat.GuidBlockOffset+16*ref.GuidIndex), io.SeekStart)
	data, err = reader.ReadBytes(rs, 16)
	if err != nil {
		return nil, err
	}
	u, err = uuid.FromBytes(data)
	if err != nil {
		return nil, err
	}
	guid := u.String()

	return &AssetPathToId{
		Uuid: uid,
		AssetId: &AssetId{
			Guid:  guid,
			SubId: ref.SubId,
		},
	}, nil
}

type LegacyAssetIdToRealAssetIdRef struct {
	LegacyGuidIndex uint32
	LegacySubId     uint32
	RealGuidIndex   uint32
	RealSubId       uint32
}

func (ref *LegacyAssetIdToRealAssetIdRef) Load(rs io.ReadSeeker, cat *AssetCatalog) (*LegacyAssetIdToRealAssetId, error) {
	rs.Seek(int64(cat.GuidBlockOffset+16*ref.LegacyGuidIndex), io.SeekStart)
	data, err := reader.ReadBytes(rs, 16)
	if err != nil {
		return nil, err
	}
	u, err := uuid.FromBytes(data)
	if err != nil {
		return nil, err
	}
	legacyGuid := u.String()

	rs.Seek(int64(cat.GuidBlockOffset+16*ref.RealGuidIndex), io.SeekStart)
	data, err = reader.ReadBytes(rs, 16)
	if err != nil {
		return nil, err
	}
	u, err = uuid.FromBytes(data)
	if err != nil {
		return nil, err
	}
	realGuid := u.String()

	return &LegacyAssetIdToRealAssetId{
		LegacyAssetId: &AssetId{
			Guid:  legacyGuid,
			SubId: ref.LegacySubId,
		},
		RealAssetId: &AssetId{
			Guid:  realGuid,
			SubId: ref.RealSubId,
		},
	}, nil
}

type LegacyAssetIdToRealAssetId struct {
	LegacyAssetId *AssetId
	RealAssetId   *AssetId
}

var signature = []byte("RASC")

func ParseAssetCatalog(r io.Reader) (*AssetCatalog, error) {
	cat := &AssetCatalog{}

	var data []byte
	var u32 uint32
	var err error

	data, err = reader.ReadBytes(r, 4)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(data, signature) {
		return nil, errors.New("wrong signature")
	}
	cat.Signature = data

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cat.Version = u32

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cat.FileSize1 = u32

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cat.Field4 = u32

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cat.GuidBlockOffset = u32

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cat.TypeBlockOffset = u32

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cat.DirBlockOffset = u32

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cat.FileNameBlockOffset = u32

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cat.FileSize2 = u32

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cat.AssetIdToInfoNumEntries = u32

	cat.AssetIdToInfo = make([]*AssetIdToInfoRef, cat.AssetIdToInfoNumEntries)

	for i := uint32(0); i < cat.AssetIdToInfoNumEntries; i++ {
		datum := &AssetIdToInfoRef{}

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		datum.Guid1Index = u32

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		datum.SubId1 = u32

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		datum.Guid2Index = u32

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		datum.SubId2 = u32

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		datum.TypeIndex = u32

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		datum.Field6 = u32

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		datum.FileSize = u32

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		datum.Field8 = u32

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		datum.DirOffset = u32

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		datum.FileNameOffset = u32

		cat.AssetIdToInfo[i] = datum
	}

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cat.Unknown1 = u32

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cat.AssetPathToIdNumEntries = u32

	cat.AssetPathToId = make([]*AssetPathToIdRef, cat.AssetPathToIdNumEntries)

	for i := uint32(0); i < cat.AssetPathToIdNumEntries; i++ {
		datum := &AssetPathToIdRef{}

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		datum.UuidIndex = u32

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		datum.GuidIndex = u32

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		datum.SubId = u32

		cat.AssetPathToId[i] = datum
	}

	u32, err = reader.ReadUint32(r, binary.LittleEndian)
	if err != nil {
		return nil, err
	}
	cat.LegacyAssetIdToRealAssetIdNumEntries = u32

	cat.LegacyAssetIdToRealAssetId = make([]*LegacyAssetIdToRealAssetIdRef, cat.LegacyAssetIdToRealAssetIdNumEntries)

	for i := uint32(0); i < cat.LegacyAssetIdToRealAssetIdNumEntries; i++ {
		datum := &LegacyAssetIdToRealAssetIdRef{}

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		datum.LegacyGuidIndex = u32

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		datum.LegacySubId = u32

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		datum.RealGuidIndex = u32

		u32, err = reader.ReadUint32(r, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		datum.RealSubId = u32

		cat.LegacyAssetIdToRealAssetId[i] = datum
	}

	return cat, nil
}
