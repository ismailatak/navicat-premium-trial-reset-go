#!/usr/bin/env bash
set -euxo pipefail

echo "::group::Setup git"
git config user.name iso-dev-bot
git config user.email ismailatakofficial+gh_bot@gmail.com
echo "::endgroup::"

echo "::group::Setup variables"
cur_version="$(cat VERSION | head -n 1)"
echo "::endgroup::"

echo "::group::Create cosign.key file"
echo "$COSIGN_KEY" > cosign.key
echo "::endgroup::"

echo "::group::Run goreleaser"
if [[ "$DRY_RUN" != 1 ]]; then
  goreleaser release --clean
else
  goreleaser release --clean --snapshot
fi
echo "::endgroup::"

echo "::group::Verify checksum with cosign"
output=$(cosign verify-blob --key cosign.pub --signature dist/checksums.txt.sig dist/checksums.txt 2>&1)
echo "$output"
if [ "$output" != "Verified OK" ]; then
  echo "Checksum verification failed"
  exit 1
fi
echo "::endgroup::"

echo "::group::Publish the release"
if [[ "$DRY_RUN" != 1 ]]; then
  gh release edit "$cur_version" --draft=false
fi
echo "::endgroup::"