package multimedia

import (
	rms_library "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-library"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"net/http"
	"sort"
	"strconv"
	"time"
)

const libraryRequestTimeout = 20 * time.Second
const itemsByPage = 5

type libraryPageContext struct {
	ui.PageContext
	Movies []*rms_library.Movie
	Pages  []int
	Page   int
	Type   string
	Sort   string
}

type playPageContext struct {
	ui.PageContext
	Movie    *rms_library.Movie
	BeginURL string
}

func parseMovieType(movType string) *rms_library.MovieType {
	tp := rms_library.MovieType_Film
	switch movType {
	case "films":
		return &tp
	case "tv-series":
		tp = rms_library.MovieType_TvSeries
		return &tp
	default:
		return nil
	}
}

func sortResults(results []*rms_library.Movie, sortId string) {
	sortFunc := func(i, j int) bool {
		return results[i].Info.Title < results[j].Info.Title
	}
	switch sortId {
	case "desc":
		sortFunc = func(i, j int) bool {
			return results[j].Info.Title < results[i].Info.Title
		}
	case "rating":
		sortFunc = func(i, j int) bool {
			return results[j].Info.Rating < results[i].Info.Rating
		}
	}

	sort.Slice(results, sortFunc)
}

func parsePageNo(p string) int {
	no, err := strconv.ParseInt(p, 10, 32)
	if err != nil {
		return 1
	}
	return int(no)
}

func paginate(movies []*rms_library.Movie, pageNo int) []*rms_library.Movie {
	idx := (pageNo - 1) * itemsByPage
	if idx >= len(movies) {
		return []*rms_library.Movie{}
	}
	res := movies[idx:]
	if len(res) > itemsByPage {
		res = res[0:itemsByPage]
	}
	return res
}

func (s *Service) libraryHandler(ctx *gin.Context) {
	movType := parseMovieType(ctx.Query("type"))
	resp, err := s.f.NewLibrary().GetMovies(ctx, &rms_library.GetMoviesRequest{Type: movType}, client.WithRequestTimeout(libraryRequestTimeout))
	if err != nil {
		logger.Errorf("Get movies failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось обратиться к сервису библиотеки")
		return
	}

	sortResults(resp.Result, ctx.Query("sort"))

	page := libraryPageContext{
		PageContext: *ui.New(),
		Movies:      resp.Result,
		Page:        parsePageNo(ctx.Query("page")),
		Sort:        ctx.Query("sort"),
		Type:        ctx.Query("type"),
	}

	pages := len(resp.Result) / itemsByPage
	if len(resp.Result)%itemsByPage != 0 {
		pages++
	}
	for i := 1; i <= pages; i++ {
		page.Pages = append(page.Pages, i)
	}
	page.Movies = paginate(page.Movies, page.Page)

	ctx.HTML(http.StatusOK, "multimedia.library.tmpl", &page)
}

func (s *Service) deleteMovieHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := s.f.NewLibrary().DeleteMovie(ctx, &rms_library.DeleteMovieRequest{ID: id})
	if err != nil {
		logger.Errorf("Delete movie failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось удалить фильм/сериал")
		return
	}
	ui.DisplayOK(ctx, "Удалено", "/multimedia/library")
}

func (s *Service) playHandler(ctx *gin.Context) {
	resp, err := s.f.NewLibrary().GetMovie(ctx, &rms_library.GetMovieRequest{ID: ctx.Param("id")})
	if err != nil {
		logger.Errorf("Get movie failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось достать информацию о фильме/сериале")
		return
	}
	page := playPageContext{
		PageContext: *ui.New(),
		Movie:       resp.Result,
		BeginURL:    resp.Result.Film.Files[0],
	}
	ctx.HTML(http.StatusOK, "multimedia.library.play.tmpl", &page)
}
