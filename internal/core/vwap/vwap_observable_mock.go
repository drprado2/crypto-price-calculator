package vwap

import "context"

type (
	VwapObservableMock struct {
		MockHandleNewVwap func(ctx context.Context, event *VwapUpdatedEvent)
	}
)

func (v *VwapObservableMock) HandleNewVwap(ctx context.Context, event *VwapUpdatedEvent) {
	if v.MockHandleNewVwap != nil {
		v.MockHandleNewVwap(ctx, event)
	}
}
