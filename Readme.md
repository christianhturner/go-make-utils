# go-make-utils

A reusable Go utility for managing developer-specific environment variables across Make-based projects.

## Features

- **Environment Management**: Manage developer-specific configuration with `local-config.json`
- **Template-Based**: Track required configuration keys in version control via `local-config.template.json`
- **Make Integration**: Automatically export configuration as Make variables
- **Sync Detection**: Detects when template changes and prompts for updates
- **Go 1.24 Tool Support**: Leverages Go's new tool directive for easy distribution
- **Production Safe**: Development-only utility that doesn't interfere with CI/CD builds

## Installation

### Prerequisites

- Go 1.24+ (recommended for optimal performance)
- Make

### Quick Start

In your project directory:

```bash
curl -sSL https://raw.githubusercontent.com/christianhturner/go-make-utils/main/bootstrap.sh | bash
```

Or add to your Makefile (see [Makefile Integration](#makefile-integration) below):

```bash
make bootstrap-tools
```

## Usage

### 1. Configure Your Project

Edit `local-config.template.json` with your project's required configuration:

```json
{
  "aws_region": "us-east-1",
  "log_level": "INFO",
  "api_endpoint": "https://api.example.com"
}
```

### 2. Update Your Makefile

Add this line near the top of your Makefile (after variable definitions but before other includes):

```makefile
-include tools/go-make-utils.mk
```

**Note**: The `-` prefix ensures the Makefile doesn't fail if the utility isn't installed (e.g., in CI/CD environments).

### 3. First-Time Setup

Run:

```bash
make ensure-config
```

You'll be prompted to provide values for each configuration key.

### 4. Use Configuration Variables

All keys from `local-config.json` are automatically available as uppercase Make variables:

```makefile
deploy: ensure-config
 aws s3 sync ./dist s3://my-bucket --region $(AWS_REGION)
 echo "Log level: $(LOG_LEVEL)"
```

## Makefile Integration

### Recommended Setup

Add these targets to your Makefile for the best developer experience:

```makefile
# Your existing variables
PROJECT_NAME := my-project
BUILD_DIR := ./build

# Include the utility (won't fail if not present)
-include tools/go-make-utils.mk

# Bootstrap target to set up development tools
.PHONY: bootstrap-tools
bootstrap-tools:
 @echo "🔧 Setting up development tools..."
 @curl -sSL https://raw.githubusercontent.com/christianhturner/go-make-utils/main/bootstrap.sh | bash
 @echo ""
 @echo "✅ Bootstrap complete! Development tools are ready."
 @echo "Run 'make ensure-config' to configure your local environment."

# Helper to check if tools are set up
.PHONY: check-dev-tools
check-dev-tools:
 @if [ ! -f "tools/go-make-utils" ]; then \
  echo "❌ Development tools not found."; \
  echo "Run 'make bootstrap-tools' to set up your development environment."; \
  exit 1; \
 fi
 @if [ ! -f "tools/go-make-utils.mk" ]; then \
  echo "❌ Development utilities Makefile not found."; \
  echo "Run 'make bootstrap-tools' to set up your development environment."; \
  exit 1; \
 fi

# Wrapper target that checks tools and runs ensure-config
.PHONY: ensure-dev-config
ensure-dev-config: check-dev-tools ensure-config

# Your development targets - add ensure-dev-config as a dependency
.PHONY: dev-build
dev-build: ensure-dev-config
 go build -o $(BUILD_DIR)/$(PROJECT_NAME) ./...

.PHONY: dev-test
dev-test: ensure-dev-config
 go test ./...

.PHONY: dev-run
dev-run: ensure-dev-config
 go run ./cmd/$(PROJECT_NAME)

# Production targets - NO dev-config dependency
.PHONY: build
build:
 go build -o $(BUILD_DIR)/$(PROJECT_NAME) ./...

.PHONY: test
test:
 go test ./...
```

### Usage Workflow

**For new developers:**

```bash
# Clone the repo
git clone <your-repo>
cd <your-repo>

# Set up development tools (one-time)
make bootstrap-tools

# Configure your local environment (one-time)
make ensure-config

# Now use development targets
make dev-build
make dev-test
```

**For CI/CD:**

```bash
# Production builds work without any dev tools
make build
make test
```

### Update .gitignore

The bootstrap script automatically updates your `.gitignore`, but if you're setting up manually, add:

```gitignore
# Development configuration (go-make-utils)
local-config.json
.env.mk
tools/.go-make-utils.cache
tools/go-make-utils.mk
```

## How It Works

1. **Template File** (`local-config.template.json`): Tracked in version control, defines required configuration keys
2. **Local Config** (`local-config.json`): Git-ignored, contains actual values for your environment
3. **Environment File** (`.env.mk`): Auto-generated, exports variables for Make
4. **Wrapper Script** (`tools/go-make-utils`): Cached binary wrapper for fast execution

The tool ensures your local config stays in sync with the template, prompting for new keys and removing obsolete ones.

## Available Make Targets

### From go-make-utils.mk

- `make ensure-config`: Sync local config with template
- `make show-config`: Display all loaded configuration variables
- `make clean-config`: Remove local configuration files

### From your project Makefile

- `make bootstrap-tools`: Install go-make-utils (one-time setup)
- `make check-dev-tools`: Verify development tools are installed
- `make ensure-dev-config`: Check tools and sync config (use as dependency)

## Advanced Usage

### Environment Variable Naming

Keys in `local-config.json` are automatically converted to uppercase Make variables:

```json
{
  "aws_region": "us-west-2",
  "api_key": "secret123",
  "debug_mode": "true"
}
```

Becomes:

```makefile
$(AWS_REGION)    # us-west-2
$(API_KEY)       # secret123
$(DEBUG_MODE)    # true
```

### Conditional Defaults

You can set defaults in your Makefile that are overridden by local config:

```makefile
# Default values
AWS_REGION ?= us-east-1
LOG_LEVEL ?= INFO

# Include local config (will override defaults if present)
-include tools/go-make-utils.mk
```

### Multiple Environments

For projects with multiple environments, you can create different template files:

```bash
# Development
make ensure-config  # Uses local-config.template.json

# Staging (manual approach)
cp local-config.staging.json local-config.json
```

## Troubleshooting

### "Development tools not found"

Run `make bootstrap-tools` to install the utility.

### "No rule to make target 'ensure-config'"

The utility isn't installed. Run `make bootstrap-tools`.

### Changes to template not detected

Run `make ensure-config` manually to sync your local config with template changes.

### Tool runs slowly

The first execution compiles the tool. Subsequent runs use a cached binary and should be fast (~10-15ms).

## Development

### Building Locally

```bash
cd ~/Development/go-make-utils
go build -o config-tool.
```

### Testing

```bash
# Create test files
echo '{"test_key": "test_value"}' > local-config.template.json

# Run the tool
./config-tool ensure.

# Verify output
cat local-config.json
cat .env.mk
```

### Publishing a New Version

```bash
git tag v0.x.x
git push origin v0.x.x
```

## License

MIT

## Contributing

Contributions welcome! Please open an issue or PR at <https://github.com/christianhturner/go-make-utils>
