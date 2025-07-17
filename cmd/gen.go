package cmd

import (
	"fmt"
	"os"

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
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a claim/code configuration using AI based on examples and rulesets",
	Long:  `The 'gen' command uses the Gemini AI model to generate configurations from code examples and optional rulesets.`,
	Run: func(cmd *cobra.Command, args []string) {

		// READ API KEY
		apiKey := os.Getenv(envAPIKeyVar)
		if apiKey == "" {
			panic("GEMINI_API_KEY is not set in environment")
		}

		// READ EXAMPLES
		if examplesDir != "" {
			dirExamples, err := internal.LoadCodeExamples(examplesDir)
			if err != nil {
				panic(err)
			}
			examples = append(examples, dirExamples...)
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

		fmt.Println("Generated Prompt:", prompt)

		result, err := ai.CallGeminiAPI(apiKey, prompt)
		if err != nil {
			panic(err)
		}

		if err := internal.SaveOutput(destination, result); err != nil {
			panic(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().StringVar(&exampleFiles, "example-files", "", "Comma-separated list of example file paths")
	genCmd.Flags().StringVar(&examplesDir, "examples-dir", "", "Directory containing example code files")
	genCmd.Flags().StringVar(&rulesetEnvDir, "ruleset-env-dir", "", "Directory containing environment rulesets (optional)")
	genCmd.Flags().StringVar(&rulesetUsecaseDir, "ruleset-usecase-dir", "", "Directory containing use case rulesets (optional)")
	genCmd.Flags().StringVar(&rulesetEnvFiles, "ruleset-env-files", "", "Comma-separated list of environment ruleset files")
	genCmd.Flags().StringVar(&rulesetUsecaseFiles, "ruleset-usecase-files", "", "Comma-separated list of usecase ruleset files")
	genCmd.Flags().StringVar(&usecase, "usecase", "", "usecase context for generation")
	genCmd.Flags().StringVar(&instruction, "instruction", "", "Specific instruction to guide the AI")
	genCmd.Flags().StringVar(&destination, "destination", "", "Specific destination file to save the generated configuration (default: stdout)")
}
