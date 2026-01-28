# Briks

A customizable terminal Tetris game built in Go with Bubble Tea + Lipgloss.

## Install

```bash
go install github.com/meszmate/briks/cmd/briks@latest
```

## Build from source

```bash
git clone https://github.com/meszmate/briks.git
cd briks
make build
./briks
```

## Features

- Standard Tetris gameplay with SRS rotation and wall kicks
- 7-bag randomizer for fair piece distribution
- Ghost piece, hold piece, and next piece preview
- DAS/ARR for responsive movement
- T-Spin and combo scoring
- 8 built-in themes (default, light, dracula, nord, monokai, gruvbox, catppuccin, rainbow)
- Persistent configuration and high scores
- Fully customizable key bindings

## Controls (vim-style)

| Key | Action |
|-----|--------|
| h / Left | Move left |
| l / Right | Move right |
| j / Down | Soft drop |
| k / Up | Rotate clockwise |
| z | Rotate counter-clockwise |
| Space / Enter | Hard drop |
| c | Hold piece |
| p / Esc | Pause |

## Menu Navigation

| Key | Action |
|-----|--------|
| j / k | Navigate |
| l / Enter | Select |
| q | Quit / Back |

## License

MIT
