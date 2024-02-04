package riotdrgaon

import (
	"encoding/json"
	"io"
	"net/http"
)

func (R *RiotDragon) GetSummonersFromRiot(version string) (map[string]*Summoner, error) {

	summoners := make(map[string]*Summoner, 0)
	for _, language := range R.Languages {
		url := "https://ddragon.leagueoflegends.com/cdn/" + version + "/data/" + language + "/summoner.json"

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

		for index, dirtySummoner := range payload["data"].(map[string]interface{}) {
			summonerMap := dirtySummoner.(map[string]interface{})

			if value, ok := summoners[index]; ok {
				value.Names[language] = summonerMap["name"].(string)
				value.Description[language] = summonerMap["description"].(string)
				value.Tooltip[language] = summonerMap["tooltip"].(string)

			} else {
				summoner := &Summoner{
					Id:          summonerMap["id"].(string),
					Effect:      nil,
					Vars:        make([]string, 0),
					Key:         summonerMap["key"].(string),
					Names:       make(map[string]string),
					Description: make(map[string]string),
					Tooltip:     make(map[string]string),
				}

				if value, ok := summonerMap["effectBurn"].([]interface{}); ok {
					for _, effect := range value {
						if effect == nil {
							summoner.Effect = append(summoner.Effect, "")
						} else {
							summoner.Effect = append(summoner.Effect, effect.(string))
						}
					}
				}

				if value, ok := summonerMap["vars"].([]interface{}); ok {
					for _, vars := range value {
						varsMap := vars.(map[string]interface{})
						summoner.Vars = append(summoner.Vars, varsMap["key"].(string))
					}
				}

				if value, ok := summonerMap["cost"].([]interface{}); ok {
					summoner.Cost = int(value[0].(float64))
				}

				if value, ok := summonerMap["maxrank"].(float64); ok {
					summoner.MaxRank = int(value)
				}

				if value, ok := summonerMap["cooldown"].([]interface{}); ok {
					summoner.Cooldown = int(value[0].(float64))
				}

				if value, ok := summonerMap["summonerLevel"].(float64); ok {
					summoner.SummonerLevel = int(value)
				}

				summoner.Names[language] = summonerMap["name"].(string)
				summoner.Description[language] = summonerMap["description"].(string)
				summoner.Tooltip[language] = summonerMap["tooltip"].(string)

				summoners[index] = summoner
			}
		}
	}

	return summoners, nil
}

type Summoner struct {
	Id            string   `json:"id"`
	MaxRank       int      `json:"maxRank"`
	Cost          int      `json:"cost"`
	Cooldown      int      `json:"cooldown"`
	Effect        []string `json:"effect"`
	Vars          []string `json:"vars"`
	Key           string   `json:"key"`
	SummonerLevel int      `json:"summonerLevel"`

	Names       map[string]string `json:"names"`
	Description map[string]string `json:"description"`
	Tooltip     map[string]string `json:"tooltip"`
}
