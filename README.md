# /ˈkæf.kən/ (k2n)

the project kaeffken, or in short k2n (/keɪ tuː ɛn/ ) is a cli for generating ai based claims or inlcude statements. The generation is based on examples and rulesets. claims are user-facing Kubernetes custom resources (CRDs) that allow application teams (developers, workloads) to request infrastructure or services without knowing the underlying implementation details.

## DEV

```bash
go run main.go gen \
--usecase crossplane \
--examples-dir examples/examples \
--ruleset-env-dir examples/ruleset-env \
--ruleset-usecase-dir examples/ruleset-runner \
--instruction "give one runner claim definition for the repo flux and cluster app3. no description. see examples for schema"
```
