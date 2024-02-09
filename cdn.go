package main

import (
	"flag"
	"fmt"

	"bsuu.eu/riot-cdn/internal/riotdrgaon"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func main() {

	download := flag.Bool("download", true, "Download the data from the Riot API")
	path := flag.String("path", "data/", "Path to the data")
	cacheSize := flag.Int("cache", 50, "Cache size")

	flag.Parse()

	riotDragon := riotdrgaon.NewRiotDragon(&riotdrgaon.RiotDragonConfig{
		Download: *download,
		Path:     *path,
		Cache:    *cacheSize,
	})

	err := riotDragon.ScanLocalVersions()

	if err == nil {
		go riotdrgaon.VersionWorker(riotDragon)
	} else {
		fmt.Println(err)
	}

	app := fiber.New()

	//Middlewares

	app.Use(cache.New())
	app.Get("/metrics", monitor.New(monitor.Config{
		Title:      "CDN Metrics",
		CustomHead: "🤡",
	}))
	app.Use(logger.New())

	// Routes

	cdn := app.Group("/cdn")

	cdn.Get("/versions", riotDragon.GetVersionsHandler)
	cdn.Get("/versions/:version", riotDragon.CacheHandler, riotDragon.GetVersionHandler)

	cdn.Get("/languages", riotDragon.GetLanguagesHandler)

	cdn.Get("/champions/:version", riotDragon.CacheHandler, riotDragon.GetChampionsHandler)
	cdn.Get("/champions/:version/:champion", riotDragon.CacheHandler, riotDragon.GetChampionHandler)

	cdn.Get("/items/:version", riotDragon.CacheHandler, riotDragon.GetItemsHandler)
	cdn.Get("/items/:version/:item", riotDragon.CacheHandler, riotDragon.GetItemHandler)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(418)
	})

	// Static files

	cdn.Static("/images", "./images")

	app.Listen(":8080")
}
