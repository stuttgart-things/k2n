# /ˈkæf.kən/ (k2n)

the project kaeffken, or in short k2n (/keɪ tuː ɛn/ ) is a cli for generating ai based claims or inlcude statements. The generation is based on examples and rulesets. claims are user-facing Kubernetes custom resources (CRDs) that allow application teams (developers, workloads) to request infrastructure or services without knowing the underlying implementation details.

## DEV


## USAGE-EXAMPLES

<details><summary>VERBOSE OUTPUT OF THE PROMPT (w/o SENDING IT)</summary>

```bash
k2n gen \
--examples-dirs _examples/examples \
--ruleset-env-dir _examples/ruleset-env \
--ruleset-usecase-dir _examples/ruleset-runner \
--usecase crosssplane \
--instruction "give me a runner-claim for the repository dagger for the cluster sthings" \
--verbose=true \
--prompt-to-ai=false
```

</details>

<details><summary>PROMPT AI + OUTPUT TO STDOUT</summary>

```bash
export GEMINI_API_KEY=""

k2n gen \
--examples-dirs _examples/examples \
--ruleset-env-dir _examples/ruleset-env \
--ruleset-usecase-dir _examples/ruleset-runner \
--usecase crosssplane \
--instruction "give me a runner-claim for the repository dagger for the cluster sthings"
```

</details>

<details><summary>PROMPT AI + OUTPUT TO STDOUT</summary>

```bash
k2n gen \
--examples-dirs _examples/examples \
--ruleset-env-dir _examples/ruleset-env \
--ruleset-usecase-dir _examples/ruleset-runner \
--usecase crosssplane \
--instruction "give me a runner-claim for the repository dagger for the cluster sthings. add also a proposal for a branch name and a pr title" \
--destination "/tmp/runner/" \
--verbose=true
```

</details>


## AUTHOR

```bash
Patrick Hermann, stuttgart-things 07/2025
```

## LICENSE

Licensed under the Apache License, Version 2.0 (the "License").

You may obtain a copy of the License at [apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0).

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an _"AS IS"_ basis, without WARRANTIES or conditions of any kind, either express or implied.

See the License for the specific language governing permissions and limitations under the License.
