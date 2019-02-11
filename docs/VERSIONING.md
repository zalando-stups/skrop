## Versioning

This project uses [semantic versioning](https://semver.org/).

The patch-version (3rd digit) is bumped up automatically at every merge to master (by [_Travis CI_](/.travis.yml)).

### Increment the patch version
This is done automatically by _Travis CI_. Nothing special to do here. Example: merging when latest tag is `v3.23.291`
will automaticall tag a version `v3.23.292`.

### Increment the minor version
Since _Travis CI_ only automatically increases the patch-version, we need to manually pre-tag with the new version we
want.

Scenario:
- actual version is `v3.23.291`
- tag one of the commit on your branch with the new version you want.
  - it is **important** that the patch version be `-1`, since it will be incremented automatically by _Travis CI_.
  - `git tag v3.24.-1 && git push --tags`
- open the pull request.
- after merge, _Travis CI_ will tag automatically the right final version `v3.24.0`.
- delete the temporary manual tag
  - `git tag -d v3.24.-1 && git push --tags`

### Increment the major version
Since _Travis CI_ only automatically increases the patch-version, we need to manually pre-tag with the new version we
want.

Scenario:
- actual version is `v3.23.291`
- tag one of the commit on your branch with the new version you want.
  - it is **important** that the patch version be `-1`, since it will be incremented automatically by _Travis CI_.
  - `git tag v4.0.-1 && git push --tags`
- open the pull request.
- after merge, _Travis CI_ will tag automatically the right final version `v4.0.0`.
- delete the temporary manual tag
  - `git tag -d v4.0.-1 && git push --tags`
  