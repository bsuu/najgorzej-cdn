package riotdrgaon

import (
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (r *RiotDragon) GetLanguagesHandler(c *fiber.Ctx) error {
	return c.JSON(r.Languages)
}

func (r *RiotDragon) GetVersionsHandler(c *fiber.Ctx) error {
	return c.JSON(r.VersionsIds)
}

func (r *RiotDragon) GetVersionHandler(c *fiber.Ctx) error {
	v := c.Locals("version").(*Version)
	return c.JSON(v)
}

func (r *RiotDragon) GetChampionsHandler(c *fiber.Ctx) error {
	v := c.Locals("version").(*Version)
	return c.JSON(v.Champions)
}

func (r *RiotDragon) GetChampionHandler(c *fiber.Ctx) error {
	championId := c.Params("champion")
	chapmionsIds := strings.Split(championId, ",")
	v := c.Locals("version").(*Version)
	champion := v.GetChampion(chapmionsIds...)
	if len(champion) == 0 {
		return c.SendStatus(404)
	}
	if len(champion) == 1 {
		return c.JSON(champion[chapmionsIds[0]])
	}
	return c.JSON(champion)
}

func (r *RiotDragon) GetItemsHandler(c *fiber.Ctx) error {
	v := c.Locals("version").(*Version)
	return c.JSON(v.Items)
}

func (r *RiotDragon) GetItemHandler(c *fiber.Ctx) error {
	itemId := c.Params("item")
	itemsIds := strings.Split(itemId, ",")
	v := c.Locals("version").(*Version)
	item := v.GetItem(itemsIds...)
	if len(item) == 0 {
		return c.SendStatus(404)
	}
	if len(item) == 1 {
		return c.JSON(item[itemsIds[0]])
	}
	return c.JSON(item)
}

func (r *RiotDragon) GetSummonersHandler(c *fiber.Ctx) error {
	v := c.Locals("version").(*Version)
	return c.JSON(v.Summoners)
}

func (r *RiotDragon) GetSummonerHandler(c *fiber.Ctx) error {
	summonerId := c.Params("summoner")
	summonersIds := strings.Split(summonerId, ",")
	v := c.Locals("version").(*Version)
	summoner := v.GetSummoner(summonersIds...)
	if len(summoner) == 0 {
		return c.SendStatus(404)
	}
	if len(summoner) == 1 {
		return c.JSON(summoner[summonersIds[0]])
	}
	return c.JSON(summoner)
}

func (r *RiotDragon) GetRunesHandler(c *fiber.Ctx) error {
	v := c.Locals("version").(*Version)
	return c.JSON(v.Runes)
}

func (r *RiotDragon) GetRuneHandler(c *fiber.Ctx) error {
	runeId := c.Params("rune")
	runesIds := strings.Split(runeId, ",")
	v := c.Locals("version").(*Version)
	rune := v.GetRune(runesIds...)
	if len(rune) == 0 {
		return c.SendStatus(404)
	}
	if len(rune) == 1 {
		return c.JSON(rune[runesIds[0]])
	}
	return c.JSON(rune)
}

func (r *RiotDragon) CacheHandler(c *fiber.Ctx) error {
	versionId := c.Params("version")
	v, err := r.GetVersion(versionId)
	if err != nil {
		version, err := r.LoadLocalVersion(r.Config.Path + versionId)
		if err != nil {
			return c.SendStatus(404)
		}
		c.Locals("version", version)
		r.Versions = append(r.Versions, version)
		return c.Next()
	}
	lang := c.Locals("language")
	if lang != nil {
		v = v.ToLanguage(lang.(string))
	}

	c.Locals("version", v)
	return c.Next()
}

func (r *RiotDragon) LanguageHandler(c *fiber.Ctx) error {
	language := c.Query("lang")
	if language == "" {
		return c.Next()
	}
	index := sort.SearchStrings(r.Languages, language)
	if index < len(r.Languages) {
		c.Locals("language", language)
	}
	return c.Next()
}
