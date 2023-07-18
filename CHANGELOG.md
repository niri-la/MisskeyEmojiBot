# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Delete channels and archive threads after deciding whether or not to apply for pictograms. [`#13`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/13)

### Changed
- Changed from hard-coding such as Token to using env files [`#11`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/11)

### Fixed
- change Regex pattern [`#11`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/11)
- Numbers are replaced by _ when they are included in pictogram names [`#12`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/12)
- Examples of characters that can be entered [`#14`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/14)

### Removed

## [0.0.1] - 2023-07-17
### Added
- emoji managed functions
- added libraries
  - go-misskey (connect to misskey)
  - discord-go (connect to discord)
