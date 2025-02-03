package main

import (
	"fmt"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Attempting to run migrations...")

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cmd := exec.Command(
		"tern",
		"migrate",
		"--migrations",
		"./internal/store/pgstore/migrations",
		"--config",
		"./internal/store/pgstore/migrations/tern.conf",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Migrations command execution failed: ", err)
		fmt.Println("Output: ", string(output))
		panic(err)
	}

	fmt.Println("Migrations command executed successfully!")
}
