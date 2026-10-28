# This fork

`gtk96/magpie` builds [yetone/magpie](https://github.com/yetone/magpie) and
publishes and updates from this repository's Releases. It has three changes
of its own; everything else is upstream's.

## The commits

Each upstream version starts from a baseline commit on `main`:

- `upstream: vX.Y.Z` — the tree of that upstream tag.
- `fork: ...` — the three changes below, in one or more commits after it.

The fork's own patch is therefore `git diff <last upstream commit>..HEAD`:
one diff to carry to the next upstream version.

## The changes

- **Updates come from here.** `internal/update` reads `feed.json` and
  `notes.json` from `releases/latest/download` of this repository (`Feed`,
  `NotesFeed`, `ReleasePage`), and `AutoUpdateDisabled` keeps automatic
  checks, downloads and restarts off: `magpie update` and Settings' Check
  are the ways an update happens.
- **No partners.** `partnersSource` keeps to none unless `MAGPIE_PARTNERS`
  names a list, so the add sheet shows no "Partners · magpie's sponsors".
- **This repository's release workflow.** `.github/workflows/release.yml`
  builds the Windows and Linux apps and the terminal builds here, and
  publishes them with the `feed.json` and `notes.json` this build reads. The
  release is created as a draft, every asset is uploaded, and only then is
  it made Latest (so `releases/latest/download` never serves half a
  release), and the workflow refuses to publish unless Test passed on the
  same commit.

## Taking an upstream version

```sh
git fetch --depth 1 upstream refs/tags/v0.1.NNNN   # upstream: yetone/magpie
git read-tree --reset -u FETCH_HEAD
git commit -m "upstream: v0.1.NNNN"
git diff <last upstream commit>..<previous HEAD> | git apply -3
```

`-3` leaves the three changes as conflicts where upstream moved under them;
resolve them, then run what a release needs:

```sh
gofmt -l $(git diff --name-only HEAD)
go vet -tags nogui ./...                     # GOOS=windows and GOOS=linux
go build -tags production .                  # GOOS=windows, the desktop app
GOOS=linux go build -tags nogui .
go test -tags nogui ./...                    # -timeout 45m; CI runs it too
```

Commit the result (`fork: ...`), push `main`, trigger Test on that commit,
and release only once its run is green: the Release workflow needs a
successful Test run for the same commit.

## What this fork does not do

- **No macOS desktop app.** `AppAsset()` names `magpie-darwin-<arch>.zip`,
  and the release workflow builds only the terminal binary for macOS; a
  macOS desktop build cannot update itself from here.
- **No Stable channel.** The feed is the static `feed.json` of the newest
  release, which carries one release whatever channel is asked for: Stable
  lists the same newest release as every channel.

## Building and installing here

Go 1.26.8 as `go.mod` says, `CGO_ENABLED=0`; the Windows desktop build needs
no GTK. `make release-windows VERSION=v0.2.N` and `make release-cli
TARGETS="windows/amd64 linux/amd64" VERSION=v0.2.N` write into `dist/`.

To put a build in use:

1. Quit the running magpie from its tray icon (Windows refuses to replace a
   running exe), and check no `magpie` process is left.
2. Copy the previous `magpie.exe` aside as `magpie.exe.bak-<old version>-<date>`.
3. Copy `dist/magpie-windows-amd64.exe` over
   `E:\dev-day\magpie-deploy\bin\magpie.exe`.
4. Start it through explorer (`explorer.exe <path>`) so it outlives the shell
   that started it.
5. Check `magpie --version`, and that `http://127.0.0.1:3425/v1/models`
   answers.

The previous binary stays beside it; putting `magpie.exe.bak-...` back is the
rollback.
