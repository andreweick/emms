/*
Copyright © 2021 M. Andrew Eick

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		synchPhotoMetadata()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func synchPhotoMetadata() {
	fmt.Printf("synchPhotoMetadata\n")
}

// func sendCloudflareList() {
// 	cfconfig, err := LoadConfig(".")

// 	if err != nil {
// 		log.Fatalf("Cannont load config: %s", err.Error())
// 	}

// 	// Cloudflare List (GET https://api.cloudflare.com/client/v4/accounts/5930846a5870031c415bb26e42e38833/images/v1?page=777&per_page=10)

// 	// Create client
// 	client := &http.Client{}

// 	// Create request
// 	req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/accounts/"+cfconfig.CloudflareAccountID+"/images/v1?page=777&per_page=10", nil)

// 	if err != nil {
// 		fmt.Printf("Error on call %s\n", err.Error())
// 		panic(err)
// 	}

// 	// Headers
// 	req.Header.Add("Authorization", "Bearer "+cfconfig.CloudflareBearerToken)

// 	parseFormErr := req.ParseForm()
// 	if parseFormErr != nil {
// 		fmt.Println(parseFormErr)
// 	}

// 	// Fetch Request
// 	resp, err := client.Do(req)

// 	if err != nil {
// 		fmt.Println("Failure : ", err)
// 	}

// 	// Read Response Body
// 	respBody, err := ioutil.ReadAll(resp.Body)

// 	if err != nil {
// 		fmt.Printf("Error on call %s\n", err.Error())
// 		panic(err)
// 	}

// 	var imageList CFImages
// 	err = json.Unmarshal(respBody, &imageList)

// 	if err != nil {
// 		log.Printf("Cannot unmarshal json: %s\n", err.Error())
// 	}

// 	// Display Results
// 	fmt.Println("response Body : ", string(respBody))
// }
