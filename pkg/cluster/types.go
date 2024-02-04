package cluster

type Provider struct {
	Kind   string
	Flavor string
}

const (
	AWSClusterKind         = "AWSCluster"
	AWSClusterKindProvider = "capa"

	AWSManagedClusterKind         = "AWSManagedCluster"
	AWSManagedClusterKindProvider = "eks"

	AzureClusterKind         = "AzureCluster"
	AzureClusterKindProvider = "capz"

	AzureManagedClusterKind         = "AzureManagedCluster"
	AzureManagedClusterKindProvider = "aks"

	VCDClusterKind         = "VCDCluster"
	VCDClusterKindProvider = "cloud-director"

	VSphereClusterKind         = "VSphereCluster"
	VSphereClusterKindProvider = "vsphere"

	GCPClusterKind         = "GCPCluster"
	GCPClusterKindProvider = "gcp"

	GCPManagedClusterKind         = "GCPManagedCluster"
	GCPManagedClusterKindProvider = "gke"
)
