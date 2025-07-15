package acceptance

import (
	"testing"

	"github.com/roma-glushko/frens/cmd"
	"github.com/stretchr/testify/require"
)

func TestJournal(t *testing.T) {
	app := cmd.NewApp()

	args := []string{"frens", "journal", "init"}

	err := app.Run(args)
	require.NoError(t, err)
}
