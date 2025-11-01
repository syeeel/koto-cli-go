# Changelog

All notable changes to koto will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.9] - 2025-11-01

### Added
- **Export Feature**: Added dedicated export screen with interactive file path selection
  - Default file path with timestamp: `~/.koto/export_YYYYMMDD_HHMMSS.json`
  - Success screen displaying file location, todo count, and file size
  - Accessible via `/export` command
- **Import Feature**: Added multi-step import flow with confirmation
  - File path input with validation
  - Preview screen showing number of todos to import
  - Confirmation step before actual import
  - Success/failure summary screen
  - Accessible via `/import` command

### Changed
- **Recent Todos Display**: Increased maximum character count for task titles from 15 to 22 characters (1.5x increase)
- **Export/Import Commands**: Now use dedicated interactive screens instead of inline execution

### Removed
- **Unnecessary File**: Removed `--help` file from repository

### Documentation
- **README.md**: Translated to English for global audience
- **CHANGELOG.md**: Translated to English for global audience

### Technical
- `internal/tui/model.go`: Added `ViewModeExport` and `ViewModeImport` view modes, added export/import state fields
- `internal/tui/views.go`: Implemented `renderExportView()` and `renderImportView()` functions
- `internal/tui/update.go`: Added key handling for export/import views, implemented `handleExportEnter()` and `handleImportEnter()` helpers
- `internal/tui/commands.go`: Removed old inline export/import command handlers

## [1.0.7] - 2025-10-30

### Changed
- **Version Display Format**: Changed to a simple and readable format like "Version: 1.0.7"
- **Makefile**: Improved to automatically extract the latest version number from CHANGELOG
- **Pomodoro Timeout Sound**: Changed to a higher and longer sound (880Hz, 500ms) to be more noticeable

### Fixed
- **Version Information**: Fixed issue where "dev" was displayed, now shows the correct version number

## [1.0.6] - 2025-10-30

### Added
- **Responsive Layout**: Dynamically adjusts layout based on terminal width
- **Minimum Terminal Width Check**: Displays error message when terminal width is less than 100 characters
- **Dynamic Width Calculation System**: Added `DynamicWidths` structure to centrally manage full-screen width
- **Cross-Terminal Compatible ASCII Borders**: Custom borders using simple ASCII characters (+, -, |)
- **Dynamic Version Information System**: Automatically injects version, commit hash, and build date at build time
- **Makefile**: Added comprehensive Makefile to streamline development builds (build, test, clean, install, run, release, etc.)

### Fixed
- **Task List**: Fixed issue where selected row was displayed in 2 lines when focused
- **Startup Screen**: Fixed issue where borders were cut off (centered ASCII art and ToDo box)
- **Detail Screen**: Fixed issue where borders were cut off (all boxes now support dynamic width)
- **Pomodoro Screen**: Fixed issue where borders were cut off (progress bar and info box support dynamic width)
- **Pomodoro Screen**: Changed task information (Task ID and task name) to center display
- **Pomodoro Screen**: Improved progress bar center display
- **macOS Terminal Compatibility**: Replaced Unicode box-drawing characters with ASCII characters so they display correctly in all terminals

### Changed
- **Main List Screen**: Column widths calculated dynamically, title column changed to variable width
- **Detail Screen**: 3-column layout dynamically adjusted with proportional distribution
- **Pomodoro Screen**: All elements centered according to terminal width
- **Border Style**: Changed from RoundedBorder/NormalBorder to highly compatible ASCII borders
- **Version Display**: Display detailed version information (commit, build date) on startup screen and command line

### Technical
- `internal/tui/styles.go`: Added dynamic width calculation functions and helper functions, added ASCII border definitions
- `internal/tui/views.go`: Made all screens responsive, added minimum width check feature, replaced all borders with ASCII, improved progress bar center display
- `internal/tui/banner.go`: Exported version information variables, changed `GetVersion()` to detailed display
- `cmd/koto/main.go`: Inject version information into TUI package
- `Makefile`: Set version information with ldflags at build time, automated various development tasks
- Implemented center alignment using Lipgloss's `PlaceHorizontal`
- Replaced all Unicode box-drawing characters with ASCII characters to improve cross-terminal compatibility
- Integrated with GoReleaser so version information is automatically injected in release builds

## [1.0.3] - 2025-10-30

### Fixed
- **install.sh**: Fixed archive name format (using koto-cli-go project name)
- **install.sh**: Removed 'v' prefix from version number to match actual GoReleaser archive name

### Changed
- Fixed installation script to work properly on macOS/Linux

## [1.0.2] - 2025-10-30

### Changed
- Temporarily disabled Homebrew Tap configuration (setup in progress)
- Removed deprecated configuration (format_overrides) for GoReleaser v2 compatibility
- Simplified archive configuration

### Fixed
- Fixed 401 error in Homebrew Tap during release (disabled configuration)

## [1.0.1] - 2025-10-30

### Changed
- Updated Homebrew instructions in README to "Coming Soon 🚧" status
- Clarified that Homebrew will be available from next release (v1.0.1 or later)

## [1.0.0] - 2025-10-30

### Added
- **GoReleaser Integration**: Automated release workflow
- **Multi-Platform Build**: macOS (Intel/Apple Silicon), Linux (amd64/arm64), Windows (amd64)
- **Installation Script**: One-line installation (curl | sh)
- **Version Display**: Display version, commit, and build date with `--version` flag
- **GitHub Actions**: Automated release workflow on tag push

### Documentation
- Release guide (docs/RELEASE.md)
- Homebrew setup guide (docs/SETUP_HOMEBREW.md)
- Quick start guide (docs/QUICKSTART_RELEASE.md)

### Infrastructure
- `.goreleaser.yaml`: GoReleaser configuration
- `.github/workflows/release.yml`: Automated release workflow
- `install.sh`: User-facing installation script

---
