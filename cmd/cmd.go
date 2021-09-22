package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/arxeiss/go-expression-calculator/evaluator"
	"github.com/arxeiss/go-expression-calculator/parser"
	"github.com/arxeiss/go-expression-calculator/parser/recursivedescent"
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
	flagParser = rootCmd.Flags().StringP("parser", "p", "recursive", fmt.Sprintf(
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

		var funcs []map[string]evaluator.FunctionHandler
		if !*flagNoFuncs {
			funcs = append(funcs, evaluator.MathFunctions())
		}

		parserName := "Recursive descent"
		var p parser.Parser
		switch *flagParser {
		case "shunt-yard":
			p, err = shuntyard.NewParser(parser.DefaultTokenPriorities())
			parserName = "Shunting Yard"
		default:
			p, err = recursivedescent.NewParser(parser.DefaultTokenPriorities())
			if !*flagNoFuncs {
				funcs = append(funcs, evaluator.MathFunctionsWithVarArgs())
			}
		}
		if err != nil {
			return err
		}

		numEvaluator, err := evaluator.NewNumericEvaluator(vars, funcs...)
		if err != nil {
			return err
		}

		fmt.Printf("Welcome to the expression calculator, write '%s' to get more info\n", color.HiCyanString("help"))
		fmt.Printf("Current parser is '%s'\n", color.HiGreenString(parserName))

		controlC := false
		emptyLine := true
		promptParser := prompt.NewStandardInputParser()

		pp := prompt.New(
			func(s string) {
				controlC = false
				if s != "" && s != "exit" {
					parseLine(s, numEvaluator, p)
				}
			},
			func(d prompt.Document) []prompt.Suggest {
				if d.CursorPositionCol() == 0 && d.LastKeyStroke() == prompt.Escape {
					return []prompt.Suggest{}
				}
				s := []prompt.Suggest{
					{Text: "help", Description: "Open this help"},
					{Text: "functions", Description: "Show all available functions"},
					{Text: "variables", Description: "Show all available variables"},
					{Text: "tree", Description: "Prints AST tree"},
					{Text: "exit", Description: "Quits console"},
				}
				return prompt.FilterHasPrefix(s, d.Text, true)
			},
			prompt.OptionParser(promptParser),
			prompt.OptionAddKeyBind(prompt.KeyBind{
				Key: prompt.ControlC,
				Fn: func(b *prompt.Buffer) {
					controlC = true
				},
			}),
			prompt.OptionSetExitCheckerOnInput(func(in string, breakline bool) bool {
				if controlC && len(in) == 0 && emptyLine {
					fmt.Println("Received interrupt signal")
					return true
				}
				if breakline && in == "exit" {
					return true
				}
				emptyLine = breakline || len(in) == 0
				return false
			}),
			prompt.OptionPrefix(">>> "),
			prompt.OptionTitle("Expression calculator"),
			prompt.OptionSuggestionTextColor(prompt.Turquoise),
			prompt.OptionSuggestionBGColor(prompt.Black),
			prompt.OptionSelectedSuggestionTextColor(prompt.Turquoise),
			prompt.OptionSelectedSuggestionBGColor(prompt.DarkGray),
			prompt.OptionDescriptionTextColor(prompt.Turquoise),
			prompt.OptionDescriptionBGColor(prompt.Black),
			prompt.OptionSelectedDescriptionTextColor(prompt.Turquoise),
			prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),
		)
		pp.Run()
		fmt.Println(color.HiMagentaString("All done, good bye!"))

		return promptParser.TearDown()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	_ = rootCmd.Execute()
}
