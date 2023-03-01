<div align="center">

# Dyad

A Go module for making deep copies of arbitrary values.

[![Documentation](https://img.shields.io/badge/go.dev-documentation-007d9c?&style=for-the-badge)](https://pkg.go.dev/github.com/dogmatiq/dyad)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/dyad.svg?&style=for-the-badge&label=semver)](https://github.com/dogmatiq/dyad/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/dogmatiq/dyad/ci.yml?style=for-the-badge&branch=main)](https://github.com/dogmatiq/dyad/actions/workflows/ci.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/dyad/main.svg?style=for-the-badge)](https://codecov.io/github/dogmatiq/dyad)

</div>

Dyad makes clones of Go values. It attempts to make all types clonable, or
fallback to predictable (and configurable) behavior when non-clonable types
(such as channels) are encountered.
