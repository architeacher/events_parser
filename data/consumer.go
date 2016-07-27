package main

import (
	"flag"
	"net/http"
	httpProtocol "splash/communication/protocols/http"
	"splash/services"
	"splash/client"
	"os"
	"splash/communication"
)

func main() {

	filePath := flag.String("file.path", "./data/events.csv", "data file path")
	host := flag.String("server.host", "localhost:8090", "Server Host.")
	path := flag.String("server.path", "/server", "Server Path.")

	flag.Parse()

	serviceLocator := services.NewLocator()
	logger := serviceLocator.Logger()

	go serviceLocator.Stats()

	client := client.NewClient(httpProtocol.NewProtocol(&http.Client{}))
	data, err := client.LoadCSVFile(filePath)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(0)
	}

	clientResponse := make(chan *communication.Response, len(data))

	go func() {

		for _, event := range data {

			go func() {
				err := client.SendData(&event, host, path, &clientResponse)

				if err != nil {
					logger.Error(err.Error())
				}

				return
			}()
		}

		return
	}()

	for response := range clientResponse {

		if response != nil {

			logger.Debug(response.Body())
		}
	}

	close(clientResponse)
}