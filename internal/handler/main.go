package handler

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/suzuki-shunsuke/durl/internal/domain"
)

// Main calls a command.
func Main() {
	app := cli.NewApp()
	app.Name = "durl"
	app.Version = domain.Version
	app.Usage = "check whether dead urls are included in files"

	app.Commands = []*cli.Command{
		&initCommand,
		&checkCommand,
	}
	_ = app.Run(os.Args)
}
