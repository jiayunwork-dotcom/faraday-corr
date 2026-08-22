// Package metal holds the material data that feeds a faraday-rate
// calculation: a registry of common metals with their molar mass,
// valence and density, and a loader that turns a JSON specification file
// (the same shape consumed by the CLI `rate` subcommand) into a
// faraday.Input.
//
// The registry is deliberately small: iron and aluminium are the two
// metals named by the specification, and they double as the reference
// pair for the cross-rule that the depth rate differs between metals at
// the same current density.
package metal
