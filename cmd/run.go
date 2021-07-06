package cmd

// Copyright © 2019 oleg2807@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/olegsu/iris/pkg/app"
	"github.com/olegsu/iris/pkg/renderer"
)

var runCmdOptions struct {
	file      string
	kube      string
	inCluster bool
	template  string
	values    []string
}

const (
	rootContext       = "Values"
	defaultOutputFile = "iris-generated.yaml"
)

var runCmd = &cobra.Command{
	Use:  "run",
	Long: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		logger := buildLogger("run")
		if runCmdOptions.kube == "" {
			path := fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))
			logger.Debug("Path to kubeconfig not set, using default", "path", path)
			runCmdOptions.kube = path
		}
		if runCmdOptions.template != "" {
			f, err := os.Open(runCmdOptions.template)
			dieOnError(err)

			writer, err := os.Create(fmt.Sprintf("./%s", defaultOutputFile))
			dieOnError(err)

			valueReaders := make(map[string][]io.Reader)
			for _, valuePath := range runCmdOptions.values {
				var file *os.File
				var err error
				var valueReader io.Reader
				values := strings.Split(valuePath, "=")
				if len(values) == 1 {
					file, err = os.Open(values[0])
					dieOnError(err)
					if strings.HasSuffix(values[0], ".json") {
						valueReader, err = jsonToYaml(file)
						dieOnError(err)
					} else {
						valueReader = file
					}
					if _, ok := valueReaders[rootContext]; !ok {
						valueReaders[rootContext] = []io.Reader{valueReader}
					} else {
						valueReaders[rootContext] = append(valueReaders[rootContext], valueReader)
					}
				} else {
					file, err = os.Open(values[1])
					dieOnError(err)
					if strings.HasSuffix(values[1], ".json") {
						valueReader, err = jsonToYaml(file)
						dieOnError(err)
					} else {
						valueReader = file
					}
					if _, ok := valueReaders[values[0]]; !ok {
						valueReaders[values[0]] = []io.Reader{valueReader}
					} else {
						valueReaders[values[0]] = append(valueReaders[values[0]], valueReader)
					}
				}
				defer file.Close()
			}

			res, err := renderer.New(&renderer.Options{
				TemplateReaders: map[string]io.Reader{
					runCmdOptions.template: f,
				},
				ValueReaders: valueReaders,
				LeftDelim:    "{{",
				RightDelim:   "}}",
				Name:         path.Base(runCmdOptions.template),
			}).Render()
			dieOnError(err)
			fmt.Fprintln(writer, res.String())
			os.Exit(1)
		}
		config := app.NewApplicationOptions(runCmdOptions.file, runCmdOptions.kube, runCmdOptions.inCluster, logger)
		app.CreateApp(config)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	viper.BindEnv("file", "IRIS_FILE")
	viper.BindEnv("template", "IRIS_TEMPLATE")

	viper.BindEnv("kube", "KUBECONFIG")

	viper.BindEnv("inCluster", "IRIS_IN_CLUSTER")

	runCmd.Flags().StringVar(&runCmdOptions.file, "iris-file", viper.GetString("file"), "Iris yaml config file [$IRIS_FILE]")
	runCmd.Flags().StringVar(&runCmdOptions.template, "iris-template", viper.GetString("template"), "Iris template file [$IRIS_TEMPLATE]")
	runCmd.Flags().StringArrayVar(&runCmdOptions.values, "value", []string{}, "Path to value file (yaml or json)")
	runCmd.Flags().StringVar(&runCmdOptions.kube, "kube-config", viper.GetString("kube"), "Path to kube-config file (default is ~/.kube/config) [$KUBECONFIG]")
	runCmd.Flags().BoolVar(&runCmdOptions.inCluster, "in-cluster", viper.GetBool("inCluster"), "Set when running inside a cluster. NOTE: This option will ignore --kube-config flag [$IRIS_IN_CLUSTER]")

}
