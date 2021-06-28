package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/arxeiss/go-expression-calculator/evaluator"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	flagInitVars *bool
	flagNoFuncs  *bool
	flagParser   *string

	availableParsers = []string{"shunt-yard", "recursive"}
)

func init() {
	flagInitVars = rootCmd.Flags().BoolP("init-vars", "i", false, "Before start, initialize values")
	flagParser = rootCmd.Flags().StringP("parser", "p", "shunt-yard", fmt.Sprintf(
		"Parser to be used, available ones are: '"+strings.Join(availableParsers, "', '")+"'",
	))
	flagNoFuncs = rootCmd.Flags().Bool("no-functions", false, "Disable functions for parser")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "calculator",
	Short: "Expression Calculator",
	Long:  `Expression calculator in Go with REPL.`,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !strInStrSlice(availableParsers, *flagParser) {
			return errors.New("Invalid parser, available ones are: '" + strings.Join(availableParsers, "', '") + "'")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		vars, err := initVariables(*flagInitVars)
		if err != nil {
			return err
		}

		var funcs map[string]evaluator.FunctionHandler
		if !*flagNoFuncs {
			funcs = evaluator.MathFunctions()
		}
		numEvaluator, err := evaluator.NewNumericEvaluator(vars, funcs)
		if err != nil {
			return err
		}

		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print(color.HiBlackString(">>> "))
			expr, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			expr = strings.TrimSpace(expr)
			switch expr {
			case "exit":
				fmt.Println(color.HiMagentaString("All done, good bye!"))
				return nil
			case "func", "funcs", "functions":
			case "vars", "variables":
			default:
				parseExpression(numEvaluator, expr)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	_ = rootCmd.Execute()
}
