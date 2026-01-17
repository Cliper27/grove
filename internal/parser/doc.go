// Package parser implements loading, parsing, and validation of .gro schema files.
//
// A schema defines filesystem rules using required, allowed, and denied nodes.
// Schemas may include other schemas explicitly via the `include` directive and
// reference them by name.
//
// Included schemas form a dependency graph which must be acyclic (a DAG).
// Cyclic includes and duplicate schema names are rejected at load time.
//
// Schema loading is path-based and cached globally within the package.
// Errors returned by this package are wrapped and may be inspected using errors.Is.
package parser
