package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	variableRegex = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")
)

func strInStrSlice(slice []string, elem string) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}

func initVariables(flagInitVars bool) (map[string]float64, error) {
	if !flagInitVars {
		return nil, nil
	}
	reader := bufio.NewReader(os.Stdin)
	vars := map[string]float64{}
	for {
		fmt.Print(color.HiBlackString("> "), "Enter new variable name (keep empty to exit): ")
		color.Set(color.FgHiBlue)
		name, err := reader.ReadString('\n')
		color.Unset()
		if err != nil {
			return nil, err
		}
		name = strings.TrimSpace(name)
		if name == "" {
			break
		}
		if !variableRegex.MatchString(name) {
			fmt.Print(color.RedString(
				"Error: variable cannot start with a number and can contain only letters, numbers and underscore\n",
			))
			continue
		}

		for {
			fmt.Print(color.HiBlackString("> "))
			fmt.Printf("Enter new value for variable '%s' (use decimal dot): ", name)
			color.Set(color.FgBlue)
			strVal, err := reader.ReadString('\n')
			color.Unset()
			if err != nil {
				return nil, err
			}
			strVal = strings.TrimSpace(strVal)
			val, err := strconv.ParseFloat(strVal, 64)
			if err != nil {
				fmt.Print(color.RedString("Error: cannot parse given number, remember to use deciamal dot\n"))
				continue
			}

			if _, has := vars[name]; has {
				fmt.Print(color.YellowString("Overriding variable! "))
			}
			vars[name] = val
			fmt.Println(color.GreenString("Variable '%s' with value %f was added", name, val))
			break
		}
	}

	return vars, nil
}
