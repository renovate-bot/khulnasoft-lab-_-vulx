## vul server

Server mode

```
vul server [flags]
```

### Examples

```
  # Run a server
  $ vul server

  # Listen on 0.0.0.0:10000
  $ vul server --listen 0.0.0.0:10000

```

### Options

```
      --cache-backend string        cache backend (e.g. redis://localhost:6379) (default "fs")
      --cache-ttl duration          cache TTL when using redis as cache backend
      --clear-cache                 clear image caches without scanning
      --db-repository string        OCI repository to retrieve vul-db from (default "ghcr.io/khulnasoft-lab/vul-db")
      --download-db-only            download/update vulnerability database but don't run a scan
      --download-java-db-only       download/update Java index database but don't run a scan
      --enable-modules strings      [EXPERIMENTAL] module names to enable
  -h, --help                        help for server
      --java-db-repository string   OCI repository to retrieve vul-java-db from (default "ghcr.io/khulnasoft-lab/vul-java-db")
      --listen string               listen address in server mode (default "localhost:4954")
      --module-dir string           specify directory to the wasm modules that will be loaded (default "$HOME/.vul/modules")
      --no-progress                 suppress progress bar
      --password strings            password. Comma-separated passwords allowed. VUL_PASSWORD should be used for security reasons.
      --redis-ca string             redis ca file location, if using redis as cache backend
      --redis-cert string           redis certificate file location, if using redis as cache backend
      --redis-key string            redis key file location, if using redis as cache backend
      --redis-tls                   enable redis TLS with public certificates, if using redis as cache backend
      --registry-token string       registry token
      --reset                       remove all caches and database
      --skip-db-update              skip updating vulnerability database
      --skip-java-db-update         skip updating Java index database
      --token string                for authentication in client/server mode
      --token-header string         specify a header name for token in client/server mode (default "Vul-Token")
      --username strings            username. Comma-separated usernames allowed.
```

### Options inherited from parent commands

```
      --cache-dir string          cache directory (default "/path/to/cache")
  -c, --config string             config path (default "vul.yaml")
  -d, --debug                     debug mode
      --generate-default-config   write the default config to vul-default.yaml
      --insecure                  allow insecure server connections
  -q, --quiet                     suppress progress bar and log output
      --timeout duration          timeout (default 5m0s)
  -v, --version                   show version
```

### SEE ALSO

* [vul](vul.md)	 - Unified security scanner
