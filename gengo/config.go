package gengo

import "flag"

var (
	import_path = flag.String("import_path", "", "Specify import path/prefix for nested types")
)
