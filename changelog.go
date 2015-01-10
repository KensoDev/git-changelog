// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command tipgodoc is the beginning of the new tip.golang.org server,
// serving the latest HEAD straight from the Git oven.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	repoURL = "git@github.com:gogobot/deploy.git"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var err error

	dir := os.TempDir()
	path := filepath.Join(dir, "gogobot-deploy")

	if err = initRepository(path); err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile("./current-head", getHead(path), 0644)
	check(err)
}

func runCommand(cmd *exec.Cmd) {
	_, err := runErr(cmd)
	check(err)
}

func getHead(path string) []byte {
	cmd := exec.Command("git", "fetch")
	cmd.Dir = path
	runCommand(cmd)

	cmd = exec.Command("git", "reset", "--hard", "HEAD")
	cmd.Dir = path
	runCommand(cmd)

	cmd = exec.Command("git", "clean", "-d", "-f", "-x")
	cmd.Dir = path
	runCommand(cmd)

	cmd = exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = path
	out, err := runErr(cmd)
	check(err)
	return out
}

func initRepository(path string) error {
	if _, err := os.Stat(filepath.Join(path, ".git")); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		_, err := runErr(exec.Command("git", "clone", repoURL, path))
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func runErr(cmd *exec.Cmd) (stdout []byte, stderr error) {
	out, err := cmd.CombinedOutput()
	if err != nil {
		if len(out) == 0 {
			return nil, err
		}
		return nil, fmt.Errorf("%s\n%v", out, err)
	}
	return out, err
}