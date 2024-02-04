package riotdrgaon

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
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

		var payload map[string]interface{}

		err = json.Unmarshal(body, &payload)

		if err != nil {
			continue
		}

		for index, dirtyItem := range payload["data"].(map[string]interface{}) {
			itemMap := dirtyItem.(map[string]interface{})

			if item, ok := items[index]; ok {
				if value, ok := itemMap["name"].(string); ok {
					item.Names[language] = value
				}

				if value, ok := itemMap["description"].(string); ok {
					item.Descriptions[language] = value
				}

				if value, ok := itemMap["plaintext"].(string); ok {
					item.Plaintext[language] = value
				}

			} else {
				// Not existing
				item := &Item{
					Id:           index,
					Stats:        make(map[string]int),
					Names:        make(map[string]string),
					Descriptions: make(map[string]string),
					Plaintext:    make(map[string]string),

					Into: make([]string, 0),
					From: make([]string, 0),
					Tags: make([]string, 0),
				}

				if value, ok := itemMap["stats"].(map[string]interface{}); ok {
					for k, v := range value {
						if stat, ok := v.(float64); ok {
							item.Stats[k] = int(stat)
						}
					}
				}

				if value, ok := itemMap["into"].([]interface{}); ok {
					for _, v := range value {
						if v == nil {
							continue
						}
						item.Into = append(item.Into, v.(string))
					}
				}

				if value, ok := itemMap["from"].([]interface{}); ok {
					for _, v := range value {
						if v == nil {
							continue
						}
						item.From = append(item.From, v.(string))
					}
				}

				if value, ok := itemMap["tags"].([]interface{}); ok {
					for _, v := range value {
						if v == nil {
							continue
						}
						item.Tags = append(item.Tags, v.(string))
					}
				}

				if value, ok := itemMap["depth"]; ok {
					item.Depth = int(value.(float64))
				}

				item.Gold = Gold{
					Base:        int(itemMap["gold"].(map[string]interface{})["base"].(float64)),
					Total:       int(itemMap["gold"].(map[string]interface{})["total"].(float64)),
					Sell:        int(itemMap["gold"].(map[string]interface{})["sell"].(float64)),
					Purchasable: itemMap["gold"].(map[string]interface{})["purchasable"].(bool),
				}

				if value, ok := itemMap["name"].(string); ok {
					item.Names[language] = value
				}

				if value, ok := itemMap["description"].(string); ok {
					item.Descriptions[language] = value
				}

				if value, ok := itemMap["plaintext"].(string); ok {
					item.Plaintext[language] = value
				}

				items[index] = item
			}

		}
	}

	return items, nil
}

type Item struct {
	Id string `json:"id"`

	Into []string `json:"into"` // jak coś to w co item jest budowany
	From []string `json:"from"` // jak coś to z czego item jest budowany

	Gold  Gold           `json:"gold"`  // ile golda kosztuje item
	Tags  []string       `json:"tags"`  // jakie tagi ma item
	Stats map[string]int `json:"stats"` // jakie staty daje item

	Descriptions map[string]string `json:"descriptions"`
	Names        map[string]string `json:"names"`
	Plaintext    map[string]string `json:"plaintext"`

	Depth int `json:"depth"` // ile itemów jest potrzebnych do zbudowania tego itemu
}

type Gold struct {
	Base        int  `json:"base"`
	Total       int  `json:"total"`
	Sell        int  `json:"sell"`
	Purchasable bool `json:"purchasable"`
}
