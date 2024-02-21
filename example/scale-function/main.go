package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/openfaas/faas-provider/types"
	"github.com/openfaas/go-sdk"
)

func main() {
	// NOTE: You can have any name for environment variables. below defined variables names are not standard names
	username := os.Getenv("OPENFAAS_USERNAME")
	password := os.Getenv("OPENFAAS_PASSWORD")

	gatewayURL, _ := url.Parse(os.Getenv("OPENFAAS_GATEWAY_URL"))
	auth := &sdk.BasicAuth{
		Username: username,
		Password: password,
	}

	client := sdk.NewClient(gatewayURL, auth, http.DefaultClient)

	status, err := client.Deploy(context.Background(), types.FunctionDeployment{
		Service:    "env-store-test",
		Image:      "ghcr.io/openfaas/alpine:latest",
		Namespace:  "openfaas-fn",
		EnvProcess: "env",
		Labels: &map[string]string{
			"purpose": "test",
		},
	})

	// non 200 status value will have some error
	if err != nil {
		log.Printf("Status: %d Deploy Failed: %s", status, err)
	}

	fmt.Println("Wait for 15 seconds....")
	time.Sleep(15 * time.Second)
	fn, err := client.GetFunction(context.Background(), "env-store-test", "openfaas-fn")
	if err != nil {
		log.Printf("Get Failed: %s", err)
	}
	fmt.Printf("Function Replica: %d \n", fn.Replicas)

	// scale functions
	err = client.ScaleFunction(context.Background(), "env-store-test", "openfaas-fn", uint64(2))
	// non 200 status value will have some error
	if err != nil {
		log.Printf("Scale Failed: %s", err)
	}

	fmt.Println("Wait for 15 seconds....")
	time.Sleep(15 * time.Second)
	fn, err = client.GetFunction(context.Background(), "env-store-test", "openfaas-fn")
	if err != nil {
		log.Printf("Get Failed: %s", err)
	}
	fmt.Printf("Function Replica: %d \n", fn.Replicas)

	// delete function
	err = client.DeleteFunction(context.Background(), "env-store-test", "openfaas-fn")
	// non 200 status value will have some error
	if err != nil {
		log.Printf("Delete Failed: %s", err)
	}
}
