package bridge

import (
	"day-planner/internal/models"
	"day-planner/internal/service"
)

type Bridge struct {
	day      *service.DayService
	tasks    *service.TaskService
	tags     *service.TagService
	calendar *service.CalendarService
}

func New(day *service.DayService, tasks *service.TaskService, tags *service.TagService, calendar *service.CalendarService) *Bridge {
	return &Bridge{day: day, tasks: tasks, tags: tags, calendar: calendar}
}

func (b *Bridge) GetDayEntry(date string) (models.DayEntry, error) {
	return b.day.GetDayEntry(date)
}

func (b *Bridge) SaveDayPlan(date string, planText string) error {
	return b.day.SaveDayPlan(date, planText)
}

func (b *Bridge) SaveDayResult(date string, resultText string) error {
	return b.day.SaveDayResult(date, resultText)
}

func (b *Bridge) CreateTask(input models.CreateTaskInput) (models.Task, error) {
	return b.tasks.CreateTask(input)
}

func (b *Bridge) UpdateTask(id int64, input models.UpdateTaskInput) (models.Task, error) {
	return b.tasks.UpdateTask(id, input)
}

func (b *Bridge) DeleteTask(id int64) error {
	return b.tasks.DeleteTask(id)
}

func (b *Bridge) ToggleTaskCompleted(id int64) (models.Task, error) {
	return b.tasks.ToggleTaskCompleted(id)
}

func (b *Bridge) ListTasks(filter models.TaskFilter) ([]models.Task, error) {
	return b.tasks.ListTasks(filter)
}

func (b *Bridge) CreateTag(name string, color string) (models.Tag, error) {
	return b.tags.CreateTag(name, color)
}

func (b *Bridge) UpdateTag(id int64, name string, color string) (models.Tag, error) {
	return b.tags.UpdateTag(id, name, color)
}

func (b *Bridge) DeleteTag(id int64) error {
	return b.tags.DeleteTag(id)
}

func (b *Bridge) ListTags() ([]models.Tag, error) {
	return b.tags.ListTags()
}

func (b *Bridge) CreateCalendarEvent(input models.CreateCalendarEventInput) (models.CalendarEvent, error) {
	return b.calendar.CreateCalendarEvent(input)
}

func (b *Bridge) UpdateCalendarEvent(id int64, input models.CreateCalendarEventInput) (models.CalendarEvent, error) {
	return b.calendar.UpdateCalendarEvent(id, input)
}

func (b *Bridge) DeleteCalendarEvent(id int64) error {
	return b.calendar.DeleteCalendarEvent(id)
}

func (b *Bridge) ToggleCalendarEventCompleted(id int64) (models.CalendarEvent, error) {
	return b.calendar.ToggleCalendarEventCompleted(id)
}

func (b *Bridge) ListCalendarEvents(startDate string, endDate string) ([]models.CalendarEvent, error) {
	return b.calendar.ListCalendarEvents(startDate, endDate)
}
