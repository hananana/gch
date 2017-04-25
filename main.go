package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		repos, err := repos()
		if err != nil {
			return err
		}

		return nil
	}
	app.Run(os.Args)
}

func repos() ([]string, error) {
	out, err := exec.Command("ghq", "list").Output()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	lines := make([]string, 0, 100)
	for scanner.Scan() {
		append(lines, scanner.Text())
	}

	return lines, nil
}
