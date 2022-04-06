package capacity

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseCapacity(t *testing.T) {
	capacity, err := ParseCapacities("-1gb1M2MB5b")
	require.NoError(t, err)
	require.Equal(t, capacity.String(), "-1GB3MB5B")
	require.Equal(t, int(capacity), "-1076887557")
}
