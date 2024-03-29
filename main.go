package main

import (
	"flag"
	"os"
	"splash/app"
	"splash/services"
	"strconv"
)

func main() {

	serviceLocator := services.NewLocator()
	logger := serviceLocator.Logger()

	configPath := flag.String("config.path", "config/main.json", "Base config.")

	flag.Parse()

	config, err := serviceLocator.LoadConfig(configPath)

	if err != nil {

		logger.Critical(err.Error())
		os.Exit(0)
	}

	mainConfig := parseConfigs(&config)

	dispatcher := app.NewDispatcher(mainConfig)
	dispatcher.Run()

	serviceLocator.BlockIndefinitely()
}

func parseConfigs(config *map[string]interface{}) map[string]map[string]string {

	baseConfig := map[string]map[string]string{}

	for key, value := range *config {
		switch key {
		case "servers":

			data := value.(map[string]interface{})
			for id, server := range data {
				baseConfig[id] = map[string]string{}
				serverData := server.(map[string]interface{})

				for k, base := range serverData {

					switch base.(type) {
					case float64:
						baseConfig[id][k] = strconv.Itoa(int(base.(float64)))
					case string:
						baseConfig[id][k] = base.(string)
						break
					}
				}
			}
			break
		}
	}

	return baseConfig
}
