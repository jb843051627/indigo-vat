package regression

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "io"
    "net/http"
    "net/http/httptest"
    "strings"
    "sync"
    "testing"
    "time"

    "github.com/jb843051627/indigo-vat/internal/codec"
    "github.com/jb843051627/indigo-vat/internal/clock"
    "github.com/jb843051627/indigo-vat/internal/httpapi"
    "github.com/jb843051627/indigo-vat/internal/metrics"
    "github.com/jb843051627/indigo-vat/internal/model"
    "github.com/jb843051627/indigo-vat/internal/queue"
    "github.com/jb843051627/indigo-vat/internal/service"
    "github.com/jb843051627/indigo-vat/internal/store"
    "github.com/jb843051627/indigo-vat/internal/validation"
    "github.com/jb843051627/indigo-vat/internal/worker"
)

var (
    _ = errors.Is
    _ = fmt.Sprintf
    _ = io.Discard
    _ = http.MethodGet
    _ = httptest.NewRecorder
    _ = strings.Builder{}
    _ = sync.Mutex{}
    _ = clock.NewFixed
    _ = codec.Encode
    _ = httpapi.New
    _ = metrics.New
    _ = queue.New
    _ = service.ErrNotFound
    _ = store.Open
    _ = validation.IsID
    _ = worker.Retry
    _ = sql.ErrNoRows
    _ = time.Now
)

func testDB(t *testing.T) (*store.DB, *service.Service) {
    t.Helper()
    db, err := store.Open(t.TempDir() + "/indigo.db")
    if err != nil { t.Fatal(err) }
    t.Cleanup(func() { _ = db.Close() })
    return db, service.New(db)
}

func seedCycle(t *testing.T, svc *service.Service) model.Cycle {
    t.Helper()
    vat, err := svc.CreateVat(context.Background(), model.VatInput{Name: "Vat", Site: "Shed", Capacity: 100, TargetPH: 10, TargetTemp: 30, Timezone: "UTC"})
    if err != nil { t.Fatal(err) }
    recipe, err := svc.CreateRecipe(context.Background(), model.RecipeInput{Name: "Recipe", TargetHue: 220, HueTolerance: 5, MinMinutes: 1, MaxMinutes: 20})
    if err != nil { t.Fatal(err) }
    if _, err := svc.AddStage(context.Background(), recipe.ID, model.StageInput{Name: "Ferment", Minutes: 1, PHMin: 8, PHMax: 12, TempMin: 20, TempMax: 40}); err != nil { t.Fatal(err) }
    recipe, err = svc.PublishRecipe(context.Background(), recipe.ID)
    if err != nil { t.Fatal(err) }
    cycle, err := svc.StartCycle(context.Background(), model.CycleInput{VatID: vat.ID, RecipeID: recipe.ID})
    if err != nil { t.Fatal(err) }
    return cycle
}
func TestBug20_OpenCriticalAlertBlocksRelease(t *testing.T) {
    db, svc := testDB(t)
    cycle := seedCycle(t, svc)
    matured, err := svc.AdvanceCycle(context.Background(), cycle.ID, model.CycleMatured, cycle.Revision, "matured")
    if err != nil { t.Fatal(err) }
    if err := db.PutInspection(context.Background(), model.Inspection{ID: "inspection-pass", CycleID: matured.ID, Result: model.InspectionPass, Score: 90, CreatedAt: time.Now()}); err != nil { t.Fatal(err) }
    if err := db.PutAlert(context.Background(), model.Alert{ID: "alert-open", CycleID: matured.ID, Level: model.AlertCritical, Code: "ph", State: model.AlertOpen, CreatedAt: time.Now()}); err != nil { t.Fatal(err) }
    if _, err := svc.ReleaseCycle(context.Background(), matured.ID, matured.Revision); !errors.Is(err, service.ErrNotReady) { t.Fatalf("want not ready, got %v", err) }
}
