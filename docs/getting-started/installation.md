# Installation

## RHEL/CentOS


=== "Repository"
    Add repository setting to `/etc/yum.repos.d`.

    ``` bash
    $ sudo vim /etc/yum.repos.d/vul.repo
    [vul]
    name=Vul repository
    baseurl=https://khulnasoft-lab.github.io/vul-repo/rpm/releases/$releasever/$basearch/
    gpgcheck=0
    enabled=1
    $ sudo yum -y update
    $ sudo yum -y install vul
    ```

=== "RPM"

    ``` bash
    rpm -ivh https://github.com/khulnasoft-lab/vul/releases/download/{{ git.tag }}/vul_{{ git.tag[1:] }}_Linux-64bit.rpm
    ```

## Debian/Ubuntu

=== "Repository"
    Add repository setting to `/etc/apt/sources.list.d`.

    ``` bash
    sudo apt-get install wget apt-transport-https gnupg lsb-release
    wget -qO - https://khulnasoft-lab.github.io/vul-repo/deb/public.key | sudo apt-key add -
    echo deb https://khulnasoft-lab.github.io/vul-repo/deb $(lsb_release -sc) main | sudo tee -a /etc/apt/sources.list.d/vul.list
    sudo apt-get update
    sudo apt-get install vul
    ```

=== "DEB"

    ``` bash
    wget https://github.com/khulnasoft-lab/vul/releases/download/{{ git.tag }}/vul_{{ git.tag[1:] }}_Linux-64bit.deb
    sudo dpkg -i vul_{{ git.tag[1:] }}_Linux-64bit.deb
    ```



## Arch Linux
Package vul-bin can be installed from the Arch User Repository.

=== "pikaur"

    ``` bash
    pikaur -Sy vul-bin
    ```

=== "yay"

    ``` bash
    yay -Sy vul-bin
    ```

## Homebrew

You can use homebrew on macOS and Linux.

```bash
brew install khulnasoft-lab/vul/vul
```

## Nix/NixOS

You can use nix on Linux or macOS and on others unofficially.

Note that vul is currently only in the unstable channels.

```bash
nix-env --install vul
```

Or through your configuration on NixOS or with home-manager as usual


## Install Script
This script downloads Vul binary based on your OS and architecture.

```bash
curl -sfL https://raw.githubusercontent.com/khulnasoft-lab/vul/main/contrib/install.sh | sh -s -- -b /usr/local/bin {{ git.tag }}
```

## Binary

Download the archive file for your operating system/architecture from [here](https://github.com/khulnasoft-lab/vul/releases/tag/{{ git.tag }}). 
Unpack the archive, and put the binary somewhere in your `$PATH` (on UNIX-y systems, /usr/local/bin or the like).
Make sure it has execution bits turned on.

## From source

```bash
mkdir -p $GOPATH/src/github.com/khulnasoft-lab
cd $GOPATH/src/github.com/khulnasoft-lab
git clone --depth 1 --branch {{ git.tag }} https://github.com/khulnasoft-lab/vul
cd vul/cmd/vul/
export GO111MODULE=on
go install
```

## Docker
### Docker Hub
Replace [YOUR_CACHE_DIR] with the cache directory on your machine.

```bash
docker pull khulnasoft/vul:{{ git.tag[1:] }}
```

Example:

=== "Linux"

    ``` bash
    docker run --rm -v [YOUR_CACHE_DIR]:/root/.cache/ khulnasoft/vul:{{ git.tag[1:] }} [YOUR_IMAGE_NAME]
    ```

=== "macOS"

    ``` bash
    yay -Sy vul-bin
    docker run --rm -v $HOME/Library/Caches:/root/.cache/ khulnasoft/vul:{{ git.tag[1:] }} python:3.4-alpine
    ```

If you would like to scan the image on your host machine, you need to mount `docker.sock`.

```bash
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
    -v $HOME/Library/Caches:/root/.cache/ khulnasoft/vul:{{ git.tag[1:] }} python:3.4-alpine
```

Please re-pull latest `khulnasoft/vul` if an error occurred.

<details>
<summary>Result</summary>

```bash
2019-05-16T01:20:43.180+0900    INFO    Updating vulnerability database...
2019-05-16T01:20:53.029+0900    INFO    Detecting Alpine vulnerabilities...

python:3.4-alpine3.9 (alpine 3.9.2)
===================================
Total: 1 (UNKNOWN: 0, LOW: 0, MEDIUM: 1, HIGH: 0, CRITICAL: 0)

+---------+------------------+----------+-------------------+---------------+--------------------------------+
| LIBRARY | VULNERABILITY ID | SEVERITY | INSTALLED VERSION | FIXED VERSION |             TITLE              |
+---------+------------------+----------+-------------------+---------------+--------------------------------+
| openssl | CVE-2019-1543    | MEDIUM   | 1.1.1a-r1         | 1.1.1b-r1     | openssl: ChaCha20-Poly1305     |
|         |                  |          |                   |               | with long nonces               |
+---------+------------------+----------+-------------------+---------------+--------------------------------+
```

</details>

### GitHub Container Registry

The same image is hosted on [GitHub Container Registry][registry] as well.

```bash
docker pull ghcr.io/khulnasoft-lab/vul:{{ git.tag[1:] }}
```


### Amazon ECR Public

The same image is hosted on [Amazon ECR Public][ecr] as well.

```bash
docker pull public.ecr.aws/khulnasoft-lab/vul:{{ git.tag[1:] }}
```
## Helm
### Installing from the the Khulnasoft Chart Repository

```
helm repo add khulnasoft-lab https://khulnasoft-lab.github.io/helm-charts/
helm repo update
helm search repo vul
helm install my-vul khulnasoft-lab/vul
```

### Installing the Chart

To install the chart with the release name `my-release`:

```
helm install my-release .
```

The command deploys Vul on the Kubernetes cluster in the default configuration. The [Parameters][helm]
section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`.

[ecr]: https://gallery.ecr.aws/khulnasoft-lab/vul
[registry]: https://github.com/orgs/khulnasoft-lab/packages/container/package/vul
[helm]: https://github.com/khulnasoft-lab/vul/tree/{{ git.tag }}/helm/vul
