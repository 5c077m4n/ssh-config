#!/usr/bin/env bash
# Bump latest git tag patch semver and push. e.g. v1.2.3 -> v1.2.4
# Usage: scripts/bump-tag.sh [major|minor|patch] [tag-prefix]
set -euo pipefail

level="${1:-patch}"
prefix="${2:-v}"

case "$level" in
major | minor | patch) ;;
*)
	echo "error: invalid level '$level' (use major|minor|patch)" >&2
	exit 1
	;;
esac

last=$(git describe --tags --abbrev=0 --match="${prefix}*" 2>/dev/null || true)

if [ "$last" = "" ]; then
	tag="${prefix}0.1.0"
	echo "no existing ${prefix}* tags, creating first tag $tag"
else
	ver="${last#"$prefix"}"
	IFS='.' read -r major minor patch <<<"$ver"
	# strip any pre-release/suffix from patch (e.g. 3-rc1 -> 3)
	patch="${patch%%-*}"
	case "$level" in
	major)
		major=$((major + 1))
		minor=0
		patch=0
		;;
	minor)
		minor=$((minor + 1))
		patch=0
		;;
	patch) patch=$((patch + 1)) ;;
	esac
	tag="${prefix}${major}.${minor}.${patch}"
	echo "bumping $last -> $tag"
fi

if git rev-parse --verify --quiet "refs/tags/$tag" >/dev/null; then
	echo "error: tag $tag already exists" >&2
	exit 1
fi

git tag "$tag"
echo "created tag $tag"

if [ "${PUSH:-1}" = "1" ]; then
	git push origin "$tag"
	echo "pushed $tag to origin"
fi
