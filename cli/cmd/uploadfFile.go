package cmd

import (
	"bond/common"
	"client"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var uploadFileCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload file to a remote device",
	Long:  `upload file to a remote device: bond upload -d <DEVICE-ID> -f <FILE-PATH>`,
	Run: func(cmd *cobra.Command, args []string) {
		device, err := cmd.Flags().GetString("device")
		defer func() {
			if err != nil {
				fmt.Println(err)
				return
			}
		}()
		if err != nil {
			return
		}
		address, err := cmd.Flags().GetString("address")
		if err != nil {
			return
		}
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			return
		}
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			return
		}
		uuid := uuid.New()
		id := uuid.String()
		cli, err := client.Setup(id, address, port, nil)
		if err != nil {
			return
		}
		// 从file中读取数据
		data, err := os.ReadFile(file)
		if err != nil {
			return
		}
		request := client.FileUploadRequest{
			FileInfo: client.FileInfo{
				FileName: filepath.Base(file),
				FileSize: len(data),
			},
			FileData: &data,
		}
		d, err := request.Encode()
		if err != nil {
			return
		}
		fmt.Printf("upload file %s to device %s\n", file, device)
		response, err := cli.RequestResponse(common.UploadFilePipeName, device, d)
		if err != nil {
			return
		}
		res := new(client.FileUploadResponse)
		err = json.Unmarshal(response, &res)
		if err != nil {
			return
		}
		if res.Code == "S0001" {
			fmt.Println("upload success")
			return
		}
		fmt.Println(res.Message)
		return
	},
}

func init() {
	uploadFileCmd.Flags().StringP("device", "d", "", "device id")
	uploadFileCmd.Flags().StringP("address", "a", "192.168.233.52", "address of server")
	uploadFileCmd.Flags().IntP("port", "p", 8080, "port of server")
	rootCmd.AddCommand(uploadFileCmd)
}
