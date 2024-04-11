package expr

import (
	"github.com/new-world-tools/new-world-tools/datasheet"
	"log"
	"strconv"
	"strings"
)

func GetConsumablePotencies(dsStore *datasheet.Store) (map[string]float64, map[string]float64, error) {
	var f64 float64
	var err error
	var ds *datasheet.DataSheet

	files, err := dsStore.GetByType("StatusEffectData")
	if err != nil {
		return nil, nil, err
	}

	statusEffectPotency := map[string]float64{}
	for _, dsFile := range files {
		ds, err = dsFile.GetData()
		if err != nil {
			return nil, nil, err
		}

		for _, row := range ds.Rows {
			id, err := ds.GetCellValueByColumnName(row, "StatusID")
			if err != nil {
				break
			}

			val, err := ds.GetCellValueByColumnName(row, "PotencyPerLevel")
			if err != nil {
				break
			}

			if val != "" {
				f64, err = strconv.ParseFloat(val, 64)
				if err == nil {
					statusEffectPotency[strings.ToLower(id)] = f64
				}
			}
		}
	}

	files, err = dsStore.GetByType("ConsumableItemDefinitions")
	if err != nil {
		return nil, nil, err
	}

	consumablePotency := map[string]float64{}
	for _, dsFile := range files {
		ds, err = dsFile.GetData()
		if err != nil {
			return nil, nil, err
		}

		for _, row := range ds.Rows {
			id, err := ds.GetCellValueByColumnName(row, "ConsumableID")
			if err != nil {
				break
			}

			val, err := ds.GetCellValueByColumnName(row, "AddStatusEffects")
			if err != nil {
				break
			}

			statusEffects := strings.Split(val, "+")
			i := 0
			for _, statusEffect := range statusEffects {
				f64, ok := statusEffectPotency[strings.ToLower(statusEffect)]
				if ok {
					consumablePotency[strings.ToLower(id)] = f64
					i++
				}
			}
			if i > 1 {
				log.Printf("consumable %q has more than one status effect with potency", id)
			}
		}
	}

	return statusEffectPotency, consumablePotency, nil
}
