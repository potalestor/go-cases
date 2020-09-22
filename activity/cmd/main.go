package main

import (
	"activity"
	"os"

	"github.com/nikandfor/cli"
)

func main() {
	cli.App = cli.Command{
		Name: "Trade Finances",
		Flags: []*cli.Flag{
			cli.NewFlag("activity_handler", 0, "user activity: 0 - none, 1 - syslog 2 - file. Default: 0"),
			cli.NewFlag("activity_tag", "cbg", "user activity tag. Default: cbg"),
			cli.NewFlag("activity_name", "local2", "syslog tag or filename of activity. Default: local2"),
		},
		Before: func(c *cli.Command) error {
			activity.TheConfig = &activity.Config{
				c.Int("activity_handler"),
				c.String("activity_name"),
				c.String("activity_tag"),
			}
			return nil
		},
	}

	cli.RunAndExit(os.Args)

}
