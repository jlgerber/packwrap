package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/jlgerber/packwrap"
)

// pawList - List the packages in the system.
func pawList() error {
	usage := `Usage: paw list [options] 

paw list - list the packages and the paths to their respective manifests.
 
Options:
   -l, --loglevel=<level>
   -d, --debug
   -q, --quiet
   -h, --help
   `
	args, _ := docopt.Parse(usage, nil, true, "", false)

	processCommonArgs(args)

	lst := packwrap.GetPackageList()
	fmt.Println()
	for _, pack := range lst {
		fmt.Println(pack)
	}
	return nil
}
