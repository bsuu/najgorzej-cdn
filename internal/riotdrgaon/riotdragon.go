package riotdrgaon

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
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
	fmt.Printf("\r[%s] Status: Saving to local files        ", version)

	v := &Version{
		Id:        version,
		Items:     items,
		Champions: champions,
		Summoners: summoners,
	}

	r.SaveLocalVersions(v)
	fmt.Printf("\r[%s] Status: Download complete            \n", version)
	return v, nil
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

	if _, err := os.Stat(config.Path); os.IsNotExist(err) {
		err := os.Mkdir(config.Path, 0755)
		if err != nil {
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

	Config *RiotDragonConfig `json:"config"`
}

type Version struct {
	Id        string               `json:"id"`
	Items     map[string]*Item     `json:"items"`
	Champions map[string]*Champion `json:"champions"`
	Summoners map[string]*Summoner `json:"summoners"`
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
