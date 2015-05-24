package packwrap

import (
	"fmt"
	"log"
	"os"
	"testing"
)

// setup = create a fake houdini manifest
func setup(m *testing.M, manifestPath string, manifestName string, contents []byte) string {

	manifest := manifestPath + "/" + manifestName

	err := os.MkdirAll(manifestPath, 0777)
	if err != nil {
		log.Fatalf("Mkdir %q %s", manifestPath, err)
	}

	fh, err := os.Create(manifest)
	if err != nil {
		log.Fatalf("os.Create failed %s %s", manifest, err)
	}

	_, err = fh.Write(contents)
	if err != nil {
		log.Fatal("unable to write manifest")
	}

	return manifest
}

// teardown - remove the fake houdini manifest
func teardown(m *testing.M, manifest string) {
	// remove manifest
	if err := os.Remove(manifest); err != nil {
		log.Fatal("unable to remove manifest file")
	}

}

// TestMain - Run setup, execute tests and teardown, then exit
func TestMain(m *testing.M) {

	manifestContents := []byte(
		`{ 
		"name":"houdini"
		"version":"14.0.335"
		"basepath":"/Library/Frameworks/Houdini.framework/Versions/14.0.335/Resources"
		"environ": {
			"HFS":"$$basepath",
			"H": "${HFS}",
			"HB":"${H}/bin",
			"HDSO":"${H}/../Libraries",
			"HD":"${H}/demo",
			"HH":"${H}/houdini",
			"HHC":"${HH}/config",
			"HT":"${H}/toolkit",
			"HSB":"${HH}/sbin",
			"TEMP":"/tmp"
			"JAVA_HOME":"/Library/Java/Home",
			"HOUDINI_MAJOR_RELEASE":"14",
			"HOUDINI_MINOR_RELEASE":"0",
			"HOUDINI_BUILD_VERSION":"335"
			"HOUDINI_VERSION":"${HOUDINI_MAJOR_RELEASE}.${HOUDINI_MINOR_RELEASE}.${HOUDINI_BUILD_VERSION}",
			"HOUDINI_BUILD_KERNEL":"XXX_BUILD_KERNEL_XXX",
			"HOUDINI_BUILD_PLATFORM":"XXX_BUILD_PLATFORM_XXX",
			"HOUDINI_BUILD_COMPILER":"XXX_BUILD_COMPILER_XXX"
		}
	 }`)

	// your func
	testpath := "/var/tmp/houdini/14.0.335"
	testmanifest := "manifest.json"

	manifest := setup(m, testpath, testmanifest, manifestContents)

	retCode := m.Run()

	// remove manifest
	teardown(m, manifest)

	// call with result of m.Run()
	os.Exit(retCode)
}

func TestEnviron_GetManifestPathSearchPathFor(t *testing.T) {
	origValue := os.Getenv(Envvar_manifestPath)
	os.Setenv(Envvar_manifestPath, "/var/tmp/manifest")

	app := "houdini"
	ver := "14.0.335"
	spath := GetManifestSearchPathFor(app, ver)

	if spath != "/var/tmp/manifest:/packages/manifest" {
		t.Errorf("Incorrect Search Path:%s", spath)
	}
	//restore environment
	os.Setenv(Envvar_manifestPath, origValue)
}

func TestEnviron_GetManifestFor(t *testing.T) {
	testpath := "/var/tmp"
	os.Setenv(Envvar_manifestPath, testpath)
	app := "houdini"
	ver := "14.0.335"
	manifest, err := GetManifestFor(app, ver)
	if err != nil {
		t.Error(err)
	}
	if manifest != fmt.Sprintf("%s/%s/%s/manifest.json", testpath, app, ver) {
		t.Errorf("manifest path incorrect:%s", manifest)
	}
}
