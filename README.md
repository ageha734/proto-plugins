# proto-plugins

A community-maintained registry of [proto](https://moonrepo.dev/proto) TOML plugins for CLI tools.

## Available Plugins

|Plugin|Tool|Homepage|
|---|---|---|
|actionlint|GitHub Actions linter|[rhysd/actionlint](https://github.com/rhysd/actionlint)|
|argo|Argo Workflows CLI|[argoproj/argo-workflows](https://github.com/argoproj/argo-workflows)|
|commitlint|Commit message linter|[KeisukeYamashita/commitlint](https://github.com/KeisukeYamashita/commitlint)|
|dprint|Code formatter|[dprint/dprint](https://github.com/dprint/dprint)|
|ghalint|GitHub Actions linter|[suzuki-shunsuke/ghalint](https://github.com/suzuki-shunsuke/ghalint)|
|hadolint|Dockerfile linter|[hadolint/hadolint](https://github.com/hadolint/hadolint)|
|helm|Kubernetes package manager|[helm/helm](https://github.com/helm/helm)|
|helmfile|Helm declarative spec|[helmfile/helmfile](https://github.com/helmfile/helmfile)|
|hyperfine|CLI benchmarking tool|[sharkdp/hyperfine](https://github.com/sharkdp/hyperfine)|
|kubeconform|Kubernetes manifest validator|[yannh/kubeconform](https://github.com/yannh/kubeconform)|
|kubectl|Kubernetes CLI|[kubernetes/kubectl](https://github.com/kubernetes/kubectl)|
|kubectx|Kubernetes context switcher|[ahmetb/kubectx](https://github.com/ahmetb/kubectx)|
|kustomize|Kubernetes configuration|[kubernetes-sigs/kustomize](https://github.com/kubernetes-sigs/kustomize)|
|lefthook|Git hooks manager|[evilmartians/lefthook](https://github.com/evilmartians/lefthook)|
|pinact|GitHub Actions SHA pinner|[suzuki-shunsuke/pinact](https://github.com/suzuki-shunsuke/pinact)|
|shellcheck|Shell script linter|[koalaman/shellcheck](https://github.com/koalaman/shellcheck)|
|shfmt|Shell script formatter|[mvdan/sh](https://github.com/mvdan/sh)|
|task|Task runner|[go-task/task](https://github.com/go-task/task)|
|terraform-docs|Terraform documentation|[terraform-docs/terraform-docs](https://github.com/terraform-docs/terraform-docs)|
|terragrunt|Terraform wrapper|[gruntwork-io/terragrunt](https://github.com/gruntwork-io/terragrunt)|
|tflint|Terraform linter|[terraform-linters/tflint](https://github.com/terraform-linters/tflint)|
|tilt|Local Kubernetes dev|[tilt-dev/tilt](https://github.com/tilt-dev/tilt)|
|trivy|Security scanner|[aquasecurity/trivy](https://github.com/aquasecurity/trivy)|
|zizmor|GitHub Actions security|[woodruffw/zizmor](https://github.com/woodruffw/zizmor)|

## Usage

Add a plugin to your `.prototools`:

```toml
[plugins]
actionlint = "github://ageha734/proto-plugins/toml/actionlint.toml"
```

Then install:

```bash
proto install actionlint
```

## Requirements

- [proto](https://moonrepo.dev/proto) >= 0.57.3

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and contribution guidelines.

## License

MIT
