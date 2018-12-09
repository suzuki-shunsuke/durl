git status --porcelain || exit 1
# mainly ignore change of package-lock.json
git checkout -- . || exit 1
if [ -n "$CIRCLE_TAG" ]; then
  goreleaser release || exit 1
else
  TAG=0.1.0-alpha
  git tag $TAG || exit 1
  goreleaser release --skip-publish || exit 1
  git tag -d $TAG || exit 1
fi
