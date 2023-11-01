package cctv

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
)

type Service struct {
	f servicemgr.ServiceFactory
}

func New(f servicemgr.ServiceFactory) *Service {
	return &Service{f: f}
}
