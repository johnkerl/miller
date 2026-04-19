#!/usr/bin/env bash
#
# tools/release.sh -- automate the bulk of docs/src/how-to-release.md.
#
# Usage:
#   ./tools/release.sh <tag> pre-release [--notes-file FILE] [--branch NAME] [--yes] [--dry-run]
#   ./tools/release.sh <tag> docs        [--branch NAME] [--yes] [--dry-run]
#   ./tools/release.sh <tag> afterwork   [--branch NAME] [--yes] [--dry-run]
#
# <tag> may be "v6.18.0" or "6.18.0"; it is normalized to both forms internally.
#
# Subcommands:
#   pre-release
#     Phase 0 preflight, phase 1 version bumps + `make dev` + commit + push,
#     phase 2 release tarball + sha256, phase 3 SRPM build, phase 4 GitHub
#     pre-release create + asset upload.  Stops after the pre-release is
#     created.  The operator must then wait for CI/goreleaser, verify a
#     downloaded binary, and flip the release from pre-release to public by
#     hand before running `docs`.
#
#   docs
#     Phase 5. Creates (or reuses) the `<VERSION>` docs branch, edits
#     `docs/mkdocs.yml`, commits and (prompted) pushes.  Refuses to run while
#     the GitHub release is still marked pre-release.  ReadTheDocs admin steps
#     remain manual -- the script prints the URLs to visit.
#
#   afterwork
#     Phases 6 and 7. Back on the main branch, flips
#     `pkg/version/version.go` to `<VERSION>-dev`, runs `make dev`, commits
#     and (prompted) pushes.  Prints the brew/macports/readthedocs/social
#     reminders with a pre-filled `brew bump-formula-pr` command.
#
# Every mutating action is echoed before execution, every phase is framed by a
# banner, and every phase begins with an idempotency check so the script is
# safe to re-run after a partial failure.  `rpmbuild` is mandatory for
# `pre-release` and the script aborts in phase 0 if missing; `rpmlint` is
# optional (phase 3 runs it if installed, otherwise skips the lint step).

set -euo pipefail

# ============================================================================
# Globals populated by arg-parsing
# ============================================================================
TAG=""            # e.g. v6.18.0
VERSION=""        # e.g. 6.18.0
SUBCOMMAND=""     # pre-release | docs | afterwork
BRANCH="main"     # main branch to commit version bumps against
YES="no"          # --yes skips confirmation prompts
DRY_RUN="no"      # --dry-run echoes commands without executing
NOTES_FILE=""     # --notes-file: release notes body (required for new releases)

REPO_ROOT=""      # absolute path to the miller repo root (cwd at script start)

# ============================================================================
# Terminal helpers
# ============================================================================
if [ -t 1 ]; then
  C_RED=$'\033[1;31m'
  C_YEL=$'\033[1;33m'
  C_GRN=$'\033[1;32m'
  C_CYN=$'\033[1;36m'
  C_DIM=$'\033[2m'
  C_OFF=$'\033[0m'
else
  C_RED=""; C_YEL=""; C_GRN=""; C_CYN=""; C_DIM=""; C_OFF=""
fi

banner() {
  echo
  echo "${C_CYN}================================================================${C_OFF}"
  echo "${C_CYN}== $*${C_OFF}"
  echo "${C_CYN}================================================================${C_OFF}"
}

log()  { echo "${C_GRN}[release]${C_OFF} $*"; }
note() { echo "${C_DIM}[release]${C_OFF} $*"; }
warn() { echo "${C_YEL}[release][warn]${C_OFF} $*" >&2; }
die()  { echo "${C_RED}[release][error]${C_OFF} $*" >&2; exit 1; }

# run_cmd echoes a command and then runs it (unless --dry-run was passed).
run_cmd() {
  echo "${C_DIM}\$ $*${C_OFF}"
  if [ "$DRY_RUN" = "yes" ]; then
    return 0
  fi
  "$@"
}

# confirm prompts unless --yes is set.  Returns 0 on y/Y, dies on n/N.
confirm() {
  local prompt="$1"
  if [ "$YES" = "yes" ]; then
    note "--yes: auto-confirming: $prompt"
    return 0
  fi
  local reply=""
  read -r -p "${C_YEL}?${C_OFF} $prompt [y/N] " reply
  case "$reply" in
    y|Y|yes|YES) return 0 ;;
    *) die "user declined: $prompt" ;;
  esac
}

# ============================================================================
# Argument parsing
# ============================================================================
usage() {
  sed -n '3,35p' "$0" | sed 's/^# \{0,1\}//'
  exit "${1:-1}"
}

parse_args() {
  local positional=()
  while [ $# -gt 0 ]; do
    case "$1" in
      -h|--help) usage 0 ;;
      --yes) YES="yes"; shift ;;
      --dry-run) DRY_RUN="yes"; shift ;;
      --notes-file)
        [ $# -ge 2 ] || die "--notes-file requires an argument"
        NOTES_FILE="$2"; shift 2 ;;
      --branch)
        [ $# -ge 2 ] || die "--branch requires an argument"
        BRANCH="$2"; shift 2 ;;
      --) shift; while [ $# -gt 0 ]; do positional+=("$1"); shift; done ;;
      -*) die "unknown flag: $1 (see --help)" ;;
      *) positional+=("$1"); shift ;;
    esac
  done

  [ "${#positional[@]}" -eq 2 ] || die "expected two positional arguments: <tag> <subcommand>. Got ${#positional[@]}. See --help."

  local raw_tag="${positional[0]}"
  SUBCOMMAND="${positional[1]}"

  case "$raw_tag" in
    v*) TAG="$raw_tag"; VERSION="${raw_tag#v}" ;;
    *)  TAG="v$raw_tag"; VERSION="$raw_tag" ;;
  esac

  [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] \
    || die "version '$VERSION' is not a simple MAJOR.MINOR.PATCH triple"

  case "$SUBCOMMAND" in
    pre-release|docs|afterwork) ;;
    *) die "unknown subcommand '$SUBCOMMAND'; expected pre-release | docs | afterwork" ;;
  esac

  if [ -n "$NOTES_FILE" ] && [ ! -f "$NOTES_FILE" ]; then
    die "--notes-file '$NOTES_FILE' does not exist"
  fi
}

# ============================================================================
# Phase 0 -- preflight
# ============================================================================
preflight_common() {
  banner "PHASE 0: preflight -- $SUBCOMMAND ($TAG)"

  REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || true)"
  [ -n "$REPO_ROOT" ] || die "not inside a git working tree"
  cd "$REPO_ROOT"
  log "repo root: $REPO_ROOT"
  log "tag:       $TAG"
  log "version:   $VERSION"
  log "subcmd:    $SUBCOMMAND"
  log "branch:    $BRANCH"
  [ "$DRY_RUN" = "yes" ] && warn "dry-run mode: no mutating commands will execute"

  for tool in go make git sed gh awk; do
    command -v "$tool" >/dev/null 2>&1 || die "required tool '$tool' not found on PATH"
  done

  if ! gh auth status >/dev/null 2>&1; then
    die "gh is not authenticated; run 'gh auth login' first"
  fi

  # The three main files we touch must be where we expect them.
  [ -f pkg/version/version.go ]    || die "pkg/version/version.go not found; wrong repo root?"
  [ -f miller.spec ]               || die "miller.spec not found"
  [ -f docs/mkdocs.yml ]           || die "docs/mkdocs.yml not found"
  [ -x create-release-tarball ]    || die "create-release-tarball not found or not executable"

  # Working tree must be clean.  (git diff-index is quiet if clean.)
  if ! git diff-index --quiet HEAD -- 2>/dev/null; then
    git status --short
    die "working tree has uncommitted changes; commit or stash first"
  fi

  log "preflight (common) ok"
}

preflight_pre_release_extras() {
  # rpmbuild is mandatory: the SRPM is a required release artifact.
  # rpmlint is optional: if present we run it as a lint; if absent we skip.
  command -v rpmbuild >/dev/null 2>&1 || die "rpmbuild not found on PATH; SRPMs are required -- install rpm-build"
  command -v rpmlint  >/dev/null 2>&1 || warn "rpmlint not found on PATH; phase 3 will skip the lint step"
  command -v shasum   >/dev/null 2>&1 || die "shasum not found on PATH"

  # create-release-tarball needs a tar that supports --transform; on macOS
  # that means gnu-tar installed as gtar.
  if [ "$(uname)" = "Darwin" ]; then
    if ! command -v gtar >/dev/null 2>&1 \
         && [ ! -x /usr/local/bin/gtar ] \
         && [ ! -x /opt/homebrew/bin/gtar ]; then
      die "gtar (gnu-tar) not found; 'brew install gnu-tar' first"
    fi
  fi

  # Current on-disk version.go must carry a -dev marker so we know we're
  # bumping forward from the last released version.
  local current
  current="$(read_version_string)"
  case "$current" in
    "$VERSION")
      note "version.go already set to $VERSION; phase 1 will skip its edit"
      ;;
    *-dev)
      note "version.go currently '$current'; will bump to '$VERSION'"
      ;;
    *)
      die "version.go currently '$current' -- expected something-dev or '$VERSION'"
      ;;
  esac

  # Fetch to be sure we know about origin's current tag list, then check that
  # the release tag does not already exist on origin (we haven't published yet).
  run_cmd git fetch --tags origin >/dev/null
  if git ls-remote --tags origin "refs/tags/$TAG" | grep -q "refs/tags/$TAG"; then
    # Not fatal; the GitHub-release step is idempotent and may be resuming.
    warn "tag $TAG already exists on origin; phase 4 will treat the release as a resume"
  fi

  log "preflight (pre-release) ok"
}

preflight_docs_extras() {
  # The release must exist and must no longer be pre-release.
  if ! gh release view "$TAG" >/dev/null 2>&1; then
    die "gh release '$TAG' does not exist -- run 'pre-release' first"
  fi
  local is_pre
  is_pre="$(gh release view "$TAG" --json isPrerelease --jq .isPrerelease 2>/dev/null || echo "true")"
  if [ "$is_pre" != "false" ]; then
    die "release $TAG is still marked pre-release -- verify the CI artifacts, flip it to public on GitHub, then retry"
  fi
  log "preflight (docs) ok -- release $TAG is public"
}

preflight_afterwork_extras() {
  local current
  current="$(read_version_string)"
  case "$current" in
    "$VERSION"|"$VERSION-dev") ;;
    *) die "version.go is '$current' -- expected '$VERSION' or '$VERSION-dev'" ;;
  esac
  log "preflight (afterwork) ok"
}

# ============================================================================
# Small file-editing helpers
# ============================================================================
read_version_string() {
  # Extracts the quoted string assigned to STRING in pkg/version/version.go.
  awk -F'"' '/^var[[:space:]]+STRING/ { print $2; exit }' pkg/version/version.go
}

write_version_string() {
  local new="$1"
  local tmp
  tmp="$(mktemp)"
  awk -v new="$new" '
    /^var[[:space:]]+STRING/ {
      sub(/"[^"]*"/, "\"" new "\"")
    }
    { print }
  ' pkg/version/version.go > "$tmp"
  mv "$tmp" pkg/version/version.go
}

# Prints today's date in miller.spec %changelog format: "Sun Apr 19 2026"
# (day number unpadded, matching recent entries).
spec_changelog_date() {
  local wday mon day year
  wday="$(date +%a)"
  mon="$(date +%b)"
  day="$(date +%d | sed 's/^0//')"
  year="$(date +%Y)"
  echo "$wday $mon $day $year"
}

spec_has_version_entry() {
  grep -q "^\* .* - ${VERSION}-1$" miller.spec
}

spec_current_version() {
  awk '/^Version:/ { print $2; exit }' miller.spec
}

prepend_spec_changelog_entry() {
  local dt
  dt="$(spec_changelog_date)"
  local tmp
  tmp="$(mktemp)"
  awk -v dt="$dt" -v ver="$VERSION" '
    BEGIN { inserted = 0 }
    {
      print
      if (!inserted && $0 ~ /^%changelog[[:space:]]*$/) {
        printf "* %s John Kerl <kerl.john.r@gmail.com> - %s-1\n", dt, ver
        printf "- %s release\n", ver
        printf "\n"
        inserted = 1
      }
    }
    END {
      if (!inserted) {
        exit 2
      }
    }
  ' miller.spec > "$tmp" || die "failed to insert %changelog entry (no %changelog block in miller.spec?)"
  mv "$tmp" miller.spec
}

# ============================================================================
# Phase 1 -- bump version files, make dev, commit, push
# ============================================================================
phase_1_bump_versions() {
  banner "PHASE 1: bump pkg/version/version.go and miller.spec to $VERSION"

  local current
  current="$(read_version_string)"
  if [ "$current" = "$VERSION" ]; then
    note "pkg/version/version.go already set to '$VERSION' (idempotent: skipping edit)"
  else
    log "editing pkg/version/version.go: $current -> $VERSION"
    if [ "$DRY_RUN" = "yes" ]; then
      note "(dry-run) would rewrite version.go STRING to $VERSION"
    else
      write_version_string "$VERSION"
    fi
  fi

  local spec_ver
  spec_ver="$(spec_current_version)"
  if [ "$spec_ver" = "$VERSION" ]; then
    note "miller.spec Version: already $VERSION (idempotent: skipping edit)"
  else
    log "editing miller.spec Version: $spec_ver -> $VERSION"
    if [ "$DRY_RUN" = "yes" ]; then
      note "(dry-run) would set miller.spec Version: $VERSION"
    else
      # Use a literal delimiter that cannot appear in a version string.
      sed -i.bak "s|^Version:.*|Version: $VERSION|" miller.spec
      rm -f miller.spec.bak
    fi
  fi

  if spec_has_version_entry; then
    note "miller.spec %changelog already has an entry for ${VERSION}-1 (idempotent: skipping)"
  else
    log "prepending miller.spec %changelog entry for ${VERSION}-1"
    if [ "$DRY_RUN" = "yes" ]; then
      note "(dry-run) would prepend %changelog entry dated '$(spec_changelog_date)'"
    else
      prepend_spec_changelog_entry
    fi
  fi

  banner "PHASE 1: make dev"
  run_cmd make dev

  banner "PHASE 1: commit + push"
  if git diff-index --quiet HEAD -- 2>/dev/null; then
    note "no staged or unstaged changes after make dev -- commit likely already present"
  else
    run_cmd git add pkg/version/version.go miller.spec
    # `make dev` may have regenerated docs and manpages; include them.
    # Pick up anything else make dev touched, but limit to tracked files only
    # to avoid sweeping in stray files.
    run_cmd git add -u
    run_cmd git commit -m "Prepare ${VERSION} release"
  fi

  confirm "push branch '$BRANCH' to origin?"
  run_cmd git push origin "$BRANCH"

  log "phase 1 complete"
}

# ============================================================================
# Phase 2 -- release tarball + sha256
# ============================================================================
phase_2_release_tarball() {
  banner "PHASE 2: build release tarball"

  local tgz="miller-${VERSION}.tar.gz"
  local built_ver=""
  if [ -x ./mlr ]; then
    built_ver="$(./mlr --bare-version 2>/dev/null || true)"
  fi

  if [ -f "$tgz" ] && [ "$built_ver" = "$VERSION" ]; then
    note "$tgz already present and ./mlr reports $VERSION (idempotent: skipping make release_tarball)"
  else
    run_cmd make release_tarball
  fi

  [ "$DRY_RUN" = "yes" ] || [ -f "$tgz" ] || die "$tgz was not produced"

  banner "PHASE 2: sha256"
  local sha_file="${tgz}.sha256"
  if [ "$DRY_RUN" = "yes" ]; then
    note "(dry-run) would compute sha256 of $tgz"
  else
    shasum -a 256 "$tgz" > "$sha_file"
    log "wrote $sha_file:"
    cat "$sha_file"
  fi

  log "phase 2 complete"
}

# ============================================================================
# Phase 3 -- SRPM (mandatory)
# ============================================================================
phase_3_srpm() {
  banner "PHASE 3: build SRPM"

  local tgz="miller-${VERSION}.tar.gz"
  [ -f "$tgz" ] || [ "$DRY_RUN" = "yes" ] || die "$tgz not found -- phase 2 should have produced it"

  local rpm_top="$HOME/rpmbuild"
  run_cmd mkdir -p "$rpm_top/SPECS" "$rpm_top/SOURCES" "$rpm_top/SRPMS"

  # Check for an already-built SRPM for this version before rebuilding.
  local existing
  existing="$(ls -1 "$rpm_top"/SRPMS/miller-"${VERSION}"-1*.src.rpm 2>/dev/null | head -n1 || true)"
  if [ -n "$existing" ]; then
    note "existing SRPM found (idempotent: reusing): $existing"
    SRPM_PATH="$existing"
    log "phase 3 complete"
    return 0
  fi

  run_cmd cp miller.spec "$rpm_top/SPECS/"
  run_cmd cp "$tgz"      "$rpm_top/SOURCES/"

  banner "PHASE 3: rpmlint miller.spec (optional)"
  if command -v rpmlint >/dev/null 2>&1; then
    if [ "$DRY_RUN" != "yes" ]; then
      rpmlint "$rpm_top/SPECS/miller.spec" || warn "rpmlint produced warnings; review above"
    else
      note "(dry-run) would run rpmlint"
    fi
  else
    note "rpmlint not installed -- skipping lint step"
  fi

  banner "PHASE 3: rpmbuild -bs"
  run_cmd rpmbuild -bs "$rpm_top/SPECS/miller.spec"

  if [ "$DRY_RUN" = "yes" ]; then
    SRPM_PATH="$rpm_top/SRPMS/miller-${VERSION}-1.DRYRUN.src.rpm"
    note "(dry-run) pretend SRPM path: $SRPM_PATH"
  else
    SRPM_PATH="$(ls -1 "$rpm_top"/SRPMS/miller-"${VERSION}"-1*.src.rpm | head -n1)"
    [ -n "$SRPM_PATH" ] && [ -f "$SRPM_PATH" ] || die "rpmbuild reported success but no SRPM found under $rpm_top/SRPMS"
    log "built SRPM: $SRPM_PATH"
  fi

  log "phase 3 complete"
}

# ============================================================================
# Phase 4 -- GitHub pre-release create + upload assets
# ============================================================================
phase_4_github_release() {
  banner "PHASE 4: create GitHub pre-release $TAG"

  local tgz="miller-${VERSION}.tar.gz"
  [ -f "$tgz" ] || [ "$DRY_RUN" = "yes" ] || die "$tgz missing"
  [ -n "${SRPM_PATH:-}" ] || die "SRPM_PATH not set -- phase 3 did not run?"

  if gh release view "$TAG" >/dev/null 2>&1; then
    note "gh release $TAG already exists (idempotent: will only upload missing assets)"
  else
    [ -n "$NOTES_FILE" ] || die "gh release $TAG does not exist and --notes-file was not provided"
    confirm "create GitHub pre-release $TAG from notes file '$NOTES_FILE'?"
    run_cmd gh release create "$TAG" \
      --prerelease \
      --title "Miller $VERSION" \
      --notes-file "$NOTES_FILE"
  fi

  banner "PHASE 4: upload assets"
  # --clobber so that a resumed run can re-upload a replaced artifact.
  run_cmd gh release upload "$TAG" "$tgz" --clobber
  if [ "$DRY_RUN" != "yes" ] || [ -f "${SRPM_PATH}" ]; then
    run_cmd gh release upload "$TAG" "$SRPM_PATH" --clobber
  fi

  banner "PHASE 4: follow-ups"
  cat <<EOF
${C_YEL}Pre-release $TAG is created. Next, by hand:${C_OFF}
  1. Wait for the goreleaser workflow to finish attaching binaries.
     gh run list --limit 5
  2. Download one binary from the release page, then:
       # macOS only:
       xattr -d com.apple.quarantine ./mlr
       ./mlr version   # should print $VERSION
  3. If that looks good, flip the release from pre-release to public on
     the GitHub release page.
  4. Then run:
       $0 $TAG docs
EOF

  log "phase 4 complete"
}

# ============================================================================
# Phase 5 -- docs branch + mkdocs.yml
# ============================================================================
phase_5_docs_branch() {
  banner "PHASE 5: docs branch $VERSION"

  run_cmd git fetch origin >/dev/null

  local local_exists="no" remote_exists="no"
  git rev-parse --verify --quiet "refs/heads/$VERSION" >/dev/null && local_exists="yes"
  git ls-remote --exit-code --heads origin "$VERSION" >/dev/null 2>&1 && remote_exists="yes"

  if [ "$local_exists" = "yes" ]; then
    note "local branch '$VERSION' already exists (idempotent: checking it out)"
    run_cmd git checkout "$VERSION"
  elif [ "$remote_exists" = "yes" ]; then
    note "remote branch '$VERSION' already exists; tracking it locally"
    run_cmd git checkout -b "$VERSION" "origin/$VERSION"
  else
    log "creating branch '$VERSION' from current HEAD"
    run_cmd git checkout -b "$VERSION"
  fi

  banner "PHASE 5: edit docs/mkdocs.yml"
  local already_new="no"
  grep -q "^site_name: Miller $VERSION Documentation\$" docs/mkdocs.yml && already_new="yes"

  if [ "$already_new" = "yes" ]; then
    note "docs/mkdocs.yml already says 'Miller $VERSION Documentation' (idempotent: skipping)"
  else
    grep -q '^site_name: Miller Dev Documentation$' docs/mkdocs.yml \
      || die "docs/mkdocs.yml does not have 'site_name: Miller Dev Documentation' to replace"
    log "rewriting site_name in docs/mkdocs.yml"
    if [ "$DRY_RUN" = "yes" ]; then
      note "(dry-run) would set site_name to 'Miller $VERSION Documentation'"
    else
      sed -i.bak "s|^site_name: Miller Dev Documentation\$|site_name: Miller $VERSION Documentation|" docs/mkdocs.yml
      rm -f docs/mkdocs.yml.bak
    fi
  fi

  banner "PHASE 5: commit + push"
  if git diff-index --quiet HEAD -- docs/mkdocs.yml 2>/dev/null; then
    note "no changes to docs/mkdocs.yml to commit (may already be on the branch)"
  else
    run_cmd git add docs/mkdocs.yml
    run_cmd git commit -m "docs: Miller $VERSION Documentation site_name"
  fi

  confirm "push branch '$VERSION' to origin?"
  run_cmd git push -u origin "$VERSION"

  banner "PHASE 5: manual ReadTheDocs steps"
  cat <<EOF
${C_YEL}Now do these by hand on ReadTheDocs:${C_OFF}
  - Admin page:  https://readthedocs.org/projects/miller
  - Versions tab: scroll to 'Activate a version' and activate $VERSION.
  - Admin > Advanced Settings: set Default Version AND Default Branch to $VERSION, Save.
  - Builds tab: build $VERSION and latest if they are not already building.

Then verify:
  - https://miller.readthedocs.io/en/$VERSION   (should exist)
  - https://miller.readthedocs.io/en/latest     (hard-reload; should show 'Miller $VERSION Documentation')

When that's all good:
  $0 $TAG afterwork
EOF

  log "phase 5 complete"
}

# ============================================================================
# Phase 6 -- afterwork: version.go back to -dev
# ============================================================================
phase_6_afterwork_bump() {
  banner "PHASE 6: afterwork -- version.go back to ${VERSION}-dev"

  run_cmd git checkout "$BRANCH"
  run_cmd git pull --ff-only origin "$BRANCH"

  local current
  current="$(read_version_string)"
  local target="${VERSION}-dev"

  if [ "$current" = "$target" ]; then
    note "version.go already '$target' (idempotent: skipping edit and make dev)"
  else
    [ "$current" = "$VERSION" ] \
      || die "version.go is '$current' -- expected '$VERSION' before afterwork bump"
    log "editing pkg/version/version.go: $current -> $target"
    if [ "$DRY_RUN" = "yes" ]; then
      note "(dry-run) would set version.go STRING to $target"
    else
      write_version_string "$target"
    fi

    banner "PHASE 6: make dev"
    run_cmd make dev

    banner "PHASE 6: commit + push"
    if git diff-index --quiet HEAD -- 2>/dev/null; then
      note "no changes to commit after afterwork bump"
    else
      run_cmd git add pkg/version/version.go
      run_cmd git add -u
      run_cmd git commit -m "Post-${VERSION} release: back to ${target}"
    fi

    confirm "push branch '$BRANCH' to origin?"
    run_cmd git push origin "$BRANCH"
  fi

  log "phase 6 complete"
}

# ============================================================================
# Phase 7 -- reminders
# ============================================================================
phase_7_reminders() {
  banner "PHASE 7: distro/social reminders (no automation)"

  local tgz="miller-${VERSION}.tar.gz"
  local sha_file="${tgz}.sha256"
  local sha=""
  if [ -f "$sha_file" ]; then
    sha="$(awk '{ print $1 }' "$sha_file")"
  elif [ -f "$tgz" ]; then
    sha="$(shasum -a 256 "$tgz" | awk '{ print $1 }')"
  fi

  cat <<EOF
${C_YEL}Homebrew -- submit a version upgrade:${C_OFF}
  https://github.com/Homebrew/homebrew-core/blob/HEAD/CONTRIBUTING.md#to-submit-a-version-upgrade-for-the-foo-formula

EOF

  if [ -n "$sha" ]; then
    cat <<EOF
  brew bump-formula-pr --force --strict miller \\
    --url https://github.com/johnkerl/miller/releases/download/${TAG}/${tgz} \\
    --sha256 ${sha}

EOF
  else
    cat <<EOF
  (Could not find ${sha_file} or ${tgz} to pre-fill sha256.
   Download the release tarball and run:
     shasum -a 256 ${tgz})

EOF
  fi

  cat <<EOF
${C_YEL}MacPorts Portfile:${C_OFF}
  https://github.com/macports/macports-ports/blob/master/textproc/miller/Portfile

${C_YEL}Other distros / contacts:${C_OFF}
  https://github.com/johnkerl/miller/blob/main/README-versions.md

${C_YEL}ReadTheDocs sanity check:${C_OFF}
  https://miller.readthedocs.io/en/${VERSION}
  https://miller.readthedocs.io/en/latest

${C_YEL}Social-media updates.${C_OFF}
EOF

  log "phase 7 complete"
}

# ============================================================================
# Dispatch
# ============================================================================
main() {
  parse_args "$@"
  preflight_common

  case "$SUBCOMMAND" in
    pre-release)
      preflight_pre_release_extras
      phase_1_bump_versions
      phase_2_release_tarball
      phase_3_srpm
      phase_4_github_release
      log "pre-release subcommand complete -- see the 'Next, by hand' block above."
      ;;
    docs)
      preflight_docs_extras
      phase_5_docs_branch
      log "docs subcommand complete."
      ;;
    afterwork)
      preflight_afterwork_extras
      phase_6_afterwork_bump
      phase_7_reminders
      log "afterwork subcommand complete. Release $TAG is done."
      ;;
  esac
}

main "$@"
