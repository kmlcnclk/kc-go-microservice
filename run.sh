#!/bin/bash

set -e

# Add air to PATH automatically if not already
if ! grep -q 'export PATH="$HOME/go/bin:$PATH"' ~/.zshrc; then
  echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
  echo "✅ Added Go bin to PATH in ~/.zshrc"
fi

source ~/.zshrc

# Create .air.toml for each service
for dir in gateway orders notification-service; do
  echo "📁 Setting up $dir/.air.toml"
  mkdir -p "$dir/tmp"
  cat > "$dir/.air.toml" <<EOF
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o tmp/main ./main.go"
bin = "tmp/main"
full_bin = "tmp/main"
include_ext = ["go"]
exclude_dir = ["tmp", "vendor"]

[log]
time = true
EOF
done

# Create Procfile
cat > Procfile <<EOF
gateway-service: cd services && cd gateway-service && air
order-service: cd services && cd order-service && air
product-service: cd services && cd product-service && air
notification-service: cd services && cd notification-service && air
EOF

echo "✅ Created Procfile"

# Start docker-compose
echo "🐳 Starting Docker containers..."
docker-compose up -d

# Check/install overmind
if ! command -v overmind &> /dev/null; then
  echo "🔧 Installing overmind..."
  brew install overmind
fi

echo "🚀 Starting services with overmind..."
overmind start
