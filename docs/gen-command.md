# Gen Command

The `gen` command is the core code generation engine of k2n. It uses AI to generate configurations from code examples and optional rulesets.

## Usage

```bash
k2n gen [flags]
```

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--instruction` | string | | Natural language instruction for AI (required) |
| `--usecase` | string | | Context/technology for generation |
| `--examples-dirs` | string | | Comma-separated directories with example files |
| `--example-files` | string | | Comma-separated example file paths |
| `--example-file-ext` | string | `.yaml,.tf` | Allowed file extensions |
| `--ruleset-env-dir` | string | | Directory with environment rulesets |
| `--ruleset-usecase-dir` | string | | Directory with use-case rulesets |
| `--ruleset-env-files` | string | | Comma-separated environment ruleset files |
| `--ruleset-usecase-files` | string | | Comma-separated use-case ruleset files |
| `--destination` | string | stdout | Output: stdout, file path, or directory |
| `--ai-provider` | string | openrouter | AI provider: `openrouter` or `gemini` |
| `--ai-model` | string | | Model name for the AI provider |
| `--ai-base-url` | string | | Base URL for OpenRouter API |
| `--verbose`, `-v` | bool | false | Enable verbose output |
| `--prompt-to-ai`, `-p` | bool | true | Send prompt to AI |

## How It Works

1. **Load examples** from directories or file paths
2. **Load rulesets** (environment and use-case specific constraints)
3. **Build prompt** combining role, rules, examples, and instruction
4. **Call AI** provider with the constructed prompt
5. **Output** the result to stdout, file, or directory

## Examples

### Generate with examples and rulesets

```bash
k2n gen \
  --examples-dirs _examples/examples \
  --ruleset-env-dir _examples/ruleset-env \
  --ruleset-usecase-dir _examples/ruleset-runner \
  --usecase crossplane \
  --instruction "generate a runner claim for the dagger repository"
```

### Preview the prompt without calling AI

```bash
k2n gen \
  --examples-dirs _examples/examples \
  --usecase crossplane \
  --instruction "generate a deployment for nginx" \
  --verbose=true \
  --prompt-to-ai=false
```

### Save output to a directory

```bash
k2n gen \
  --examples-dirs _examples/examples \
  --instruction "generate helm values for a web application" \
  --destination /tmp/output/
```

## Examples and Rulesets

### Examples

Example files serve as few-shot learning material for the AI. Place your reference configurations in a directory and point to them with `--examples-dirs` or `--example-files`.

### Rulesets

Rulesets are constraint files that guide the AI generation:

- **Environment rulesets** (`--ruleset-env-dir`): Environment-specific rules like versions, namespaces, or cluster configurations
- **Use-case rulesets** (`--ruleset-usecase-dir`): Technology-specific rules like Crossplane runner specifications

## Output Modes

- **stdout** (default): Print generated output to terminal
- **Single file**: `--destination /tmp/output.yaml` saves all content to one file
- **Directory**: `--destination /tmp/output/` parses the AI output by `---` delimiter and saves each file separately
