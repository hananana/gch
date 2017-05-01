package main

import (
	"bufio"
	"fmt"
	"github.com/mitchellh/colorstring"
	"github.com/urfave/cli"
	"gopkg.in/kyokomi/emoji.v1"
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
		status, err := gitStatus(paths)
		if err != nil {
			return err
		}

		err = outPut(paths, status)
		if err != nil {
			return err
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

	dotPath := os.Getenv("HOME")
	dotPath = filepath.Join(dotPath, ".dotfiles")
	_, err = os.Stat(dotPath)
	if err == nil {
		lines = append(lines, dotPath)
	}

	return lines, nil
}

func gitStatus(paths []string) ([]string, error) {
	status := make([]string, 0, 100)
	for _, v := range paths {
		err := os.Chdir(v)
		if err != nil {
			return nil, err
		}
		out, err := exec.Command("git", "status").Output()
		if err != nil {
			return nil, err
		}
		status = append(status, string(out))
	}
	return status, nil
}

func outPut(paths []string, status []string) error {
	for k, v := range status {
		path := paths[k]
		if isOk(v) {
			fmt.Println(colorstring.Color("[blue]" + path + "[reset]" + emoji.Sprint(" ok!:sparkles:")))
		} else {
			fmt.Println(colorstring.Color("[yellow]" + path + "[reset]" + emoji.Sprint(" check me!:pig:")))
		}
	}
	return nil
}

func isOk(status string) bool {
	return isClean(status)
}

func isClean(status string) bool {
	return strings.Index(status, "working tree clean") >= 0
}

func isPushed(status string) bool {
	return strings.Index(status, "branch is ahead of") != -1
}

func isUpToDate(status string) bool {
	return strings.Index(status, "up-to-date") >= 0
}

func isNoDiff(status string) bool {
	return strings.Index(status, "nothing") >= 0
}
