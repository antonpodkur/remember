# remember

Zero-friction note-taking CLI for developers. Capture thoughts instantly without leaving the terminal.

## Install

```bash
go install github.com/antonpodkur/remember@latest
```

## Usage

```bash
# Add a note
remember ideas "use redis for caching"

# Interactive input
remember ideas

# Pipe content
echo "deploy fix" | remember todo

# Quick access to last modified note
remember latest "another thought"

# Open note in $EDITOR
remember open ideas
remember open latest

# List all notes
remember list

# Fuzzy search
remember search idea

# View note without editor
remember cat ideas
remember cat latest

# Show last N entries
remember tail ideas -n 3

# Export minified for LLM context
remember export ideas
# Output: 2024-01-15 14:30|use redis for caching

# Append from clipboard
remember clipboard ideas
```

## Storage

Notes are stored as markdown files in `~/.remember/` with timestamped entries:

```markdown
## 2024-01-15 14:30

use redis for caching

## 2024-01-15 15:45

another thought
```

## License

GPL-3.0
