package riotdrgaon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	riot_model "bsuu.eu/riot-cdn/internal/riotdrgaon/model"
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

		var payload riot_model.RiotChampionFull

		err = json.Unmarshal(body, &payload)

		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, dirtyChampion := range payload.Data {
			if champion, ok := champions[dirtyChampion.Id]; ok {
				champion.Titles[language] = dirtyChampion.Title
				champion.Blurbs[language] = dirtyChampion.Blurb
				champion.Lores[language] = dirtyChampion.Lore
				champion.Passive.Name[language] = dirtyChampion.Passive.Name
				champion.Passive.Description[language] = dirtyChampion.Passive.Description

				for i, spell := range dirtyChampion.Spells {
					champion.Skills[i].Names[language] = spell.Name
					champion.Skills[i].Descriptions[language] = spell.Description
					champion.Skills[i].Tooltips[language] = spell.Tooltip
				}
			} else {

				champion := &Champion{
					Id:     dirtyChampion.Id,
					Key:    dirtyChampion.Key,
					Name:   dirtyChampion.Name,
					Tags:   dirtyChampion.Tags,
					Titles: make(map[string]string),
					Blurbs: make(map[string]string),
					Lores:  make(map[string]string),
					Skills: make([]*Skill, 0),
					Passive: &Passive{
						Name:        make(map[string]string),
						Description: make(map[string]string),
					},
				}

				for _, spell := range dirtyChampion.Spells {

					skill := &Skill{
						Id:           spell.Id,
						Cooldown:     spell.Cooldown,
						Cost:         spell.Cost,
						Range:        spell.Range,
						Maxrank:      spell.Maxrank,
						Names:        make(map[string]string),
						Descriptions: make(map[string]string),
						Tooltips:     make(map[string]string),
					}

					skill.Names[language] = spell.Name
					skill.Descriptions[language] = spell.Description
					skill.Tooltips[language] = spell.Tooltip

					champion.Skills = append(champion.Skills, skill)
				}

				champion.Info = Info{
					Attack:     dirtyChampion.Info.Attack,
					Defense:    dirtyChampion.Info.Defense,
					Magic:      dirtyChampion.Info.Magic,
					Difficulty: dirtyChampion.Info.Difficulty,
				}

				champion.Stats = Stats{
					Hp:                   dirtyChampion.Stats.Hp,
					Hpperlevel:           dirtyChampion.Stats.Hpperlevel,
					Mp:                   dirtyChampion.Stats.Mp,
					Mpperlevel:           dirtyChampion.Stats.Mpperlevel,
					Movespeed:            dirtyChampion.Stats.Movespeed,
					Armor:                dirtyChampion.Stats.Armor,
					Armorperlevel:        dirtyChampion.Stats.Armorperlevel,
					Spellblock:           dirtyChampion.Stats.Spellblock,
					Spellblockperlevel:   dirtyChampion.Stats.Spellblockperlevel,
					Attackrange:          dirtyChampion.Stats.Attackrange,
					Hpregen:              dirtyChampion.Stats.Hpregen,
					Hpregenperlevel:      dirtyChampion.Stats.Hpregenperlevel,
					Mpregen:              dirtyChampion.Stats.Mpregen,
					Mpregenperlevel:      dirtyChampion.Stats.Mpregenperlevel,
					Crit:                 dirtyChampion.Stats.Crit,
					Critperlevel:         dirtyChampion.Stats.Critperlevel,
					Attackdamage:         dirtyChampion.Stats.Attackdamage,
					Attackdamageperlevel: dirtyChampion.Stats.Attackdamageperlevel,
					Attackspeedperlevel:  dirtyChampion.Stats.Attackspeedperlevel,
					Attackspeed:          dirtyChampion.Stats.Attackspeed,
				}

				champion.Titles[language] = dirtyChampion.Title
				champion.Blurbs[language] = dirtyChampion.Blurb
				champion.Lores[language] = dirtyChampion.Lore
				champion.Passive.Name[language] = dirtyChampion.Passive.Name
				champion.Passive.Description[language] = dirtyChampion.Passive.Description

				champions[dirtyChampion.Id] = champion
			}

		}

	}
	return champions, nil

}

func (c *Champion) ToLanguage(language string) *Champion {
	return &Champion{
		Id:     c.Id,
		Key:    c.Key,
		Name:   c.Name,
		Tags:   c.Tags,
		Titles: map[string]string{language: c.Titles[language]},
		Blurbs: map[string]string{language: c.Blurbs[language]},
		Lores:  map[string]string{language: c.Lores[language]},
		Passive: &Passive{
			Name:        map[string]string{language: c.Passive.Name[language]},
			Description: map[string]string{language: c.Passive.Description[language]},
		},
		Info:  c.Info,
		Stats: c.Stats,
		Skills: func() []*Skill {
			skills := make([]*Skill, 0)
			for _, skill := range c.Skills {
				skills = append(skills, skill.ToLanguage(language))
			}
			return skills
		}(),
	}
}

func (s *Skill) ToLanguage(language string) *Skill {
	return &Skill{
		Id:           s.Id,
		Cooldown:     s.Cooldown,
		Cost:         s.Cost,
		Range:        s.Range,
		Maxrank:      s.Maxrank,
		Names:        map[string]string{language: s.Names[language]},
		Descriptions: map[string]string{language: s.Descriptions[language]},
		Tooltips:     map[string]string{language: s.Tooltips[language]},
	}
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
	Id       string    `json:"id"`
	Cooldown []float64 `json:"cooldown"`
	Cost     []float64 `json:"cost"`
	Range    []float64 `json:"range"`
	Maxrank  int       `json:"maxrank"`

	// Moze trzeba dodać effect nie wiem, po co jest ale w sumie nie ma tego w riot dragonie

	Names        map[string]string `json:"names"`
	Descriptions map[string]string `json:"descriptions"`
	Tooltips     map[string]string `json:"tooltips"`
}
