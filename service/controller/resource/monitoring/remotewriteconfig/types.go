package remotewriteconfig

import "math"

type PrometheusAgentShardingStrategy struct {
	// Configures the number of series needed to add a new shard. Computation is number of series / ShardScaleUpSeriesCount
	ShardScaleUpSeriesCount float64
	// Percentage of needed series based on ShardScaleUpSeriesCount to scale down agents
	ShardScaleDownPercentage float64
}

// We want to start with 1 prometheus-agent for each 1M time series with a scale down 20% threshold.
func (pass PrometheusAgentShardingStrategy) ComputeShards(currentShardCount int, timeSeries float64) int {
	shardScaleDownThreshold := pass.ShardScaleDownPercentage * pass.ShardScaleUpSeriesCount
	desiredShardCount := int(math.Ceil(timeSeries / pass.ShardScaleUpSeriesCount))

	// Compute Scale Down
	if currentShardCount > desiredShardCount {
		// We get the rest of a division of timeSeries by shardStep and we compare it with the scale down threshold
		if math.Mod(timeSeries, pass.ShardScaleUpSeriesCount) > pass.ShardScaleUpSeriesCount-shardScaleDownThreshold {
			desiredShardCount = currentShardCount
		}
	}

	// We always have a minimum of 1 agent, even if there is no worker node
	if desiredShardCount <= 0 {
		return 1
	}
	return desiredShardCount
}
