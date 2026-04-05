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
	"github.com/stuttgart-things/k2n/internal/talk"
)

var (
	talkAPIURL      string
	talkAuthToken   string
	talkInstruction string
	talkDestination string
	talkProvider    string
	talkModel       string
	talkBaseURL     string
	talkVerbose     bool
)

var talkCmd = &cobra.Command{
	Use:   "talk",
	Short: "AI-powered conversation for rendering Crossplane claims via claim-machinery-api",
	Long: `The 'talk' command connects to a running claim-machinery-api instance,
discovers available claim templates, and uses AI to match your natural language
request to the right template and parameters. The result is a rendered Crossplane
claim in YAML format.`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.PrintBanner()

		// Validate API URL
		if talkAPIURL == "" {
			talkAPIURL = os.Getenv("CLAIM_API_URL")
		}
		if talkAPIURL == "" {
			fmt.Fprintln(os.Stderr, "Error: --api-url or CLAIM_API_URL is required")
			os.Exit(1)
		}

		// Validate instruction
		if talkInstruction == "" {
			fmt.Fprintln(os.Stderr, "Error: --instruction is required")
			os.Exit(1)
		}

		// Read API key
		apiKey := os.Getenv(envAPIKeyVar)
		if apiKey == "" {
			fmt.Fprintln(os.Stderr, "Error: AI_API_KEY environment variable is required")
			os.Exit(1)
		}

		// Auth token for claim-machinery-api
		if talkAuthToken == "" {
			talkAuthToken = os.Getenv("CLAIM_API_TOKEN")
		}

		// Setup AI provider
		providerConfig := &ai.ProviderConfig{APIKey: apiKey}
		if talkProvider != "" {
			providerConfig.Type = ai.ProviderType(talkProvider)
		} else if envProvider := os.Getenv("AI_PROVIDER"); envProvider != "" {
			providerConfig.Type = ai.ProviderType(envProvider)
		} else {
			providerConfig.Type = ai.ProviderOpenRouter
		}

		switch providerConfig.Type {
		case ai.ProviderOpenRouter:
			if talkModel != "" {
				providerConfig.Model = talkModel
			} else if envModel := os.Getenv("AI_MODEL"); envModel != "" {
				providerConfig.Model = envModel
			} else {
				providerConfig.Model = "openai/gpt-3.5-turbo"
			}
			if talkBaseURL != "" {
				providerConfig.BaseURL = talkBaseURL
			} else if envURL := os.Getenv("AI_BASE_URL"); envURL != "" {
				providerConfig.BaseURL = envURL
			} else {
				providerConfig.BaseURL = "https://openrouter.ai/api/v1/chat/completions"
			}
		case ai.ProviderGemini:
			// No additional config needed
		}

		// Print config
		internal.PrintEnvTable(map[string]string{
			"CLAIM_API_URL": talkAPIURL,
			"AI_PROVIDER":   string(providerConfig.Type),
			"AI_MODEL":      providerConfig.Model,
			"INSTRUCTION":   talkInstruction,
			"DESTINATION":   talkDestination,
		})

		// Step 1: Fetch templates from claim-machinery-api
		client := talk.NewClient(talkAPIURL, talkAuthToken)

		var templates []talk.ClaimTemplate
		fmt.Println()
		ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel1()
		if err := spinner.New().
			Context(ctx1).
			Title("Fetching claim templates...").
			Action(func() {
				var fetchErr error
				templates, fetchErr = client.ListTemplates()
				if fetchErr != nil {
					fmt.Fprintf(os.Stderr, "\nError fetching templates: %v\n", fetchErr)
					os.Exit(1)
				}
			}).Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(templates) == 0 {
			fmt.Println("No claim templates found on the API.")
			return
		}
		fmt.Printf("Found %d claim template(s)\n\n", len(templates))

		// Step 2: Build prompt and call AI
		systemPrompt := talk.BuildSystemPrompt(templates)
		fullPrompt := talk.BuildUserPrompt(systemPrompt, talkInstruction)

		if talkVerbose {
			fmt.Println("--- PROMPT ---")
			fmt.Println(fullPrompt)
			fmt.Println("--- END PROMPT ---")
		}

		var aiOutput string
		ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel2()
		if err := spinner.New().
			Context(ctx2).
			Title(fmt.Sprintf("Asking %s AI to select template and parameters...", string(providerConfig.Type))).
			Action(func() {
				var callErr error
				aiOutput, callErr = ai.CallAI(providerConfig, fullPrompt)
				if callErr != nil {
					fmt.Fprintf(os.Stderr, "\nError calling AI: %v\n", callErr)
					os.Exit(1)
				}
			}).Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if talkVerbose {
			fmt.Println("--- AI RESPONSE ---")
			fmt.Println(aiOutput)
			fmt.Println("--- END AI RESPONSE ---")
		}

		// Step 3: Parse AI response
		aiResp, parseErr := talk.ParseAIResponse(aiOutput)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "Error parsing AI response: %v\n", parseErr)
			os.Exit(1)
		}

		if aiResp.TemplateName == "" {
			fmt.Printf("AI could not match your request to a template.\nReason: %s\n", aiResp.Explanation)
			return
		}

		fmt.Printf("Selected template: %s\n", aiResp.TemplateName)
		fmt.Printf("Explanation: %s\n", aiResp.Explanation)
		if talkVerbose {
			fmt.Printf("Parameters: %v\n", aiResp.Parameters)
		}
		fmt.Println()

		// Step 4: Order the claim
		var orderResp *talk.OrderResponse
		ctx3, cancel3 := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel3()
		if err := spinner.New().
			Context(ctx3).
			Title("Rendering claim via claim-machinery-api...").
			Action(func() {
				var orderErr error
				orderResp, orderErr = client.OrderClaim(aiResp.TemplateName, aiResp.Parameters, "k2n-talk")
				if orderErr != nil {
					fmt.Fprintf(os.Stderr, "\nError ordering claim: %v\n", orderErr)
					os.Exit(1)
				}
			}).Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Step 5: Output the rendered YAML
		fmt.Println("\nClaim rendered successfully!")
		if err := internal.SaveOutput(talkDestination, orderResp.Rendered); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving output: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(talkCmd)
	talkCmd.Flags().StringVar(&talkAPIURL, "api-url", "", "Base URL of the claim-machinery-api (or CLAIM_API_URL env var)")
	talkCmd.Flags().StringVar(&talkAuthToken, "api-token", "", "Auth token for claim-machinery-api (or CLAIM_API_TOKEN env var)")
	talkCmd.Flags().StringVar(&talkInstruction, "instruction", "", "Natural language description of the claim you want")
	talkCmd.Flags().StringVar(&talkDestination, "destination", "", "Output destination: stdout (default), file path, or directory")
	talkCmd.Flags().StringVar(&talkProvider, "ai-provider", "", "AI provider: openrouter or gemini (default from AI_PROVIDER env)")
	talkCmd.Flags().StringVar(&talkModel, "ai-model", "", "AI model name (default from AI_MODEL env)")
	talkCmd.Flags().StringVar(&talkBaseURL, "ai-base-url", "", "Base URL for OpenRouter API (default from AI_BASE_URL env)")
	talkCmd.Flags().BoolVarP(&talkVerbose, "verbose", "v", false, "Enable verbose output (show prompts and raw AI responses)")
}
