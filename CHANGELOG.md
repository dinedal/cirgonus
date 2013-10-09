# 0.2.0 (10/9/2013)

All of these features are covered in the README documentation.

* Cirgonus no longer polls on each hit -- it does so on a tick value then
  serves requests from cache. You can adjust the frequency at which it polls by
  tweaking the "PollInterval" json configuration, which defaults to 60 seconds
  for `cirgonus generate`.
* Now logging to syslog -- you can adjust the facility at which it logs to by
  tweaking `Facility` in the json configuration which defaults to `daemon` and
  `LogLevel` for scoping messages, which defaults to `info`.
* Result publishing lets you push metrics to cirgonus which will then be picked
  up by circonus as a metric -- great for custom tooling!
* Makefile now to build releases and copies of `cirgonus` and `cstat`

# 0.1.0 (09/22/2013)

* First release!
