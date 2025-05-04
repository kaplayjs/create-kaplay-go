package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/kaplayjs/create-kaplay/config"
	"github.com/kaplayjs/create-kaplay/templates"
	"github.com/spf13/cobra"
)

// Types
type KAFile struct {
	Path     string `json:"path"`
	Type     string `json:"type"`
	Strategy string `json:"strategy"`
	URL      string `json:"url"`
}

type KATemplateData struct {
	Files []KAFile `json:"files"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "create-kaplay [flags] <appName>",
	Long: `Create a new KAPLAY game in no time! 🦖`,
	Run: func(cmd *cobra.Command, args []string) {
		isList := cmd.Flags().Lookup(("list")).Value.String() == "true"
		template := cmd.Flags().Lookup(("template")).Value.String()
		version := cmd.Flags().Lookup(("version")).Value.String()

		if isList {
			listTemplates()
			return
		}

		if len(args) < 1 {
			cmd.Help()
			os.Exit(1)
		}

		if template == "" {
			template = "vite"
		}

		if version == "" {
			version = config.DefaultVersion
		}

		appName := args[0]

		print("Creating a new KAPLAY game...\n")

		fmt.Println("Clonning template:", template)
		repoURL := parseGitURL(template)

		cloneCmd := exec.Command("git", "clone", repoURL, appName)

		if err := cloneCmd.Run(); err != nil {
			fmt.Println("Error at repo cloning: ", err)
			return
		}

		// Delete the .git folder
		err := os.RemoveAll(appName + "/.git")

		if err != nil {
			fmt.Println("Error deleting .git folder:", err)
			return
		}

		// Template data
		templateData := map[string]string{
			"version": version,
			"title":   appName,
		}

		// Process katemplate.json
		file, err := os.Open(appName + "/katemplate.json")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Decoding JSON
		decoder := json.NewDecoder(file)
		var data KATemplateData
		err = decoder.Decode(&data)
		if err != nil {
			panic(err)
		}

		// Parse json.files
		for _, tpl := range data.Files {
			switch expression := tpl.Type; expression {
			case "folder":
				// Create the folder
				err := os.MkdirAll(appName+"/"+tpl.Path, 0755)
				if err != nil {
					panic(err)
				}

				fmt.Printf("- created → %s\n", tpl.Path)
			case "file":
				switch tpl.Strategy {
				case "template":
					// Content
					content, err := os.ReadFile(appName + "/" + tpl.Path)
					if err != nil {
						fmt.Println("Error reading file:", err)
						return
					}

					// Parse the template
					newContent := templates.ParseTemplate(string(content), templateData)

					// Write the new content to the file
					err = os.WriteFile(appName+"/"+tpl.Path, []byte(newContent), 0644)

					if err != nil {
						fmt.Println("Error writing file:", err)
						return
					}

					fmt.Printf("- %sed → %s\n", tpl.Strategy, tpl.Path)
				case "fetch":
					{
						// Fetch the file
						fetchURL := templates.ParseTemplate(string(tpl.URL), templateData)
						res, err := http.Get(fetchURL)

						if err != nil {
							panic(err)
						}

						defer res.Body.Close()

						// Check if the response status is OK
						body, err := io.ReadAll(res.Body)

						if err != nil {
							panic(err)
						}

						// Write the new content to the file
						err = os.WriteFile(appName+"/"+tpl.Path, body, 0644)

						if err != nil {
							panic(err)
						}

						fmt.Printf("- fetched → %s\n", tpl.Path)
					}
				}
			}
		}

		fmt.Println("New game created:", appName)
	},
}

func parseGitURL(input string) string {
	if val, ok := templates.DefaultTemplates[input]; ok {
		return val.Url
	}

	// Si ya es una URL completa, la usamos directo
	if strings.HasPrefix(input, "http") || strings.HasSuffix(input, ".git") {
		return input
	}

	// Si es tipo github.com/user/repo → lo convertimos a https
	if strings.HasPrefix(input, "github.com/") {
		return "https://" + input + ".git"
	}

	return input
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.create-kaplay.yaml)")

	// Template option
	rootCmd.Flags().BoolP("help", "h", false, "Show this help message")
	rootCmd.Flags().StringP("template", "r", "", "Template to use for the new project")
	rootCmd.Flags().BoolP("list", "l", false, "List all default templates")
	rootCmd.Flags().StringP("version", "v", "3001", "Version of KAPLAY to use")
}

// #region Flags
func listTemplates() {
	fmt.Println("Default templates:")

	for _, value := range templates.DefaultTemplates {
		fmt.Printf("- %s: %s\n", value.Name, value.Description)
	}
}

// #endregion
