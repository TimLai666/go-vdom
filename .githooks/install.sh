#!/bin/bash
#
# Installation script for go-vdom git hooks
# Run this script to set up pre-commit hooks for the repository
#

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

echo ""
echo "=========================================="
echo "  go-vdom Git Hooks Installation"
echo "=========================================="
echo ""

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    print_error "Not a git repository. Please run this script from the repository root."
    exit 1
fi

# Check if we're in the correct directory
if [ ! -f "go.mod" ] || [ ! -d ".githooks" ]; then
    print_error "Must be run from go-vdom repository root"
    exit 1
fi

# Create .git/hooks directory if it doesn't exist
mkdir -p .git/hooks

# Install pre-commit hook
print_info "Installing pre-commit hook..."

if [ -f ".git/hooks/pre-commit" ]; then
    print_warning "Existing pre-commit hook found"
    read -p "Do you want to overwrite it? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_warning "Skipping pre-commit hook installation"
    else
        cp .githooks/pre-commit .git/hooks/pre-commit
        chmod +x .git/hooks/pre-commit
        print_success "Pre-commit hook installed (overwritten)"
    fi
else
    cp .githooks/pre-commit .git/hooks/pre-commit
    chmod +x .git/hooks/pre-commit
    print_success "Pre-commit hook installed"
fi

# Build template linter
print_info "Building template linter..."
cd tools/template-linter

if go build -o template-linter; then
    print_success "Template linter built successfully"
else
    print_error "Failed to build template linter"
    cd ../..
    exit 1
fi

cd ../..

# Test the hooks
echo ""
print_info "Testing hooks installation..."

if [ -x ".git/hooks/pre-commit" ]; then
    print_success "Pre-commit hook is executable"
else
    print_error "Pre-commit hook is not executable"
    exit 1
fi

# Verify template linter
if [ -f "tools/template-linter/template-linter" ]; then
    print_success "Template linter binary found"
else
    print_error "Template linter binary not found"
    exit 1
fi

# Optional: Configure git to use .githooks directory for all hooks
echo ""
read -p "Do you want to configure git to use .githooks directory for all future hooks? (y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    git config core.hooksPath .githooks
    print_success "Git configured to use .githooks directory"
    print_info "Note: Hooks in .githooks must be executable (chmod +x)"
else
    print_info "You can manually configure this later with:"
    echo "    git config core.hooksPath .githooks"
fi

echo ""
echo "=========================================="
print_success "Git hooks installation complete!"
echo "=========================================="
echo ""
print_info "The following checks will run before each commit:"
echo "  • Code formatting (go fmt)"
echo "  • Static analysis (go vet)"
echo "  • Template linting (custom linter)"
echo ""
print_warning "To bypass hooks (not recommended), use: git commit --no-verify"
echo ""

exit 0
