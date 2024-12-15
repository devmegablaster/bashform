package constants

import (
	"github.com/devmegablaster/bashform/internal/styles"
)

const logoArt = `
   ___           __   ____             
  / _ )___ ____ / /  / __/__  ______ _ 
 / _  / _ /(_-</ _ \/ _// _ \/ __/  ' \
/____/\_,_/___/_//_/_/  \___/_/ /_/_/_/
`

var Logo string = styles.Logo.Render(logoArt)
