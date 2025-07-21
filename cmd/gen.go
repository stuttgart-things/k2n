package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/huh/spinner"

	"github.com/spf13/cobra"
	"github.com/stuttgart-things/k2n/internal"
	"github.com/stuttgart-things/k2n/internal/ai"
)

const (
	envAPIKeyVar = "GEMINI_API_KEY" // pragma: allowlist secret
)

var (
	examplesDir         string
	exampleFiles        string
	rulesetEnvDir       string
	rulesetUsecaseDir   string
	usecase             string
	instruction         string
	examples            []string
	err                 error
	rulesetEnvFiles     string
	rulesetUsecaseFiles string
	destination         string
	verbose             bool
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a claim/code configuration using AI based on examples and rulesets",
	Long:  `The 'gen' command uses the Gemini AI model to generate configurations from code examples and optional rulesets.`,
	Run: func(cmd *cobra.Command, args []string) {

		allFlags := map[string]string{
			"EXAMPLES-DIR":          examplesDir,
			"EXAMPLE-FILES":         exampleFiles,
			"RULESET-ENV-DIR":       rulesetEnvDir,
			"RULESET-USECASE-DIR":   rulesetUsecaseDir,
			"USECASE":               usecase,
			"INSTRUCTION":           instruction,
			"RULESET-ENV-FILES":     rulesetEnvFiles,
			"RULESET-USECASE-FILES": rulesetUsecaseFiles,
			"DESTINATION":           destination,
		}

		internal.PrintBanner()
		internal.PrintEnvTable(allFlags)

		// READ API KEY
		apiKey := os.Getenv(envAPIKeyVar)
		if apiKey == "" {
			panic("GEMINI_API_KEY is not set in environment")
		}

		// READ EXAMPLES
		if examplesDir != "" {
			dirs := internal.SplitAndTrimPaths(examplesDir)
			for _, dir := range dirs {
				dirExamples, err := internal.LoadCodeExamples(dir)
				if err != nil {
					panic(fmt.Errorf("failed to load examples from dir %s: %w", dir, err))
				}
				examples = append(examples, dirExamples...)
			}
		}
		if exampleFiles != "" {
			paths := internal.SplitAndTrimPaths(exampleFiles)
			fmt.Println("Example file paths:", paths)

			fileExamples, err := internal.LoadExampleFiles(paths)
			if err != nil {
				panic(err)
			}
			examples = append(examples, fileExamples...)
		}

		if len(examples) == 0 {
			fmt.Println("No examples provided. Proceeding without examples.")
		} else {
			examples = internal.DeduplicateStrings(examples)
		}

		envRules, _ := internal.LoadRulesetsIfExists(rulesetEnvDir)
		usecaseRules, _ := internal.LoadRulesetsIfExists(rulesetUsecaseDir)

		if rulesetEnvFiles != "" {
			files := internal.SplitAndTrimPaths(rulesetEnvFiles)
			fileRules, err := internal.LoadExampleFiles(files)
			if err != nil {
				panic(err)
			}
			envRules = append(envRules, fileRules...)
		}

		if rulesetUsecaseFiles != "" {
			files := internal.SplitAndTrimPaths(rulesetUsecaseFiles)
			fileRules, err := internal.LoadExampleFiles(files)
			if err != nil {
				panic(err)
			}
			usecaseRules = append(usecaseRules, fileRules...)
		}

		finalInstruction := instruction
		if finalInstruction == "" {
			finalInstruction = fmt.Sprintf("Generate a %s configuration. Only return one file definition, no description.", usecase)
		}

		prompt := internal.BuildPrompt(examples, envRules, usecaseRules, usecase, finalInstruction)

		if verbose {
			fmt.Println(prompt)
		}

		// ASK GEMINI AI
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		var result string
		spinnerErr := spinner.New().
			Context(ctx).
			Title("CALLING GEMINI AI...🚀").
			Action(func() {
				res, err := ai.CallGeminiAPI(apiKey, prompt)
				if err != nil {
					fmt.Println("Error calling Gemini API:", err)
					return
				}
				result = res
			}).
			Run()

		if spinnerErr != nil {
			panic(spinnerErr)
		}

		if err := internal.SaveOutput(destination, result); err != nil {
			panic(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().StringVar(&exampleFiles, "example-files", "", "Comma-separated list of example file paths")
	genCmd.Flags().StringVar(&examplesDir, "examples-dirs", "", "Comma-separated list of directories containing example code files")
	genCmd.Flags().StringVar(&rulesetEnvDir, "ruleset-env-dir", "", "Directory containing environment rulesets (optional)")
	genCmd.Flags().StringVar(&rulesetUsecaseDir, "ruleset-usecase-dir", "", "Directory containing use case rulesets (optional)")
	genCmd.Flags().StringVar(&rulesetEnvFiles, "ruleset-env-files", "", "Comma-separated list of environment ruleset files")
	genCmd.Flags().StringVar(&rulesetUsecaseFiles, "ruleset-usecase-files", "", "Comma-separated list of usecase ruleset files")
	genCmd.Flags().StringVar(&usecase, "usecase", "", "usecase context for generation")
	genCmd.Flags().StringVar(&instruction, "instruction", "", "Specific instruction to guide the AI")
	genCmd.Flags().StringVar(&destination, "destination", "", "Destination for generated files: stdout (default), a file (combined content), or a directory (separate files)")
	genCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}
