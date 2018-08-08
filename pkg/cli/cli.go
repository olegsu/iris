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
					Name:   "kube-config",
					Usage:  "Path to kube-config file",
					EnvVar: "KUBECONFIG",
					Value:  fmt.Sprintf("%s/.kube/config", os.Getenv("HOME")),
				},
				cli.BoolFlag{
					Name:  "in-cluster",
					Usage: "Set when running inside a cluster. NOTE: This option will ignore --kube-config flag",
				},
				cli.StringFlag{
					Name:  "iris-cm",
					Usage: "Name of configmap with iris config yaml. NOTE: This options will ignore --iris-file path",
				},
				cli.StringFlag{
					Name:  "iris-cm-namespace",
					Usage: "Namespace in which to look the configmap",
					Value: "default",
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
		c.String("iris-cm"),
		c.String("iris-cm-namespace"),
	)
	fmt.Println("Started")
	app.CreateApp(config)
	return nil
}
