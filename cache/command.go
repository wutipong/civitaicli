package cache

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "cache",
		Usage: "Managing cache",
		Commands: []*cli.Command{
			{
				Name:  "purge",
				Usage: "purge the cache content",
				Action: func(ctx context.Context, c *cli.Command) error {
					path, err := CacheLocation()
					if err != nil {
						return fmt.Errorf("unable to get cache location: %w", err)
					}

					err = os.RemoveAll(path)

					if err != nil {
						return fmt.Errorf("unable to remove cache directory: %w", err)
					}
					return nil
				},
			},
			{
				Name:  "location",
				Usage: "return cache location",
				Action: func(ctx context.Context, c *cli.Command) error {
					path, err := CacheLocation()
					if err != nil {
						return fmt.Errorf("unable to get cache location: %w", err)
					}

					fmt.Println(path)

					return nil
				},
			},
		},
	}
}
