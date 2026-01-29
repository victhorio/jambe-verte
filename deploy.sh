#!/bin/bash

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="jv-server"
PACKAGE_NAME="jv-deploy-$(date +%Y%m%d-%H%M%S).tar.gz"
BUILD_DIR="build"
DEPLOY_DIR="$BUILD_DIR/deploy"

echo -e "${GREEN}=== JV Server Deployment Script ===${NC}"

# Step 1: Clean and create build directories
echo -e "${YELLOW}[1/6] Preparing build directories...${NC}"
rm -rf "$BUILD_DIR"
mkdir -p "$DEPLOY_DIR"

# Step 2: Build the binary for arm64
echo -e "${YELLOW}[2/6] Building $BINARY_NAME for arm64...${NC}"
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o "$DEPLOY_DIR/$BINARY_NAME" ./cmd/jv-server

if [ ! -f "$DEPLOY_DIR/$BINARY_NAME" ]; then
    echo -e "${RED}Error: Failed to build binary${NC}"
    exit 1
fi

# Step 3: Copy required directories
echo -e "${YELLOW}[3/6] Copying required files...${NC}"
cp -r templates "$DEPLOY_DIR/"
cp -r static "$DEPLOY_DIR/"
cp -r content "$DEPLOY_DIR/"
cp package.json tailwind.config.js "$DEPLOY_DIR/"

# Step 4: Create deployment package
echo -e "${YELLOW}[4/6] Creating deployment package...${NC}"
cd "$BUILD_DIR"
COPYFILE_DISABLE=1 tar -czf "$PACKAGE_NAME" --exclude='._*' --exclude='.DS_Store' --no-xattrs --no-mac-metadata deploy/
cd ..

# Step 5: Transfer to server
echo -e "${YELLOW}[5/6] Transferring package to server...${NC}"
PACKAGE_PATH="$BUILD_DIR/$PACKAGE_NAME"

# Check if JV_SERVER_IP is set
if [ -z "${JV_SERVER_IP:-}" ]; then
    echo -e "${RED}Error: JV_SERVER_IP environment variable is not set${NC}"
    echo -e "Please set it with: export JV_SERVER_IP=your.server.ip"
    exit 1
fi

# Transfer using rsync
rsync -avz -e "ssh -i ~/.ssh/id_ed25519" "$PACKAGE_PATH" "root@$JV_SERVER_IP:/tmp/"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Package transferred successfully!${NC}"
else
    echo -e "${RED}Error: Failed to transfer package${NC}"
    exit 1
fi

# Step 6: Clean up and report
echo -e "${YELLOW}[6/6] Cleaning up...${NC}"
rm -rf "$DEPLOY_DIR"

# Final report
PACKAGE_SIZE=$(du -h "$PACKAGE_PATH" | cut -f1)

echo -e "${GREEN}✓ Deployment completed successfully!${NC}"
echo -e "Package: ${GREEN}$PACKAGE_NAME${NC} (${PACKAGE_SIZE})"
echo -e "Location on server: ${GREEN}/tmp/$PACKAGE_NAME${NC}"
echo
echo -e "Next steps on the server:"
echo -e "> ${YELLOW}cd /tmp && tar -xzf $PACKAGE_NAME${NC}"
echo -e "> ${YELLOW}rm -rf ~/deploy && mv deploy ~/${NC}"
echo -e "> ${YELLOW}cd ~/deploy && bun install${NC}"
echo -e "> ${YELLOW}systemctl restart jv-server${NC}"