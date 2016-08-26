package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"splash/client"
	"splash/communication"
	httpProtocol "splash/communication/protocols/http"
	"splash/services"
)

func main() {

	filePath := flag.String("file.path", "./data/events.csv", "data file path")
	host := flag.String("server.host", "localhost:8090", "Server Host.")
	path := flag.String("server.path", "/server", "Server Path.")

	flag.Parse()

	// Todo: This should be requested from the server
	authorizationToken := "princing"

	serviceLocator := services.NewLocator()
	logger := serviceLocator.Logger()

	client := client.NewClient(httpProtocol.NewProtocol(&http.Transport{}))

	data, err := client.LoadCSVFile(filePath)

	patchesCount := len(data)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(0)
	}

	clientResponse := make(chan *communication.Response, len(data))

	defer close(clientResponse)

	go func() {
		for index, event := range data {

			response, err := client.SendData(index, &event, host, path, &authorizationToken, patchesCount)

			if err != nil {
				fmt.Println(index, len(event), err.Error())
			}

			clientResponse <- response
		}
	}()

	for response := range clientResponse {
		if response != nil {
			fmt.Println(response.Body())
		}
	}
}
