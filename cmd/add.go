/*
Copyright © 2022 M. Andrew Eick

*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// addPhoto returns the Cloudflare UID and an error
func addPhoto(path string) (string, error) {
	CloudflareUID := "0"

	return CloudflareUID, nil
}

// COBRA below!
// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		photoname := args[0]
		cfid, err := addPhoto(photoname)

		if err != nil {
			log.Printf("Error: %s\n", err.Error())
		}

		log.Printf("Uploaded %s UID: %s, AccountID: %s\n", photoname, cfid, emmsCfg.CloudflareAccountID)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
