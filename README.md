# /ˈkæf.kən/ (k2n)

the project kaeffken, or in short k2n (/keɪ tuː ɛn/ ) is a cli for generating ai based claims or inlcude statements. The generation is based on examples and rulesets. claims are user-facing Kubernetes custom resources (CRDs) that allow application teams (developers, workloads) to request infrastructure or services without knowing the underlying implementation details.

## DEV

```bash
# EXAMPLES FOLDER
go run main.go gen \
--usecase crossplane \
--examples-dir examples/examples \
--ruleset-env-dir examples/ruleset-env \
--ruleset-usecase-dir examples/ruleset-runner \
--instruction "give one runner claim definition for the repo flux and cluster app3. no description. see examples for schema"
```

```bash
# EXAMPLES FILES
go run main.go gen \
--usecase crossplane \
--example-files examples/nginx-git.yaml,examples/nginx-local.yaml \
--ruleset-env-dir examples/ruleset-env \
--ruleset-usecase-dir examples/ruleset-runner \
--instruction "give one runner claim definition for the repo flux and cluster app3. no description. see examples for schema"
```


TEST

```bash
# EXAMPLES FILES
go run main.go gen \
    --examples-dir ./examples \
    --example-files example1.yaml, example2.yaml \
    --ruleset-env-dir ./env-rulesets \
    --ruleset-env-files ./env1.yaml, ./env2.yaml \
    --ruleset-usecase-dir ./usecase-rulesets \
    --ruleset-usecase-files ./usecase1.yaml, ./usecase2.yaml \
    --instruction "Generate XYZ" \
    --usecase myusecase
```
