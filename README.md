# prc

**Prometheus Rule Converter**

This program will convert recording rules into `PrometheusRule` type custom resource definitions

Example usage would be:

```
prc convert --from-file=examples/kube-recording-rules.yaml -r k8sPromRules
```
