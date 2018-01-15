# Font Control

A font manager for keeping your font files up-to date

1. Distributed, plain & simple repository model
   1. There's a default, main repository of fonts that are free
      (as in free to distribute) that the working group hosts and maintains.
   2. Anyone can setup a repository; e.g. host fonts internally at your company
      and restrict access to people who have licenses for the fonts.
   3. A repository is [just a set of files](#font-repository)
      served over HTTP(S) and should be serviceable from any HTTP server
      (no dynamic content generation required.) E.g. AWS S3, GitHub, etc.
2. A small and portable program (`fontctrl`) that manages fonts on a computer
   and serves as the client for accessing repositories.
   1. Can be configured to operate on any number of repositories via a simple
      JSON configuration file (`~/.fontctrl.json` on Posix systems, )
   2. Memorizing the user's preference of fonts and any versions limitations
      (e.g. `"inter-ui = 2.*"`)
   3. Keeps font files up to date according to the user's preference and
      availability in configured repositories, by

Long-term goal is to make it easy for designers to keep all their font files
up to date, including updates to fonts they have purchased licenses for.
Operation should be as automatic as possible, optimally operating without user
intervention (e.g. via scheduled invocation or as a service.)


Currently supported systems:

- macOS


## Font Repository

A font repository is a source for fonts served over HTTP.

Paths:

```txt
/index.json
# machine-readable index of the state of the repo

/<fontname>/<fontname>-<version>.zip
# archive containing the font files for <version> of <fontname>

/<fontname>/<fontname>-<version>.json
# description of <version> of <fontname>
```

Shape of `/index.json`:

```json
{
  "fonts": {
    "<fontname>": {
      "name":     "<family_name>",
      "versions": [ "<version>" ]
    }
  }
}
```

Shape of `/<fontname>/<fontname>-<version>.json`:

```json
{
  "version":     "<version>",
  "checksum":    "<sha1-checksum>",
  "name":        "<family_name>",
  "styles":      [ "<style>" ],

  "archive_url": "<archive_url>",
  "description": "<description>",
  "info_url":    "<info_url>",
  "authors":     [ "<author>" ],
  "license":     "<license>"
}
```

Parameters:

- `<version>` should be a
  [SemVer](https://github.com/semver/semver/blob/master/semver.md) formatted
  version string. E.g. "1.1.141+2013"
  If `<version>` is not in the SemVer format, it is interpreted as an opaque
  identifier. I.e. `2.003 < 2.004 => false` but `"2.003" == "2.003" => true`.
- `<fontname>` should be a short name for the font using only the following
  characters: `A-Z`, `a-z`, `0-9`, `-`, `_`, `.` (regexp `[A-Za-z0-9_\-\.]+`)
  E.g. "inter-ui"
- `<family_name>` should be the human-readable name of the font family.
  This should match the `typoFamilyName` record of the font files'
  `name` tables. E.g. "Inter UI"
- `<sha1-checksum>` should be the hexadecimal representation of the
  SHA-1 checksum of the zip archive.
- `<style>` should be the same name as in the respective font file's
  `typoSubfamilyName` record of the `name` table. E.g. "Medium Italic".

Optional parameters:

- `<archive_url>` URL pointing to a font-file archive in an external location.
  Note that `<sha1-checksum>` must match the archive file even if it's served
  from an external location.
- `<description>` should be a human-readable description of the typeface.
- `<info_url>` a well-formed URL pointing to a resource with more information
  about the typeface.
- `<author>` should be the name (and optionally an electronic address) of a
  person or entity that is the (co-)author of the typeface.
- `<license>` should either be a copyright statement, a complete end-user
  license or a url to a complete end-user license for the font files.


## Building & developing

[Posix]

First-time setup:

```
$ ln -s /abs/path/to/misc/fontctrl.json ~/.fontctrl.json
$ ./init.sh
```

Build:

```txt
$ client/build.sh
$ ./build/fontctrl version
```


