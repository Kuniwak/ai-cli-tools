How to Contribute
=================

How to Release
--------------

```console
$ edit ./version/version.go
$ git add ./version/version.go
$ git commit -m "Ready to be vx.y.z"
$ git tag x.y.z
$ git push && git push --tags

$ goreleaser release --clean --skip publish && open ./dist && gh repo view -w

$ # Create new release with x.y.z and add checksums.txt and *.tar.gz in ./dist
```
