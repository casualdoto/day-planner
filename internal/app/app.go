package app

import (
	"day-planner/internal/bridge"
	"day-planner/internal/db"
	"day-planner/internal/repository"
	"day-planner/internal/service"
)

func Run() error {
	conn, err := db.Open("")
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := db.Migrate(conn); err != nil {
		return err
	}

	dayRepo := repository.NewDayRepository(conn)
	taskRepo := repository.NewTaskRepository(conn)
	tagRepo := repository.NewTagRepository(conn)
	calendarRepo := repository.NewCalendarRepository(conn)

	bridge := bridge.New(
		service.NewDayService(dayRepo),
		service.NewTaskService(taskRepo),
		service.NewTagService(tagRepo),
		service.NewCalendarService(calendarRepo),
	)

	return RunServer(bridge)
}
