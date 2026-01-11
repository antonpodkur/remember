---
name: remember-notes
description: Save and retrieve developer notes using the remember CLI. Use when user wants to note something, save findings, check notes, or retrieve from their notes.
allowed-tools: Bash
---

# Remember Notes

A skill for managing developer notes during Claude Code sessions using the `remember` CLI.

## Commands Available

```bash
remember <name> "content"    # Add entry to a note
remember latest "content"    # Add to most recently modified note
remember list                # List all notes
remember cat <name>          # View full note content
remember export <name>       # Get minified output (for context)
remember tail <name> -n N    # View last N entries
```

## When to Add Notes

When the user says things like:
- "note this", "note that", "save this to notes"
- "add to my notes", "remember this"
- "save to [notename]"

### Adding Notes

1. If user specifies a note name, use: `remember <name> "content"`
2. If no name specified, use: `remember latest "content"`
3. For new topics, suggest a descriptive name

Example:
```bash
remember ideas "consider using Redis for caching"
remember latest "fixed the authentication bug"
```

## When to Read Notes

When the user says things like:
- "check my notes", "what's in my notes"
- "get from notes", "show my [notename] note"
- "what did I note about..."

### Reading Notes

1. To show full content: `remember cat <name>`
2. To get minified for context: `remember export <name>`
3. To show recent entries: `remember tail <name> -n 5`
4. To list all notes: `remember list`

Example:
```bash
remember list                    # See available notes
remember cat ideas               # Full content
remember export latest           # Minified latest note
remember tail ideas -n 3         # Last 3 entries
```

## Output Format

The `export` command outputs minified format ideal for LLM context:
```
2024-01-15 14:30|content here on single line
2024-01-15 15:45|another entry
```

The `cat` command outputs raw markdown with timestamps as headers.
