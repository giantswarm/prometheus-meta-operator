package certificates

// The Kubernetes Secrets we currently use for prometheus (e.g: $CLUSTER_ID-prometheus) are held in the default namespace.
// We want to run the tenant cluster Prometheus servers in a per-tenant cluster namespace ($CLUSTER_ID-prometheus). prometheus-operator does not support referencing a secret in a different namespace.
// So, we need to copy the secret default/$CLUSTER_ID-prometheus to the tenant cluster prometheus namespace - e.g: $CLUSTER_ID-prometheus/cluster-certificates
