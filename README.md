# K8s-Auditor

Kubernetes-Auditor is an open source project for facilitating global config search and store events which can be queried and audited. 

----

## Issues in K8s which k8s-Auditor tackle

- There is no way to fetch changes in config-map and secrets natively. 

- K8s events are short-lived and vanishes after certain time. This makes debugging harder for unexpected issues.


---

## Features

- Config Map: Any change in config maps would be captured and queried. This would serve as a single place for config search for the cluster.

- Events: Kubernetes events would be stored with proper metadata in order to be meaningful and could be used for auditing purposes.

---

## Future Scope

- Config Map: This can be extended to test for changes in config maps before it can be rolled out.

- Secrets: This can be extended for secrets as well. Ofcourse, we would need some layer that could mask important data which would go into the query engine.
