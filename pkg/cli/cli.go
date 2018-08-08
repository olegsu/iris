package cli

import (
	"fmt"
	"os"

	"github.com/olegsu/iris/pkg/app"
	"github.com/urfave/cli"
)

func CreateApplication() *cli.App {
	app := cli.NewApp()
	app.Name = "Iris"
	app.Email = "oleg2807@gmail.com"
	app.Description = "Watch on Kubernetes event, filter and send them as standard wehbook you any system"
	app.Version = "0.0.1"
	setupCommands(app)
	return app
}

func setupCommands(app *cli.App) {
	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "Start server",
			Action: run,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "iris-file",
					Usage: "Iris yaml config file",
				},
				cli.StringFlag{
					Name:  "kube-config",
					Usage: "Path to kube-config file",
					Value: fmt.Sprintf("%s/.kube/config", os.Getenv("HOME")),
				},
			},
		},
	}
}

func run(c *cli.Context) error {
	fmt.Println("Started")
	app.CreateApp(c.String("iris-file"), c.String("kube-config"))
	return nil
}
