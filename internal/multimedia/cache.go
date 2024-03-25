package multimedia

import (
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
)

func (s *Service) movieFromCache(id string) (*rms_library.FoundMovie, bool) {
	mov := &rms_library.FoundMovie{}
	val, ok := s.cache.Load(id)
	if ok {
		mov, ok = val.(*rms_library.FoundMovie)
	}

	return mov, ok
}
