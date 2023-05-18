package remotewriteconfig

import "math"

const (
	shardNodeStep            = 20
	shardScaleDownPercentage = float64(0.25)
	shardScaleDownThreshold  = shardScaleDownPercentage * shardNodeStep
)

// We want to start with 1 prometheus-agent every 20 nodes with 25% threshold to scale down.
func computeShards(currentShardCount int, nodeCount int) int {
	desiredShardCount := int(math.Ceil(float64(nodeCount) / shardNodeStep))

	// Compute Scale Down
	if currentShardCount > desiredShardCount {
		// We get the rest of a division of nodeCount by ShardNodeStep (215 % 20 = 15 and we check if the threshold is passed)
		if float64(nodeCount%shardNodeStep) > shardNodeStep-shardScaleDownThreshold {
			desiredShardCount = currentShardCount
		}
	}

	// We always have a minimum of 1 agent, even if there is no worker node
	if desiredShardCount <= 0 {
		return 1
	}
	return desiredShardCount
}
