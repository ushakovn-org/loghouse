package main

import (
  "context"
  "log"

  "github.com/gammazero/workerpool"
  "github.com/ushakovn-org/loghouse/internal/app/log_house"
  "github.com/ushakovn-org/loghouse/internal/config"
  "github.com/ushakovn-org/loghouse/internal/pkg/storage"
  "github.com/ushakovn/boiler/pkg/app"
)

func main() {
  ctx := context.Background()
  boiler := app.NewApp()

  addr := config.NewProvider(ctx, config.ClickhouseAddress).Watch(ctx)

  logs, err := storage.NewStorage(ctx, storage.Config{
    Addr: []string{addr.Provide().String()},
  })
  if err != nil {
    log.Fatalf("storage.NewStorage: %v", err)
  }
  workers := config.Get(ctx, config.LoghouseWorkersCount).Int()
  pool := workerpool.New(workers)

  loghouse := log_house.NewLogHouse(log_house.Config{
    Pool:    pool,
    Storage: logs,
  })
  boiler.Run(loghouse)
}
