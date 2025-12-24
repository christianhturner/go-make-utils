# go-make-utils

A reusable Go utility for managing developer-specific environment variables across Make-based projects.

## Features

- **Environment Management**: Manage developer-specific configuration with `local-config.json`
- **Template-Based**: Track required configuration keys in version control via `local-config.template.json`
- **Make Integration**: Automatically export configuration as Make variables
- **Sync Detection**: Detects when template changes and prompts for updates
- **Go 1.24 Tool Support**: Leverages Go's new tool directive for easy distribution

## Installation

### Prerequisites

- Go 1.24+ (recommended for optimal performance)
- Make

### Quick Start

In your project directory:

```bash
curl -sSL https://raw.githubusercontent.com/christianhturner/go-make-utils/main/bootstrap.sh | bash
```

Or manually:

```bash
# Add as a Go tool
go get -tool github.com/christianhturner/go-make-utils@latest

# Download and run bootstrap script
curl -sSL https://raw.githubusercontent.com/christianhturner/go-make-utils/main/bootstrap.sh -o bootstrap.sh
chmod +x bootstrap.sh
./bootstrap.sh
```

## Usage

### 1. Configure Your Project

Edit `local-config.template.json` with your project's required configuration:

```json
{
  "ada_profile": "YourProfileHere",
  "aws_region": "us-east-1",
  "log_level": "INFO"
}
```

### 2. Update Your Makefile

Add this line near the top of your Makefile:

```makefile
include .go-make-utils.mk
```

Add `ensure-config` as a dependency to your build targets:

```makefile
build: ensure-config
	go build ./...

test: ensure-config
	go test ./...
```

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
	echo "Using profile: $(ADA_PROFILE)"
```

## How It Works

1. **Template File** (`local-config.template.json`): Tracked in version control, defines required configuration keys
2. **Local Config** (`local-config.json`): Git-ignored, contains actual values for your environment
3. **Environment File** (`.env.mk`): Auto-generated, exports variables for Make

The tool ensures your local config stays in sync with the template, prompting for new keys and removing obsolete ones.

## Available Make Targets

- `make ensure-config`: Sync local config with template
- `make show-config`: Display all loaded configuration variables
- `make clean-config`: Remove local configuration files

## Development

### Building Locally

```bash
go build -o config-tool .
```

### Testing

```bash
# Create test files
echo '{"test_key": "test_value"}' > local-config.template.json

# Run the tool
./config-tool ensure .

# Verify output
cat local-config.json
cat .env.mk
```

## License

MIT

## Contributing

Contributions welcome! Please open an issue or PR.
