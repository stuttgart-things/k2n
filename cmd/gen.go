// Package cmd provides the command-line interface for generating configurations using AI.
//
// Copyright © 2025 PATRICK HERMANN
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
	envAPIKeyVar = "AI_API_KEY" // pragma: allowlist secret
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
	exampleFileExt      string
	promptToAI          bool
	generatedResult     string
	aiprovider          string
	aiproviderModel     string
	aiproviderBaseURL   string
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
			"PROMPT-TO-AI":          fmt.Sprintf("%t", promptToAI),
			"VERBOSE":               fmt.Sprintf("%t", verbose),
		}

		internal.PrintBanner()
		internal.PrintEnvTable(allFlags)

		// READ API KEY
		apiKey := os.Getenv(envAPIKeyVar)
		if apiKey == "" {
			panic("AI_API_KEY is not set in environment")
		}

		// SETUP PROVIDER CONFIGURATION
		providerConfig := &ai.ProviderConfig{
			APIKey: apiKey,
		}

		// Use flag values or environment variables for provider settings
		if aiprovider != "" {
			providerConfig.Type = ai.ProviderType(aiprovider)
		} else if envProvider := os.Getenv("AI_PROVIDER"); envProvider != "" {
			providerConfig.Type = ai.ProviderType(envProvider)
		} else {
			providerConfig.Type = ai.ProviderOpenRouter
		}

		// Configure provider-specific settings
		switch providerConfig.Type {
		case ai.ProviderOpenRouter:
			if aiproviderModel != "" {
				providerConfig.Model = aiproviderModel
			} else if envModel := os.Getenv("AI_MODEL"); envModel != "" {
				providerConfig.Model = envModel
			} else {
				providerConfig.Model = "openai/gpt-3.5-turbo"
			}
			if aiproviderBaseURL != "" {
				providerConfig.BaseURL = aiproviderBaseURL
			} else if envURL := os.Getenv("AI_BASE_URL"); envURL != "" {
				providerConfig.BaseURL = envURL
			} else {
				providerConfig.BaseURL = "https://openrouter.ai/api/v1/chat/completions"
			}
		case ai.ProviderGemini:
			// Gemini doesn't require additional configuration
		}

		// Add AI environment variables to flags display
		allFlags["AI_API_KEY"] = "***" // Don't expose actual key
		allFlags["AI_PROVIDER"] = string(providerConfig.Type)
		if providerConfig.Model != "" {
			allFlags["AI_MODEL"] = providerConfig.Model
		}
		if providerConfig.BaseURL != "" {
			allFlags["AI_BASE_URL"] = providerConfig.BaseURL
		}

		fmt.Println("\n📋 AI Configuration:")
		internal.PrintEnvTable(map[string]string{
			"AI_API_KEY":  "***",
			"AI_PROVIDER": string(providerConfig.Type),
			"AI_MODEL":    providerConfig.Model,
			"AI_BASE_URL": providerConfig.BaseURL,
		})

		// READ EXAMPLES
		if examplesDir != "" {
			dirs := internal.SplitAndTrimPaths(examplesDir)
			for _, dir := range dirs {
				dirExamples, err := internal.LoadCodeExamplesWithExtensions(dir, internal.SplitAndTrimExts(exampleFileExt))
				if err != nil {
					panic(fmt.Errorf("failed to load examples from dir %s: %w", dir, err))
				}
				examples = append(examples, dirExamples...)
			}
		}
		if exampleFiles != "" {
			paths := internal.SplitAndTrimPaths(exampleFiles)
			fmt.Println("Example file paths:", paths)

			fileExamples, err := internal.LoadExampleFilesWithExtensions(paths, internal.SplitAndTrimExts(exampleFileExt))
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

		if promptToAI && instruction != "" {

			// CALL AI PROVIDER
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			spinnerErr := spinner.New().
				Context(ctx).
				Title(fmt.Sprintf("CALLING %s AI...🚀", string(providerConfig.Type))).
				Action(func() {
					res, err := ai.CallAI(providerConfig, prompt)
					if err != nil {
						fmt.Printf("ERROR CALLING %s API: %v\n", string(providerConfig.Type), err)
						return
					}
					generatedResult = res
				}).
				Run()

			if spinnerErr != nil {
				panic(spinnerErr)
			}

			if err := internal.SaveOutput(destination, generatedResult); err != nil {
				panic(err)
			}
		} else if promptToAI && instruction == "" {
			fmt.Println("⚠️  No instruction provided. Skipping AI call. Use --instruction to prompt the AI.")
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
	genCmd.Flags().BoolVarP(&promptToAI, "prompt-to-ai", "p", true, "Prompt the AI with the generated content (default true)")
	genCmd.Flags().StringVar(&exampleFileExt, "example-file-ext", ".yaml,.tf", "Comma-separated list of allowed example file extensions (e.g., .yaml,.tf)")
	genCmd.Flags().StringVar(&aiprovider, "ai-provider", "", "AI provider: openrouter or gemini (default: gemini, can also use AI_PROVIDER env var)")
	genCmd.Flags().StringVar(&aiproviderModel, "ai-model", "", "Model name for the AI provider (e.g., openai/gpt-4 for OpenRouter, can also use AI_MODEL env var)")
	genCmd.Flags().StringVar(&aiproviderBaseURL, "ai-base-url", "", "Base URL for OpenRouter API (can also use AI_BASE_URL env var)")
}
