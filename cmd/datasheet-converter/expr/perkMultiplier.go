package expr

import (
	"fmt"
	"github.com/new-world-tools/new-world-tools/datasheet"
	"strconv"
	"strings"
)

func GetPerkMultipliers(dsStore *datasheet.Store) (map[string]*Scaling, map[string]float64, error) {
	var err error
	var ds *datasheet.DataSheet

	files, err := dsStore.GetByType("PerkData")
	if err != nil {
		return nil, nil, err
	}

	perkScaling := map[string]*Scaling{}
	for _, dsFile := range files {
		ds, err = dsFile.GetData()
		if err != nil {
			return nil, nil, err
		}

		for _, row := range ds.Rows {
			id, err := ds.GetCellValueByColumnName(row, "PerkID")
			if err != nil {
				break
			}
			scaling, err := ds.GetCellValueByColumnName(row, "ScalingPerGearScore")
			if err != nil {
				break
			}
			sc, _ := ParseScaling(scaling)
			if sc != nil {
				perkScaling[strings.ToLower(id)] = sc
			}
		}
	}

	return perkScaling, nil, nil
}

func ParseScaling(str string) (*Scaling, error) {
	if str == "" {
		return nil, nil
	}

	scalingData := &Scaling{
		Points: make([]*ScalingPoint, 0),
	}

	var err error

	points := strings.Split(str, ",")
	for _, point := range points {
		parts := strings.Split(point, ":")
		if len(parts) > 2 || len(parts) < 1 {
			return nil, fmt.Errorf("%q strange scaling", str)
		}

		var scaling float64
		var gearScore int64
		if len(parts) == 1 {
			scaling, err = strconv.ParseFloat(parts[0], 64)
			if err != nil {
				return nil, fmt.Errorf("strconv.ParseFloat: %s", err)
			}
			gearScore = 100
		}
		if len(parts) == 2 {
			gearScore, err = strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("strconv.ParseInt: %s", err)
			}
			scaling, err = strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return nil, fmt.Errorf("strconv.ParseFloat: %s", err)
			}
		}
		scalingData.Points = append(scalingData.Points, &ScalingPoint{
			GearScore: gearScore,
			Scaling:   scaling,
		})
	}

	return scalingData, nil
}

type Scaling struct {
	Points []*ScalingPoint
}

func (scaling *Scaling) GetScaling(gearScore int64) float64 {
	var result float64 = 1
	for i, point := range scaling.Points {
		if point.GearScore > gearScore {
			continue
		}

		min := point.GearScore
		var max int64
		if len(scaling.Points) == (i + 1) {
			max = gearScore
		} else {
			max = scaling.Points[i+1].GearScore
		}
		result = result + float64(max-min)*point.Scaling
	}

	return result
}

type ScalingPoint struct {
	GearScore int64
	Scaling   float64
}
