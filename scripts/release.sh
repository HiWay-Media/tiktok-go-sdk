#!/usr/bin/env bash
#
# release.sh — tag the current commit as the release described at the top of
# CHANGELOG.md.
#
# The changelog is the source of truth for the version: you write the section,
# the script derives "vX.Y.Z" from it, runs the same gates as CI and creates the
# annotated tag with that section as the tag message. Pushing stays manual (the
# project rule is that the agent never pushes):
#
#   ./scripts/release.sh              # gate + create the tag
#   ./scripts/release.sh --dry-run    # show what would happen
#   ./scripts/release.sh --check      # only verify CHANGELOG/tag consistency
#   ./scripts/release.sh --check v0.3.0
#   ./scripts/release.sh --notes      # print the top changelog section
#   ./scripts/release.sh --notes v0.3.0
#
# Then: git push && git push --tags
set -euo pipefail

cd "$(dirname "$0")/.."

CHANGELOG="CHANGELOG.md"

die() {
	echo "release: $*" >&2
	exit 1
}

# top_version prints the version of the first "## X.Y.Z" section of the changelog.
top_version() {
	awk '/^## [0-9]+\.[0-9]+\.[0-9]+$/ { print $2; exit }' "$CHANGELOG"
}

# section VERSION prints the body of that changelog section.
section() {
	awk -v want="## $1" '
		$0 == want { inside = 1; next }
		inside && /^## / { exit }
		inside { print }
	' "$CHANGELOG"
}

# higher A B prints the greater of two semver versions.
higher() { printf '%s\n%s\n' "$1" "$2" | sort -V | tail -1; }

mode="tag"
want_tag=""
case "${1:-}" in
--check)
	mode="check"
	want_tag="${2:-}"
	;;
--notes)
	mode="notes"
	want_tag="${2:-}"
	;;
--dry-run) mode="dry-run" ;;
"") ;;
*) die "unknown argument: $1 (use --check, --notes or --dry-run)" ;;
esac

[ -f "$CHANGELOG" ] || die "$CHANGELOG not found"

version="$(top_version)"
[ -n "$version" ] || die "no '## X.Y.Z' section found at the top of $CHANGELOG"
tag="v$version"

body="$(section "$version")"
# A release with empty notes is worse than no release: it hides what changed.
[ -n "$(echo "$body" | tr -d '[:space:]')" ] || die "the $version section of $CHANGELOG is empty"

if [ "$mode" = "notes" ]; then
	v="${want_tag#v}"
	[ -n "$v" ] || v="$version"
	notes="$(section "$v")"
	[ -n "$(echo "$notes" | tr -d '[:space:]')" ] || die "no changelog section for $v"
	echo "$notes"
	exit 0
fi

if [ "$mode" = "check" ]; then
	# In CI the tag being built must match the top of the changelog, otherwise
	# the published release notes describe a different version than the tag.
	if [ -n "$want_tag" ] && [ "$want_tag" != "$tag" ]; then
		die "tag $want_tag does not match the top of $CHANGELOG ($tag). Update the changelog or retag."
	fi
	echo "release: $tag matches the top of $CHANGELOG"
	exit 0
fi

# --- tag mode ---------------------------------------------------------------

git rev-parse --git-dir >/dev/null 2>&1 || die "not a git repository"

if git rev-parse -q --verify "refs/tags/$tag" >/dev/null; then
	die "tag $tag already exists: bump the version at the top of $CHANGELOG"
fi

last="$(git tag --sort=-v:refname | head -1 || true)"
if [ -n "$last" ]; then
	if [ "$(higher "${last#v}" "$version")" = "${last#v}" ]; then
		die "version $version is not greater than the last tag $last"
	fi
fi

if [ -n "$(git status --porcelain)" ]; then
	die "working tree not clean: commit first, the tag must point at the release commit"
fi

echo "release: gate (go vet, go test)"
go vet ./...
go test ./...

echo "release: $tag on $(git rev-parse --short HEAD)${last:+ (previous: $last)}"
if [ "$mode" = "dry-run" ]; then
	echo "--- tag message ---"
	printf 'Release %s\n\n%s\n' "$version" "$body"
	echo "--- (dry-run: nothing created) ---"
	exit 0
fi

git tag -a "$tag" -m "$(printf 'Release %s\n\n%s' "$version" "$body")"
echo "release: created $tag — push with: git push && git push --tags"
