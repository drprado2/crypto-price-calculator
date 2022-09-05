package coinbase

import "context"

type (
	ControllerMock struct {
		MockHandle func(ctx context.Context, message []byte) error
	}
)

func (cm *ControllerMock) Handle(ctx context.Context, message []byte) error {
	if cm.MockHandle != nil {
		return cm.MockHandle(ctx, message)
	}

	return nil
}
