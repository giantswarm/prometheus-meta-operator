package remotewriteconfig

import "testing"

func TestShardComputationScaleUp(t *testing.T) {
	expected := 1
	result := computeShards(0, 20)
	if result != expected {
		t.Errorf(`expected computeShards(0, 20) to be %d, got %d`, expected, result)
	}

	expected = 2
	result = computeShards(0, 21)
	if result != expected {
		t.Errorf(`expected computeShards(0, 21) to be %d, got %d`, expected, result)
	}

	expected = 3
	result = computeShards(0, 41)
	if result != expected {
		t.Errorf(`expected computeShards(0, 41) to be %d, got %d`, expected, result)
	}
}

func TestShardComputationReturnsAtLeast1Shart(t *testing.T) {
	expected := 1
	result := computeShards(0, 0)
	if result != expected {
		t.Errorf(`expected computeShards(0, 0) to be %d, got %d`, expected, result)
	}

	expected = 1
	result = computeShards(0, -5)
	if result != expected {
		t.Errorf(`expected computeShards(0, -5) to be %d, got %d`, expected, result)
	}
}

func TestShardComputationScaleDown(t *testing.T) {
	expected := 2
	result := computeShards(1, 21)
	if result != expected {
		t.Errorf(`expected computeShards(1, 21) to be %d, got %d`, expected, result)
	}

	expected = 2
	result = computeShards(2, 19)
	if result != expected {
		t.Errorf(`expected computeShards(2, 19) to be %d, got %d`, expected, result)
	}

	expected = 2
	result = computeShards(2, 16)
	if result != expected {
		t.Errorf(`expected computeShards(2, 16) to be %d, got %d`, expected, result)
	}

	// threshold hit
	expected = 1
	result = computeShards(2, 15)
	if result != expected {
		t.Errorf(`expected computeShards(2, 15) to be %d, got %d`, expected, result)
	}
}
