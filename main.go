package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/go-spanapi/v4/apitools"
)

func main() {
	// It's always a good idea to leave authentication tokens out of the source
	// code so we use a command line parameter here.
	token := ""
	collectionID := ""
	flag.StringVar(&token, "token", "", "API token for the Span service")
	flag.StringVar(&collectionID, "collection-id", "", "The collection to query")

	flag.Parse()

	if token == "" {
		fmt.Println("Missing token parameter")
		flag.PrintDefaults()
		return
	}

	config := spanapi.NewConfiguration()

	// Set this to true to list the requests and responses in the client. It can
	// be useful if you are wondering what is happening.
	config.Debug = false

	client := spanapi.NewAPIClient(config)

	// In the Real World this context should also have a context.WithTimeout
	// call to ensure we time out if there's no response.
	ctx, done := apitools.ContextWithAuthAndTimeout(token, 30*time.Second)
	defer done()

	req := client.CollectionsApi.RetrieveCollection(ctx, collectionID)

	collection, _, err := req.Execute()
	if err != nil {
		fmt.Println("Error retrieving collection: ", err.Error())
		return
	}

	fmt.Println("Data from collection ", collection.CollectionId)
	fmt.Println("======================================")

	// This will retrieve the last 10 payloads from the service.
	items, _, err := client.CollectionsApi.
		ListCollectionData(ctx, *collection.CollectionId).
		Limit(10).
		Execute()
	if err != nil {
		fmt.Println("Error retrieving data: ", err.Error())
	}
	for _, data := range items.Data {
		// Print the payload. The payload is base64 encoded so you have decode
		// it if you want the raw data with base64.StdEncoding.DecodeString()
		fmt.Println("Device ID: ", *data.Device.DeviceId, " Payload: ", *data.Payload, " Transport: ", *data.Transport)
	}
}
