# Config

``` bash
NAME:
   vul config - scan config files

USAGE:
   vul config [command options] dir

OPTIONS:
   --template value, -t value                     output template [$VUL_TEMPLATE]
   --format value, -f value                       format (table, json, template) (default: "table") [$VUL_FORMAT]
   --severity value, -s value                     severities of vulnerabilities to be displayed (comma separated) (default: "UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL") [$VUL_SEVERITY]
   --output value, -o value                       output file name [$VUL_OUTPUT]
   --exit-code value                              Exit code when vulnerabilities were found (default: 0) [$VUL_EXIT_CODE]
   --skip-policy-update                           skip updating built-in policies (default: false) [$VUL_SKIP_POLICY_UPDATE]
   --reset                                        remove all caches and database (default: false) [$VUL_RESET]
   --clear-cache, -c                              clear image caches without scanning (default: false) [$VUL_CLEAR_CACHE]
   --ignorefile value                             specify .vulignore file (default: ".vulignore") [$VUL_IGNOREFILE]
   --timeout value                                timeout (default: 5m0s) [$VUL_TIMEOUT]
   --skip-files value                             specify the file paths to skip traversal [$VUL_SKIP_FILES]
   --skip-dirs value                              specify the directories where the traversal is skipped [$VUL_SKIP_DIRS]
   --policy value, --config-policy value          specify paths to the Rego policy files directory, applying config files [$VUL_POLICY]
   --data value, --config-data value              specify paths from which data for the Rego policies will be recursively loaded [$VUL_DATA]
   --policy-namespaces value, --namespaces value  Rego namespaces (default: "users") [$VUL_POLICY_NAMESPACES]
   --file-patterns value                          specify file patterns [$VUL_FILE_PATTERNS]
   --include-successes                            include successes of misconfigurations (default: false) [$VUL_INCLUDE_SUCCESSES]
   --help, -h                                     show help (default: false)
```
