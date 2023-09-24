Vul has several sub commands, image, fs, repo, client and server.

``` bash
NAME:
   vul - A simple and comprehensive vulnerability scanner for containers

USAGE:
   vul [global options] command [command options] target

VERSION:
   dev

COMMANDS:
   image, i          scan an image
   filesystem, fs    scan local filesystem
   repository, repo  scan remote repository
   client, c         client mode
   server, s         server mode
   config, conf      scan config files
   plugin, p         manage plugins
   help, h           Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --quiet, -q        suppress progress bar and log output (default: false) [$VUL_QUIET]
   --debug, -d        debug mode (default: false) [$VUL_DEBUG]
   --cache-dir value  cache directory (default: "/Users/teppei/Library/Caches/vul") [$VUL_CACHE_DIR]
   --help, -h         show help (default: false)
   --version, -v      print the version (default: false)

```
