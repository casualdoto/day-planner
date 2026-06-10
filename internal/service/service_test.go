package service_test

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"day-planner/internal/db"
	"day-planner/internal/models"
	"day-planner/internal/repository"
	"day-planner/internal/service"
)

type services struct {
	day      *service.DayService
	tasks    *service.TaskService
	tags     *service.TagService
	calendar *service.CalendarService
}

func newServices(t *testing.T) services {
	t.Helper()
	databaseURL := os.Getenv("DAY_PLANNER_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("set DAY_PLANNER_TEST_DATABASE_URL to run Postgres integration tests")
	}
	conn, err := openTestDatabase(t, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	if err := db.Migrate(conn); err != nil {
		t.Fatal(err)
	}
	return services{
		day:      service.NewDayService(repository.NewDayRepository(conn)),
		tasks:    service.NewTaskService(repository.NewTaskRepository(conn)),
		tags:     service.NewTagService(repository.NewTagRepository(conn)),
		calendar: service.NewCalendarService(repository.NewCalendarRepository(conn)),
	}
}

func openTestDatabase(t *testing.T, databaseURL string) (*sql.DB, error) {
	t.Helper()
	schema := fmt.Sprintf("test_%d", time.Now().UnixNano())

	admin, err := db.Open(databaseURL)
	if err != nil {
		return nil, err
	}

	if _, err := admin.Exec(`CREATE SCHEMA ` + schema); err != nil {
		_ = admin.Close()
		return nil, err
	}
	t.Cleanup(func() {
		_, _ = admin.Exec(`DROP SCHEMA IF EXISTS ` + schema + ` CASCADE`)
		_ = admin.Close()
	})

	testURL, err := withSearchPath(databaseURL, schema)
	if err != nil {
		return nil, err
	}
	return db.Open(testURL)
}

func withSearchPath(databaseURL string, schema string) (string, error) {
	parsed, err := url.Parse(databaseURL)
	if err != nil {
		return "", err
	}
	query := parsed.Query()
	query.Set("search_path", schema)
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func TestDayPlanAndResult(t *testing.T) {
	s := newServices(t)
	if err := s.day.SaveDayPlan("2026-06-10", "Plan"); err != nil {
		t.Fatal(err)
	}
	if err := s.day.SaveDayResult("2026-06-10", "Result"); err != nil {
		t.Fatal(err)
	}
	day, err := s.day.GetDayEntry("2026-06-10")
	if err != nil {
		t.Fatal(err)
	}
	if day.PlanText != "Plan" || day.ResultText != "Result" {
		t.Fatalf("unexpected day entry: %+v", day)
	}
}

func TestTaskTagImportanceAndCompletion(t *testing.T) {
	s := newServices(t)
	tag, err := s.tags.CreateTag("Work", "#2563eb")
	if err != nil {
		t.Fatal(err)
	}
	task, err := s.tasks.CreateTask(models.CreateTaskInput{
		Title:      "Write MVP",
		Date:       "2026-06-10",
		Importance: 3,
		TagIDs:     []int64{tag.ID},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(task.Tags) != 1 {
		t.Fatalf("expected assigned tag, got %+v", task.Tags)
	}
	task, err = s.tasks.ToggleTaskCompleted(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !task.Completed || task.CompletedAt == nil {
		t.Fatalf("expected completed task, got %+v", task)
	}
	byTag, err := s.tasks.ListTasks(models.TaskFilter{TagID: &tag.ID})
	if err != nil {
		t.Fatal(err)
	}
	if len(byTag) != 1 {
		t.Fatalf("expected one task by tag, got %d", len(byTag))
	}
	importance := 3
	byImportance, err := s.tasks.ListTasks(models.TaskFilter{Importance: &importance})
	if err != nil {
		t.Fatal(err)
	}
	if len(byImportance) != 1 {
		t.Fatalf("expected one task by importance, got %d", len(byImportance))
	}
	if err := s.tasks.DeleteTask(task.ID); err != nil {
		t.Fatal(err)
	}
	remaining, err := s.tasks.ListTasks(models.TaskFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(remaining) != 0 {
		t.Fatalf("expected deleted task, got %d", len(remaining))
	}
}

func TestCalendarEventsByRange(t *testing.T) {
	s := newServices(t)
	_, err := s.calendar.CreateCalendarEvent(models.CreateCalendarEventInput{
		Title:     "Release",
		StartDate: "2026-06-10",
	})
	if err != nil {
		t.Fatal(err)
	}
	events, err := s.calendar.ListCalendarEvents("2026-06-01", "2026-06-30")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one event, got %d", len(events))
	}
}
