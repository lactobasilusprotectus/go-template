package jwt

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)

	mock := NewMockJwtInterface(ctrl)

	mock.EXPECT().GenerateToken(gomock.Any(), gomock.Any())
	mock.EXPECT().ExtractToken(gomock.Any(), gomock.Any())

	_, _ = mock.GenerateToken(nil, JwtData{})
	_, _ = mock.ExtractToken(nil, "")
}
