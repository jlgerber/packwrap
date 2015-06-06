package packwrap

import (
	"fmt"
	"os"
	"reflect"
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
		"name":"houdini",
		"schema":1,
		"major":14,
		"minor":0,
		"micro":335,
		"basepath":"/Library/Frameworks/Houdini.framework/Versions/14.0.335/Resources",
		"environ": [
			"HFS","$$basepath",
			"H", "${HFS}",
			"HB","${H}/bin",
			"HDSO","${H}/../Libraries",
			"HD","${H}/demo",
			"HH","${H}/houdini",
			"HHC","${HH}/config",
			"HT","${H}/toolkit",
			"HSB","${HH}/sbin",
			"TEMP","/tmp",
			"JAVA_HOME","/Library/Java/Home",
			"HOUDINI_MAJOR_RELEASE","14",
			"HOUDINI_MINOR_RELEASE","0",
			"HOUDINI_BUILD_VERSION","335",
			"HOUDINI_VERSION","${HOUDINI_MAJOR_RELEASE}.${HOUDINI_MINOR_RELEASE}.${HOUDINI_BUILD_VERSION}",
			"HOUDINI_BUILD_KERNEL","XXX_BUILD_KERNEL_XXX",
			"HOUDINI_BUILD_PLATFORM","XXX_BUILD_PLATFORM_XXX",
			"HOUDINI_BUILD_COMPILER","XXX_BUILD_COMPILER_XXX"
		]
	 }`)

	// your func
	testpath := "/var/tmp/houdini/14.0.335"
	testmanifest := "manifest.json"
	os.Setenv("TESTMANIFEST", testpath+"/"+testmanifest)

	manifest := setup(m, testpath, testmanifest, manifestContents)

	retCode := m.Run()

	// remove manifest
	teardown(m, manifest)

	// call with result of m.Run()
	os.Exit(retCode)
}

func TestManifest_GetManifestPathSearchPath(t *testing.T) {
	origValue := os.Getenv(ENVVAR_MANIFESTPATH)
	os.Setenv(ENVVAR_MANIFESTPATH, "/var/tmp/manifest")

	//app := "houdini"
	//ver := "14.0.335"
	spath := GetManifestSearchPath()

	if spath[0] != "/var/tmp/manifest" || spath[1] != "/packages/manifest" {
		t.Errorf("Incorrect Search Path:%s", spath)
	}
	//restore environment
	os.Setenv(ENVVAR_MANIFESTPATH, origValue)
}

func TestManifest_GetManifestLocationFor(t *testing.T) {
	testpath := "/var/tmp"
	os.Setenv(ENVVAR_MANIFESTPATH, testpath)
	app := "houdini"
	ver := "14.0.335"
	manifest, err := GetManifestLocationFor(app, ver)
	if err != nil {
		t.Error(err)
	}
	if manifest != fmt.Sprintf("%s/%s/%s/manifest.json", testpath, app, ver) {
		t.Errorf("manifest path incorrect:%s", manifest)
	}
}

func TestManifest_NewManifestFromJsonByteSlice(t *testing.T) {
	md := []byte(`{"name":"houdini",
		"major":14,
		"minor":0,
		"micro":335,
		"url":"http://blabla",
		"environ":["H","/foo/bar/bla","HBIN","${H}/bin"]}`)
	_, err := NewManifestFromJsonByteSlice(md)
	if err != nil {
		t.Error(err)

	}

}

func getManifestStringField(v *Manifest, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

func getManifestIntField(v *Manifest, field string) uint16 {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return uint16(f.Uint())
}

func TestManifest_NewManifestFromJsonFile(t *testing.T) {
	jsonFile := "/packages/manifest/houdini/14.0.335/manifest.json"
	manifest, err := NewManifestFromJsonFile(jsonFile)
	if err != nil {
		t.Error(err)

	}

	strtests := map[string]string{
		"Name":     "houdini",
		"Basepath": "/Library/Frameworks/Houdini.framework/Versions/14.0.335/Resources",
	}

	for key, strtest := range strtests {
		if strtest != getManifestStringField(manifest, key) {
			t.Errorf("%s does not match", key)
		}
	}

	inttests := map[string]uint16{
		"Major": 14,
		"Minor": 0,
		"Micro": 335,
	}

	for key, inttest := range inttests {
		tmp := getManifestIntField(manifest, key)
		if tmp != inttest {
			t.Errorf("%s does not match", key)
		}
	}

	envtests := map[string]string{
		"HFS":                    "$$basepath",
		"H":                      "${HFS}",
		"HB":                     "${H}/bin",
		"HDSO":                   "${H}/../Libraries",
		"HD":                     "${H}/demo",
		"HH":                     "${H}/houdini",
		"HHC":                    "${HH}/config",
		"HT":                     "${H}/toolkit",
		"HSB":                    "${HH}/sbin",
		"TEMP":                   "/tmp",
		"JAVA_HOME":              "/Library/Java/Home",
		"HOUDINI_MAJOR_RELEASE":  "14",
		"HOUDINI_MINOR_RELEASE":  "0",
		"HOUDINI_BUILD_VERSION":  "335",
		"HOUDINI_VERSION":        "${HOUDINI_MAJOR_RELEASE}.${HOUDINI_MINOR_RELEASE}.${HOUDINI_BUILD_VERSION}",
		"HOUDINI_BUILD_KERNEL":   "XXX_BUILD_KERNEL_XXX",
		"HOUDINI_BUILD_PLATFORM": "XXX_BUILD_PLATFORM_XXX",
		"HOUDINI_BUILD_COMPILER": "XXX_BUILD_COMPILER_XXX",
	}
	_ = envtests
	// TODO convert test to deal with manifest change

	// for key, strtest := range envtests {
	// 	if strtest != manifest.Environ[key] {
	// 		t.Errorf("%s does not match", key)
	// 	}
	// }

}

func TestManifest_ReplaceLocalVars(t *testing.T) {

	jsonFile := os.Getenv("TESTMANIFEST")
	println(jsonFile)
	manifest, err := NewManifestFromJsonFile(jsonFile)
	if err != nil {
		t.Error(err)
	}
	val := manifest.ReplaceLocalVars("$$basepath/fff")
	fmt.Println("$$basepath/foo", val)

	val = manifest.ReplaceLocalVars("foo/$$basepath")
	fmt.Println("foo/$$basepath", val)

	val = manifest.ReplaceLocalVars("foo/$$basepath/bar")
	fmt.Println("foo/$$basepath/bar", val)

}
