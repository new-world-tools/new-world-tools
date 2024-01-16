package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/new-world-tools/new-world-tools/asset"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	inputPtr := flag.String("input", "", "directory or .pak path")
	assetInfoOutputPtr := flag.String("asset-info-output", "./asset-info-output.csv", ".csv path")
	flag.Parse()

	input, err := filepath.Abs(filepath.Clean(*inputPtr))
	if err != nil {
		log.Fatalf("filepath.Abs: %s", err)
	}

	_, err = os.Stat(input)
	if os.IsNotExist(err) {
		log.Fatalf("'%s' does not exist", input)
	}

	assetInfoOutput := *assetInfoOutputPtr

	f, err := os.Open(input)
	if err != nil {
		log.Fatalf("os.Open: %s", err)
	}
	defer f.Close()

	log.Printf("Parsing the catalog...")
	cat, err := asset.ParseAssetCatalog(f)
	if err != nil {
		log.Fatalf("asset.ParseAssetCatalog: %s", err)
	}

	log.Printf("Loading info...")
	assetMap := make(map[string]*asset.AssetInfo, cat.AssetIdToInfoNumEntries)
	for _, ref := range cat.AssetIdToInfo {
		assetInfo, err := ref.Load(f, cat)
		if err != nil {
			log.Fatalf("ref.Load: %s", err)
		}
		assetMap[assetInfo.AssetId.String()] = assetInfo
	}

	keys := make([]string, len(assetMap))
	var i int
	for key, _ := range assetMap {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	log.Printf("Writing a csv...")
	of, err := os.Create(assetInfoOutput)
	if err != nil {
		log.Fatalf("os.Create: %s", err)
	}
	defer of.Close()

	w := csv.NewWriter(of)

	err = w.Write([]string{
		"Guid",
		"SubId",
		"Asset Type",
		"Size (bytes)",
		"Relative Path",
	})
	if err != nil {
		log.Fatalf("w.Write: %s", err)
	}

	var assetInfo *asset.AssetInfo
	for i, key := range keys {
		assetInfo = assetMap[key]
		err = w.Write([]string{
			fmt.Sprintf("%s", assetInfo.AssetId.Guid),
			fmt.Sprintf("%d", assetInfo.AssetId.SubId),
			fmt.Sprintf("%s", assetInfo.AssetType),
			fmt.Sprintf("%d", assetInfo.SizeBytes),
			fmt.Sprintf("%s", assetInfo.RelativePath),
		})
		if err != nil {
			log.Fatalf("w.Write: %s", err)
		}
		if i%1000 == 0 {
			w.Flush()
		}
	}

	w.Flush()
	log.Printf("Finish")
}
