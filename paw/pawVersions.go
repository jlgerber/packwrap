package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/jlgerber/packwrap"
)

// pawVersions - Lists package versions for a named package supplied as
// the first arugment.
func pawVersions(manifestLocator *packwrap.ManifestLocator,
	manifestReaderFactory *packwrap.ManifestReaderFactory) error {

	usage := `Usage: paw versions [options] <command> [<args>...]

paw versions - list the versions for the provided command.
 
Options:
   -l, --loglevel=<level>
   -d, --debug
   -q, --quiet
   -h, --help
   `
	args, _ := docopt.Parse(usage, nil, true, "", false)

	if args["<command>"] == nil {
		log.Fatal("pawVersions - Need to provide a command to look up versions for")
	}
	command := args["<command>"].(string)

	versions := manifestLocator.GetPackageVersions(command)
	if versions == nil {
		log.Info("pawVersions - No Package Versions Found for ", command)
		return nil
	}
	fmt.Println()
	for _, version := range versions {
		fmt.Println(version)
	}

	return nil
}
