# Dependency Licence Review

This document records the licence review for all direct dependencies in the Crux project.

## Review Process

All dependencies must have their licences reviewed and approved before being added to the project. This ensures:

1. Legal compliance with open source licences
2. No incompatible licence restrictions
3. Clear understanding of obligations

## Acceptable Licences

- MIT
- Apache 2.0
- BSD-3-Clause
- BSD-2-Clause
- ISC

## Unacceptable Licences (without explicit approval)

- GPL (any version)
- AGPL (any version)
- LGPL (any version)

## Current Dependencies

### github.com/spf13/cobra v1.10.2

- **Licence**: Apache-2.0
- **Purpose**: CLI framework for building command-line applications
- **Review Date**: 2026-03-11
- **Reviewed By**: Platform Team
- **Status**: ✅ Approved
- **Notes**: Apache 2.0 is compatible with our project. Widely used in Go ecosystem.
- **Licence URL**: https://github.com/spf13/cobra/blob/main/LICENSE.txt

### github.com/stretchr/testify v1.11.1

- **Licence**: MIT
- **Purpose**: Testing toolkit with assertions and mocking
- **Review Date**: 2026-03-11
- **Reviewed By**: Platform Team
- **Status**: ✅ Approved
- **Notes**: MIT licence is permissive and compatible. Standard testing library in Go.
- **Licence URL**: https://github.com/stretchr/testify/blob/master/LICENSE

## Indirect Dependencies

Indirect dependencies are automatically pulled in by direct dependencies. They are also subject to licence review:

### github.com/spf13/pflag v1.0.9

- **Licence**: BSD-3-Clause
- **Purpose**: Drop-in replacement for Go's flag package (used by cobra)
- **Status**: ✅ Approved (indirect via cobra)

### github.com/inconshreveable/mousetrap v1.1.0

- **Licence**: Apache-2.0
- **Purpose**: Windows console handling (used by cobra)
- **Status**: ✅ Approved (indirect via cobra)

## Review History

| Date | Dependency | Version | Action | Reviewer |
|---|---|---|---|---|
| 2026-03-11 | github.com/spf13/cobra | v1.10.2 | Added | Platform Team |
| 2026-03-11 | github.com/stretchr/testify | v1.11.1 | Added | Platform Team |

## Future Additions

When adding new dependencies, update this document with:
1. Package name and version
2. Licence type
3. Purpose
4. Review date and reviewer
5. Approval status
6. Link to licence file

This ensures we maintain a clear audit trail of all dependency licences.
