# k2n (/ˈkæf.kən/)

An AI-powered CLI for generating Kubernetes and infrastructure configurations from examples and rulesets.

## Overview

k2n leverages artificial intelligence to transform natural language instructions into properly formatted configuration files, including:

- Kubernetes manifests and CRDs
- Crossplane claims and compositions
- Helm values
- Terraform configurations

## Commands

| Command | Description |
|---------|-------------|
| `k2n gen` | Generate configurations using AI based on examples and rulesets |
| `k2n talk` | AI-powered conversational claim rendering via claim-machinery-api |
| `k2n version` | Show version information |

## Quick Start

```bash
# Set AI provider
export AI_PROVIDER="openrouter"
export AI_MODEL="openai/gpt-4"
export AI_API_KEY="your-api-key"

# Generate a claim from examples
k2n gen \
  --examples-dirs ./examples \
  --usecase crossplane \
  --instruction "generate a runner claim for the dagger repository"

# Or talk to claim-machinery-api
k2n talk \
  --api-url http://localhost:8080 \
  --instruction "I need a 50Gi volume in namespace production"
```

## Interactive Mode

Run `k2n` without arguments to launch the interactive TUI menu with guided configuration for both `gen` and `talk` commands.

## AI Providers

k2n supports multiple AI providers:

- **OpenRouter** - Access to multiple models (GPT-4, Deepseek, etc.)
- **Google Gemini** - Google's Gemini API

See [AI Providers](ai-providers.md) for configuration details.

## Related Documentation

- [Gen Command](gen-command.md)
- [Talk Command](talk-command.md)
- [AI Providers](ai-providers.md)
- [Architecture](architecture.md)
