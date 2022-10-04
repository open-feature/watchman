# watchman

![goversion](https://img.shields.io/github/go-mod/go-version/open-feature/watchman/main)
![version](https://img.shields.io/badge/version-pre--alpha-green)
![licence](https://img.shields.io/github/license/open-feature/watchman)
![size](https://img.shields.io/github/repo-size/open-feature/watchman)

<img src="images/doggo.png" width="250px;" />

A general purpose admission controller that uses open-feature to control it's capabilities.

This is a demonstration project that is designed to familiarise the user with how open-feature can be used within Kubernetes to control infrastructure. 

This project demonstrates:

- [open-feature golang SDK](https://github.com/open-feature/golang-sdk)
- [flagD](https://github.com/open-feature/flagd)
- [open-feature operator](https://github.com/open-feature/open-feature-operator)

_But how does it work?_

<img src="images/003.png" width="850px;">

## Run the demo

- Setup a Kubernetes cluster
- Install open-feature-operator as described [here](https://github.com/open-feature/open-feature-operator#deploy-the-latest-release)
- `kubectl create namespace watchman-system`
- `kubectl create -f config/samples/flags.yaml`
- `kubectl create -f config/webhook/certificate.yaml`
- `IMG=cnskunkworks/watchman:0.0.3 make deploy`

You should see that you have watchman injected with flagd.

<img src="images/001.png" width="450px;">

This now means you can turn on or off pod admission with the watchman `FeatureFlagConfiguration` custom resource named `watchman-flags/watchman-system`

<img src="images/002.png" width="450px;">

The end result should be that you can turn on or off the admission controller with the `FeatureFlagConfiguration` custom resource.

```
❯ kubectl apply -f config/samples/pod.yaml
Error from server (flag set to block admission): error when creating "config/samples/pod.yaml": admission webhook "mutate-demo.openfeature.dev" denied the request: flag set to block admission
```
