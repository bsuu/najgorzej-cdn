package riot_model

type Season struct {
	Id     int    `json:"id"`
	Season string `json:"season"`
}

type Queue struct {
	QueueId     int    `json:"queueId"`
	Map         string `json:"map"`
	Description string `json:"description"`
	Notes       string `json:"notes"`
}

type Map struct {
	MapId   int    `json:"mapId"`
	MapName string `json:"mapName"`
	Notes   string `json:"notes"`
}

type GameMode struct {
	GameMode    string `json:"gameMode"`
	Description string `json:"description"`
}

type GameType struct {
	GameType    string `json:"gameType"`
	Description string `json:"description"`
}
