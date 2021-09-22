package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"

	"github.com/arxeiss/go-expression-calculator/evaluator"
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

func prettyPrintVariables(vars []evaluator.VariableTuple) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Value"})
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)

	for _, v := range vars {
		table.Append([]string{
			color.HiBlueString(v.Name),
			fmt.Sprintf("%.8f", v.Value),
		})
	}

	fmt.Println(color.GreenString("All variables:"))
	table.Render()
}

func prettyPrintFunctions(funcs []evaluator.FunctionTuple) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Description"})
	table.SetAutoWrapText(false)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})

	for _, v := range funcs {
		table.Append([]string{
			color.HiBlueString(v.Name) + "(" + formatFuncArgs(v.Function) + ")",
			v.Function.Description,
		})
	}

	fmt.Println(color.GreenString("All variables:"))
	table.Render()
}

func formatFuncArgs(v evaluator.FunctionHandler) string {
	p := ""

	i := 0
	for ; i < len(v.ArgsNames); i++ {
		if i >= v.MinArguments {
			break
		}
		p += color.HiGreenString(v.ArgsNames[i])
		if i < v.MinArguments-1 || v.MaxArguments == 0 {
			p += color.HiBlackString(", ")
		}
	}

	if v.MinArguments > 0 && v.MaxArguments == 0 {
		p += color.HiGreenString(v.ArgsNames[i]) + color.HiRedString("...")
	} else {
		varP := ""
		for b := len(v.ArgsNames) - 1; b >= i; b-- {
			inn := color.HiBlackString("[")
			if b > i || len(p) > 0 {
				inn += color.HiBlackString(", ")
			}
			varP = inn + color.HiGreenString(v.ArgsNames[i]) + varP + color.HiBlackString("]")
		}
		p += varP
	}

	return p
}
