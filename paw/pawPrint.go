package main

import (
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/jlgerber/packwrap"
)

// pawPrint - Given the name and version of a particular executable, find its
// manifest and print its contents.
func pawPrint(manifestLocator *packwrap.ManifestLocator,
	manifestReaderFactory *packwrap.ManifestReaderFactory) error {
	usage := `Usage: paw print [options] <command> <version>

paw print - print the contents of a specific manifest.
 
Options:
   -l, --loglevel=<level>
   -d, --debug
   -q, --quiet
   -h, --help
   `
	args, _ := docopt.Parse(usage, nil, true, "", false)

	processCommonArgs(args)

	if args["<command>"] == nil || args["<version>"] == nil {
		err := errors.New("pawPrint - wrong number of arguements. paw print <package> <version>")
		return err
	}

	command := args["<command>"].(string)
	version := args["<version>"].(string)

	manifest, err := manifestLocator.GetManifestLocationFor(command, version)
	if err != nil {
		log.Errorf("manifestLocator.GetManifestLocationFor(\"%s\", \"%s\") failed", command, version)
		return err
	}

	contents, err := readLines(manifest)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println(manifest)
	fmt.Println("")

	for _, ln := range contents {
		fmt.Println(ln)
	}
	return nil
}
