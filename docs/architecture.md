# Architecture

## Project Structure

```
k2n/
├── main.go                       # Entry point
├── cmd/
│   ├── root.go                   # Root command, interactive menu
│   ├── gen.go                    # Gen command
│   ├── talk.go                   # Talk command
│   └── version.go                # Version command
├── internal/
│   ├── ai/
│   │   ├── provider.go           # Provider abstraction and factory
│   │   ├── gemini.go             # Google Gemini implementation
│   │   └── openrouter.go         # OpenRouter implementation
│   ├── menu/
│   │   └── interactive.go        # Interactive TUI menu
│   ├── talk/
│   │   ├── client.go             # claim-machinery-api HTTP client
│   │   └── conversation.go       # AI conversation logic and prompt building
│   ├── examples.go               # Example file loading
│   ├── ruleset.go                # Ruleset loading
│   ├── prompt.go                 # Prompt construction for gen
│   ├── output.go                 # Output handling (stdout, file, directory)
│   └── print.go                  # Terminal UI (banner, tables)
├── _examples/                    # Example files and rulesets
├── docs/                         # MkDocs documentation
├── catalog-info.yaml             # Backstage catalog entry
└── mkdocs.yml                    # MkDocs configuration
```

## Component Overview

### CLI Layer (`cmd/`)

Built with [Cobra](https://github.com/spf13/cobra). Handles flag parsing, environment variable resolution, and orchestrates the workflow for each command.

### AI Provider Layer (`internal/ai/`)

Pluggable provider architecture with a common `AIProvider` interface:

```go
type AIProvider interface {
    Call(apiKey, prompt string) (string, error)
}
```

New providers can be added by implementing this interface and registering them in the factory.

### Talk Layer (`internal/talk/`)

Two components:

- **Client**: HTTP client for the claim-machinery-api REST API (list templates, get template, order claim)
- **Conversation**: Builds AI prompts from template metadata and parses structured JSON responses

### Gen Pipeline

```
Examples + Rulesets → BuildPrompt() → AI Provider → SaveOutput()
```

### Talk Pipeline

```
claim-machinery-api → BuildSystemPrompt() → AI Provider → ParseAIResponse() → OrderClaim() → SaveOutput()
```

### Interactive Menu (`internal/menu/`)

Built with [Charmbracelet Huh](https://github.com/charmbracelet/huh) for terminal forms. Provides step-by-step configuration wizards for both `gen` and `talk` commands.

## Dependencies

| Dependency | Purpose |
|------------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/charmbracelet/huh` | Interactive terminal forms |
| `github.com/charmbracelet/huh/spinner` | Loading spinners |
| `github.com/pterm/pterm` | Terminal styling and tables |
| `go.hein.dev/go-version` | Version information |
