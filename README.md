# faraday-corr

Uniform-corrosion Faraday rate calculator. Given a metal's corrosion
current density `i_corr`, molar mass `M`, valence `n` and density `ρ`, it
computes the mass loss rate and the penetration depth rate of the metal
under uniform corrosion.

This is the electrochemical-equivalent kernel, not a corrosion database
and not a maintenance ledger: the output is the rate pair

```
mdot  = M · i_corr / (n · F)
CR    = K · (M / n) · i_corr / ρ
```

where `F` is the Faraday constant pinned to the CODATA value
(96485.33212 C/mol) and `K` is the pinned unit-chain constant that folds
μA/cm² and g/cm³ into mm/y (`K = 10 · 1e-6 · seconds-per-year / F`). The
two formulas must agree: a doubling of the current doubles both rates, a
doubling of the valence halves them, and doubling the density leaves the
mass rate unchanged while halving the depth rate.

## What it computes

For one input specification it reports:

- mass loss rate `mdot` in g/(cm²·s)
- annual mass loss in g/(cm²·y), also g/(m²·y) and mdd
- corrosion depth rate in mm/y and μm/y, also mils per year (mpy)
- total corrosion current `I = i_corr · A` when an exposed area is given
- cumulative mass loss over the area and exposure time

The `reverse` subcommand inverts the kernel (solve `i_corr` from a target
rate); `life` converts a depth rate into a corrosion-allowance lifetime;
`schedule` tabulates the linear-in-time cumulative loss; `chain` walks
the unit chain and verifies that the two depth-rate paths close.

## Usage

Run a rate calculation from a JSON specification:

```bash
go run . rate example/fe-seawater.json
```

Example output for iron in seawater at 10 μA/cm²:

```
Corrosion rate:    0.1159 mm/y (115.9 um/y)
Annual mass loss:  0.09126 g/(cm^2 y)
```

Other subcommands:

```bash
go run . metals                    # list the built-in metal registry
go run . validate example/fe-seawater.json
go run . reverse example/fe-seawater.json   # needs a target_cr_mm_y field
go run . schedule example/fe-seawater.json
go run . life example/fe-seawater.json --allowance 0.5 --thickness 8
go run . chain example/fe-seawater.json
go run . json example/fe-seawater.json
go run . template Al               # print a ready-to-edit JSON spec
go run . help
```

### Specification files

Input files are JSON with fixed units:

| field           | unit        | notes                                   |
|-----------------|-------------|-----------------------------------------|
| `i_corr`        | μA/cm²      | corrosion current density, must be ≥ 0  |
| `M`             | g/mol       | molar mass, must be > 0                 |
| `n`             | electrons   | valence, must be > 0                    |
| `rho`           | g/cm³       | density, must be > 0                    |
| `area`          | cm²         | optional exposed area                   |
| `duration_y`    | years       | optional exposure time                  |
| `metal`         | symbol/name | optional; fills M/n/ρ from the registry |

`example/fe-seawater.json` and `example/al-seawater.json` are checked-in
specs. Invalid input (negative current, zero valence, non-positive
density, malformed JSON) produces an error on stderr and a non-zero exit
code.

## Cross rules

These relationships are pinned by the test suite and printed by the
`invariants` subcommand:

- `i = 0` ⇒ `CR = 0`
- `i × 2` ⇒ `mdot × 2` and `CR × 2`
- `n × 2` ⇒ both rates halve
- `ρ × 2` ⇒ mass rate unchanged, depth rate halves
- Fe vs Al at the same `i_corr` ⇒ different `CR`
- time × 2 ⇒ cumulative mass loss × 2

## Key conventions

- `F = 96485.33212` C/mol (CODATA), never an approximation.
- The year is 365 days; `K = 315.36 / F`.
- All physics live in `internal/faraday`; `internal/metal` holds the
  registry and JSON loading; `internal/report` renders output.
- The built-in registry covers iron (Fe → Fe²⁺, M = 55.845, ρ = 7.874)
  and aluminium (Al → Al³⁺, M = 26.9815, ρ = 2.70).

## Build & test

```bash
go build ./...
go test ./...
```
