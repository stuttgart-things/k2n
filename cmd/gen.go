package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stuttgart-things/k2n/internal"
	"github.com/stuttgart-things/k2n/internal/ai"
)

const (
	envAPIKeyVar = "GEMINI_API_KEY"
)

var (
	examplesDir       string
	rulesetEnvDir     string
	rulesetUsecaseDir string
	usecase           string
	instruction       string
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a claim/code configuration using AI based on examples and rulesets",
	Long:  `The 'gen' command uses the Gemini AI model to generate configurations from code examples and optional rulesets.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := os.Getenv(envAPIKeyVar)
		if apiKey == "" {
			panic("GEMINI_API_KEY is not set in environment")
		}

		var examples []string
		var err error

		if examplesDir != "" {
			examples, err = internal.LoadCodeExamples(examplesDir)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("No examples directory provided. Proceeding without examples.")
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

	genCmd.Flags().StringVar(&examplesDir, "examples-dir", "", "Directory containing example code files")
	genCmd.Flags().StringVar(&rulesetEnvDir, "ruleset-env-dir", "", "Directory containing environment rulesets (optional)")
	genCmd.Flags().StringVar(&rulesetUsecaseDir, "ruleset-usecase-dir", "", "Directory containing use case rulesets (optional)")
	genCmd.Flags().StringVar(&usecase, "usecase", "", "usecase context for generation")
	genCmd.Flags().StringVar(&instruction, "instruction", "", "Specific instruction to guide the AI")
}
