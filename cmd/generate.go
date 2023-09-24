package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"context"

	yaml "github.com/ha36d/infy/pkg/vaultutility"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm/utils"
)

// generateCmd represents the copy command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate tar.gz to kv data",
	Long:  `generate tar.gz to kv data`,
	Run: func(cmd *cobra.Command, args []string) {

		dstaddr := viper.GetString("addr")

		destination, err := yaml.VaultClient(dstaddr, dsttoken)

		if err != nil {
			log.Fatal(err)
		}

		ctx := context.Background()

		resp, err := destination.System.MountsListSecretsEngines(ctx)
		if err != nil {
			log.Fatal(err)
		}

		osPath, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}

		myFile, err := os.Open(generate)
		if err != nil {
			log.Fatalf("unable to read file: %v", err)
		}
		defer myFile.Close()

		if err = targz.Untar(fmt.Sprintf("%s/%s", osPath, "vault-generate"), myFile); err != nil {
			log.Fatalf("unable to untar the file: %v", err)
		}

		versions := []string{"1", "2"}

		for engine, property := range resp.Data {
			engineVersion := "1"
			engineProperty := property.(map[string]interface{})
			if engineProperty["options"] != nil {
				engineOption := engineProperty["options"].(map[string]interface{})
				if engineOption["version"] != nil {
					engineVersion = engineOption["version"].(string)
				}
			}
			if utils.Contains(versions, engineVersion) && engineProperty["type"] == "kv" && (dstengine == "" || utils.Contains(strings.Split(dstengine, ","), strings.TrimSuffix(engine, "/"))) {
				saveSecretToKv(destination, ctx, engine)
			}
		}

		defer os.RemoveAll(fmt.Sprintf("%s/%s", osPath, "vault-generate"))
		log.Println("Job Finished!")

	},
}

func saveSecretToKv(destination *vault.Client, ctx context.Context, engine string) {

	verbose = viper.GetBool("verbose")

	osPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	filepath.Walk(fmt.Sprintf("%s/%s/%s", osPath, "vault-generate", engine), func(path string, info os.FileInfo, err error) error {

		var payload map[string]interface{}

		if err != nil {
			log.Fatalf(err.Error())
		}
		if !info.IsDir() {
			fileName := info.Name()
			content, err := ioutil.ReadFile(path)
			err = json.Unmarshal(content, &payload)
			if err != nil {
				log.Println(err)
			}

			_, err = destination.Secrets.KvV2Write(ctx, fmt.Sprintf("%s%s", strings.TrimPrefix(filepath.Dir(path), fmt.Sprintf("%s/%s/%s", osPath, "vault-generate", strings.TrimSuffix(engine, "/"))), strings.TrimSuffix(fileName, filepath.Ext(fileName))), schema.KvV2WriteRequest{
				Data: payload,
			}, vault.WithMountPath(engine))
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})

}

func init() {

	generateCmd.Flags().StringP("generate", "f", "", "Backup file path")
	viper.BindPFlag("generate", generateCmd.Flags().Lookup("generate"))

	rootCmd.AddCommand(generateCmd)
}
