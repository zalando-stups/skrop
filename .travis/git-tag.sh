#!/usr/bin/env bash

. docker/version.sh

echo "TRAVIS_BRANCH=${TRAVIS_BRANCH}"
echo "TRAVIS_PULL_REQUEST=${TRAVIS_PULL_REQUEST}"
echo "GITHUB_AUTH=${GITHUB_AUTH}"

if [ "${TRAVIS_BRANCH}_${TRAVIS_PULL_REQUEST}" = "master_false" ]; then
  echo "Merge to 'master'. Tagging patch version up."
  git config --global user.email "builds@travis-ci.com"
  git config --global user.name "Travis CI"
  echo "Creating tag for version: ${NEXT_PATCH_VERSION}"
  git tag ${NEXT_PATCH_VERSION} -a -m "Generated tag from TravisCI for build ${TRAVIS_BUILD_NUMBER}"
  git push -q --tags https://$GITHUB_AUTH@github.com/zalando-stups/skrop
else
  echo "Not a merge to 'master'. Don't release a new version."
fi
