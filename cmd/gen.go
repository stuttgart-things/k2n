package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stuttgart-things/k2n/internal"
	"github.com/stuttgart-things/k2n/internal/ai"
)

const (
	envAPIKeyVar = "GEMINI_API_KEY" // pragma: allowlist secret
)

var (
	examplesDir       string
	exampleFiles      string
	rulesetEnvDir     string
	rulesetUsecaseDir string
	usecase           string
	instruction       string
	examples          []string
	err               error
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
			paths := strings.Split(exampleFiles, ",")
			for i := range paths {
				paths[i] = strings.TrimSpace(paths[i])
			}
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

		fmt.Println(result)
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().StringVar(&exampleFiles, "example-files", "", "Comma-separated list of example file paths")
	genCmd.Flags().StringVar(&examplesDir, "examples-dir", "", "Directory containing example code files")
	genCmd.Flags().StringVar(&rulesetEnvDir, "ruleset-env-dir", "", "Directory containing environment rulesets (optional)")
	genCmd.Flags().StringVar(&rulesetUsecaseDir, "ruleset-usecase-dir", "", "Directory containing use case rulesets (optional)")
	genCmd.Flags().StringVar(&usecase, "usecase", "", "usecase context for generation")
	genCmd.Flags().StringVar(&instruction, "instruction", "", "Specific instruction to guide the AI")
}
