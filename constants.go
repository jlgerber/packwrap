package packwrap

// environment variable name for the manifest search path
const ENVVAR_MANIFESTPATH string = "PACKWRAP_MANIFEST_PATH"

// manifest file extension
const Extension string = ".json"

// default location for manifest files
const DEFAULT_MANIFEST_LOCATION string = "/packages/manifest"

// default golang template string. This should encode
// major minor and micro...
const VERSION_TEMPLATE_STRING string = "{{.Major}}.{{.Minor}}.{{.Micro}}"

const MANIFEST_TEMPLATE_STRING string = `
{ 
	"schema": 1, // do not remove! 

	// The name of the executable
	"name":"",

	// major version - an integer is expected
	"major": , 
	
	// minor version - an integer is expected
	"minor": ,
	
	// micro version - an integer is expected
	"micro": ,
	
	// path to executable
	"basepath":"", 

	// The environ list consists of key value pairs, 
	// in order, key first, then value. We use a list rather than a map because
	// we need to preserve the order of the values. You may either use 
	// shell variables ( $<name> or ${<name>} eg $fred or ${fred}), or manifest 
	// variables, denoted by $$. Manifest variables may refer to any key in 
	// the manifest ( other than environ of course). 
	// eg 
	// "environ" : [
	//  	"maya","$$basepath/$$name"
	// ]
	"environ" : [
		// "key1", "value1",
		// "key2", "value2",
	]
}
`
