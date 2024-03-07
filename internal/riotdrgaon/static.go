package riotdrgaon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	riot_model "bsuu.eu/riot-cdn/internal/riotdrgaon/model"
)

func GetGameTypesFromRiot() (map[string]*riot_model.GameType, error) {
	gameTypes := make(map[string]*riot_model.GameType, 0)

	url := "https://static.developer.riotgames.com/docs/lol/gameTypes.json"

	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var payload []*riot_model.GameType
	err = json.Unmarshal(body, &payload)

	if err != nil {
		return nil, err
	}

	for _, gameType := range payload {
		gameTypes[gameType.GameType] = gameType
	}

	return gameTypes, nil
}

func GetQueuesFromRiot() (map[int]*riot_model.Queue, error) {
	queues := make(map[int]*riot_model.Queue, 0)

	url := "https://static.developer.riotgames.com/docs/lol/queues.json"

	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var payload []*riot_model.Queue
	err = json.Unmarshal(body, &payload)

	if err != nil {
		return nil, err
	}

	for _, queue := range payload {
		queues[queue.QueueId] = queue
	}

	return queues, nil
}

func GetSeasonsFromRiot() (map[int]*riot_model.Season, error) {
	seasons := make(map[int]*riot_model.Season, 0)

	url := "https://static.developer.riotgames.com/docs/lol/seasons.json"

	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var payload []riot_model.Season
	err = json.Unmarshal(body, &payload)
	fmt.Println(payload)

	if err != nil {
		return nil, err
	}

	for _, season := range payload {
		seasons[season.Id] = &season
	}

	return seasons, nil
}

func GetMapsFromRiot() (map[int]*riot_model.Map, error) {
	maps := make(map[int]*riot_model.Map, 0)

	url := "https://static.developer.riotgames.com/docs/lol/maps.json"

	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var payload []*riot_model.Map
	err = json.Unmarshal(body, &payload)

	if err != nil {
		return nil, err
	}

	for _, m := range payload {
		maps[m.MapId] = m
	}

	return maps, nil
}

func GetGameModesFromRiot() (map[string]*riot_model.GameMode, error) {
	gameModes := make(map[string]*riot_model.GameMode, 0)

	url := "https://static.developer.riotgames.com/docs/lol/gameModes.json"

	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var payload []*riot_model.GameMode
	err = json.Unmarshal(body, &payload)

	if err != nil {
		return nil, err
	}

	for _, gameMode := range payload {
		gameModes[gameMode.GameMode] = gameMode
	}

	return gameModes, nil
}
