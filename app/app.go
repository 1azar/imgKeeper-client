package app

import (
	"fmt"
	"imgKeeper-client/service"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	serverAddr string
	filePath   string
	batchSize  int
	method     string
	rootCmd    = &cobra.Command{
		Use:   "transfer_client",
		Short: "Sending/downloading files via gRPC",
		Run: func(cmd *cobra.Command, args []string) {
			clientService := service.New(serverAddr, filePath, batchSize, method)
			if err := clientService.TransferFile(); err != nil {
				log.Fatal(err)
			}

			//switch method {
			//case service.UploadFile:
			//	if err := clientService.SendFile(); err != nil {
			//		log.Fatal(err)
			//	}
			//case service.DownloadFile:
			//	if err := clientService.DownloadFile(); err != nil {
			//		log.Fatal(err)
			//	}
			//case service.ListFiles:
			//	panic("implement me")
			//}
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&serverAddr, "addr", "a", "", "server address")
	rootCmd.Flags().StringVarP(&filePath, "file", "f", "", "file path")
	rootCmd.Flags().StringVarP(&filePath, "method", "m", "none", "method (upload, download, list)")
	rootCmd.Flags().IntVarP(&batchSize, "batch", "b", 1024*1024, "batch size for sending")
	if err := rootCmd.MarkFlagRequired("file"); err != nil {
		log.Fatal(err)
	}
	if err := rootCmd.MarkFlagRequired("addr"); err != nil {
		log.Fatal(err)
	}
	if err := rootCmd.MarkFlagRequired("method"); err != nil {
		log.Fatal(err)
	}
}
