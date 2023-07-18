# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Fixed

### Removed

## [1.0.1] - 2023-07-19

### Fixed
- fix dockerfile [`#29`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/29)

## [1.0.0] - 2023-07-18

### Added
- Delete channels and archive threads after deciding whether or not to apply for pictograms. [`#13`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/13)
- User feedback implementation.　[`#15`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/15)
- add emoji note function. [`#21`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/21)

### Changed
- Changed from hard-coding such as Token to using env files [`#11`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/11)

### Fixed
- change Regex pattern [`#11`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/11)
- Numbers are replaced by _ when they are included in pictogram names [`#12`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/12)
- Examples of characters that can be entered [`#14`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/14)
- Fixed not being able to press a button after making a request. [`#20`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/20)
- add interaction message. [`#23`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/23)
- add emoji length check. [`#26`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/26)

## [0.0.1] - 2023-07-17
### Added
- emoji managed functions
- added libraries
  - go-misskey (connect to misskey)
  - discord-go (connect to discord)
