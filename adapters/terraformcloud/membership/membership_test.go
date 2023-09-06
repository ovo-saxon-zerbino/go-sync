package membership_test

import (
	"context"
	"log"

	"github.com/hashicorp/go-tfe"

	gosync "github.com/ovotech/go-sync"
	"github.com/ovotech/go-sync/adapters/terraformcloud/membership"
)

func ExampleNew() {
	client, err := tfe.NewClient(&tfe.Config{Token: "my-org-token"})
	if err != nil {
		log.Fatal(err)
	}

	adapter := membership.New(client, "my-org")

	gosync.New(adapter)
}

func ExampleInit() {
	ctx := context.Background()

	adapter, err := membership.Init(ctx, map[gosync.ConfigKey]string{
		membership.Token:        "my-org-token",
		membership.Organisation: "ovotech",
	})
	if err != nil {
		log.Fatal(err)
	}

	gosync.New(adapter)
}
