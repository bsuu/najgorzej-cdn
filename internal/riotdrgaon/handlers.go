package riotdrgaon

import (
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
	v := c.Locals("version").(*Version)
	champion, err := v.GetChampion(championId)
	if err != nil {
		return c.SendStatus(404)
	}
	return c.JSON(champion)
}

func (r *RiotDragon) GetItemsHandler(c *fiber.Ctx) error {
	v := c.Locals("version").(*Version)
	return c.JSON(v.Items)
}

func (r *RiotDragon) GetItemHandler(c *fiber.Ctx) error {
	itemId := c.Params("item")
	v := c.Locals("version").(*Version)
	item, err := v.GetItem(itemId)
	if err != nil {
		return c.SendStatus(404)
	}
	return c.JSON(item)
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

	c.Locals("version", v)
	return c.Next()
}
