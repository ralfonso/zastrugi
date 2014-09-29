zastrugi - key/value based token replacement
============================================

overview
----------------------------------------------------
zastrugi is a build preprocessor that scans source and configuration files for tokens to replace. It uses [etcd](https://github.com/coreos/etcd) as its key/value datastore. This enables developers to keep sensitive data such as API keys, database connection strings, usernames, passwords, and SSH keys out of their repository and instead store them in a centralized location. zastrugi supports namespacing of keys to enable its use for multiple projects.

components
-----------------------------------------------------

### zastrugi
`zastrugi` is the command-line tool to apply token substitution to the specified files.

By default, `zastrugi` scans the current directory for a `.zastrugi` config file. This file should contain configuration in the following JSON format:

```javascript
{
	"namespace":  "project-name",
	"searchPaths": ["site/config.py", "site/config/*.py"]
}
```

`zastrugi` is run as follows, from the root of the project on which to perform replacement.

`zastrugi -datasource=http://etcd.url:4001`

### zastrugi-web
`zastrugi-web` is an upcoming web-based tool for managing keys in the etcd datastore.

-----------------------

TODO
-----------------------------

* documentation
* build instructions
* the completion of zastrugi-web
* allow complete files to be stored in etcd
* configurable delimiters
* support multiple backends (a very long term goal)