package packwrap

// environment variable name for the manifest search path
const Envvar_manifestPath string = "PACKWRAP_MANIFEST_PATH"

// manifest file extension
const Extension string = ".json"

// default location for manifest files
const ManifestPath string = "/packages/manifest"

// default golang template string. This should encode
// major minor and micro...
const VERSION_TEMPLATE_STRING string = "{{.Major}}.{{.Minor}}.{{.Micro}}"
