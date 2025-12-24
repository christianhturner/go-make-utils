go-make-utils

A reusable Go utility for managing developer-specific environment variables across Make-based projects.

Features
• Environment Management: Manage developer-specific configuration with local-config.json
• Template-Based: Track required configuration keys in version control via local-config.template.json
• Make Integration: Automatically export configuration as Make variables
• Sync Detection: Detects when template changes and prompts for updates
• Go 1.24 Tool Support: Leverages Go's new tool directive for easy distribution
• Production Safe: Development-only utility that doesn't interfere with CI/CD builds

Installation

Prerequisites
• Go 1.24+ (recommended for optimal performance)
• Make

Quick Start

In your project directory:

Or add to your Makefile (see Makefile Integration below):

Usage
Configure Your Project

Edit local-config.template.json with your project's required configuration:
Update Your Makefile

Add this line near the top of your Makefile (after variable definitions but before other includes):

Note: The - prefix ensures the Makefile doesn't fail if the utility isn't installed (e.g., in CI/CD environments).
First-Time Setup

Run:

You'll be prompted to provide values for each configuration key.
Use Configuration Variables

All keys from local-config.json are automatically available as uppercase Make variables:

Makefile Integration

Recommended Setup

Add these targets to your Makefile for the best developer experience:

Usage Workflow

For new developers:

For CI/CD:

Update .gitignore

The bootstrap script automatically updates your .gitignore, but if you're setting up manually, add:

How It Works
Template File (local-config.template.json): Tracked in version control, defines required configuration keys
Local Config (local-config.json): Git-ignored, contains actual values for your environment
Environment File (.env.mk): Auto-generated, exports variables for Make
Wrapper Script (tools/go-make-utils): Cached binary wrapper for fast execution

The tool ensures your local config stays in sync with the template, prompting for new keys and removing obsolete ones.

Available Make Targets

From go-make-utils.mk
• make ensure-config: Sync local config with template
• make show-config: Display all loaded configuration variables
• make clean-config: Remove local configuration files

From your project Makefile
• make bootstrap-tools: Install go-make-utils (one-time setup)
• make check-dev-tools: Verify development tools are installed
• make ensure-dev-config: Check tools and sync config (use as dependency)

Advanced Usage

Environment Variable Naming

Keys in local-config.json are automatically converted to uppercase Make variables:

Becomes:

Conditional Defaults

You can set defaults in your Makefile that are overridden by local config:

Multiple Environments

For projects with multiple environments, you can create different template files:

Troubleshooting

"Development tools not found"

Run make bootstrap-tools to install the utility.

"No rule to make target 'ensure-config'"

The utility isn't installed. Run make bootstrap-tools.

Changes to template not detected

Run make ensure-config manually to sync your local config with template changes.

Tool runs slowly

The first execution compiles the tool. Subsequent runs use a cached binary and should be fast (~10-15ms).

Development

Building Locally

Testing

Publishing a New Version

License

MIT

Contributing

Contributions welcome! Please open an issue or PR at <https://github.com/christianhturner/go-make-utils>
