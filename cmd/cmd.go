package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/arxeiss/go-expression-calculator/evaluator"
	"github.com/arxeiss/go-expression-calculator/parser"
	"github.com/arxeiss/go-expression-calculator/parser/shuntyard"
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
	Long:  `Expression calculator in Go with REPL. Write 'help' to REPL console to get more info.`,
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
		p, err := shuntyard.NewParser(parser.DefaultTokenPriorities())
		if err != nil {
			return err
		}

		fmt.Printf("Welcome to the expression calculator, write '%s' to get more info\n", color.HiCyanString("help"))
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
			case "help":
				fmt.Printf(
					"%s\n   %s - %s\n   %s - %s\n   %s - %s\n   %s - %s\n   %s - %s\n",
					"Write directly any expression to evaluate, or one of those commands:",
					color.HiYellowString("functions  "), "Show all available functions",
					color.HiYellowString("variables  "), "Prints all variables with values",
					color.HiYellowString("help       "), "Show this help",
					color.HiYellowString("tree {expr}"), "Write tree and then expression to print AST tree",
					color.HiYellowString("exit       "), "Quit this REPL",
				)
				continue
			case "func", "funcs", "functions":
				funcs := numEvaluator.FunctionList()
				if len(funcs) == 0 {
					fmt.Println(color.YellowString("There are no defined functions"))
					continue
				}
				fmt.Println(color.GreenString("All functions:"))
				for _, f := range funcs {
					fmt.Printf("%s: %s\n", color.HiBlueString(f.Name), f.Description)
				}
			case "vars", "variables":
				vars := numEvaluator.VariableList()
				if len(vars) == 0 {
					fmt.Println(color.YellowString("There are no variables now"))
					continue
				}
				fmt.Println(color.GreenString("All variables:"))
				for _, f := range vars {
					fmt.Printf("%s: %f\n", color.HiBlueString(f.Name), f.Value)
				}
			default:
				parseExpression(numEvaluator, p, expr)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	_ = rootCmd.Execute()
}
