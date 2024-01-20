package core

import (
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSyncCharger(t *testing.T) {
	tc := []struct {
		status                      api.ChargeStatus
		expected, actual, corrected bool
	}{
		{api.StatusA, false, false, false},
		{api.StatusC, false, false, true}, // disabled but charging
		{api.StatusA, false, true, true},
		{api.StatusA, true, false, false},
		{api.StatusA, true, true, true},
	}

	ctrl := gomock.NewController(t)

	for _, tc := range tc {
		t.Logf("%+v", tc)

		charger := api.NewMockCharger(ctrl)
		charger.EXPECT().Enabled().Return(tc.actual, nil).AnyTimes()

		if tc.status == api.StatusC {
			charger.EXPECT().Enable(tc.corrected).Times(1)
		}

		lp := &Loadpoint{
			log:     util.NewLogger("foo"),
			clock:   clock.New(),
			charger: charger,
			status:  tc.status,
			enabled: tc.expected,
		}

		require.NoError(t, lp.syncCharger())
		assert.Equal(t, tc.corrected, lp.enabled)
	}
}
