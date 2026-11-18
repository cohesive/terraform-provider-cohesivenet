# Changelog

## 1.1.0

SECURITY:

* The provider no longer disables TLS verification process-wide by modifying
  `http.DefaultTransport`. Each client now uses its own transport. (CWE-295)

FEATURES:

* provider: new `verify_cert` option to control TLS certificate verification.
  It defaults to `false` to preserve existing behaviour (new controllers
  present a self-signed certificate that must be accepted before a trusted
  certificate can be uploaded). Set `verify_cert = true` (or
  `CN_VERIFY_CERT=true`) once a trusted certificate is installed on the
  controller. A warning is shown when verification is disabled only because of
  the default; set the option explicitly to acknowledge it. **The default will
  change to `true` in 2.0.**

DEPRECATIONS:

* The `vns3 {}` block is deprecated and will be removed in a future release.
  Configure host and credentials in the provider block instead.
