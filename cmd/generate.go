package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// generateCmd represents the copy command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the infrastructure file",
	Long:  `Generate the infrastructure yaml file from an interactive CLI`,
	Run: func(cmd *cobra.Command, args []string) {

		generateYaml()

	},
}

type promptContent struct {
	errorMsg string
	label    string
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

func promptGetInput(pc promptContent) string {
	validate := func(input string) error {
		if input == "" {
			return errors.New(pc.errorMsg)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Input: %s\n", result)

	return result
}

func promptGetChoice(pc promptContent, items []string) string {
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label:    pc.label,
			Items:    items,
			AddLabel: "Other",
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Input: %s\n", result)

	return result
}

func generateYaml() {
	namePromptContent := promptContent{
		"Service",
		"What is the name of your Service?",
	}
	name := promptGetInput(namePromptContent)

	teamPromptContent := promptContent{
		"Team",
		"What is the name of your team?",
	}
	team := promptGetInput(teamPromptContent)

	envPromptContent := promptContent{
		"Environment",
		"What is the name of your environment?",
	}
	environment := promptGetInput(envPromptContent)

	cloudPromptContent := promptContent{
		"Cloud Provider",
		fmt.Sprintf("What cloud do you need for %s?", name),
	}
	cloud := promptGetChoice(cloudPromptContent, []string{"aws", "gcp", "azure"})

	accountPromptContent := promptContent{
		"Account Name",
		fmt.Sprintf("What is your account/project/subscription in %s?", cloud),
	}
	account := promptGetInput(accountPromptContent)

	regionPromptContent := promptContent{
		"Region",
		fmt.Sprintf("What is the region of your Service in %s?", account),
	}
	region := promptGetInput(regionPromptContent)

	isComponent := "Yes"
	var components []map[string]any
	var componentStr string

	for isComponent == "Yes" {
		isComponentPromptContent := promptContent{
			"New Component",
			fmt.Sprintf("Do you want to add component to %s?", name),
		}
		isComponent = promptGetChoice(isComponentPromptContent, []string{"Yes", "No"})

		if isComponent == "No" {
			break
		}

		component := make(map[string]any)
		componentPromptContent := promptContent{
			"Component Name",
			fmt.Sprintf("What component do you need for %s?", name),
		}
		componentStr = promptGetChoice(componentPromptContent, []string{"Storage", "Compute"})
		if componentStr == "Storage" {
			storageNamePromptContent := promptContent{
				"Storage Name",
				fmt.Sprintf("What is the name of the storage for %s?", name),
			}
			storageName := promptGetInput(storageNamePromptContent)
			component[componentStr] = make(map[string]string)
			component[componentStr].(map[string]string)["name"] = storageName
		} else if componentStr == "Compute" {
			computeNamePromptContent := promptContent{
				"Compute Name",
				fmt.Sprintf("What is the name of the VM for %s?", name),
			}
			computeName := promptGetInput(computeNamePromptContent)
			component[componentStr] = make(map[string]string)
			component[componentStr].(map[string]string)["name"] = computeName
		}
		components = append(components, component)
	}

	C.Name = name
	C.Team = team
	C.Env = environment
	C.Cloud = cloud
	C.Account = account
	C.Region = region
	C.Components = components

	yamlData, err := yaml.Marshal(&C)

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}

	if cfgFile == "" {
		home, _ := os.UserHomeDir()
		cfgFile = filepath.Join(home, ".infy.yaml")
	}

	err = os.WriteFile(cfgFile, yamlData, 0600)
	if err != nil {
		panic("Unable to write data into the file")
	}
}
