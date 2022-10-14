/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bond/common"
	"client"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"log"
)

type onMetadataPush struct {
	register   func([]byte)
	unregister func([]byte)
}

func (p *onMetadataPush) ServiceRegister(input []byte) {
	log.Println("ServiceRegister")
	p.register(input)
}

func (p *onMetadataPush) ServiceUnRegister(input []byte) {
	p.unregister(input)
}

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "get log from a remote device, then print in terminal",
	Long: `subscribe log from a remote device:

bond console -d <DEVICE-ID>.`,
	Run: func(cmd *cobra.Command, args []string) {
		device, err := cmd.Flags().GetString("device")
		if err != nil {
			fmt.Println(err)
			return
		}
		address, err := cmd.Flags().GetString("address")
		if err != nil {
			fmt.Println(err)
			return
		}
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			fmt.Println(err)
			return
		}
		uuid := uuid.New()
		id := uuid.String()
		cli, err := client.Setup(id, address, port, nil)
		if err != nil {
			log.Println("connect error: " + err.Error())
			panic(err)
		}
		err = cli.RegisterConsole()
		if err != nil {
			return
		}

		errorChan := make(chan any)
		cli.OnMetadataPush(&onMetadataPush{
			register: func(input []byte) {
				api, err := common.NewApiServiceFormData(input)
				if err != nil {
					errorChan <- struct{}{}
					return
				}
				if api.Name == common.ConsoleListenerPipeName && api.From == device {
					request := fmt.Sprintf(" {\"user\":\"%s\"}", id)
					response, err := cli.RequestResponse(common.ConsoleListenerPipeName, api.From, []byte(request))
					if err != nil {
						log.Println(err)
						errorChan <- struct{}{}
						return
					}
					log.Println(string(response))
				}
			},
			unregister: func(bytes []byte) {

			},
		})

		err = cli.SoundOutService(device, common.ConsoleListenerPipeName)
		if err != nil {
			return
		}

		defer func(cli *client.Entity) {
			err := cli.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}(cli)
		fmt.Printf("subscribed %s 's log\n", device)
		<-errorChan
	},
}

func init() {

	consoleCmd.Flags().StringP("device", "d", "", "device id")
	consoleCmd.Flags().StringP("address", "a", "192.168.233.52", "address of server")
	consoleCmd.Flags().IntP("port", "p", 7878, "port of server")
	rootCmd.AddCommand(consoleCmd)
}
