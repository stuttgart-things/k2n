# /ˈkæf.kən/ (k2n)

the project kaeffken, or in short k2n (/keɪ tuː ɛn/ ) is a cli for generating ai based claims or inlcude statements. The generation is based on examples and rulesets. claims are user-facing Kubernetes custom resources (CRDs) that allow application teams (developers, workloads) to request infrastructure or services without knowing the underlying implementation details.

## DEV



## USAGE

```bash
export GEMINI_API_KEY=""

k2n gen \
--examples-dirs examples/examples \
--ruleset-env-dir examples/ruleset-env \
--ruleset-usecase-dir examples/ruleset-runner \
--usecase crosssplane \
--instruction "give me a runner-claim for the repository dagger for the cluster sthings"
```
