# SSH Manager (sshman)

A comprehensive SSH key management tool built in Go that simplifies generating, managing, and organizing SSH keys for different Git providers and services.

## Features

- **SSH Key Generation**: Create ED25519 and RSA SSH key pairs with custom purposes
- **Provider Support**: Built-in configurations for GitHub, GitLab, Bitbucket
- **SSH Agent Management**: Add, remove, list, and clear keys from SSH agent
- **Automatic SSH Config**: Automatically updates SSH config with host aliases and key associations
- **Key Listing**: Display all SSH keys with their status (loaded/not loaded in agent)
- **Key Deletion**: Remove keys and clean up from agent and filesystem
- **Purpose-based Organization**: Organize keys by purpose (work, personal, etc.)

## Installation

### Prerequisites

- Go 1.24.4 or later

### Build from Source

```bash
git clone https://github.com/residwi/sshman.git
cd sshman
go build -o sshman
```

### Install to PATH

```bash
go install github.com/residwi/sshman@latest
```

## Quick Start

1. **Create an SSH key for GitHub:**

```bash
sshman create github --email your@email.com --purpose work
```

1. **List all SSH keys:**

```bash
sshman list
```

1. **Manage SSH agent:**

```bash
sshman agent list
sshman agent add id_ed25519_github_work
```

## Usage

### Creating SSH Keys

#### For Git Providers

Create keys for popular Git providers with automatic configuration:

```bash
# GitHub
sshman create github --email your@email.com --purpose work

# GitLab
sshman create gitlab --email your@email.com --purpose personal

# Bitbucket
sshman create bitbucket --email your@email.com --purpose project
```

#### For Generic Hosts

Create keys for custom hosts:

```bash
sshman create generic --email your@email.com --user myuser --hostname example.com --purpose server
```

#### Key Type Options

```bash
# ED25519 (default, recommended)
sshman create github --email your@email.com -t ed25519

# RSA
sshman create github --email your@email.com -t rsa
```

### Managing SSH Agent

#### List Keys in Agent

```bash
sshman agent list
```

#### Add Key to Agent

```bash
sshman agent add id_ed25519_github_work
```

#### Remove Key from Agent

```bash
sshman agent remove id_ed25519_github_work
```

#### Clear All Keys from Agent

```bash
sshman agent clear
```

### Listing SSH Keys

View all SSH keys with their status:

```bash
sshman list
```

Output example:

```output
NAME                    TYPE     STATUS      PATH
id_ed25519_github_work  ED25519  Loaded      ~/.ssh/id_ed25519_github_work
id_rsa_personal         RSA      Not Loaded  ~/.ssh/id_rsa_personal
```

### Deleting SSH Keys

Remove SSH key pairs and clean up:

```bash
sshman delete id_ed25519_github_work
```

This will:

- Remove the key from SSH agent (if loaded)
- Delete the private key file
- Delete the public key file
- Notify about manual SSH config cleanup

## Examples

### Workflow: Setting up Multiple Git Accounts

1. **Create work GitHub key:**

```bash
sshman create github --email work@company.com --purpose work
```

1. **Create personal GitHub key:**

```bash
sshman create github --email personal@gmail.com --purpose personal
```

1. **List keys to verify:**

```bash
sshman list
```

1. **Clone repositories using host aliases:**

```bash
git clone git@github-work:company/project.git
git clone git@github-personal:username/personal-project.git
```

### Workflow: Server Access

1. **Create key for server access:**

```bash
sshman create generic --email admin@company.com --user root --hostname server.company.com --purpose production
```

1. **Add to agent:**

```bash
sshman agent add id_ed25519_generic_production
```

1. **Connect to server:**

```bash
ssh server.company.com-production
```

## How It Works

### SSH Config Integration

When creating keys, sshman automatically adds entries to your SSH config (`~/.ssh/config`):

```sshconfig
Host github-work
    User git
    Hostname github.com
    IdentityFile ~/.ssh/id_ed25519_github_work
```

This allows you to use the host alias in Git operations:

```bash
git clone git@github-work:username/repository.git
```

### Key Naming Convention

Keys are named using the pattern: `id_{type}_{provider}_{purpose}`

Examples:

- `id_ed25519_github_work`
- `id_rsa_gitlab_personal`
- `id_ed25519_bitbucket_project`

### Provider Configurations

Built-in provider configurations:

| Provider  | User | Hostname      |
| --------- | ---- | ------------- |
| GitHub    | git  | github.com    |
| GitLab    | git  | gitlab.com    |
| Bitbucket | git  | bitbucket.org |

## Development

### Building

```bash
make all
```

### Testing

```bash
make test
```

### Code Quality

```bash
make lint
make fmt
```

### Coverage

```bash
make coverage
```

### Available Make commands

- `make all`: Run all checks (format, lint, test, coverage)
- `make test`: Run tests with coverage
- `make lint`: Run go vet
- `make fmt`: Format code
- `make coverage`: Show test coverage
- `make mocks`: Generate mocks for testing
