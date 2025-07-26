# Product Requirements Document: Enhanced `grab update` Command

## Executive Summary

Enhance the existing `grab update` command to accept an optional package name argument, enabling users to update specific packages instead of updating all packages in their configuration file.

## Background

Currently, the `grab update` command updates all packages defined in the user's configuration file (`~/.grab/config.yml`) with their latest upstream versions. This enhancement will provide users with more granular control by allowing them to update individual packages when needed.

## Goals

- Provide users with the ability to update specific packages
- Maintain backward compatibility with existing functionality
- Preserve consistent user experience and output formatting

## Requirements

### Functional Requirements

#### Command Interface
```bash
# Existing behavior (unchanged)
grab update

# New behavior
grab update [package-name]
```

#### Behavior Specification

**Default behavior (no arguments):**
- Updates all packages defined in `~/.grab/config.yml`
- Maintains existing functionality without changes

**Single package behavior:**
- Updates only the specified package in the user's configuration
- Performs exact string match against package names in `~/.grab/config.yml`
- Uses case-sensitive matching
- Only operates on packages already present in user's configuration

#### Error Handling

- **Package not found**: Command must terminate with descriptive error message if the specified package name does not exist in `~/.grab/config.yml`
- **Invalid input**: Standard error handling for malformed package names
- **Configuration errors**: Existing error handling patterns remain unchanged

### Technical Requirements

#### Configuration Scope
- Only operates on packages present in user's `~/.grab/config.yml`
- Does not add new packages to configuration
- Does not interact with packages outside user's config

#### Validation
- Validate package name exists in user configuration before processing
- Implement fail-fast behavior with clear error messaging

#### Output Consistency
- Maintain identical output format between single-package and all-package updates
- Preserve existing logging patterns, progress indicators, and status reporting
- No changes to success/failure messaging format

#### Backward Compatibility
- Zero breaking changes to existing `grab update` command behavior
- Maintain existing command interface and output for current users
- Preserve all existing functionality and performance characteristics

## Success Criteria

1. Users can successfully update individual packages using `grab update [package-name]`
2. Existing `grab update` behavior remains completely unchanged
3. Clear error messages are displayed for invalid package names
4. Output formatting remains consistent across both usage patterns
5. No regression in existing functionality or performance

## Implementation Notes

- Command implementation should leverage existing update logic in `cmd/update.go`
- Package validation should occur early in the command execution flow
- Error messages should follow established patterns in the codebase
- Consider extending existing updater functionality in `pkg/updater.go`

## Timeline

To be determined based on development capacity and priority.