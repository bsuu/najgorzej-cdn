package riotdrgaon

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	riot_model "bsuu.eu/riot-cdn/internal/riotdrgaon/model"
)

func (r *RiotDragon) DownloadVersion(version string) (*Version, error) {

	fmt.Printf("\r[%s] Status: Downloading items", version)
	items, err := r.GetItemsFromRiot(version)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\r[%s] Status: Downloading champions, skills", version)
	champions, err := r.GetChampionsFromRiot(version)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\r[%s] Status: Downloading summoners        ", version)
	summoners, err := r.GetSummonersFromRiot(version)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\r[%s] Status: Downloading runes             ", version)
	runes, err := r.GetRunesFromRiot(version)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\r[%s] Status: Saving to local files        ", version)

	v := &Version{
		Id:        version,
		Items:     items,
		Champions: champions,
		Summoners: summoners,
		Runes:     runes,
	}

	r.SaveLocalVersions(v)
	fmt.Printf("\r[%s] Status: Download complete            \n", version)
	return v, nil
}

func (r *RiotDragon) DownloadStatic() error {
	seasons, err := GetSeasonsFromRiot()
	if err != nil {
		return err
	}

	queues, err := GetQueuesFromRiot()
	if err != nil {
		return err
	}

	maps, err := GetMapsFromRiot()
	if err != nil {
		return err
	}

	gameModes, err := GetGameModesFromRiot()
	if err != nil {
		return err
	}

	gameTypes, err := GetGameTypesFromRiot()
	if err != nil {
		return err
	}

	r.Static = &Static{
		Seasons:   seasons,
		Queues:    queues,
		Maps:      maps,
		GameModes: gameModes,
		GameTypes: gameTypes,
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (r *RiotDragon) LoadLocalVersions() ([]*Version, []string, error) {
	var notFound []string
	var versions []*Version
	for _, version := range r.VersionsIds {
		if fileExists(r.Config.Path + version) {
			version, err := r.LoadLocalVersion(r.Config.Path + version)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if len(versions) <= r.Config.Cache {
				versions = append(versions, version)
			}

		} else {
			notFound = append(notFound, version)
		}
	}
	return versions, notFound, nil
}

func (r *RiotDragon) LoadLocalVersion(path string) (*Version, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	decoder := gob.NewDecoder(file)

	var version Version

	err = decoder.Decode(&version)

	if err != nil {
		return nil, err
	}

	return &version, nil
}

func (r *RiotDragon) SaveLocalVersions(version *Version) error {
	file, err := os.Create(r.Config.Path + version.Id)

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := gob.NewEncoder(file)
	// encoder := json.NewEncoder(file)

	err = encoder.Encode(version)

	if err != nil {
		return err
	}

	return nil
}

func (r *RiotDragon) GetVersion(id string) (*Version, error) {
	for _, v := range r.Versions {
		if v.Id == id {
			return v, nil
		}
	}
	return nil, errors.New("version not found")
}

func (r *RiotDragon) ExistGameVersion(id string) bool {
	for _, v := range r.VersionsIds {
		if v == id {
			return true
		}
	}
	return false
}

func (r *RiotDragon) GetGameVersions() ([]string, error) {
	url := "https://ddragon.leagueoflegends.com/api/versions.json"
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var versions []string
	err = json.Unmarshal(body, &versions)
	if err != nil {
		return nil, err
	}

	return versions, nil
}

func (r *RiotDragon) FindMissingVersions() []string {
	missingVersions := []string{}
	for _, version := range r.VersionsIds {
		found := false
		for _, localVersion := range r.Versions {
			if version == localVersion.Id {
				found = true
				break
			}
		}
		if !found {
			missingVersions = append(missingVersions, version)
		}
	}
	return missingVersions
}

func (r *RiotDragon) ScanLocalVersions() error {
	versions, notFound, err := r.LoadLocalVersions()

	if err != nil {
		return err
	}

	r.Versions = versions

	if r.Config.Download {
		for _, version := range notFound {
			v, err := r.DownloadVersion(version)
			if err != nil {
				continue
			}
			if len(r.Versions) < r.Config.Cache {
				fmt.Println("Dodaje do cashu wersje: ", version)
				r.Versions = append(r.Versions, v)
			}
		}
	}
	return nil
}

func (r *RiotDragon) GetGameLanguages() ([]string, error) {
	url := "https://ddragon.leagueoflegends.com/cdn/languages.json"
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var languages []string
	err = json.Unmarshal(body, &languages)
	if err != nil {
		return nil, err
	}

	return languages, nil
}

func NewRiotDragon(config *RiotDragonConfig) *RiotDragon {
	riotDragon := &RiotDragon{
		Config: config,
	}
	versionsIds, err := riotDragon.GetGameVersions()

	if err != nil {
		fmt.Println(err)
	}

	riotDragon.VersionsIds = versionsIds

	languages, err := riotDragon.GetGameLanguages()

	if err != nil {
		fmt.Println(err)
	}

	riotDragon.Languages = languages

	err = riotDragon.DownloadStatic()

	if err != nil {
		fmt.Println(err)
	}

	if _, err := os.Stat(config.Path); os.IsNotExist(err) {
		err := os.MkdirAll(config.Path, 0755)
		if err != nil {
			fmt.Println("Error creating directory")
			fmt.Println(err)
		}
	}

	return riotDragon
}

type RiotDragonConfig struct {
	Cache    int
	Path     string
	Download bool
}

type RiotDragon struct {
	VersionsIds []string `json:"versions"`
	Languages   []string `json:"languages"`

	Versions []*Version `json:"gameVersions"`

	Static *Static `json:"static"`

	Config *RiotDragonConfig `json:"config"`
}

type Static struct {
	Seasons   map[int]*riot_model.Season      `json:"seasons"`
	Queues    map[int]*riot_model.Queue       `json:"queues"`
	Maps      map[int]*riot_model.Map         `json:"maps"`
	GameModes map[string]*riot_model.GameMode `json:"gameModes"`
	GameTypes map[string]*riot_model.GameType `json:"gameTypes"`
}

type Version struct {
	Id        string                   `json:"id"`
	Items     map[string]*Item         `json:"items"`
	Champions map[string]*Champion     `json:"champions"`
	Summoners map[string]*Summoner     `json:"summoners"`
	Runes     map[string]*RuneReforged `json:"runes"`
}

func (v *Version) ToLanguage(language string) *Version {
	items := make(map[string]*Item, 0)
	champions := make(map[string]*Champion, 0)
	summoners := make(map[string]*Summoner, 0)
	runes := make(map[string]*RuneReforged, 0)

	for id, item := range v.Items {
		items[id] = item.ToLanguage(language)
	}

	for id, champion := range v.Champions {
		champions[id] = champion.ToLanguage(language)
	}

	for id, summoner := range v.Summoners {
		summoners[id] = summoner.ToLanguage(language)
	}

	for id, rune := range v.Runes {
		runes[id] = rune.ToLanguage(language)
	}

	return &Version{
		Id:        v.Id,
		Items:     items,
		Champions: champions,
		Summoners: summoners,
		Runes:     runes,
	}
}

func (v *Version) GetItem(ids ...string) map[string]*Item {
	items := make(map[string]*Item, 0)
	for _, id := range ids {
		if item, ok := v.Items[id]; ok {
			items[id] = item
		}
	}
	return items
}

func (v *Version) GetChampion(ids ...string) map[string]*Champion {
	champions := make(map[string]*Champion, 0)
	for _, id := range ids {
		if champion, ok := v.Champions[id]; ok {
			champions[id] = champion
		}
	}
	return champions
}

func (v *Version) GetSummoner(ids ...string) map[string]*Summoner {
	summoners := make(map[string]*Summoner, 0)
	for _, id := range ids {
		if summoner, ok := v.Summoners[id]; ok {
			summoners[id] = summoner
		}
	}
	return summoners
}

func (v *Version) GetRune(ids ...string) map[string]*RuneReforged {
	runes := make(map[string]*RuneReforged, 0)
	for _, id := range ids {
		if rune, ok := v.Runes[id]; ok {
			runes[id] = rune
		}
	}
	return runes
}
