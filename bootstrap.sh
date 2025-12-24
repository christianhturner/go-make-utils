#!/bin/bash
set -e

echo "🔧 Setting up go-make-utils development tools..."

# Check Go version (macOS compatible)
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go1\.//' | cut -d. -f1)
if [ "$GO_VERSION" -lt 24 ]; then
  echo "⚠️  Warning: Go 1.24+ recommended for optimal tool support"
  echo "Current version: $(go version)"
fi

# Add the go-make-utils as a Go tool
echo "📦 Adding go-make-utils..."
GOPROXY=direct GONOSUMDB=github.com/christianhturner/go-make-utils go get -tool github.com/christianhturner/go-make-utils@latest

# Create tools directory and wrapper script
mkdir -p tools

echo "🎁 Creating go-make-utils wrapper..."
cat >tools/go-make-utils <<'EOF'
#!/bin/bash
set -e

WD=$(dirname "$0")
WD=$(cd "$WD"; pwd)

function hash() {
  {
    go list -m -f '{{.Path}}@{{.Version}}' github.com/christianhturner/go-make-utils 2>/dev/null || echo "local"
    go env | grep -v GOGCCFLAGS
  } | shasum -a 256 | cut -f1 -d' '
}

function findpath() {
  local path="$1"
  local key="$2"
  while IFS= read -r line; do
    fk="$(echo "$line" | cut -f1 -d=)"
    if [[ "${fk}" == "${key}" ]]; then
      echo "$line" | cut -f2 -d=
      break
    fi
  done <"${path}"
}

key="$(hash)"
cache_file="${WD}/.go-make-utils.cache"

if [[ -f "${cache_file}" ]]; then
  bin="$(findpath "${cache_file}" "${key}")"
  if [ -n "${bin}" ] && [ ! -f "${bin}" ]; then
    # Remove the stale entry (macOS compatible)
    grep -v "${key}=" "${cache_file}" > "${cache_file}.tmp" 2>/dev/null || true
    mv "${cache_file}.tmp" "${cache_file}" 2>/dev/null || rm -f "${cache_file}"
    bin=""
  fi
fi

if [[ -z "${bin}" ]]; then
    bin="$(go tool -n go-make-utils)"
    echo "${key}=${bin}" >> "${cache_file}"
fi

exec "${bin}" "$@"
EOF

chmod +x tools/go-make-utils

# Create template if it doesn't exist
if [ ! -f "local-config.template.json" ]; then
  echo "📄 Creating local-config.template.json..."
  cat >local-config.template.json <<EOF
{
  "example_key": "example_value",
  "another_key": "another_value"
}
EOF
  echo "⚠️Please edit local-config.template.json with your project's configuration keys"
fi

# Update .gitignore
if [ -f ".gitignore" ]; then
  if ! grep -q "local-config.json" .gitignore 2>/dev/null; then
    echo "" >>.gitignore
    echo "# Development configuration (go-make-utils)" >>.gitignore
    echo "local-config.json" >>.gitignore
    echo ".env.mk" >>.gitignore
    echo "tools/.go-make-utils.cache" >>.gitignore
    echo "✅ Updated .gitignore"
  fi
else
  cat >.gitignore <<EOF
# Development configuration (go-make-utils)
local-config.json
.env.mk
tools/.go-make-utils.cache
EOF
  echo "✅ Created .gitignore"
fi

# Download Makefile.include
echo "📥 Downloading Makefile.include..."
curl -sSL https://raw.githubusercontent.com/christianhturner/go-make-utils/main/Makefile.include -o .go-make-utils.mk

echo ""
echo "✅ go-make-utils setup complete!"
echo ""
echo "Next steps:"
echo "  1. Edit local-config.template.json with your project's configuration keys"
echo "  2. Add this line to your Makefile:"
echo "     include .go-make-utils.mk"
echo "  3. Add 'ensure-config' as a dependency to your build targets"
echo "  4. Run 'make ensure-config' to configure your environment"
