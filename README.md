# alidns-ingress-updater
A k8s app to auto update alidns from ingress

## Install

### Install via Kubectl
```bash
$ kubectl apply -f https://github.com/xzzpig/alidns-ingress-updater/raw/main/deploy/bundle.yaml
```


# Cookbook
1. Create AliDnsAccount
```yaml
apiVersion: network.xzzpig.com/v1
kind: AliDnsAccount
metadata:
  name: alidnsaccount-sample
spec:
  accessKeyId: YourKeyIdHere
  accessKeySecret: YourKeySecretHere
  domainName: yourdomain.com
```

2. Create/Modify Ingress
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example
  annotations:
    xzzpig.com/alidns-ignore: "false" # (optional)whether ignore this ingress to update
spec:
  rules:
  - host: example.yourdomain.com
    http:
      paths:
      - backend:
          service:
            name: example
            port:
              number: 5000
        path: /
        pathType: Prefix
```
3. Enjoy!
> This app will auto match AliDnsAccount.spec.domainName to Ingress.spec.rules.host and modify alidns record(`example`) to clusters public ip

## Reference
### Ingress Annotations
| name | describe | type |
| ---- | -------- | ---- |
| xzzpig.com/alidns-ignore | whether ignore this ingress to update | "true"/"false" |