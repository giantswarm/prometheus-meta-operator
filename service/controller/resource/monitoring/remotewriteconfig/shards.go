package remotewriteconfig

import "math"

const (
	shardStep                = float64(1_000_000)
	shardScaleDownPercentage = float64(0.20)
	shardScaleDownThreshold  = shardScaleDownPercentage * shardStep
)

// We want to start with 1 prometheus-agent for each 1M time series with a scale down 25% threshold.
func computeShards(currentShardCount int, timeSeries float64) int {
	desiredShardCount := int(math.Ceil(timeSeries / shardStep))

	// Compute Scale Down
	if currentShardCount > desiredShardCount {
		// We get the rest of a division of timeSeries by shardStep and we compare it with the scale down threshold
		if math.Mod(timeSeries, shardStep) > shardStep-shardScaleDownThreshold {
			desiredShardCount = currentShardCount
		}
	}

	// We always have a minimum of 1 agent, even if there is no worker node
	if desiredShardCount <= 0 {
		return 1
	}
	return desiredShardCount
}
