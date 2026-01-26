package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/reveever/gocookie"
)

const envFileName = ".lcode"

func loadEnv() (LEETCODE_SESSION string, LEETCODE_CSRF_TOKEN string) {
	godotenv.Load()
	home, _ := os.UserHomeDir()
	envPath := filepath.Join(home, envFileName)
	godotenv.Load(envPath)
	return os.Getenv("LEETCODE_SESSION"), os.Getenv("LEETCODE_CSRF_TOKEN")
}

func openInEditor(fileName string) {
	tmuxEnv := os.Getenv("TMUX")

	if tmuxEnv != "" {
		fmt.Printf("\nOpening %s in a new Tmux pane...\n", fileName)
		cmd := exec.Command("tmux", "new-window", "-d", fmt.Sprintf("nvim %s", fileName))
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error opening tmux pane:", err)
		}
	} else {
		// fmt.Printf("\nOpening %s in Neovim...\n", fileName)
		// cmd := exec.Command("nvim", fileName)
		// cmd.Stdin = os.Stdin
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// cmd.Run()
	}
}

func initEnvFile() {
	home, _ := os.UserHomeDir()
	envPath := filepath.Join(home, envFileName)
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Println("Authenticating..")
		file, err := os.Create(envPath)
		if err != nil {
			log.Fatal("Error creating .lcode file:", err)
		}
		defer file.Close()
		fmt.Fprintln(file, "# LeetCode CLI environment variables")
		fmt.Fprintln(file, "LEETCODE_SESSION=")
		fmt.Fprintln(file, "LEETCODE_CSRF_TOKEN=")
	}
	session := os.Getenv("LEETCODE_SESSION")
	if session == "" {
		fmt.Println("No session found. Extracting from browser...")

		cookies, err := gocookie.GetCookies(gocookie.Chrome, gocookie.WithDomainSuffix("leetcode.com"))
		if err != nil {
			log.Fatal(err)
		}

		for _, cookie := range cookies {
			if cookie.Name == "LEETCODE_SESSION" {
				os.Setenv("LEETCODE_SESSION", cookie.Value)
			}
			if cookie.Name == "csrftoken" {
				os.Setenv("LEETCODE_CSRF_TOKEN", cookie.Value)
			}
		}

		envData := map[string]string{
			"LEETCODE_SESSION":    os.Getenv("LEETCODE_SESSION"),
			"LEETCODE_CSRF_TOKEN": os.Getenv("LEETCODE_CSRF_TOKEN"),
		}

		err = godotenv.Write(envData, envPath)
		if err != nil {
			log.Fatal("Could not save session file")
		}
		fmt.Println("Successfully logged in and saved session!")
	} else {
		err := godotenv.Load(envPath)
		if err != nil {
			log.Fatal("Error loading .lcode file:", err)
		}
		fmt.Println("Using existing session", envPath)
	}
}
