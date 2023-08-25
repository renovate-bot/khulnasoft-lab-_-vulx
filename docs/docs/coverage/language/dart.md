# Dart

Vul supports [Dart][dart].

The following scanners are supported.

| Package manager         | SBOM  | Vulnerability | License |
|-------------------------| :---: | :-----------: |:-------:|
| [Dart][dart-repository] |   ✓   |       ✓       |    -    |

The following table provides an outline of the features Vul offers.


| Package manager         | File         | Transitive dependencies | Dev dependencies | [Dependency graph][dependency-graph] | Position |
|-------------------------|--------------|:-----------------------:|:----------------:|:------------------------------------:|:--------:|
| [Dart][dart-repository] | pubspec.lock |            ✓            |     Included     |                  -                   |    -     |

## Dart
In order to detect dependencies, Vul searches for `pubspec.lock`.

Vul marks indirect dependencies, but `pubspec.lock` file doesn't have options to separate root and dev transitive dependencies.
So Vul includes all dependencies in report.

[dart]: https://dart.dev/
[dart-repository]: https://pub.dev/
[dependency-graph]: ../../configuration/reporting.md#show-origins-of-vulnerable-dependencies
