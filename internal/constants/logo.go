package constants

import "github.com/charmbracelet/lipgloss"

const asciiArt = `
   ___           __   ____             
  / _ )___ ____ / /  / __/__  ______ _ 
 / _  / _ /(_-</ _ \/ _// _ \/ __/  ' \
/____/\_,_/___/_//_/_/  \___/_/ /_/_/_/
`

var Logo string = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Render(asciiArt)
