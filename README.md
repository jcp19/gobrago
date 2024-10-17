# gobrago
Expands Gobra commands from json configs.

## Example
File `gobra-cfg.json`:
```js
{
	"jar_path": "/gobra.jar",
	"jvm_options": [
		"-Xss1g",
		"-Xmx4g"
	]
}
```

File `scion-path-cfg.json`:
```js
{
	"includes": [
		"./",
		"./verification/dependencies"
	],
	"mce_mode": "on",
	"module": "github.com/scionproto/scion",
	"more_joins": "off",
	"parallelize_branches": true,
	"pkg_path": "/pkg/slayers/path/",
	"print_vpr": true,
	"other": ["--disableNL"]
}
```

Running `gobrago gobra-cfg.json scion-path-cfg.json` generates the command
```
java -Xss1g -Xmx4g -jar /gobra.jar --assumeInjectivityOnInhale
	--backend SILICON --checkConsistency --onlyFilesWithHeader
	-I / /verification/dependencies --mceMode=on
	-m github.com/scionproto/scion --moreJoins off
	-p pkg/slayers/path --parallelizeBranches
	--printVpr --requireTriggers --disableNL
```

Passing the flag `--run` to gobrago will make it execute the command as well.

Note: gobrago expands all relative paths in the gobra (verification job) configuration relative to the gobra (verification job) configuration file.
