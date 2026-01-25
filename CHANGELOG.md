# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## [0.1.0] - TBD

### Added
- Initial release
- Core endpoint definition structs
- Scalar UI integration
- Basic schema generation from Go structs
- net/http handler support
- Authentication schemes (Bearer, API Key, Basic, Cookie)
- Struct tag parsing (swagger, example, format, description)

### Features
- `openapi.New()` - Create new docs instance
- `docs.Add()` - Add endpoint definition
- `docs.AddAll()` - Add multiple endpoints
- `docs.Mount()` - Mount handlers on mux
- `openapi.Body()` - Request body helper
- `openapi.Response()` - Response helper
- `openapi.PathParam()` / `QueryParam()` / `HeaderParam()` - Parameter helpers
- `openapi.BearerAuth()` / `APIKeyAuth()` / `BasicAuth()` - Auth scheme helpers
