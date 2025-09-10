
# gotree

A tiny, dependency-free Go CLI that prints a text-based directory tree (similar to `tree`). By default it visualizes the current working directory. You can also pass a specific path.

## Install

```bash
go install github.com/CalebBellmyer/gotree/cmd/gotree@latest

```

Or build locally:

```bash
git clone <this repo>
cd cmd/gotree
go build -o gotree
```

## Usage

```bash
# Visualize current directory
$ gotree

# Visualize a specific directory
$ gotree /path/to/dir

# Limit depth to 2 levels
$ gotree -depth 2

# Include dotfiles
$ gotree -all

# Show directories only
$ gotree -dirs-only

# Show file sizes
$ gotree -size
```

**Examples**

```
$ gotree -depth 2
myproject
├── cmd
│   └── gotree
├── go.mod
└── pkg
    └── tree
```

## Notes
- Depth `0` means unlimited depth.
- Hidden entries (prefix `.`) are omitted unless `-all` is set.
- Sorting puts directories first, then files, each group alphabetically.
- Exit codes: `0` on success, `1` on errors (bad path, permission, etc.).
