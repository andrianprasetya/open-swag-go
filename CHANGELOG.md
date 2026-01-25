# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## [0.3.0] - TBD

### Added
- Example generator (`pkg/examples`) with smart field name heuristics
- Code snippet generator (`pkg/snippets`) for curl, JavaScript, Go, Python, PHP, Ruby
- Version diff tool (`pkg/versioning`) for breaking change detection
- Migration guide generation for breaking changes

## [0.2.0] - 2025-01-25

### Added
- Framework adapters (Chi, Gin, Echo, Fiber)
- `HandlerHTTP()` and `SpecHandlerHTTP()` methods for http.Handler compatibility
- `MountGroup()` functions for router groups (Gin, Echo, Fiber)

## [0.1.0] - 2025-01-25

### Added
- Initial release
- Core endpoint definition structs
- Scalar UI integration
- Basic schema generation from Go structs
- net/http handler support
- Authentication schemes (Bearer, API Key, Basic, Cookie)
- Struct tag parsing (swagger, example, format, description)

### Features
- `openswag.New()` - Create new docs instance
- `docs.Add()` - Add endpoint definition
- `docs.AddAll()` - Add multiple endpoints
- `docs.Mount()` - Mount handlers on mux
- `openswag.Body()` - Request body helper
- `openswag.Response()` - Response helper
- `openswag.PathParam()` / `QueryParam()` / `HeaderParam()` - Parameter helpers
- `openswag.BearerAuth()` / `APIKeyAuth()` / `BasicAuth()` - Auth scheme helpers
