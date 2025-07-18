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
--examples-dirs /home/sthings/projects/stuttgart-things/terraform/builds/labda-dagger-vm, /home/sthings/projects/stuttgart-things/terraform/builds/labda-maverick-vm \
--ruleset-env-dir /home/sthings/projects/ai/terraform/ruleset-terraformvm \
--ruleset-usecase-dir /home/sthings/projects/ai/terraform/ruleset-terraformvm \
--usecase terraform \
--instruction "give me one terraformconfig for a medium vm with a random name (movie reference) + one ansible playbook with baseos profile. no description. see examples for reference " \
--destination "/tmp/krock3/bla9/" # "" = stdout;  /tmp/allinone.yaml = allinonefile; /tmp/new-folder = new folder + single files
```
