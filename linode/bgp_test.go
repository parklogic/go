package linode

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBGPPeerAddresses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		region  string
		wantErr bool
	}{
		{name: "known region", region: "us-east"},
		{name: "unknown region", region: "missing", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			require := require.New(t)
			assert := assert.New(t)

			got, err := GetBGPPeerAddresses(test.region)
			if test.wantErr {
				require.ErrorIs(err, ErrInvalidRegion)
				assert.Nil(got)

				return
			}

			require.NoError(err)

			dcID := DatacenterIDFromRegion[test.region]
			want := make([]string, len(got))

			for i := range want {
				want[i] = fmt.Sprintf("%s:%d:34::%d", PeerAddressPrefix, dcID, i+1)
			}

			assert.Equal(want, got)
		})
	}
}
