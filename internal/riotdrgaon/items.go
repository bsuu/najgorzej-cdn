package riotdrgaon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	riot_model "bsuu.eu/riot-cdn/internal/riotdrgaon/model"
)

func (r *RiotDragon) GetItem(versionId string, id string) (*Item, error) {
	version, err := r.GetVersion(versionId)

	if err != nil {
		return nil, err
	}

	if item, ok := version.Items[id]; ok {
		return item, nil
	}

	return nil, errors.New("item not found")
}

func (r *RiotDragon) GetItems(versionId string) (map[string]*Item, error) {
	version, err := r.GetVersion(versionId)
	if err != nil {
		return nil, err
	}
	return version.Items, nil
}

func (r *RiotDragon) GetItemsFromRiot(version string) (map[string]*Item, error) {

	items := make(map[string]*Item, 0)
	for _, language := range r.Languages {
		url := "https://ddragon.leagueoflegends.com/cdn/" + version + "/data/" + language + "/item.json"

		response, err := http.Get(url)

		if err != nil {
			return nil, err
		}

		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)

		if err != nil {
			return nil, err
		}

		var payload riot_model.RiotItem

		err = json.Unmarshal(body, &payload)

		if err != nil {
			fmt.Println(err)
			continue
		}

		for key, item := range payload.Data {
			if it, ok := items[key]; ok {
				it.Names[language] = item.Name
				it.Descriptions[language] = item.Description
				it.Plaintext[language] = item.Plaintext
			} else {
				items[key] = &Item{
					Id:   key,
					Into: item.Into,
					From: item.From,
					Gold: Gold{
						Base:        item.Gold.Base,
						Total:       item.Gold.Total,
						Sell:        item.Gold.Sell,
						Purchasable: item.Gold.Purchasable,
					},
					Tags:         item.Tags,
					Stats:        item.Stats,
					Descriptions: make(map[string]string, 0),
					Names:        make(map[string]string, 0),
					Plaintext:    make(map[string]string, 0),
					Depth:        item.Depth,
				}

				items[key].Names[language] = item.Name
				items[key].Descriptions[language] = item.Description
				items[key].Plaintext[language] = item.Plaintext
			}

		}

	}

	return items, nil
}

func (i *Item) ToLanguage(language string) *Item {
	return &Item{
		Id:           i.Id,
		Into:         i.Into,
		From:         i.From,
		Gold:         i.Gold,
		Tags:         i.Tags,
		Stats:        i.Stats,
		Descriptions: map[string]string{language: i.Descriptions[language]},
		Names:        map[string]string{language: i.Names[language]},
		Plaintext:    map[string]string{language: i.Plaintext[language]},
		Depth:        i.Depth,
	}
}

type Item struct {
	Id string `json:"id"`

	Into []string `json:"into"` // jak coś to w co item jest budowany
	From []string `json:"from"` // jak coś to z czego item jest budowany

	Gold  Gold               `json:"gold"`  // ile golda kosztuje item
	Tags  []string           `json:"tags"`  // jakie tagi ma item
	Stats map[string]float64 `json:"stats"` // jakie staty daje item

	Descriptions map[string]string `json:"descriptions"`
	Names        map[string]string `json:"names"`
	Plaintext    map[string]string `json:"plaintext"`

	Depth float64 `json:"depth"` // ile itemów jest potrzebnych do zbudowania tego itemu
}

type Gold struct {
	Base        int  `json:"base"`
	Total       int  `json:"total"`
	Sell        int  `json:"sell"`
	Purchasable bool `json:"purchasable"`
}
