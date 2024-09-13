package core

import (
	"errors"
	"fmt"
	"mpwt/pkg/log"
	"strings"
)

type TerminalConfig struct {
	Maximize        bool
	Direction       string
	Columns         int
	OpenInNewWindow bool
	Commands        []string
}

const (
	Horizontal = "horizontal"
	Vertical   = "vertical"
	Maximize   = "maximize"
)

var flagsMap = map[string]string{
	Horizontal: "-H",
	Vertical:   "-V",
	Maximize:   "-M",
}

// calculatePaneSize takes an integer (number of panes) and returns a slice of float64
func calculatePaneSize(n int) ([]float64, error) {
	if n < 1 {
		return nil, errors.New("length must be greater than 0")
	}

	results := make([]float64, n)

	for i := 0; i < n; i++ {
		num := float64(n - i)
		denom := float64(n - i + 1)
		results[i] = num / denom
	}

	return results, nil
}

// generateCommand takes an array of strings and concat it with a separator
func generateCommand(cmd []string) string {
	return strings.Join(cmd, " ")
}

func OpenWt(t *TerminalConfig) error {
	wtCmd := []string{"wt"}

	// Append maximize flag to command
	if t.Maximize {
		wtCmd = append(wtCmd, flagsMap[Maximize])
	}

	// Split commands into even groups
	cmdsLength := len(t.Commands)
	size := (cmdsLength + t.Columns - 1) / t.Columns
	splitCmds := make([][]string, 0, t.Columns)

	for i := 0; i < cmdsLength; i += size {
		end := i + size
		if end > len(t.Commands) {
			end = len(t.Commands)
		}

		splitCmds = append(splitCmds, t.Commands[i:end])
	}

	log.Debug(fmt.Sprintf("Data processing - wtCmd: %s", generateCommand(wtCmd)))
	log.Debug(fmt.Sprintf("Data processing - splitCmds: %s", splitCmds))

	// Pop and append first command from first cmds group to final windows terminal command
	wtCmd = append(wtCmd, fmt.Sprintf("cmd /k %s;", splitCmds[0][0]))
	splitCmds[0] = splitCmds[0][1:]

	// Reverse general direction when creating tree
	treeDirection := flagsMap[Horizontal]
	if t.Direction == Horizontal {
		treeDirection = flagsMap[Vertical]
	}

	// Calculate size of each tree
	treeSizes, err := calculatePaneSize(t.Columns)
	if err != nil {
		return fmt.Errorf("failed to calculate sizes for trees: %v", err)
	}

	for i := 1; i < len(splitCmds); i++ {
		// Pop and append first command from the rest of cmds group to final windows terminal command
		wtCmd = append(wtCmd, fmt.Sprintf("sp %s -s %.2f cmd /k %s;", treeDirection, treeSizes[i], splitCmds[i][0]))
		splitCmds[i] = splitCmds[i][1:]
	}

	log.Debug(fmt.Sprintf("Tree formation - wtCmd: %s", generateCommand(wtCmd)))
	log.Debug(fmt.Sprintf("Tree formation - splitCmds: %s", splitCmds))

	for i := len(splitCmds) - 1; i >= 0; i-- {
		// Calculate size of each leaf
		sizes, err := calculatePaneSize(len(splitCmds[i]))
		if err != nil {
			return fmt.Errorf("failed to calculate sizes for leaf nodes: %v", err)
		}

		// Form leaf command
		for idx, cmd := range splitCmds[i] {
			leafCmd := fmt.Sprintf("sp %s -s %.2f cmd /k %s;", flagsMap[t.Direction], sizes[idx], cmd)
			log.Debug(fmt.Sprintf("Leaf formation - leafCmd: %s", leafCmd))
			wtCmd = append(wtCmd, leafCmd)
		}

		// Move to the first tree after finish the current
		wtCmd = append(wtCmd, fmt.Sprintf("mf %s%s", map[string]string{Horizontal: "left", Vertical: "up"}[t.Direction], map[bool]string{true: ";", false: ""}[i != 0]))
	}

	log.Debug(generateCommand(wtCmd))
	return nil
}
