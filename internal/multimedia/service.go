package multimedia

import (
	"sync"

	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
)

type Service struct {
	f     servicemgr.ServiceFactory
	cache sync.Map
}

func New(f servicemgr.ServiceFactory) *Service {
	return &Service{
		f: f,
	}
}
