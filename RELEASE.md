# Release

The release process is based around tags:

- Create a new tag in [the semver format](https://semver.org/) `vX.Y.Z` (e.g. `v0.0.1`, `v0.0.2-alpha.2`):

```bash
git tag v0.0.2-alpha.1
git push --tags
```

- This should trigger a GitHub Action that builds the release and uploads it to the `Releases` section of the repository.