package cctv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/RacoonMediaServer/rms-packages/pkg/schedule"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"github.com/RacoonMediaServer/rms-web/internal/ui"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"go-micro.dev/v4/logger"
)

const maxIntervals = 10

var timePointRegex = regexp.MustCompile(`(\d\d):(\d\d)`)

func (s *Service) getSchedulesHandler(ctx *gin.Context) {
	page := struct {
		ui.PageContext
		Schedules []scheduleItem
	}{PageContext: *ui.New()}

	resp, err := s.f.NewCctvSchedules().GetSchedulesList(ctx, &empty.Empty{})
	if err != nil {
		logger.Errorf("Get schedules failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось связаться с системой видеонаблюдения")
		return
	}

	page.Schedules = make([]scheduleItem, len(resp.Result))
	for i, sched := range resp.Result {
		page.Schedules[i] = scheduleItem{Name: sched.Name, ID: sched.Id}
	}

	ctx.HTML(http.StatusOK, "cctv.schedules.tmpl", page)
}

func (s *Service) getNewScheduleHandler(ctx *gin.Context) {
	page := struct {
		ui.PageContext
		Intervals []struct{}
	}{
		PageContext: *ui.New(),
		Intervals:   make([]struct{}, maxIntervals),
	}
	ctx.HTML(http.StatusOK, "cctv.schedules.new.tmpl", page)
}

func getPositionalName(fieldName string, number int) string {
	return fmt.Sprintf("%s_%d", fieldName, number)
}

func parseTimePoint(tm string) (schedule.TimePoint, bool) {
	matches := timePointRegex.FindStringSubmatch(tm)
	if len(matches) < 3 {
		return schedule.TimePoint{}, false
	}
	hours, _ := strconv.ParseUint(matches[1], 10, 16)
	minutes, _ := strconv.ParseUint(matches[2], 10, 16)
	return schedule.TimePoint{Hours: int(hours), Minutes: int(minutes)}, true
}

func parseScheduleForm(ctx *gin.Context) (schedule.Representation, bool) {
	sched := schedule.Representation{}
	for i := 0; i < maxIntervals; i++ {
		if ctx.PostForm(getPositionalName("enabled", i)) != "on" {
			continue
		}
		begin, ok := parseTimePoint(ctx.PostForm(getPositionalName("begin", i)))
		if !ok {
			return sched, false
		}
		end, ok := parseTimePoint(ctx.PostForm(getPositionalName("end", i)))
		if !ok {
			return sched, false
		}
		interval := schedule.Interval{
			Begin:     begin,
			End:       end,
			IsHoliday: ctx.PostForm(getPositionalName("is_holiday", i)) == "on",
		}
		sched.Intervals = append(sched.Intervals, interval)
	}
	return sched, sched.IsValid()
}

func (s *Service) postNewScheduleHandler(ctx *gin.Context) {
	sched, ok := parseScheduleForm(ctx)
	if !ok {
		ui.DisplayError(ctx, http.StatusBadRequest, "Неверное задано расписание")
		return
	}

	req := rms_cctv.Schedule{Name: ctx.PostForm("name"), Content: sched.String()}
	_, err := s.f.NewCctvSchedules().AddSchedule(ctx, &req)
	if err != nil {
		logger.Errorf("Create schedule failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось создать расписание")
		return
	}
	ui.DisplayOK(ctx, "Расписание создано", "/cctv/schedules")
}

func (s *Service) getScheduleHandler(ctx *gin.Context) {
	page := struct {
		ui.PageContext
		ID        string
		Name      string
		Intervals []periodItem
	}{PageContext: *ui.New()}

	id := ctx.Param("id")
	sched, err := s.f.NewCctvSchedules().GetSchedule(ctx, &rms_cctv.GetScheduleRequest{Id: id})
	if err != nil {
		logger.Errorf("Fetch schedule failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось получить детали расписания")
		return
	}

	var parsed schedule.Representation
	err = json.Unmarshal([]byte(sched.Content), &parsed)
	if err != nil {
		logger.Errorf("Parse schedule failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось разобрать формат расписания")
		return
	}

	page.ID = id
	page.Name = sched.Name
	page.Intervals = make([]periodItem, maxIntervals)
	count := len(parsed.Intervals)
	if count > maxIntervals {
		count = maxIntervals
	}
	for i := 0; i < count; i++ {
		page.Intervals[i].Enabled = true
		page.Intervals[i].Interval = parsed.Intervals[i]
	}

	ctx.HTML(http.StatusOK, "cctv.schedules.edit.tmpl", page)
}

func (s *Service) postScheduleHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	sched, ok := parseScheduleForm(ctx)
	if !ok {
		ui.DisplayError(ctx, http.StatusBadRequest, "Неверное задано расписание")
		return
	}

	req := rms_cctv.Schedule{
		Id:      id,
		Name:    ctx.PostForm("name"),
		Content: sched.String(),
	}

	_, err := s.f.NewCctvSchedules().ModifySchedule(ctx, &req)
	if err != nil {
		logger.Errorf("Modify schedule failed: %s", err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось изменить расписание")
		return
	}
	ui.DisplayOK(ctx, "Расписание изменено", "/cctv/schedules")
}

func (s *Service) deleteScheduleHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := s.f.NewCctvSchedules().DeleteSchedule(ctx, &rms_cctv.DeleteScheduleRequest{Id: id})
	if err != nil {
		logger.Errorf("Delete schedule %s failed: %s", id, err)
		ui.DisplayError(ctx, http.StatusInternalServerError, "Не удалось удалить расписание")
		return
	}
	ui.DisplayOK(ctx, "Расписание удалено", "/cctv/schedules")
}
