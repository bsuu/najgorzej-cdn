package riotdrgaon

import (
	"encoding/json"
	"io"
	"net/http"

	riot_model "bsuu.eu/riot-cdn/internal/riotdrgaon/model"
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

		var payload riot_model.RiotSummoner

		err = json.Unmarshal(body, &payload)

		if err != nil {
			continue
		}

		for key, summoner := range payload.Data {
			if summ, ok := summoners[key]; ok {
				summ.Names[language] = summoner.Name
				summ.Description[language] = summoner.Description
				summ.Tooltip[language] = summoner.Tooltip
			} else {

				summoners[key] = &Summoner{
					Id:            summoner.Id,
					MaxRank:       summoner.Maxrank,
					Cost:          int(summoner.Cost[0]),
					Cooldown:      int(summoner.Cooldown[0]),
					Effect:        summoner.EffectBurn,
					Vars:          summoner.EffectBurn,
					Key:           summoner.Key,
					SummonerLevel: summoner.Summonerlevel,
					Names:         make(map[string]string, 0),
					Description:   make(map[string]string, 0),
					Tooltip:       make(map[string]string, 0),
				}

				summoners[key].Names[language] = summoner.Name
				summoners[key].Description[language] = summoner.Description
				summoners[key].Tooltip[language] = summoner.Tooltip
			}
		}
	}

	return summoners, nil
}

func (s *Summoner) ToLanguage(language string) *Summoner {
	return &Summoner{
		Id:            s.Id,
		MaxRank:       s.MaxRank,
		Cost:          s.Cost,
		Cooldown:      s.Cooldown,
		Effect:        s.Effect,
		Vars:          s.Vars,
		Key:           s.Key,
		SummonerLevel: s.SummonerLevel,
		Names:         map[string]string{language: s.Names[language]},
		Description:   map[string]string{language: s.Description[language]},
		Tooltip:       map[string]string{language: s.Tooltip[language]},
	}
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
