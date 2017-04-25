package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		paths, err := repoPaths()
		if err != nil {
			return err
		}

		for _, v := range paths {
			fmt.Println(v)
		}

		return nil
	}
	app.Run(os.Args)
}

func repoPaths() ([]string, error) {
	out, err := exec.Command("ghq", "list").Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	lines := make([]string, 0, 100)
	gopath := os.Getenv("GOPATH")
	for scanner.Scan() {
		path := filepath.Join(gopath, "src", scanner.Text())
		lines = append(lines, path)
	}

	return lines, nil
}
