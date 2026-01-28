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

## Controls

| Key | Action |
|-----|--------|
| Left / A | Move left |
| Right / D | Move right |
| Down / S | Soft drop |
| Up / W | Rotate clockwise |
| Z | Rotate counter-clockwise |
| Space | Hard drop |
| C | Hold piece |
| P / Esc | Pause |

## License

MIT
