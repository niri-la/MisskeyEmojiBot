# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Fixed

### Removed

## [1.4.0] - 2024-08-25
### Changed
- refactor: reorganize file structure for better clarity and maintainability [`#97`](https://github.com/niri-la/MisskeyEmojiBot/pull/97)

### Fixed
- fix: emoji name input was not converted to lowercase [`#91`](https://github.com/niri-la/MisskeyEmojiBot/pull/91)

## [1.3.5] - 2024-02-24
### Changed
- add sendDirectMessage emoji.Name [`#86`](https://github.com/niri-la/MisskeyEmojiBot/pull/86)

## [1.3.4] - 2024-01-14
### Changed
- Change message. [`#82`](https://github.com/niwaniwa/MisskeyEmojiBot/issues/82)

### Fixed
- Change the UUID version to v4. [`#18`](https://github.com/niwaniwa/MisskeyEmojiBot/issues/18)

## [1.3.3] - 2024-01-13
### Added
- Add the URL of the application source from DM. [`#79`](https://github.com/niwaniwa/MisskeyEmojiBot/issues/79)
- Add request cancel button. [`8e50ab0`](https://github.com/niwaniwa/MisskeyEmojiBot/commit/8e50ab092e383b31a2d52a64122ebf1c1fe5848e)
- Delete threads of applications that have not been completed within a certain period of time. [`#71`](https://github.com/niwaniwa/MisskeyEmojiBot/issues/71)

### Changed
- Increase auto-hide time of threads. [`#75`](https://github.com/niwaniwa/MisskeyEmojiBot/issues/75)

### Fixed
- Fixed a problem where the part where pictograms had to be specified was blank. [`#76`](https://github.com/niwaniwa/MisskeyEmojiBot/issues/76)

## [1.3.2] - 2023-07-22

### Changed
- Change rate-limit [`#69`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/69)

## [1.3.1] - 2023-07-22

### Fixed
- fixed message when request completed [`#62`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/62)
- fixed nsfw button [`#66`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/66)
- fixed license field [`#67`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/67)

## [1.3.0] - 2023-07-22

### Changed
- Change channel to thread [`#60`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/60)

### Fixed
- Fixed responce state [`#59`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/59)
- Fixed properly reflect values [`61`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/61)

## [1.2.0] - 2023-07-22

### Changed
- Separation of licenses and remarks [`#54`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/55)

### Fixed
- Change logger level timing [`#57`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/57)

## [1.1.0] - 2023-07-22

### Added
- Added Emoji request abort function [`#36`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/36)
- Added user permission checker [`#42`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/42)
- Multilingual implementation [`#45`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/45)
- Added feature to delete thread messages [`#51`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/51)

### Changed
- fmt to logrus [`#33`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/33)
- Change request flow [`#49`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/49)

### Fixed
- Fixed DM bug. [`#37`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/37)
- Fixed request input alias [`#50`](https://github.com/niwaniwa/MisskeyEmojiBot/pull/50)

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
