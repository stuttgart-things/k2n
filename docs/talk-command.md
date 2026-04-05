# Talk Command

The `talk` command provides AI-powered conversational claim rendering via the [claim-machinery-api](https://github.com/stuttgart-things/claim-machinery-api). Describe what you need in natural language, and k2n will select the right template, fill in the parameters, and render the claim.

## Usage

```bash
k2n talk [flags]
```

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--api-url` | string | | Base URL of the claim-machinery-api (or `CLAIM_API_URL` env var) |
| `--api-token` | string | | Auth token for claim-machinery-api (or `CLAIM_API_TOKEN` env var) |
| `--instruction` | string | | Natural language description of the claim you want |
| `--destination` | string | stdout | Output: stdout, file path, or directory |
| `--ai-provider` | string | | AI provider: `openrouter` or `gemini` |
| `--ai-model` | string | | AI model name |
| `--ai-base-url` | string | | Base URL for OpenRouter API |
| `--verbose`, `-v` | bool | false | Show prompts and raw AI responses |

## How It Works

```
┌─────────┐     ┌──────────────────────┐     ┌─────────┐     ┌──────────────────────┐
│  User    │────>│  k2n talk            │────>│   AI    │────>│  claim-machinery-api │
│ "I need  │     │  1. Fetch templates  │     │ Select  │     │  Render claim YAML   │
│  a 50Gi  │     │  2. Build prompt     │     │ template│     │  via KCL             │
│  volume" │     │  3. Call AI          │     │ + params│     │                      │
└─────────┘     └──────────────────────┘     └─────────┘     └──────────────────────┘
```

1. **Fetch templates** from the claim-machinery-api (`GET /api/v1/claim-templates`)
2. **Build prompt** that includes all available templates with their parameter schemas
3. **Call AI** to match the user's natural language instruction to a template and fill parameters
4. **Parse AI response** to extract the selected template name and parameter values
5. **Order the claim** via `POST /api/v1/claim-templates/{name}/order`
6. **Output** the rendered YAML

## Examples

### Basic usage with OpenRouter

```bash
export AI_PROVIDER="openrouter"
export AI_MODEL="openai/gpt-4"
export AI_API_KEY="sk-or.."

k2n talk \
  --api-url http://localhost:8080 \
  --instruction "I need a 50Gi persistent volume in namespace production with fast-ssd storage class"
```

### Using Gemini

```bash
export AI_PROVIDER="gemini"
export AI_API_KEY="your-gemini-api-key"

k2n talk \
  --api-url http://localhost:8080 \
  --instruction "create a harbor project for team-alpha"
```

### Save rendered claim to file

```bash
k2n talk \
  --api-url http://localhost:8080 \
  --instruction "give me a vsphere vm with 4 cpus and 8gb ram" \
  --destination /tmp/vm-claim.yaml
```

### Verbose mode for debugging

```bash
k2n talk \
  --api-url http://localhost:8080 \
  --instruction "create a flux kustomization for the app repo" \
  --verbose
```

This shows the full AI prompt, the raw AI response, and the parsed parameters before rendering.

### Using environment variables

```bash
export CLAIM_API_URL="http://localhost:8080"
export CLAIM_API_TOKEN="optional-auth-token"
export AI_PROVIDER="openrouter"
export AI_MODEL="openai/gpt-4"
export AI_API_KEY="sk-or.."

k2n talk --instruction "I need storage for my application"
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `CLAIM_API_URL` | Base URL of the claim-machinery-api |
| `CLAIM_API_TOKEN` | Optional auth token for the API |
| `AI_API_KEY` | API key for the AI provider |
| `AI_PROVIDER` | AI provider: `openrouter` or `gemini` |
| `AI_MODEL` | Model name for the AI provider |
| `AI_BASE_URL` | Custom base URL for OpenRouter |

## Prerequisites

- A running instance of [claim-machinery-api](https://github.com/stuttgart-things/claim-machinery-api) with loaded claim templates
- An AI provider API key (OpenRouter or Gemini)
