package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/frivas/go-git-swap/internal/config"
	"github.com/frivas/go-git-swap/internal/git"
	"github.com/frivas/go-git-swap/internal/model"
	"github.com/frivas/go-git-swap/internal/validator"
)

var (
	green  = color.New(color.FgGreen).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Printf("%s: %v\n", red("Error initializing"), err)
		os.Exit(1)
	}

	gitConfig := git.New()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n" + blue("Git Profile Manager"))
		fmt.Println("-------------------")

		// Show active profile
		activeProfile := cfg.GetActiveProfile()
		if activeProfile != "" {
			fmt.Printf("%s: %s\n", yellow("Active Profile"), green(activeProfile))
		} else {
			fmt.Printf("%s: %s\n", yellow("Active Profile"), red("None"))
		}

		// Show available profiles
		profiles := cfg.ListProfiles()
		fmt.Printf("\n%s:\n", blue("Available Profiles"))
		if len(profiles) == 0 {
			fmt.Println(red("No profiles found"))
		} else {
			for i, profile := range profiles {
				marker := " "
				if profile == activeProfile {
					marker = "*"
				}
				fmt.Printf("%s %d) %s\n", marker, i+1, profile)
			}
		}

		// Show menu
		fmt.Printf("\n%s:\n", blue("Commands"))
		fmt.Println("1) Select profile")
		fmt.Println("2) Create new profile")
		fmt.Println("3) Delete profile")
		fmt.Println("4) Exit")

		fmt.Print("\nEnter command (1-4): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			if len(profiles) == 0 {
				fmt.Println(red("\nNo profiles available to select"))
				continue
			}
			fmt.Print("Enter profile number: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			profileNum := 0
			fmt.Sscanf(input, "%d", &profileNum)

			if profileNum < 1 || profileNum > len(profiles) {
				fmt.Println(red("Invalid profile number"))
				continue
			}

			selectedProfile := profiles[profileNum-1]
			profile, exists := cfg.GetProfile(selectedProfile)
			if !exists {
				fmt.Printf("%s: Profile %s not found\n", red("Error"), selectedProfile)
				continue
			}

			if err := gitConfig.SetGlobalConfig(profile.Email, profile.Username, profile.FullName, os.ExpandEnv(profile.SigningKey)); err != nil {
				fmt.Printf("%s: %v\n", red("Error updating Git configuration"), err)
				continue
			}

			cfg.SetActiveProfile(profile.Name)
			fmt.Printf("%s: %s\n", green("Successfully activated profile"), profile.Name)

		case "2":
			var newProfile model.GitProfile
			fmt.Print("Enter profile name (identifier): ")
			newProfile.Name, _ = reader.ReadString('\n')
			newProfile.Name = strings.TrimSpace(newProfile.Name)

			fmt.Print("Enter Git user name (full name): ")
			newProfile.FullName, _ = reader.ReadString('\n')
			newProfile.FullName = strings.TrimSpace(newProfile.FullName)

			fmt.Print("Enter Git username: ")
			newProfile.Username, _ = reader.ReadString('\n')
			newProfile.Username = strings.TrimSpace(newProfile.Username)

			fmt.Print("Enter email: ")
			newProfile.Email, _ = reader.ReadString('\n')
			newProfile.Email = strings.TrimSpace(newProfile.Email)

			if err := validator.ValidateEmail(newProfile.Email); err != nil {
				fmt.Printf("%s: %v\n", red("Error"), err)
				continue
			}

			fmt.Print("Enter SSH public key path (optional, press Enter to skip): ")
			pubKeyPath, _ := reader.ReadString('\n')
			pubKeyPath = strings.TrimSpace(pubKeyPath)

			if pubKeyPath != "" {
				pubKeyPath = os.ExpandEnv(pubKeyPath)
				if strings.HasPrefix(pubKeyPath, "~/") {
					homeDir, err := os.UserHomeDir()
					if err != nil {
						fmt.Printf("%s: %v\n", red("Error getting home directory"), err)
						continue
					}
					pubKeyPath = filepath.Join(homeDir, pubKeyPath[2:])
				}
				_, err := os.ReadFile(pubKeyPath)
				if err != nil {
					fmt.Printf("%s: %v\n", red("Error reading SSH public key"), err)
					continue
				}
				newProfile.SigningKey = pubKeyPath
			}

			if err := cfg.AddProfile(newProfile); err != nil {
				fmt.Printf("%s: %v\n", red("Error adding profile"), err)
				continue
			}

			fmt.Printf("%s: %s\n", green("Successfully created profile"), newProfile.Name)

		case "3":
			if len(profiles) == 0 {
				fmt.Println(red("\nNo profiles available to delete"))
				continue
			}
			fmt.Print("Enter profile number to delete: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			profileNum := 0
			fmt.Sscanf(input, "%d", &profileNum)

			if profileNum < 1 || profileNum > len(profiles) {
				fmt.Println(red("Invalid profile number"))
				continue
			}

			profileToDelete := profiles[profileNum-1]
			if profileToDelete == activeProfile {
				fmt.Print(yellow("Warning: This is your active profile! Are you sure? (y/N): "))
				confirm, _ := reader.ReadString('\n')
				confirm = strings.ToLower(strings.TrimSpace(confirm))
				if confirm != "y" {
					continue
				}
			}

			if err := cfg.RemoveProfile(profileToDelete); err != nil {
				fmt.Printf("%s: %v\n", red("Error deleting profile"), err)
				continue
			}

			fmt.Printf("%s: %s\n", green("Successfully deleted profile"), profileToDelete)

		case "4":
			fmt.Println(green("Goodbye!"))
			return

		default:
			fmt.Println(red("Invalid choice"))
		}
	}
}
