package main

import (
	"flag"
	"net/http"
	httpProtocol "splash/communication/protocols/http"
	"splash/services"
	"splash/client"
	"os"
	"splash/communication"
	"fmt"
)

func main() {

	filePath := flag.String("file.path", "./data/events.csv", "data file path")
	host := flag.String("server.host", "localhost:8090", "Server Host.")
	path := flag.String("server.path", "/server", "Server Path.")

	flag.Parse()

	serviceLocator := services.NewLocator()
	logger := serviceLocator.Logger()

	go serviceLocator.Stats()

	client := client.NewClient(httpProtocol.NewProtocol(&http.Transport{
	}))

	data, err := client.LoadCSVFile(filePath)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(0)
	}

	clientResponse := make(chan *communication.Response, len(data))

	defer close(clientResponse)

	go func() {
		for index, event := range data {
			go func() {
				if index == 100 {
					return
				}
				response, err := client.SendData(&event, host, path)

				if err != nil {
					fmt.Println(index, len(event), err.Error())
				}

				clientResponse <- response
			}()
		}
	}()

	for response := range clientResponse {

		if response != nil {

			fmt.Println(response.Body())
		}
	}
}