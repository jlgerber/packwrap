package main

import (
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/jlgerber/packwrap"
)

// pawEnv - given a string slice of arguments to the env command, lad
// the package and print it out
func pawEnv(manifestLocator *packwrap.ManifestLocator,
	manifestReaderFactory *packwrap.ManifestReaderFactory) error {
	usage := `Usage: paw env [options] <command> <version> [<args>...]

paw env - execute the supplied command and version. Alternatively, you may supply an auxiliary
command to run (cmd) in the versioned command's environment, which it gets from the manifest.
 
Options:
   -l, --loglevel=<level>
   -d, --debug
   -q, --quiet
   -h, --help
   `

	args, _ := docopt.Parse(usage, nil, true, "", false)
	log.Debug("printEnv -cargs", args)

	if args["<command>"] == nil || args["<version>"] == nil {
		err := errors.New("pawEnv - wrong number of arguments. paw env <package> <version>")
		return err
	}

	command := args["<command>"].(string)
	version := args["<version>"].(string)

	// find manifest
	location, err := manifestLocator.GetManifestLocationFor(command, version)
	if err != nil {
		log.Errorf("manifestLocator.GetManifestLocationFor(%s,%s) failed", command, version)
		return err
	}
	// get instance of manifest given a valid location
	manifest, err := manifestReaderFactory.NewManifestFor(location)
	if err != nil {
		log.Errorf("manifestReaderFactory.NewManifestFor(%s) failed to return manifest", location)
		return err
	}

	for _, value := range manifest.Getenv() {
		fmt.Println(value)
	}
	return nil

}
