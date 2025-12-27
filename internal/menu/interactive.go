package menu

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

type K2NConfig struct {
	Command             string
	Instruction         string
	ExampleFiles        string
	ExamplesDirs        string
	Destination         string
	UseCase             string
	RulesetEnvFiles     string
	RulesetUseCaseFiles string
	AIProvider          string
	AIModel             string
	Verbose             bool
	PromptToAI          bool
}

// ShowInteractiveMenu displays the main menu when k2n is run without arguments
func ShowInteractiveMenu(rootCmd *cobra.Command) error {
	config := &K2NConfig{
		PromptToAI: true,
		AIProvider: getEnvOrDefault("AI_PROVIDER", "gemini"),
		AIModel:    getEnvOrDefault("AI_MODEL", ""),
	}

	// Main menu loop
	for {
		if err := showMainMenu(config); err != nil {
			return err
		}

		switch config.Command {
		case "gen":
			if err := showGenMenu(config); err != nil {
				return err
			}
			// After configuration, show what would be executed
			if err := showExecutionConfirmation(config, rootCmd); err != nil {
				return err
			}
			return nil
		case "help":
			rootCmd.Help()
			return nil
		case "exit":
			fmt.Println("\n👋 Goodbye!")
			return nil
		}
	}
}

func showMainMenu(config *K2NConfig) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🚀 K2N - AI-Based Code Generation").
				Description("What would you like to do?").
				Options(
					huh.NewOption("🎨 Generate Code (gen)", "gen"),
					huh.NewOption("❓ Show Help", "help"),
					huh.NewOption("🚪 Exit", "exit"),
				).
				Value(&config.Command),
		),
	).WithTheme(huh.ThemeCharm()).Run()
}

func showGenMenu(config *K2NConfig) error {
	// Step 1: Basic Configuration
	if err := showBasicConfig(config); err != nil {
		return err
	}

	// Step 2: AI Configuration
	if err := showAIConfig(config); err != nil {
		return err
	}

	// Step 3: Examples Configuration
	if err := showExamplesConfig(config); err != nil {
		return err
	}

	// Step 4: Rulesets Configuration (optional)
	var wantsRulesets bool
	if err := huh.NewConfirm().
		Title("Do you want to configure rulesets?").
		Description("Rulesets help guide the AI generation").
		Value(&wantsRulesets).
		Run(); err != nil {
		return err
	}

	if wantsRulesets {
		if err := showRulesetsConfig(config); err != nil {
			return err
		}
	}

	// Step 5: Advanced Options
	var wantsAdvanced bool
	if err := huh.NewConfirm().
		Title("Configure advanced options?").
		Value(&wantsAdvanced).
		Run(); err != nil {
		return err
	}

	if wantsAdvanced {
		if err := showAdvancedOptions(config); err != nil {
			return err
		}
	}

	return nil
}

func showBasicConfig(config *K2NConfig) error {
	fmt.Println("\n📋 Basic Configuration")
	fmt.Println(strings.Repeat("─", 50))

	return huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Instruction for AI").
				Description("What do you want the AI to generate?").
				Placeholder("Generate a Kubernetes deployment for a web app...").
				CharLimit(500).
				Validate(func(s string) error {
					if len(s) == 0 {
						return fmt.Errorf("instruction is required")
					}
					return nil
				}).
				Value(&config.Instruction),

			huh.NewInput().
				Title("Use Case").
				Description("Context for generation (optional)").
				Placeholder("kubernetes-deployment").
				Value(&config.UseCase),

			huh.NewSelect[string]().
				Title("Destination").
				Description("Where should the output go?").
				Options(
					huh.NewOption("📺 Standard Output (stdout)", "stdout"),
					huh.NewOption("📄 Single File", "file"),
					huh.NewOption("📁 Directory (separate files)", "directory"),
				).
				Value(&config.Destination),
		),
	).WithTheme(huh.ThemeCharm()).Run()
}

func showAIConfig(config *K2NConfig) error {
	fmt.Println("\n🤖 AI Configuration")
	fmt.Println(strings.Repeat("─", 50))

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("AI Provider").
				Options(
					huh.NewOption("Google Gemini", "gemini"),
					huh.NewOption("OpenRouter", "openrouter"),
				).
				Value(&config.AIProvider),

			huh.NewInput().
				Title("AI Model").
				Description("Leave empty for default").
				Placeholder("nousresearch/hermes-3-llama-3.1-405b:free").
				Value(&config.AIModel),

			huh.NewConfirm().
				Title("Prompt to AI?").
				Description("Send the generated content to AI for processing").
				Value(&config.PromptToAI),
		),
	).WithTheme(huh.ThemeCharm()).Run()
}

func showExamplesConfig(config *K2NConfig) error {
	fmt.Println("\n📚 Examples Configuration")
	fmt.Println(strings.Repeat("─", 50))

	var hasExamples bool
	if err := huh.NewConfirm().
		Title("Do you have example files?").
		Description("Examples help the AI understand your desired output format").
		Value(&hasExamples).
		Run(); err != nil {
		return err
	}

	if !hasExamples {
		return nil
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Example Files").
				Description("Comma-separated file paths").
				Placeholder("/path/to/example1.yaml,/path/to/example2.yaml").
				Value(&config.ExampleFiles),

			huh.NewInput().
				Title("Examples Directories").
				Description("Comma-separated directory paths").
				Placeholder("/path/to/examples/,/another/path/").
				Value(&config.ExamplesDirs),
		),
	).WithTheme(huh.ThemeCharm()).Run()
}

func showRulesetsConfig(config *K2NConfig) error {
	fmt.Println("\n📐 Rulesets Configuration")
	fmt.Println(strings.Repeat("─", 50))

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Environment Ruleset Files").
				Description("Comma-separated paths (optional)").
				Placeholder("/path/to/env-rules.yaml").
				Value(&config.RulesetEnvFiles),

			huh.NewInput().
				Title("Use Case Ruleset Files").
				Description("Comma-separated paths (optional)").
				Placeholder("/path/to/usecase-rules.yaml").
				Value(&config.RulesetUseCaseFiles),
		),
	).WithTheme(huh.ThemeCharm()).Run()
}

func showAdvancedOptions(config *K2NConfig) error {
	fmt.Println("\n⚙️  Advanced Options")
	fmt.Println(strings.Repeat("─", 50))

	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Enable verbose output?").
				Description("Show detailed logging information").
				Value(&config.Verbose),
		),
	).WithTheme(huh.ThemeCharm()).Run()
}

func showExecutionConfirmation(config *K2NConfig, rootCmd *cobra.Command) error {
	fmt.Println("\n" + strings.Repeat("═", 70))
	fmt.Println("📋 Command Summary")
	fmt.Println(strings.Repeat("═", 70))

	// Build the command string
	cmdArgs := buildCommandArgs(config)
	cmdString := "k2n gen " + strings.Join(cmdArgs, " ")

	fmt.Println("\n🔧 Generated Command:")
	fmt.Println("  " + cmdString)

	fmt.Println("\n📊 Configuration:")
	fmt.Printf("  Instruction:     %s\n", truncate(config.Instruction, 50))
	fmt.Printf("  Use Case:        %s\n", getOrEmpty(config.UseCase))
	fmt.Printf("  Destination:     %s\n", getOrEmpty(config.Destination))
	fmt.Printf("  AI Provider:     %s\n", config.AIProvider)
	fmt.Printf("  AI Model:        %s\n", getOrEmpty(config.AIModel))
	fmt.Printf("  Example Files:   %s\n", getOrEmpty(config.ExampleFiles))
	fmt.Printf("  Examples Dirs:   %s\n", getOrEmpty(config.ExamplesDirs))
	fmt.Printf("  Verbose:         %v\n", config.Verbose)
	fmt.Println(strings.Repeat("═", 70))

	var execute bool
	if err := huh.NewConfirm().
		Title("Execute this command?").
		Affirmative("Yes, run it!").
		Negative("No, exit").
		Value(&execute).
		Run(); err != nil {
		return err
	}

	if execute {
		fmt.Println("\n🚀 Executing command...\n")
		// Set the args and execute the gen command
		os.Args = append([]string{"k2n", "gen"}, cmdArgs...)
		return rootCmd.Execute()
	}

	fmt.Println("\n👋 Command not executed. Goodbye!")
	return nil
}

func buildCommandArgs(config *K2NConfig) []string {
	var args []string

	if config.Instruction != "" {
		args = append(args, "--instruction", fmt.Sprintf("%q", config.Instruction))
	}
	if config.UseCase != "" {
		args = append(args, "--usecase", config.UseCase)
	}
	if config.Destination != "" && config.Destination != "stdout" {
		args = append(args, "--destination", config.Destination)
	}
	if config.AIProvider != "" && config.AIProvider != "gemini" {
		args = append(args, "--ai-provider", config.AIProvider)
	}
	if config.AIModel != "" {
		args = append(args, "--ai-model", config.AIModel)
	}
	if config.ExampleFiles != "" {
		args = append(args, "--example-files", config.ExampleFiles)
	}
	if config.ExamplesDirs != "" {
		args = append(args, "--examples-dirs", config.ExamplesDirs)
	}
	if config.RulesetEnvFiles != "" {
		args = append(args, "--ruleset-env-files", config.RulesetEnvFiles)
	}
	if config.RulesetUseCaseFiles != "" {
		args = append(args, "--ruleset-usecase-files", config.RulesetUseCaseFiles)
	}
	if config.Verbose {
		args = append(args, "-v")
	}
	if !config.PromptToAI {
		args = append(args, "-p=false")
	}

	return args
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func getOrEmpty(s string) string {
	if s == "" {
		return "⊘ not set"
	}
	return s
}

func getEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
