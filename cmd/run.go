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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/olegsu/iris/pkg/app"
)

var runCmdOptions struct {
	file      string
	kube      string
	inCluster bool
}

var runCmd = &cobra.Command{
	Use:  "run",
	Long: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		if runCmdOptions.kube == "" {
			runCmdOptions.kube = fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))
		}
		config := app.NewApplicationOptions(runCmdOptions.file, runCmdOptions.kube, runCmdOptions.inCluster)
		app.CreateApp(config)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	viper.BindEnv("file", "IRIS_FILE")

	viper.BindEnv("kube", "KUBECONFIG")

	viper.BindEnv("inCluster", "KUBECONFIG")

	runCmd.Flags().StringVar(&runCmdOptions.file, "iris-file", viper.GetString("file"), "Iris yaml config file [$IRIS_FILE]")
	runCmd.Flags().StringVar(&runCmdOptions.kube, "kube-config", viper.GetString("kube"), "Path to kube-config file (default is ~/.kube/config) [$KUBECONFIG]")
	runCmd.Flags().BoolVar(&runCmdOptions.inCluster, "in-cluster", viper.GetBool("inCluster"), "Set when running inside a cluster. NOTE: This option will ignore --kube-config flag [$IRIS_IN_CLUSTER]")

}
