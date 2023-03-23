package multimedia

import (
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"net/http"
	"time"
)

const moviesSearchLimit = 10
const searchTimeout = 30 * time.Second

type searchPageContext struct {
	ui.PageContext
	Query  string
	Movies []*rms_library.FoundMovie
}

func (s *Service) getSearchHandler(ctx *gin.Context) {
	page := searchPageContext{
		PageContext: *ui.New(),
		Query:       ctx.Query("q"),
	}
	if page.Query != "" {
		resp, err := s.f.NewLibrary().SearchMovie(ctx, &rms_library.SearchMovieRequest{Text: page.Query, Limit: moviesSearchLimit}, client.WithRequestTimeout(searchTimeout))
		if err != nil {
			logger.Errorf("Search movies failed: %s", err)
			ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось обратиться к сервису поиска медиа")
			return
		}
		page.Movies = resp.Movies
	}
	ctx.HTML(http.StatusOK, "multimedia.search.tmpl", &page)
}
