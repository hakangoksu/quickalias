#!/bin/bash

set -e # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
GITHUB_REPO="hakangoksu/quickalias"
REPO_NAME="quickalias"
BINARY_NAME="qq"
INSTALL_DIR="/usr/local/bin"
GLOBAL_CONFIG_DIR="/etc/quickalias"
COMPLETION_DIR="/usr/share/bash-completion/completions"
ZSH_COMPLETION_DIR="/usr/share/zsh/site-functions"
TEMP_DIR=$(mktemp -d) # Temporary directory for cloning

# Print colored output
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if running as root
check_root() {
    if [[ $EUID -eq 0 ]]; then
        print_error "Please don't run this script as root. It will ask for sudo when needed."
        exit 1
    fi
}

# Detect distribution
detect_distro() {
    if [[ -f /etc/os-release ]]; then
        . /etc/os-release
        DISTRO=$ID
    elif [[ -f /etc/arch-release ]]; then
        DISTRO="arch"
    elif [[ -f /etc/debian_version ]]; then
        DISTRO="debian"
    elif [[ -f /etc/redhat-release ]]; then
        DISTRO="rhel"
    else
        DISTRO="unknown"
    fi

    print_info "Detected distribution: $DISTRO"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go first."
        case $DISTRO in
            "arch")
                print_info "On Arch Linux, run: sudo pacman -S go"
                ;;
            "ubuntu"|"debian")
                print_info "On Ubuntu/Debian, run: sudo apt install golang-go"
                ;;
            *)
                print_info "Please install Go from https://golang.org/dl/"
                ;;
        esac
        exit 1
    fi

    GO_VERSION=$(go version | cut -d' ' -f3)
    print_success "Go is installed: $GO_VERSION"
}

# Install dependencies based on distribution
install_dependencies() {
    case $DISTRO in
        "arch")
            print_info "Installing dependencies on Arch Linux..."
            if ! pacman -Qi base-devel &> /dev/null; then
                sudo pacman -S --needed base-devel
            fi
            print_success "Dependencies installed"
            ;;
        "ubuntu"|"debian")
            print_warning "Ubuntu/Debian support will be added in future versions"
            print_info "For now, ensure you have: build-essential"
            ;;
        *)
            print_warning "Unknown distribution. Skipping dependency installation."
            ;;
    esac
}

# Clone the repository
clone_repo() {
    print_info "Cloning QuickAlias repository into $TEMP_DIR..."
    if ! command -v git &> /dev/null; then
        print_error "Git is not installed. Please install Git first."
        exit 1
    fi
    git clone "https://github.com/$GITHUB_REPO.git" "$TEMP_DIR/$REPO_NAME"
    if [[ ! -d "$TEMP_DIR/$REPO_NAME" ]]; then
        print_error "Failed to clone repository."
        exit 1
    fi
    print_success "Repository cloned successfully."
    cd "$TEMP_DIR/$REPO_NAME"
}

# Build the binary
build_binary() {
    print_info "Building QuickAlias binary..."

    if [[ ! -f "main.go" ]]; then
        print_error "main.go not found in the cloned repository. Build failed."
        exit 1
    fi

    go build -ldflags="-s -w" -o "$BINARY_NAME"

    if [[ ! -f "$BINARY_NAME" ]]; then
        print_error "Build failed"
        exit 1
    fi

    print_success "Binary built successfully"
}

# Install binary
install_binary() {
    print_info "Installing binary to $INSTALL_DIR..."

    sudo install -m 755 "$BINARY_NAME" "$INSTALL_DIR/"

    # Verify installation
    if command -v qq &> /dev/null; then
        print_success "Binary installed successfully"
        print_info "Version: $(qq version)"
    else
        print_error "Binary installation failed"
        exit 1
    fi
}

# Create global config directory
create_global_config() {
    print_info "Creating global configuration directory..."

    sudo mkdir -p "$GLOBAL_CONFIG_DIR"
    sudo chmod 755 "$GLOBAL_CONFIG_DIR"

    # Create empty global aliases file
    if [[ ! -f "$GLOBAL_CONFIG_DIR/aliases.json" ]]; then
        echo "[]" | sudo tee "$GLOBAL_CONFIG_DIR/aliases.json" > /dev/null
        sudo chmod 644 "$GLOBAL_CONFIG_DIR/aliases.json"
    fi

    print_success "Global configuration directory created"
}

# Install shell completions
install_completions() {
    print_info "Installing shell completions..."

    # Create bash completion
    sudo mkdir -p "$COMPLETION_DIR"
    cat << 'EOF' | sudo tee "$COMPLETION_DIR/qq" > /dev/null
_qq_completion() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    case ${prev} in
        qq)
            opts="add set remove remove-global list search control status setup init config version help"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        config)
            opts="reset backup export import"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        remove|remove-global)
            # Complete with existing alias names
            local aliases=$(qq list 2>/dev/null | grep -E "^\s+\w+" | awk '{print $1}' | tr -d '→')
            COMPREPLY=( $(compgen -W "${aliases}" -- ${cur}) )
            return 0
            ;;
    esac
}

complete -F _qq_completion qq
EOF

    # Create zsh completion if zsh is installed
    if command -v zsh &> /dev/null; then
        sudo mkdir -p "$ZSH_COMPLETION_DIR"
        cat << 'EOF' | sudo tee "$ZSH_COMPLETION_DIR/_qq" > /dev/null
#compdef qq

_qq() {
    local context state line

    _arguments -C \
        '1: :->commands' \
        '*: :->args' \
        && return 0

    case $state in
        commands)
            _values 'qq commands' \
                'add[Add user-level alias]' \
                'set[Add global alias]' \
                'remove[Remove user alias]' \
                'remove-global[Remove global alias]' \
                'list[List all aliases]' \
                'search[Search aliases]' \
                'control[Show status]' \
                'status[Show status]' \
                'setup[Setup shell integration]' \
                'init[Initialize aliases]' \
                'config[Configuration commands]' \
                'version[Show version]' \
                'help[Show help]'
            ;;
        args)
            case $words[2] in
                config)
                    _values 'config commands' \
                        'reset[Reset configuration]' \
                        'backup[Show backups]' \
                        'export[Export aliases]' \
                        'import[Import aliases]'
                    ;;
            esac
            ;;
    esac
}

_qq "$@"
EOF
    fi

    print_success "Shell completions installed"
}

# Create man page
create_man_page() {
    print_info "Creating man page..."

    sudo mkdir -p /usr/share/man/man1

    cat << 'EOF' | sudo tee /usr/share/man/man1/qq.1 > /dev/null
.TH QQ 1 "2024" "QuickAlias 1.0.0" "User Commands"
.SH NAME
qq \- Cross-shell command alias manager
.SH SYNOPSIS
.B qq
[\fIOPTION\fR] [\fICOMMAND\fR] [\fIARGS\fR...]
.SH DESCRIPTION
QuickAlias (qq) is a cross-shell command alias manager that allows you to create, manage, and synchronize command aliases across different shell environments.
.SH COMMANDS
.TP
.B add \fIalias\fR "\fIcommand\fR"
Add a user-level alias
.TP
.B set \fIalias\fR "\fIcommand\fR"
Add a global alias (requires sudo)
.TP
.B remove \fIalias\fR
Remove a user alias
.TP
.B remove-global \fIalias\fR
Remove a global alias
.TP
.B list [\fIkeyword\fR]
List all aliases or filter by keyword
.TP
.B search \fIkeyword\fR
Search in alias names and commands
.TP
.B control
Show status and conflicts
.TP
.B setup
Setup shell integration
.TP
.B init
Initialize aliases (used by shell)
.TP
.B config \fIsubcommand\fR
Configuration management
.TP
.B version
Show version information
.TP
.B help
Show help information
.SH FILES
.TP
.I ~/.config/quickalias/aliases.json
User aliases configuration
.TP
.I ~/.config/quickalias/config.json
User configuration
.TP
.I /etc/quickalias/aliases.json
Global aliases configuration
.SH AUTHOR
QuickAlias development team
.SH SEE ALSO
.BR bash (1),
.BR zsh (1),
.BR fish (1)
EOF

    # Update man database
    sudo mandb &> /dev/null || true

    print_success "Man page created"
}

# Create desktop entry (for GUI applications that might want to use it)
create_desktop_entry() {
    print_info "Creating desktop entry..."

    sudo mkdir -p /usr/share/applications

    cat << 'EOF' | sudo tee /usr/share/applications/quickalias.desktop > /dev/null
[Desktop Entry]
Name=QuickAlias
Comment=Cross-shell command alias manager
Exec=qq
Icon=utilities-terminal
Terminal=true
Type=Application
Categories=System;TerminalEmulator;
Keywords=alias;command;shell;terminal;
EOF

    print_success "Desktop entry created"
}

# Setup for Arch Linux package management
setup_arch_integration() {
    if [[ "$DISTRO" == "arch" ]]; then
        print_info "Setting up Arch Linux integration..."

        # Add to PATH if not already there
        if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
            print_warning "$INSTALL_DIR not in PATH. This is unusual but should be fine."
        fi

        print_success "Arch Linux integration complete"
    fi
}

# Post-installation setup
post_install() {
    print_info "Running post-installation setup..."

    # Clean up temporary build artifacts
    if [[ -d "$TEMP_DIR" ]]; then
        print_info "Cleaning up temporary files..."
        rm -rf "$TEMP_DIR"
    fi

    print_success "Installation completed successfully!"
    echo
    print_info "Next steps:"
    echo "  1. Run: ${GREEN}qq setup${NC}   (to configure shell integration)"
    echo "  2. Restart your terminal or run: ${GREEN}source ~/.bashrc${NC}"
    echo "  3. Start adding aliases: ${GREEN}qq add ll 'ls -la'${NC}"
    echo "  4. View help: ${GREEN}qq help${NC}"
    echo "  5. Check man page: ${GREEN}man qq${NC}"
    echo
    print_info "Shell completions are available for bash and zsh"
    print_info "Global aliases require sudo: ${GREEN}sudo qq set <alias> \"<command>\"${NC}"
}

# Uninstall function
uninstall() {
    print_info "Uninstalling QuickAlias..."

    # Remove binary
    if [[ -f "$INSTALL_DIR/$BINARY_NAME" ]]; then
        sudo rm "$INSTALL_DIR/$BINARY_NAME"
        print_success "Binary removed"
    fi

    # Remove completions
    sudo rm -f "$COMPLETION_DIR/qq"
    sudo rm -f "$ZSH_COMPLETION_DIR/_qq"
    print_success "Completions removed"

    # Remove man page
    sudo rm -f /usr/share/man/man1/qq.1
    sudo mandb &> /dev/null || true
    print_success "Man page removed"

    # Remove desktop entry
    sudo rm -f /usr/share/applications/quickalias.desktop
    print_success "Desktop entry removed"

    # Ask about config removal
    read -p "Remove global configuration? (y/N): " -r
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        sudo rm -rf "$GLOBAL_CONFIG_DIR"
        print_success "Global configuration removed"
    fi

    print_success "QuickAlias uninstalled successfully"
    print_warning "User configurations in ~/.config/quickalias were preserved"
}

# Main installation flow
main() {
    echo "🚀 QuickAlias Installation Script"
    echo "================================="
    echo

    # Handle uninstall
    if [[ "$1" == "--uninstall" ]]; then
        uninstall
        exit 0
    fi

    # Show help
    if [[ "$1" == "--help" || "$1" == "-h" ]]; then
        echo "Usage: $0 [OPTIONS]"
        echo
        echo "Options:"
        echo "  --uninstall    Uninstall QuickAlias"
        echo "  --help, -h     Show this help"
        echo
        echo "This script will:"
        echo "  • Clone the QuickAlias repository from GitHub"
        echo "  • Build the QuickAlias binary"
        echo "  • Install it to $INSTALL_DIR"
        echo "  • Set up shell completions"
        echo "  • Create man page"
        echo "  • Configure system integration"
        echo "  • Clean up downloaded repository files"
        exit 0
    fi

    check_root
    detect_distro

    # Only proceed with supported distributions for now
    if [[ "$DISTRO" == "arch" ]]; then
        print_warning "For Arch Linux, it's recommended to build and install QuickAlias via a PKGBUILD file."
        print_info "This script is primarily for manual installations on other distributions."
        read -p "Do you still want to proceed with this script for manual installation on Arch? (y/N): " -r
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Installation cancelled. Please consider creating a PKGBUILD for Arch Linux."
            exit 0
        fi
    elif [[ "$DISTRO" == "ubuntu" || "$DISTRO" == "debian" ]]; then
        print_warning "Ubuntu/Debian support will be added in future versions for a more guided experience."
        print_info "For now, ensure you have: build-essential (sudo apt install build-essential)"
        read -p "Continue anyway with manual installation? (y/N): " -r
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Installation cancelled."
            exit 0
        fi
    else
        print_warning "Unknown distribution. This script will attempt a generic installation."
        read -p "Continue anyway? (y/N): " -r
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Installation cancelled."
            exit 0
        fi
    fi

    clone_repo # Call this before check_go and build_binary to get the source
    check_go
    install_dependencies
    build_binary
    install_binary
    create_global_config
    install_completions
    create_man_page
    create_desktop_entry
    setup_arch_integration
    post_install
}

# Run main function with all arguments
main "$@"
