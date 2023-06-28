package settings

import (
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"go-micro.dev/v4"
)

type Service struct {
	f   servicemgr.ServiceFactory
	pub micro.Event
}

func New(f servicemgr.ServiceFactory, pub micro.Event) *Service {
	return &Service{f: f, pub: pub}
}
