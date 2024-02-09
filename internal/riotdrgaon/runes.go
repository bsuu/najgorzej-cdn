package riotdrgaon

import (
	"encoding/json"
	"io"
	"net/http"

	riot_model "bsuu.eu/riot-cdn/internal/riotdrgaon/model"
)

func (r *RiotDragon) GetRunesFromRiot(version string) (map[string]*RuneReforged, error) {

	runesReforged := make(map[string]*RuneReforged, 0)
	for _, language := range r.Languages {
		url := "https://ddragon.leagueoflegends.com/cdn/" + version + "/data/" + language + "/runesReforged.json"

		response, err := http.Get(url)

		if err != nil {
			return nil, err
		}

		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)

		if err != nil {
			return nil, err
		}

		var payload riot_model.RiotReforgedRune

		err = json.Unmarshal(body, &payload)

		if err != nil {
			continue
		}

		for _, dirtyRune := range payload {

			if value, ok := runesReforged[dirtyRune.Key]; ok {
				value.Name[language] = dirtyRune.Name

				for slotIndex, slot := range dirtyRune.Slots {
					for runeIndex, dirtyRune := range slot.Runes {
						value.Slots[slotIndex][runeIndex].Name[language] = dirtyRune.Name
						value.Slots[slotIndex][runeIndex].ShortDesc[language] = dirtyRune.ShortDesc
						value.Slots[slotIndex][runeIndex].LongDesc[language] = dirtyRune.LongDesc
					}
				}
			} else {
				rune := &RuneReforged{
					Id:    dirtyRune.Id,
					Key:   dirtyRune.Key,
					Name:  make(map[string]string),
					Slots: make([][]*Rune, 0),
				}

				rune.Name[language] = dirtyRune.Name

				for _, slot := range dirtyRune.Slots {
					runes := make([]*Rune, 0)
					for _, dirtyRune := range slot.Runes {
						rune := &Rune{
							Id:        dirtyRune.Id,
							Key:       dirtyRune.Key,
							Name:      make(map[string]string),
							ShortDesc: make(map[string]string),
							LongDesc:  make(map[string]string),
						}

						rune.Name[language] = dirtyRune.Name
						rune.ShortDesc[language] = dirtyRune.ShortDesc
						rune.LongDesc[language] = dirtyRune.LongDesc

						runes = append(runes, rune)
					}

					rune.Slots = append(rune.Slots, runes)
				}

				runesReforged[dirtyRune.Key] = rune
			}
		}
	}

	return runesReforged, nil
}

type RuneReforged struct {
	Id    int               `json:"id"`
	Key   string            `json:"key"`
	Name  map[string]string `json:"name"`
	Slots [][]*Rune         `json:"slots"`
}

type Rune struct {
	Id        int               `json:"id"`
	Key       string            `json:"key"`
	Name      map[string]string `json:"name"`
	ShortDesc map[string]string `json:"shortDesc"`
	LongDesc  map[string]string `json:"longDesc"`
}
