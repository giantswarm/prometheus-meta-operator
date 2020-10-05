# Prometheus storage

When possible on customers installations, we decided to use dynamic provisioning to create volumes to store prometheus data.

As some installations do not support dynamic provisioning (some on-prem installations), we decided to ask customers to provision volumes in the following pattern:

- 1 Volume for the control plane.
- 1 Volume per tenant cluster managed by the control plane.
- 1 extra volume that we can use in case a new cluster is created.

## Retention

Max Retention duration: **2 weeks**

Max Retention size: **95 GB**

## Volume sizing

Previous Volume Size: **20 GB**, this was not enough as some prometheus needs more

Current Volume Size: **100GB**

Reasons: https://gigantic.slack.com/archives/C01176DKNP4/p1601471521029900


## Current State of on-prem storage

Installation | Datastore | Dynamic provisioning |
------------ | --------- | -------------------- |
anubis       | iSCSI     | no                   |
amagon       | NFS       | no                   |
buffalo      | NFS       | yes                  |
dinosaur     | NFS       | no                   |
dragon       | NFS       | no                   |
puma         | NFS       | no                   |
geckon       | NFS       | no                   |
gorgoth      | NFS       | no                   |