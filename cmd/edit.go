package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/roma-glushko/frens/internal/lifedir"
	"github.com/urfave/cli/v2"
)

var EditCommand = &cli.Command{
	Name:    "edit",
	Aliases: []string{"e"},
	Usage:   "Edit life space raw files",
	Flags:   []cli.Flag{},
	Action: func(context *cli.Context) error {
		lifeDir, err := lifedir.DefaultDir()
		if err != nil {
			return err
		}

		editor := GetEditor()

		cmd := exec.Command(editor, lifeDir+"/friends.toml") // TODO: make it configurable
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("Error running editor: %s", err)
		}

		return nil
	},
}

func GetEditor() string {
	editor := os.Getenv("EDITOR")

	if editor == "" {
		editor = "vim"
	}

	return editor
}
