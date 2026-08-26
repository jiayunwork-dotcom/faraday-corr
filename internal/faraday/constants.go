package faraday

const Faraday = 96485.33212

const yearSeconds = 365 * 24 * 3600

const milliPerCm = 10.0

const microPerAmpere = 1.0e6

const K = milliPerCm * (1.0 / microPerAmpere) * yearSeconds / Faraday

const MicroAmpsPerCm2 = "uA/cm^2"
