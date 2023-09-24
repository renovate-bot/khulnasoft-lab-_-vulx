# Repository

```bash
NAME:
   vul repository - scan remote repository

USAGE:
   vul repository [command options] repo_url

OPTIONS:
   --template value, -t value  output template [$VUL_TEMPLATE]
   --format value, -f value    format (table, json, template) (default: "table") [$VUL_FORMAT]
   --input value, -i value     input file path instead of image name [$VUL_INPUT]
   --severity value, -s value  severities of vulnerabilities to be displayed (comma separated) (default: "UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL") [$VUL_SEVERITY]
   --output value, -o value    output file name [$VUL_OUTPUT]
   --exit-code value           Exit code when vulnerabilities were found (default: 0) [$VUL_EXIT_CODE]
   --skip-update               skip db update (default: false) [$VUL_SKIP_UPDATE]
   --clear-cache, -c           clear image caches without scanning (default: false) [$VUL_CLEAR_CACHE]
   --ignore-unfixed            display only fixed vulnerabilities (default: false) [$VUL_IGNORE_UNFIXED]
   --removed-pkgs              detect vulnerabilities of removed packages (only for Alpine) (default: false) [$VUL_REMOVED_PKGS]
   --vuln-type value           comma-separated list of vulnerability types (os,library) (default: "os,library") [$VUL_VULN_TYPE]
   --ignorefile value          specify .vulignore file (default: ".vulignore") [$VUL_IGNOREFILE]
   --cache-backend value       cache backend (e.g. redis://localhost:6379) (default: "fs") [$VUL_CACHE_BACKEND]
   --timeout value             timeout (default: 5m0s) [$VUL_TIMEOUT]
   --no-progress               suppress progress bar (default: false) [$VUL_NO_PROGRESS]
   --ignore-policy value       specify the Rego file to evaluate each vulnerability [$VUL_IGNORE_POLICY]
   --list-all-pkgs             enabling the option will output all packages regardless of vulnerability (default: false) [$VUL_LIST_ALL_PKGS]
   --skip-files value          specify the file path to skip traversal [$VUL_SKIP_FILES]
   --skip-dirs value           specify the directory where the traversal is skipped [$VUL_SKIP_DIRS]
   --help, -h                  show help (default: false)
```