package main

import (
	"crypto/tls"
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/spf13/cobra"
	"os"
)

func compileNatsOptions() []nats.Option {
	var natsOptions = make([]nats.Option, 0)
	if !DisableTls {
		cert, err := tls.LoadX509KeyPair(NatsCert, NatsKey)
		if err != nil {
			panic(err)
		}
		config := &tls.Config{
			// since NATS backbone should always be on a private line with self-signed certs, we just skip host verification
			InsecureSkipVerify: true,

			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
		natsOptions = append(natsOptions, nats.Secure(config))
	}
	return natsOptions
}

var RootCmd = &cobra.Command{
	Use:   "nats-client",
	Short: "Thin client for NATS server/cluster basic interaction",
	Long:  "Thin client for NATS server/cluster basic interaction",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
	},
}

var NatsUri string
var NatsCert string
var NatsKey string
var DisableTls bool

func main() {
	var testCmd = &cobra.Command{
		Use:   "test",
		Short: "Tests the connectivity to a NATS server/cluster",
		Long:  "Tests the connectivity to a NATS server/cluster",
		Run: func(cmd *cobra.Command, args []string) {
			var natsOptions = compileNatsOptions()
			natsConn, err := nats.Connect(NatsUri, natsOptions...)
			if err != nil {
				panic(err)
			}
			defer natsConn.Close()
		},
	}

	testCmd.Flags().StringVarP(&NatsUri, "nats", "n", "nats://nats:4222", "The nats cluster URI")
	testCmd.Flags().StringVarP(&NatsCert, "natsCert", "", "/etc/ssl/certs/uniform-nats.crt", "The nats cluster TLS certificate file path")
	testCmd.Flags().StringVarP(&NatsKey, "natsKey", "", "/etc/ssl/private/uniform-nats.key", "The nats cluster TLS key file path")
	testCmd.Flags().BoolVar(&DisableTls, "disableTls", false, "A flag indicating if service should disable tls encryption")

	RootCmd.AddCommand(testCmd)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
