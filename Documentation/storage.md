# Prometheus storage

When possible on customers installations, we decided to use dynamic provisioning to create volumes to store prometheus data.

As some installations do not support dynamic provisioning (some on-prem installations), we decided to ask customers to provision volumes in the following pattern:

- 1 Volume for the management cluster.
- 1 Volume per workload cluster managed by the management cluster.
- 1 extra volume that we can use in case a new cluster is created.

## Retention

Max Retention size: **95 GB**

## Volume sizing

We decided to use the same volume size accross all our clusters as to avoid having to figure out all the details regarding sizing per cluster, volume expansion and so on.

Current Volume Size: **100GB**

Reasons: 

As described in the [prometheus doc](https://www.prometheus.io/docs/prometheus/latest/storage/) and [Stack Exchange](https://devops.stackexchange.com/questions/9298/how-to-calculate-disk-space-required-by-prometheus-v2-2):

**needed_disk_space** = **retention_time_in_seconds** * **ingested_samples_per_second** ** **bytes_per_sample**

With:

    - retention_time_in_seconds   = 1209600s
    - bytes_per_sample            = rate(prometheus_tsdb_compaction_chunk_size_bytes_sum[1d]) / rate(prometheus_tsdb_compaction_chunk_samples_sum[1d]) for bytes / sample
    - ingested_samples_per_second = rate(prometheus_tsdb_head_samples_appended_total[2h])


By taking one of our biggest cluster on asgard, we found **bytes_per_sample** = _1.0811677050996586_ and **ingested_samples_per_second** = _41713_

With this, **needed_disk_space** = 1209600 * 41713 * 1.08 bytes ~= 54 GB

As our storage needs will grow, we decided to go with 100 GB for now
