zastrugi - key/value based token replacement
============================================

components
-----------------------------------------------------
`zastrugi` is the command-line tool to search a configurable list of files for tokens and replace them with values from the key/value store [etcd](https://github.com/coreos/etcd). It supports namespaces (which generally align with projects/source-sets) and in the future, will support environment/context value overrides.

`zastrugi-web` is a web-based tool for managing keys

project configuration
-----------------------------------------------------
By default, `zastrugi` scans the current directory for a `.zastrugi` file that is used for configuration. This file is in JSON format as follows.

```javascript
{
	"namespace":  "project-name",
	"searchPaths": ["site/config.py", "site/config/*.py"]
}
```

`zastrugi` is run as follows, from the root of the project on which to perform replacement.

`zastrugi -datasource=http://etcd.url:4001`