// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"github.com/riomhaire/jrpcserver/infrastructure/api/rpc"
	"github.com/riomhaire/jrpcserver/model"
	"github.com/riomhaire/jrpcserver/usecases"

	"github.com/spf13/cobra"
)

var configPort int
var configServicename string
var configConsul string
var configHost string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts up a JSON RPC Service",
	Run: func(cmd *cobra.Command, args []string) {

		config := model.APIConfig{
			ServiceName: cmd.Flag("servicename").Value.String(),
			BaseURI:     "/api/v1/rpc",
			Port:        configPort,
			Commands:    usecases.InitializeCommands(),
			Version:     usecases.Version(),
			Consul:      configConsul,
			Hostname:    configHost,
		}

		rpc.StartAPI(config) // Wont return
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&configPort, "port", "p", 3000, "Port to run on")
	serveCmd.Flags().StringVarP(&configServicename, "servicename", "s", "jrpcserver", "Name of service")
	serveCmd.Flags().StringVarP(&configConsul, "consul", "c", "", "Consul host. Empty means dont use")
	serveCmd.Flags().StringVarP(&configHost, "host", "b", "", "Interface/hostname to bind to")

}
