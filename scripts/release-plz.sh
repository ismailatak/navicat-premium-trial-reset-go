#!/usr/bin/env bash
set -euxo pipefail

echo "::group::Setup git"
git config user.name iso-dev-bot
git config user.email ismailatakofficial+gh_bot@gmail.com
echo "::endgroup::"

echo "::group::Setup variables"
cur_version="$(cat VERSION | head -n 1)"
echo "::endgroup::"

echo "::group::Check if the release exists"
if [[ $cur_version == "" ]]; then
  if [[ "$(gh release view --json tagName 2>&1)" == "release not found" ]]; then
    echo "There is no release yet"
  else
    echo "There is a release but no VERSION file"
    echo "Please create a VERSION file with the current version"
    exit 1
  fi
else
  if [[ "$(gh release view --json tagName 2>&1)" == "release not found" ]]; then
    echo "This condition is invalid. It should not logically come here. Please check!"
    exit 1
  elif ! gh release view --json tagName | jq -r '.tagName' | grep -q "$cur_version"; then
    echo "Releasing $cur_version"
    changelog="$(git cliff --tag "$cur_version" --strip all --unreleased)"
    changelog="$(echo "$changelog" | tail -n +3)"
    git tag "$cur_version" -s -m "$changelog"
    git push --tags
    gh release create "$cur_version" --title "$cur_version" --notes "$changelog" --draft
    exit 0
  fi
fi
echo "::endgroup::"

echo "::group::Check if the release is already in a PR"
bump="$(git cliff --bumped-version)"
echo $bump >VERSION
echo time $(date +%FT%TZ) >>VERSION
branch="$(echo $bump | cut -d. -f1,2)"
branch="release/$branch"

version="$(cat VERSION | head -n 1)"
git cliff --tag "$version" -o CHANGELOG.md
changelog="$(git cliff --tag "$version" --strip all --unreleased)"
changelog="$(echo "$changelog" | tail -n +3)"

sed -i.bak "s/^const currentVersion = \"v[0-9]\{1,\}\.[0-9]\{1,\}\.[0-9]\{1,\}\"$/const currentVersion = \"$version\"/" version.go
rm version.go.bak

git add \
    VERSION \
    CHANGELOG.md \
    version.go
git clean -df
git checkout -B "$branch"
git commit -m "chore(release): $version"
git push origin "$branch" --force

if [[ "$(gh pr list --label release)" == "" ]]; then
  gh pr create --title "chore(release): $version" --body "$changelog" --label "release" --head "$branch"
else
  gh pr edit --title "chore(release): $version" --body "$changelog"
fi
echo "::endgroup::"