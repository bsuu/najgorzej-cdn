package riot_model

type RiotItem struct {
	Type    string                   `json:"type"`
	Version string                   `json:"version"`
	Data    map[string]*RiotItemData `json:"data"`
}

type RiotItemData struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Plaintext   string             `json:"plaintext"`
	Colloq      string             `json:"colloq"`
	Into        []string           `json:"into"`
	From        []string           `json:"from"`
	Tags        []string           `json:"tags"`
	Maps        map[string]bool    `json:"maps"`
	Depth       float64            `json:"depth"`
	Stats       map[string]float64 `json:"stats"`

	Image *RiotImage        `json:"image"`
	Gold  *RiotItemDataGold `json:"gold"`
}

type RiotItemDataGold struct {
	Base        int  `json:"base"`
	Total       int  `json:"total"`
	Sell        int  `json:"sell"`
	Purchasable bool `json:"purchasable"`
}

type RiotSummoner struct {
	Type   string                       `json:"type"`
	Format string                       `json:"format"`
	Data   map[string]*RiotSummonerData `json:"data"`
}

type RiotSummonerData struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Tooltip     string `json:"tooltip"`
	Maxrank     int    `json:"maxrank"`

	Cooldown     []float64   `json:"cooldown"`
	CooldownBurn string      `json:"cooldownBurn"`
	Cost         []float64   `json:"cost"`
	CostBurn     string      `json:"costBurn"`
	Effect       [][]float64 `json:"effect"`
	EffectBurn   []string    `json:"effectBurn"`

	Key           string `json:"key"`
	Summonerlevel int    `json:"summonerLevel"`
}

type RiotChampionFull struct {
	Type    string                   `json:"type"`
	Format  string                   `json:"format"`
	Version string                   `json:"version"`
	Data    map[string]*RiotChampion `json:"data"`
	Keys    map[int]string           `json:"keys"`
}

type RiotChampion struct {
	Id      string `json:"id"`
	Key     string `json:"key"`
	Name    string `json:"name"`
	Title   string `json:"title"`
	Lore    string `json:"lore"`
	Blurb   string `json:"blurb"`
	Partype string `json:"partype"`

	AllyTips  []string `json:"allytips"`
	EnemyTips []string `json:"enemytips"`
	Tags      []string `json:"tags"`

	Stats   *RiotChampionStats   `json:"stats"`
	Info    *RiotChampionInfo    `json:"info"`
	Passive *RiotChampionPassive `json:"passive"`
	Image   *RiotImage           `json:"image"`
	Spells  []*RiotChampionSpell `json:"spells"`
}

type RiotChampionStats struct {
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

type RiotChampionInfo struct {
	Attack     int `json:"attack"`
	Defense    int `json:"defense"`
	Magic      int `json:"magic"`
	Difficulty int `json:"difficulty"`
}

type RiotChampionSpell struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Tooltip     string `json:"tooltip"`
	Maxrank     int    `json:"maxrank"`
	Resource    string `json:"resource"`
	Maxammo     string `json:"maxammo"`

	Cooldown     []float64 `json:"cooldown"`
	CooldownBurn string    `json:"cooldownBurn"`
	Cost         []float64 `json:"cost"`
	CostBurn     string    `json:"costBurn"`
	Range        []float64 `json:"range"`
	RangeBurn    string    `json:"rangeBurn"`

	Effect     [][]float64 `json:"effect"`
	EffectBurn []string    `json:"effectBurn"`

	Image    *RiotImage    `json:"image"`
	Leveltip *RiotLevelTip `json:"leveltip"`
}

type RiotLevelTip struct {
	Label  []string `json:"label"`
	Effect []string `json:"effect"`
}

type RiotChampionPassive struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Image       *RiotImage `json:"image"`
}

type RiotImage struct {
	Full   string `json:"full"`
	Sprite string `json:"sprite"`
	Group  string `json:"group"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	W      int    `json:"w"`
	H      int    `json:"h"`
}

type RiotReforgedRune []struct {
	Id    int            `json:"id"`
	Key   string         `json:"key"`
	Icon  string         `json:"icon"`
	Name  string         `json:"name"`
	Slots []RiotRuneSlot `json:"slots"`
}

type RiotRuneSlot struct {
	Runes []RiotRune `json:"runes"`
}

type RiotRune struct {
	Id        int    `json:"id"`
	Key       string `json:"key"`
	Icon      string `json:"icon"`
	Name      string `json:"name"`
	ShortDesc string `json:"shortDesc"`
	LongDesc  string `json:"longDesc"`
}
