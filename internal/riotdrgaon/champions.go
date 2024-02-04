package riotdrgaon

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (r *RiotDragon) GetChampion(versionId string, championId ...string) ([]*Champion, error) {
	version, err := r.GetVersion(versionId)
	if err != nil {
		return nil, err
	}

	champions := make([]*Champion, 0)
	for _, id := range championId {
		if champion, ok := version.Champions[id]; ok {
			champions = append(champions, champion)
		}
	}

	return champions, errors.New("Champion not found")
}

func (r *RiotDragon) GetChampions(versionId string) (map[string]*Champion, error) {
	version, err := r.GetVersion(versionId)
	if err != nil {
		return nil, err
	}
	return version.Champions, nil
}

func (r *RiotDragon) GetChampionsFromRiot(version string) (map[string]*Champion, error) {

	champions := make(map[string]*Champion, 0)
	for _, language := range r.Languages {
		url := "https://ddragon.leagueoflegends.com/cdn/" + version + "/data/" + language + "/championFull.json"

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

		for index, dirtyChampion := range payload["data"].(map[string]interface{}) {
			championMap := dirtyChampion.(map[string]interface{})

			if champion, ok := champions[index]; ok {
				champion.Titles[language] = championMap["title"].(string)
				champion.Blurbs[language] = championMap["blurb"].(string)
				champion.Lores[language] = championMap["lore"].(string)

				if value, ok := championMap["passive"].(map[string]interface{}); ok {
					champion.Passive.Name[language] = value["name"].(string)
					champion.Passive.Description[language] = value["description"].(string)
				}

				if value, ok := championMap["spells"].([]interface{}); ok {
					for index, skill := range value {
						skillMap := skill.(map[string]interface{})
						champion.Skills[index].Names[language] = skillMap["name"].(string)
						champion.Skills[index].Tooltips[language] = skillMap["tooltip"].(string)
						champion.Skills[index].Descriptions[language] = skillMap["description"].(string)
					}
				}

			} else {
				champion := &Champion{
					Id:   championMap["id"].(string),
					Key:  championMap["key"].(string),
					Name: championMap["name"].(string),
					Info: Info{
						Attack:     int(championMap["info"].(map[string]interface{})["attack"].(float64)),
						Defense:    int(championMap["info"].(map[string]interface{})["defense"].(float64)),
						Magic:      int(championMap["info"].(map[string]interface{})["magic"].(float64)),
						Difficulty: int(championMap["info"].(map[string]interface{})["difficulty"].(float64)),
					},
					Stats: Stats{
						Hp:                   championMap["stats"].(map[string]interface{})["hp"].(float64),
						Hpperlevel:           championMap["stats"].(map[string]interface{})["hpperlevel"].(float64),
						Mp:                   championMap["stats"].(map[string]interface{})["mp"].(float64),
						Mpperlevel:           championMap["stats"].(map[string]interface{})["mpperlevel"].(float64),
						Movespeed:            championMap["stats"].(map[string]interface{})["movespeed"].(float64),
						Armor:                championMap["stats"].(map[string]interface{})["armor"].(float64),
						Armorperlevel:        championMap["stats"].(map[string]interface{})["armorperlevel"].(float64),
						Spellblock:           championMap["stats"].(map[string]interface{})["spellblock"].(float64),
						Spellblockperlevel:   championMap["stats"].(map[string]interface{})["spellblockperlevel"].(float64),
						Attackrange:          championMap["stats"].(map[string]interface{})["attackrange"].(float64),
						Hpregen:              championMap["stats"].(map[string]interface{})["hpregen"].(float64),
						Hpregenperlevel:      championMap["stats"].(map[string]interface{})["hpregenperlevel"].(float64),
						Mpregen:              championMap["stats"].(map[string]interface{})["mpregen"].(float64),
						Mpregenperlevel:      championMap["stats"].(map[string]interface{})["mpregenperlevel"].(float64),
						Crit:                 championMap["stats"].(map[string]interface{})["crit"].(float64),
						Critperlevel:         championMap["stats"].(map[string]interface{})["critperlevel"].(float64),
						Attackdamage:         championMap["stats"].(map[string]interface{})["attackdamage"].(float64),
						Attackdamageperlevel: championMap["stats"].(map[string]interface{})["attackdamageperlevel"].(float64),
						Attackspeedperlevel:  championMap["stats"].(map[string]interface{})["attackspeedperlevel"].(float64),
					},
					Titles: map[string]string{language: championMap["title"].(string)},
					Blurbs: map[string]string{language: championMap["blurb"].(string)},
					Lores:  map[string]string{language: championMap["lore"].(string)},
					Skills: make([]*Skill, 0),
					Passive: &Passive{
						Name:        make(map[string]string),
						Description: make(map[string]string),
					},
				}

				if value, ok := championMap["stats"].(map[string]interface{})["attackspeed"].(float64); ok {
					champion.Stats.Attackspeed = value
				}

				if value, ok := championMap["spells"].([]interface{}); ok {
					for _, skill := range value {
						skillMap := skill.(map[string]interface{})
						skill := Skill{
							Id:           skillMap["id"].(string),
							Cooldown:     make([]int, 0),
							Cost:         make([]int, 0),
							Range:        make([]int, 0),
							MaxRank:      int(skillMap["maxrank"].(float64)),
							Names:        make(map[string]string),
							Tooltips:     make(map[string]string),
							Descriptions: make(map[string]string),
						}

						if value, ok := skillMap["cooldown"].([]interface{}); ok {
							for _, v := range value {
								skill.Cooldown = append(skill.Cooldown, int(v.(float64)))
							}
						}

						if value, ok := skillMap["cost"].([]interface{}); ok {
							for _, v := range value {
								skill.Cost = append(skill.Cost, int(v.(float64)))
							}
						}

						if value, ok := skillMap["range"].([]interface{}); ok {
							for _, v := range value {
								skill.Range = append(skill.Range, int(v.(float64)))
							}
						}

						skill.Names[language] = skillMap["name"].(string)
						skill.Tooltips[language] = skillMap["tooltip"].(string)
						skill.Descriptions[language] = skillMap["description"].(string)

						champion.Skills = append(champion.Skills, &skill)
					}
				}

				if value, ok := championMap["passive"].(map[string]interface{}); ok {
					champion.Passive.Name[language] = value["name"].(string)
					champion.Passive.Description[language] = value["description"].(string)
				}

				for _, tag := range championMap["tags"].([]interface{}) {
					champion.Tags = append(champion.Tags, tag.(string))
				}

				champions[index] = champion

			}
		}
	}
	return champions, nil

}

type Champion struct {
	Id      string            `json:"id"`
	Key     string            `json:"key"`
	Name    string            `json:"name"`
	Info    Info              `json:"info"`
	Tags    []string          `json:"tags"`
	Stats   Stats             `json:"stats"`
	Titles  map[string]string `json:"titles"`
	Blurbs  map[string]string `json:"blurbs"`
	Lores   map[string]string `json:"lores"`
	Passive *Passive          `json:"passive"`
	Skills  []*Skill          `json:"skills"`
}

type Passive struct {
	Name        map[string]string `json:"name"`
	Description map[string]string `json:"description"`
}

type Info struct {
	Attack     int `json:"attack"`
	Defense    int `json:"defense"`
	Magic      int `json:"magic"`
	Difficulty int `json:"difficulty"`
}

type Stats struct {
	Hp                   float64 `json:"hp"`
	Hpperlevel           float64 `json:"hpperlevel"`
	Mp                   float64 `json:"mp"`
	Mpperlevel           float64 `json:"mpperlevel"`
	Movespeed            float64 `json:"movespeed"`
	Armor                float64 `json:"armor"`
	Armorperlevel        float64 `json:"armorperlevel"`
	Spellblock           float64 `json:"spellblock"`
	Spellblockperlevel   float64 `json:"spellblockperlevel"`
	Attackrange          float64 `json:"attackrange"`
	Hpregen              float64 `json:"hpregen"`
	Hpregenperlevel      float64 `json:"hpregenperlevel"`
	Mpregen              float64 `json:"mpregen"`
	Mpregenperlevel      float64 `json:"mpregenperlevel"`
	Crit                 float64 `json:"crit"`
	Critperlevel         float64 `json:"critperlevel"`
	Attackdamage         float64 `json:"attackdamage"`
	Attackdamageperlevel float64 `json:"attackdamageperlevel"`
	Attackspeedperlevel  float64 `json:"attackspeedperlevel"`
	Attackspeed          float64 `json:"attackspeed"`
}

type Skill struct {
	Id       string `json:"id"`
	Cooldown []int  `json:"cooldown"`
	Cost     []int  `json:"cost"`
	Range    []int  `json:"range"`
	MaxRank  int    `json:"maxrank"`

	// Moze trzeba dodać effect nie wiem, po co jest ale w sumie nie ma tego w riot dragonie

	Names        map[string]string `json:"names"`
	Descriptions map[string]string `json:"descriptions"`
	Tooltips     map[string]string `json:"tooltips"`
}
