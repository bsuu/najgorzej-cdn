package main

import (
	"flag"
	"fmt"
	"os"

	"bsuu.eu/riot-cdn/internal/riotdrgaon"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/utils"
)

func main() {
	download := flag.Bool("download", false, "Download the data from the Riot API")
	path := flag.String("path", "/tmp/cdn/data/", "Path to the data")
	cacheSize := flag.Int("cache", 50, "Cache size")
	port := flag.String("port", "3002", "Port to listen on")

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

	// pprof

	app.Use(pprof.New())

	//Middlewares

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://najgorzej.lol, https://api.najgorzej.lol, https://cdn.najgorzej.lol",
		AllowOriginsFunc: func(origin string) bool {
			return os.Getenv("ENVIRONMENT") == "development"
		},
	}))

	// Initialize default config
	app.Use(compress.New())
	app.Use(cache.New(cache.Config{
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.Path() + c.Query("lang"))
		},
	}))
	app.Use(favicon.New())
	app.Use(logger.New())

	// Routes

	cdn := app.Group("/cdn")

	cdn.Get("/versions", riotDragon.GetVersionsHandler)
	cdn.Get("/versions/:version", riotDragon.LanguageHandler, riotDragon.CacheHandler, riotDragon.GetVersionHandler)

	cdn.Get("/languages", riotDragon.GetLanguagesHandler)

	cdn.Get("/champions/:version", riotDragon.LanguageHandler, riotDragon.CacheHandler, riotDragon.GetChampionsHandler)
	cdn.Get("/champions/:version/:champion", riotDragon.LanguageHandler, riotDragon.CacheHandler, riotDragon.GetChampionHandler)

	cdn.Get("/items/:version", riotDragon.LanguageHandler, riotDragon.CacheHandler, riotDragon.GetItemsHandler)
	cdn.Get("/items/:version/:item", riotDragon.LanguageHandler, riotDragon.CacheHandler, riotDragon.GetItemHandler)

	cdn.Get("/summoners/:version", riotDragon.LanguageHandler, riotDragon.CacheHandler, riotDragon.GetSummonersHandler)
	cdn.Get("/summoners/:version/:summoner", riotDragon.LanguageHandler, riotDragon.CacheHandler, riotDragon.GetSummonerHandler)

	cdn.Get("/runes/:version", riotDragon.LanguageHandler, riotDragon.CacheHandler, riotDragon.GetRunesHandler)
	cdn.Get("/runes/:version/:rune", riotDragon.LanguageHandler, riotDragon.CacheHandler, riotDragon.GetRuneHandler)

	// Routes Static

	static := cdn.Group("/static")

	static.Get("/seasons", riotDragon.StaticSesonsHandler)
	static.Get("/queues", riotDragon.StaticQueueHandler)
	static.Get("/maps", riotDragon.StaticMapHandler)
	static.Get("/gameModes", riotDragon.StaticGameModeHandler)
	static.Get("/gameTypes", riotDragon.StaticGameTypeHandler)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(418)
	})

	// Static files

	cdn.Static("/images", "/tmp/cdn/images/")

	app.Listen(":" + *port)
}
