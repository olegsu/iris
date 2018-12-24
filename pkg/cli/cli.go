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
	app.Description = "Watch on Kubernetes events, filter and send them as standard wehbook to any system"
	app.Version = os.Getenv("IRIS_VERSION")
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
					Name:   "iris-file",
					Usage:  "Iris yaml config file",
					EnvVar: "IRIS_FILE",
				},
				cli.StringFlag{
					Name:   "kube-config",
					Usage:  "Path to kube-config file",
					EnvVar: "KUBECONFIG",
					Value:  fmt.Sprintf("%s/.kube/config", os.Getenv("HOME")),
				},
				cli.BoolFlag{
					Name:   "in-cluster",
					Usage:  "Set when running inside a cluster. NOTE: This option will ignore --kube-config flag",
					EnvVar: "IRIS_IN_CLUSTER",
				},
			},
		},
	}
}

func run(c *cli.Context) error {
	config := app.NewApplicationOptions(
		c.String("iris-file"),
		c.String("kube-config"),
		c.Bool("in-cluster"),
	)
	app.CreateApp(config)
	return nil
}
