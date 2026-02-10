# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Multiple logger initialization options (Development, Production, Custom)
- `Config` struct for flexible logger configuration
- `Panic` log level support
- `Sync()` method for flushing logs
- `With()` method for creating child loggers with pre-populated fields
- `WithZapFields()` method for using zap fields directly
- `GetZapLogger()` method for direct zap.Logger access
- Comprehensive documentation and examples
- Full test coverage
- Support for multiple output paths and formats
- Godoc comments for all public APIs

### Changed
- Improved `getFields()` method with better performance
- Enhanced error handling in logger initialization
- Better variable naming (receiver -> l/f for brevity)

### Fixed
- Memory allocation optimization in field conversion

## [1.0.0] - 2024-02-10

### Added
- Initial release
- Basic logging functionality with Debug, Info, Notice, Warn, Error, Fatal levels
- Integration with uber-go/zap
- Implementation of gsr.Logger interface
- Field-based structured logging
