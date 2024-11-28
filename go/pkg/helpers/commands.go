package commands

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

func PromptRaw(prompt string) (string, error) {
	fmt.Printf("%v: ", prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %v", err)
	}
	return strings.TrimSpace(input), nil
}

func PromptWithDefault(c *cli.Context, prompt string, def string) (string, error) {
	v := c.String(prompt)

	// use arg value if supplied
	if v != "" {
		return v, nil
	}

	// use default if available
	if def != "" {
		return def, nil
	}

	// prompt for value
	return PromptRaw(prompt)
}

func Prompt(c *cli.Context, prompt string) (string, error) {
	return PromptWithDefault(c, prompt, "")
}
