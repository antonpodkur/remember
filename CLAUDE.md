# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
go build .              # Build the binary
go run main.go          # Run without building
./remember              # Run the built binary
```

## Architecture

Three-layer architecture:

```
cmd/           → CLI commands (Cobra framework)
internal/storage/  → Data persistence (~/.remember/*.md files)
internal/ui/       → Terminal UI (lipgloss styling, Bubble Tea interactive input)
```

### Command Layer (cmd/)

| File | Command | Purpose |
|------|---------|---------|
| root.go | `remember <name> [content]` | Append timestamped entry to note |
| list.go | `remember list` | List all notes |
| search.go | `remember search <query>` | Fuzzy search note names |
| open.go | `remember open [name]` | Open note in $EDITOR |
| completion.go | `remember completion [bash\|zsh\|fish]` | Generate shell completions |

### Storage Layer (internal/storage/)

- Notes stored as markdown in `~/.remember/<name>.md`
- Each entry formatted as `## YYYY-MM-DD HH:MM\n\n<content>\n`
- Note names: alphanumeric and hyphens only
- Reserved names: list, open, search, completion, help

### UI Layer (internal/ui/)

- `styles.go`: Lipgloss styles (NoteName, Success, Error, InputBox, Hint)
- `input.go`: Bubble Tea textarea for interactive input (Ctrl+D save, Esc cancel)

## Input Modes

1. **Argument mode**: `remember notes "my content"`
2. **Interactive mode**: `remember notes` → opens styled textarea
3. **Pipe mode**: `echo "content" | remember notes`

TTY detection (`os.ModeCharDevice`) determines interactive vs pipe mode.

## Adding a New Command

1. Create `cmd/<command>.go`
2. Define `var <command>Cmd = &cobra.Command{...}`
3. Register in `init()`: `rootCmd.AddCommand(<command>Cmd)`
4. For autocompletion, add `ValidArgsFunction` if needed

## Key Dependencies

- **spf13/cobra**: CLI framework
- **charmbracelet/lipgloss**: Terminal styling
- **charmbracelet/bubbletea**: Interactive TUI
- **sahilm/fuzzy**: Fuzzy string matching
