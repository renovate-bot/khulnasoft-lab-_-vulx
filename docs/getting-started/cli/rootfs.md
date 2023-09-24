# Rootfs

```bash
NAME:
   vul rootfs - scan rootfs

USAGE:
   vul rootfs [command options] dir

OPTIONS:
   --template value, -t value                     output template [$VUL_TEMPLATE]
   --format value, -f value                       format (table, json, template) (default: "table") [$VUL_FORMAT]
   --severity value, -s value                     severities of vulnerabilities to be displayed (comma separated) (default: "UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL") [$VUL_SEVERITY]
   --output value, -o value                       output file name [$VUL_OUTPUT]
   --exit-code value                              Exit code when vulnerabilities were found (default: 0) [$VUL_EXIT_CODE]
   --skip-db-update, --skip-update                skip updating vulnerability database (default: false) [$VUL_SKIP_UPDATE, $VUL_SKIP_DB_UPDATE]
   --skip-policy-update                           skip updating built-in policies (default: false) [$VUL_SKIP_POLICY_UPDATE]
   --clear-cache, -c                              clear image caches without scanning (default: false) [$VUL_CLEAR_CACHE]
   --ignore-unfixed                               display only fixed vulnerabilities (default: false) [$VUL_IGNORE_UNFIXED]
   --vuln-type value                              comma-separated list of vulnerability types (os,library) (default: "os,library") [$VUL_VULN_TYPE]
   --security-checks value                        comma-separated list of what security issues to detect (vuln,config) (default: "vuln") [$VUL_SECURITY_CHECKS]
   --ignorefile value                             specify .vulignore file (default: ".vulignore") [$VUL_IGNOREFILE]
   --cache-backend value                          cache backend (e.g. redis://localhost:6379) (default: "fs") [$VUL_CACHE_BACKEND]
   --timeout value                                timeout (default: 5m0s) [$VUL_TIMEOUT]
   --no-progress                                  suppress progress bar (default: false) [$VUL_NO_PROGRESS]
   --ignore-policy value                          specify the Rego file to evaluate each vulnerability [$VUL_IGNORE_POLICY]
   --list-all-pkgs                                enabling the option will output all packages regardless of vulnerability (default: false) [$VUL_LIST_ALL_PKGS]
   --skip-files value                             specify the file paths to skip traversal [$VUL_SKIP_FILES]
   --skip-dirs value                              specify the directories where the traversal is skipped [$VUL_SKIP_DIRS]
   --config-policy value                          specify paths to the Rego policy files directory, applying config files [$VUL_CONFIG_POLICY]
   --config-data value                            specify paths from which data for the Rego policies will be recursively loaded [$VUL_CONFIG_DATA]
   --policy-namespaces value, --namespaces value  Rego namespaces (default: "users") [$VUL_POLICY_NAMESPACES]
   --help, -h                                     show help (default: false)
```