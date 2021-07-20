package main
//https://www.twblogs.net/a/5c841957bd9eee35fc13e7ff
import (
  "context"
  "log"
  "net/http"
  "os"
  "time"

  "go.uber.org/fx"
)

// Logger構造函數
func NewLogger() *log.Logger {
  logger := log.New(os.Stdout, "" /* prefix */, 0 /* flags */)
  logger.Print("Executing NewLogger.")
  return logger
}

// http.Handler構造函數
func NewHandler(logger *log.Logger) (http.Handler, error) {
  logger.Print("Executing NewHandler.")
  return http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
    logger.Print("Got a request.")
  }), nil
}

// http.ServeMux構造函數
func NewMux(lc fx.Lifecycle, logger *log.Logger) *http.ServeMux {
  logger.Print("Executing NewMux.")

  mux := http.NewServeMux()
  server := &http.Server{
    Addr:    ":8080",
    Handler: mux,
  }

  lc.Append(fx.Hook{  // 自定義生命週期過程對應的啓動和關閉的行爲
    OnStart: func(context.Context) error {
      logger.Print("Starting HTTP server.")
      go server.ListenAndServe()
      return nil
    },
    OnStop: func(ctx context.Context) error {
      logger.Print("Stopping HTTP server.")
      return server.Shutdown(ctx)
    },
  })

  return mux
}

// 註冊http.Handler
func Register(mux *http.ServeMux, h http.Handler) {
  mux.Handle("/", h)
}

func main() {
  app := fx.New(
    fx.Provide(
      NewLogger,
      NewHandler,
      NewMux,
    ),
    fx.Invoke(Register),  // 通過invoke來完成Logger、Handler、ServeMux的創建
  )
  startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
  defer cancel()
  if err := app.Start(startCtx); err != nil {  // 手動調用Start
    log.Fatal(err)
  }

  http.Get("http://localhost:8080/")  // 具體操作

  stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
  defer cancel()
  if err := app.Stop(stopCtx); err != nil { // 手動調用Stop
    log.Fatal(err)
  }

}
