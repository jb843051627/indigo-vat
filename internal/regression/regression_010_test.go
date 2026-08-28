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
func TestBug10_CSVExportDoesNotReorderInput(t *testing.T) {
    samples := []model.Sample{{ID: "later", TakenAt: time.Unix(20, 0)}, {ID: "early", TakenAt: time.Unix(10, 0)}}
    var out strings.Builder
    if err := codec.WriteSamples(&out, samples, "UTC"); err != nil { t.Fatal(err) }
    if samples[0].ID != "later" { t.Fatal("CSV export reordered the caller slice") }
}
