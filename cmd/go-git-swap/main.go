package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
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

// Version is set by -ldflags at release time.
var Version = ""

func version() string {
	if Version != "" {
		return Version
	}
	if bi, ok := debug.ReadBuildInfo(); ok {
		return bi.Main.Version
	}
	return "dev"
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-v", "--version", "version":
			fmt.Println("go-git-swap", version())
			return
		default:
			fmt.Fprintln(os.Stderr, "usage: go-git-swap [--version]")
			os.Exit(2)
		}
	}

	cfg, err := config.New()
	if err != nil {
		fmt.Printf("%s: %v\n", red("Error initializing"), err)
		os.Exit(1)
	}

	gitConfig := git.New()
	reader := bufio.NewReader(os.Stdin)
	// readLine exits cleanly on EOF instead of looping forever.
	readLine := func() string {
		line, err := reader.ReadString('\n')
		if err != nil && line == "" {
			fmt.Println()
			os.Exit(0)
		}
		return strings.TrimSpace(line)
	}

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
		choice := readLine()

		switch choice {
		case "1":
			if len(profiles) == 0 {
				fmt.Println(red("\nNo profiles available to select"))
				continue
			}
			fmt.Print("Enter profile number: ")
			profileNum, err := strconv.Atoi(readLine())
			if err != nil || profileNum < 1 || profileNum > len(profiles) {
				fmt.Println(red("Invalid profile number"))
				continue
			}

			selectedProfile := profiles[profileNum-1]
			profile, exists := cfg.GetProfile(selectedProfile)
			if !exists {
				fmt.Printf("%s: Profile %s not found\n", red("Error"), selectedProfile)
				continue
			}

			if err := gitConfig.SetGlobalConfig(profile); err != nil {
				fmt.Printf("%s: %v\n", red("Error updating Git configuration"), err)
				continue
			}

			if err := cfg.SetActiveProfile(profile.Name); err != nil {
				fmt.Printf("%s: %v\n", red("Error saving active profile"), err)
				continue
			}
			fmt.Printf("%s: %s\n", green("Successfully activated profile"), profile.Name)

		case "2":
			var newProfile model.GitProfile
			fmt.Print("Enter profile name (identifier): ")
			newProfile.Name = readLine()
			if err := validator.ValidateProfileName(newProfile.Name); err != nil {
				fmt.Printf("%s: %v\n", red("Error"), err)
				continue
			}

			fmt.Print("Enter Git user name (full name): ")
			newProfile.FullName = readLine()
			if err := validator.ValidateFullName(newProfile.FullName); err != nil {
				fmt.Printf("%s: %v\n", red("Error"), err)
				continue
			}

			fmt.Print("Enter Git username: ")
			newProfile.Username = readLine()
			if err := validator.ValidateUsername(newProfile.Username); err != nil {
				fmt.Printf("%s: %v\n", red("Error"), err)
				continue
			}

			fmt.Print("Enter email: ")
			newProfile.Email = readLine()
			if err := validator.ValidateEmail(newProfile.Email); err != nil {
				fmt.Printf("%s: %v\n", red("Error"), err)
				continue
			}

			fmt.Print("Enter SSH public key path (optional, press Enter to skip): ")
			pubKeyPath := readLine()
			if pubKeyPath != "" {
				pubKeyPath, err := git.ExpandPath(pubKeyPath)
				if err == nil {
					_, err = git.ReadPublicKey(pubKeyPath)
				}
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
			profileNum, err := strconv.Atoi(readLine())
			if err != nil || profileNum < 1 || profileNum > len(profiles) {
				fmt.Println(red("Invalid profile number"))
				continue
			}

			profileToDelete := profiles[profileNum-1]
			if profileToDelete == activeProfile {
				fmt.Print(yellow("Warning: This is your active profile! Are you sure? (y/N): "))
				if strings.ToLower(readLine()) != "y" {
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
