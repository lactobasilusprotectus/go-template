package redis

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestNewMockInterface(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	mock := NewMockInterface(mockCtrl)

	mock.EXPECT().Get(gomock.Any())
	mock.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any())

	_, _ = mock.Get("")
	_ = mock.Set("", "", 0)

}
