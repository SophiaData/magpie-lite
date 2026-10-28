# This fork

`gtk96/magpie` builds [yetone/magpie](https://github.com/yetone/magpie) and
publishes and updates from this repository's Releases. It has three changes
of its own; everything else is upstream's.

## The commits

Each upstream version is two commits on `main`:

- `upstream: vX.Y.Z` — the tree of that upstream tag, laid on as the baseline.
- `fork: ...` — the three changes below, on top of it.

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
  publishes them with the `feed.json` and `notes.json` this build reads.

## Taking an upstream version

```sh
git fetch --depth 1 upstream refs/tags/v0.1.NNNN   # upstream: yetone/magpie
git read-tree --reset -u FETCH_HEAD
git commit -m "upstream: v0.1.NNNN"
git diff <previous fork commit>^ <previous fork commit> | git apply -3
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

Commit the result (`fork: ...`), push `main`, and tag only once the Test run
on it is green: the Release workflow starts from the tag.

## Building and installing here

Go 1.26.8 as `go.mod` says, `CGO_ENABLED=0`; the Windows desktop build needs
no GTK. `make release-windows VERSION=v0.2.N` and `make release-cli
TARGETS="windows/amd64 linux/amd64" VERSION=v0.2.N` write into `dist/`.

To put a build in use: copy the previous `magpie.exe` aside, copy
`dist/magpie-windows-amd64.exe` over `E:\dev-day\magpie-deploy\bin\magpie.exe`,
and start it through explorer (`explorer.exe <path>`) so it outlives the
shell that started it. Then check `magpie --version`, and that
`http://127.0.0.1:3425/v1/models` answers.
