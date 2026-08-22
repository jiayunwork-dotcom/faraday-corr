// Package report renders the results of a faraday-rate calculation as
// human-readable lines for the CLI. It owns the presentation only: every
// number printed here comes from faraday.Result, and no physics is
// recomputed. Unit labels are attached to every quantity so a printed
// line always carries its dimension.
package report
