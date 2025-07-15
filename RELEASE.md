# Release

The release process is based around tags:

- Create a new tag in the format `vX.Y.Z` where `X`, `Y`, and `Z` are integers (e.g., `v0.0.1`, `v0.0.2-alpha.2`):

```bash
git tag v0.0.2-alpha.1 --push origin
```

- This should trigger a GitHub Action that builds the release and uploads it to the `Releases` section of the repository.