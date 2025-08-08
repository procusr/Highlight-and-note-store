# MH-Go Development Environment Setup

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Setting up MH-Go Development Environment...${NC}"

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check Go installation
if command_exists go; then
    echo -e "${GREEN}✓ Go is installed: $(go version)${NC}"
else
    echo -e "${RED}✗ Go is not installed. Please install Go 1.21 or later.${NC}"
    exit 1
fi

# Check Docker installation
if command_exists docker; then
    echo -e "${GREEN}✓ Docker is installed${NC}"
else
    echo -e "${RED}✗ Docker is not installed. Please install Docker.${NC}"
    exit 1
fi

# Check Docker Compose
if command_exists docker-compose || docker compose version >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Docker Compose is available${NC}"
else
    echo -e "${RED}✗ Docker Compose is not available.${NC}"
    exit 1
fi

# Install Templ if not installed
if command_exists templ; then
    echo -e "${GREEN}✓ Templ is installed${NC}"
else
    echo -e "${YELLOW}Installing Templ...${NC}"
    go install github.com/a-h/templ/cmd/templ@latest
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Templ installed successfully${NC}"
    else
        echo -e "${RED}✗ Failed to install Templ${NC}"
    fi
fi

# Install Air for hot reloading (optional)
if command_exists air; then
    echo -e "${GREEN}✓ Air (hot reload) is installed${NC}"
else
    echo -e "${YELLOW}Installing Air for hot reloading...${NC}"
    go install github.com/cosmtrek/air@latest
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Air installed successfully${NC}"
    else
        echo -e "${YELLOW}⚠ Air installation failed (optional tool)${NC}"
    fi
fi

# Download Go modules
echo -e "${YELLOW}Downloading Go modules...${NC}"
go mod download
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Go modules downloaded${NC}"
else
    echo -e "${RED}✗ Failed to download Go modules${NC}"
    exit 1
fi

# Generate Templ templates
echo -e "${YELLOW}Generating Templ templates...${NC}"
templ generate
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Templates generated successfully${NC}"
else
    echo -e "${RED}✗ Failed to generate templates${NC}"
    exit 1
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo -e "${YELLOW}Creating .env file...${NC}"
    cp .env.example .env 2>/dev/null || cat > .env << EOF
DATABASE_URL=postgres://user:password@localhost:5433/mhgo?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-change-in-production
PORT=8080
ENVIRONMENT=development
EOF
    echo -e "${GREEN}✓ .env file created${NC}"
else
    echo -e "${GREEN}✓ .env file already exists${NC}"
fi

# Create Air configuration for hot reloading
if [ ! -f .air.toml ]; then
    echo -e "${YELLOW}Creating Air configuration...${NC}"
    cat > .air.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "templ generate && go build -o ./tmp/main cmd/server/main.go"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "node_modules"]
  exclude_file = []
  exclude_regex = ["_test.go", "_templ.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "templ"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
EOF
    echo -e "${GREEN}✓ Air configuration created${NC}"
fi

# Test build
echo -e "${YELLOW}Testing build...${NC}"
go build -o tmp/test-build cmd/server/main.go
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Build test successful${NC}"
    rm -f tmp/test-build
else
    echo -e "${RED}✗ Build test failed${NC}"
    exit 1
fi

echo -e "${GREEN}"
echo "=================================================="
echo "🎉 Development environment setup complete!"
echo "=================================================="
echo -e "${NC}"
echo -e "${BLUE}To start developing:${NC}"
echo -e "${YELLOW}1. Start the database:${NC} docker-compose up -d postgres"
echo -e "${YELLOW}2. Run with hot reload:${NC} air"
echo -e "${YELLOW}3. Or run normally:${NC} ./start-app.sh"
echo ""
echo -e "${BLUE}Access your app at:${NC} http://localhost:8080"
echo ""
echo -e "${BLUE}Database management:${NC}"
echo -e "${YELLOW}• Backup database:${NC} ./db-utils.sh backup"
echo -e "${YELLOW}• Restore database:${NC} ./db-utils.sh restore <backup_file>"
echo -e "${YELLOW}• List backups:${NC} ./db-utils.sh list-backups"
echo -e "${YELLOW}• Reset database:${NC} ./db-utils.sh reset"
echo ""
