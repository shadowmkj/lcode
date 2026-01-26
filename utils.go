package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func loadEnv() (LEETCODE_SESSION string, LEETCODE_CSRF_TOKEN string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("LEETCODE_SESSION"), os.Getenv("LEETCODE_CSRF_TOKEN")
}

func openInEditor(fileName string) {
	tmuxEnv := os.Getenv("TMUX")

	if tmuxEnv != "" {
		fmt.Printf("\nOpening %s in a new Tmux pane...\n", fileName)
		cmd := exec.Command("tmux", "new-window", fmt.Sprintf("nvim %s", fileName))
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error opening tmux pane:", err)
		}
	} else {
		fmt.Printf("\nOpening %s in Neovim...\n", fileName)
		cmd := exec.Command("nvim", fileName)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}
