package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/baptistegh/lakekeeper",
		Debug:   debug,
	}

	args := os.Args[1:]

	if len(args) > 0 {
		if args[0] == "version" {
			fmt.Printf("version=%s, commit=%s, date=%s\n", version, commit, date)
			return
		}
		log.Fatalf("Command does not exist: %v, the only command accepted is `version`", args)
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
