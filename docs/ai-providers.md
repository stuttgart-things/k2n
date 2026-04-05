# AI Providers

k2n supports multiple AI providers through a pluggable provider architecture. Both the `gen` and `talk` commands use the same provider configuration.

## OpenRouter

[OpenRouter](https://openrouter.ai/) provides access to multiple AI models through a unified API.

### Configuration

```bash
export AI_PROVIDER="openrouter"
export AI_API_KEY="sk-or-..."
export AI_MODEL="openai/gpt-4"                                    # optional, default: openai/gpt-3.5-turbo
export AI_BASE_URL="https://openrouter.ai/api/v1/chat/completions" # optional, this is the default
```

Or via CLI flags:

```bash
k2n gen --ai-provider openrouter --ai-model "openai/gpt-4" ...
k2n talk --ai-provider openrouter --ai-model "openai/gpt-4" ...
```

### Supported Models

Any model available on OpenRouter can be used. Some popular options:

| Model | ID | Notes |
|-------|----|-------|
| GPT-4 | `openai/gpt-4` | High quality |
| GPT-3.5 Turbo | `openai/gpt-3.5-turbo` | Default, fast |
| Deepseek R1 | `deepseek/deepseek-r1-0528:free` | Free tier |

### Custom Base URL

For enterprise deployments or self-hosted OpenRouter-compatible endpoints:

```bash
export AI_BASE_URL="https://your-internal-endpoint/v1/chat/completions"
```

## Google Gemini

Google's [Gemini](https://ai.google.dev/) API for generative AI.

### Configuration

```bash
export AI_PROVIDER="gemini"
export AI_API_KEY="your-gemini-api-key"
```

Or via CLI flags:

```bash
k2n gen --ai-provider gemini ...
k2n talk --ai-provider gemini ...
```

### Model

k2n uses the `gemini-3-pro-preview` model by default. No additional model configuration is needed.

## Priority Order

Configuration is resolved in this order (highest priority first):

1. CLI flags (`--ai-provider`, `--ai-model`, `--ai-base-url`)
2. Environment variables (`AI_PROVIDER`, `AI_MODEL`, `AI_BASE_URL`)
3. Default values (provider: `openrouter`, model: `openai/gpt-3.5-turbo`)
