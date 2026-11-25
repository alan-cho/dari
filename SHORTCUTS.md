# Shortcuts

## Neovim: moving around code

Requires `gopls` (Go language server), configured in `~/.config/nvim/init.lua`.
Use in normal mode with the cursor on a name.

| Key | Action |
|---|---|
| `gd` | Go to definition (works into the Anthropic SDK source too) |
| `Ctrl-]` | Go to definition (built-in alternative) |
| `Ctrl-o` | Jump back to where you were (repeat to keep going back) |
| `Ctrl-i` | Jump forward again |
| `K` | Hover: show the type signature and docs |
| `grr` | List every reference/usage |
| `gri` | Go to implementations (e.g. on the `Tool` interface → `ReadFile`) |
| `grn` | Rename the symbol everywhere in the project |

`gd` is a custom mapping; the rest are Neovim defaults whenever a language server is attached.

## Neovim: autocomplete

Built-in LSP completion (no plugin). The popup opens automatically as you type in insert mode.

| Key | Action |
|---|---|
| `Ctrl-n` / `Ctrl-p` | Next / previous suggestion |
| `Ctrl-y` | Accept the selected suggestion |
| `Ctrl-e` | Close the popup |
| `Ctrl-x Ctrl-o` | Open the popup manually |
