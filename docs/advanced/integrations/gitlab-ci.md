# GitLab CI

If you're a GitLab Ultimate customer, GitLab 14.0 and above include out-of-the-box integration with Vul. To enable it for your project, simply add the container scanning template to your `.gitlab-ci.yml` file. For more details, please refer to [GitLab's documentation](https://docs.gitlab.com/ee/user/application_security/container_scanning/).

If you're using an earlier version of GitLab, you can still use the new integration by copying the [contents of the 14.0 template](https://gitlab.com/gitlab-org/gitlab/blob/master/lib/gitlab/ci/templates/Security/Container-Scanning.gitlab-ci.yml) to your configuration.

Alternatively, you can always use the example configurations below.

```yaml
stages:
  - test

vul:
  stage: test
  image: docker:stable
  services:
    - name: docker:dind
      entrypoint: ["env", "-u", "DOCKER_HOST"]
      command: ["dockerd-entrypoint.sh"]
  variables:
    DOCKER_HOST: tcp://docker:2375/
    DOCKER_DRIVER: overlay2
    # See https://github.com/docker-library/docker/pull/166
    DOCKER_TLS_CERTDIR: ""
    IMAGE: vul-ci-test:$CI_COMMIT_SHA
  before_script:
    - export VUL_VERSION=$(wget -qO - "https://api.github.com/repos/khulnasoft-lab/vul/releases/latest" | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')
    - echo $VUL_VERSION
    - wget --no-verbose https://github.com/khulnasoft-lab/vul/releases/download/v${VUL_VERSION}/vul_${VUL_VERSION}_Linux-64bit.tar.gz -O - | tar -zxvf -
  allow_failure: true
  script:
    # Build image
    - docker build -t $IMAGE .
    # Build report
    - ./vul --exit-code 0 --cache-dir .vulcache/ --no-progress --format template --template "@contrib/gitlab.tpl" -o gl-container-scanning-report.json $IMAGE
    # Print report
    - ./vul --exit-code 0 --cache-dir .vulcache/ --no-progress --severity HIGH $IMAGE
    # Fail on severe vulnerabilities
    - ./vul --exit-code 1 --cache-dir .vulcache/ --severity CRITICAL --no-progress $IMAGE
  cache:
    paths:
      - .vulcache/
  # Enables https://docs.gitlab.com/ee/user/application_security/container_scanning/ (Container Scanning report is available on GitLab EE Ultimate or GitLab.com Gold)
  artifacts:
    reports:
      container_scanning: gl-container-scanning-report.json
```

[Example][example]
[Repository][repository]

### GitLab CI using Vul container

To scan a previously built image that has already been pushed into the
GitLab container registry the following CI job manifest can be used.
Note that `entrypoint` needs to be unset for the `script` section to work.
In case of a non-public GitLab project Vul additionally needs to
authenticate to the registry to be able to pull your application image.
Finally, it is not necessary to clone the project repo as we only work
with the container image.

```yaml
container_scanning:
  image:
    name: docker.io/khulnasoft/vul:latest
    entrypoint: [""]
  variables:
    # No need to clone the repo, we exclusively work on artifacts.  See
    # https://docs.gitlab.com/ee/ci/runners/README.html#git-strategy
    GIT_STRATEGY: none
    VUL_USERNAME: "$CI_REGISTRY_USER"
    VUL_PASSWORD: "$CI_REGISTRY_PASSWORD"
    VUL_AUTH_URL: "$CI_REGISTRY"
    FULL_IMAGE_NAME: $CI_REGISTRY_IMAGE:$CI_COMMIT_REF_SLUG
  script:
    - vul --version
    # cache cleanup is needed when scanning images with the same tags, it does not remove the database
    - time vul image --clear-cache
    # update vulnerabilities db
    - time vul --download-db-only --no-progress --cache-dir .vulcache/
    # Builds report and puts it in the default workdir $CI_PROJECT_DIR, so `artifacts:` can take it from there
    - time vul --exit-code 0 --cache-dir .vulcache/ --no-progress --format template --template "@/contrib/gitlab.tpl"
        --output "$CI_PROJECT_DIR/gl-container-scanning-report.json" "$FULL_IMAGE_NAME"
    # Prints full report
    - time vul --exit-code 0 --cache-dir .vulcache/ --no-progress "$FULL_IMAGE_NAME"
    # Fail on critical vulnerabilities
    - time vul --exit-code 1 --cache-dir .vulcache/ --severity CRITICAL --no-progress "$FULL_IMAGE_NAME"
  cache:
    paths:
      - .vulcache/
  # Enables https://docs.gitlab.com/ee/user/application_security/container_scanning/ (Container Scanning report is available on GitLab EE Ultimate or GitLab.com Gold)
  artifacts:
    when:                          always
    reports:
      container_scanning:          gl-container-scanning-report.json
  tags:
    - docker-runner
```

[example]: https://gitlab.com/khulnasoft-lab/vul-ci-test/pipelines
[repository]: https://github.com/khulnasoft-lab/vul-ci-test

### Gitlab CI alternative template

Depending on the edition of gitlab you have or your desired workflow, the
container scanning template may not meet your needs. As an addition to the
above container scanning template, a template for
[code climate](https://docs.gitlab.com/ee/user/project/merge_requests/code_quality.html)
has been included. The key things to update from the above examples are
the `template` and `report` type. An updated example is below.

```yaml
stages:
  - test

vul:
  stage: test
  image: docker:stable
  services:
    - name: docker:dind
      entrypoint: ["env", "-u", "DOCKER_HOST"]
      command: ["dockerd-entrypoint.sh"]
  variables:
    DOCKER_HOST: tcp://docker:2375/
    DOCKER_DRIVER: overlay2
    # See https://github.com/docker-library/docker/pull/166
    DOCKER_TLS_CERTDIR: ""
    IMAGE: vul-ci-test:$CI_COMMIT_SHA
  before_script:
    - export VUL_VERSION=$(wget -qO - "https://api.github.com/repos/khulnasoft-lab/vul/releases/latest" | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')
    - echo $VUL_VERSION
    - wget --no-verbose https://github.com/khulnasoft-lab/vul/releases/download/v${VUL_VERSION}/vul_${VUL_VERSION}_Linux-64bit.tar.gz -O - | tar -zxvf -
  allow_failure: true
  script:
    # Build image
    - docker build -t $IMAGE .
    # Build report
    - ./vul --exit-code 0 --cache-dir .vulcache/ --no-progress --format template --template "@contrib/gitlab-codeclimate.tpl" -o gl-codeclimate.json $IMAGE
  cache:
    paths:
      - .vulcache/
  # Enables https://docs.gitlab.com/ee/user/application_security/container_scanning/ (Container Scanning report is available on GitLab EE Ultimate or GitLab.com Gold)
  artifacts:
    paths:
      gl-codeclimate.json
    reports:
      codequality: gl-codeclimate.json
```

Currently gitlab only supports a single code quality report. There is an
open [feature request](https://gitlab.com/gitlab-org/gitlab/-/issues/9014)
to support multiple reports. Until this has been implemented, if you
already have a code quality report in your pipeline, you can use
`jq` to combine reports. Depending on how you name your artifacts, it may
be necessary to rename the artifact if you want to reuse the name. To then
combine the previous artifact with the output of vul, the following `jq`
command can be used, `jq -s 'add' prev-codeclimate.json vul-codeclimate.json > gl-codeclimate.json`.
