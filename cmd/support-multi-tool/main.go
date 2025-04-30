package main

import (
	"os"
	"os/signal"
	"support-multi-tool/config"
	"support-multi-tool/internal/app"
	"support-multi-tool/pkg/logger"
	"syscall"
)

func main() {
	// читаем данные из конфиг файла
	cfg := config.MustLoad()

	// инициализируем логгер и уровень логирования для окружения
	log, logFile := logger.Initial(cfg)
	if logFile != nil {
		defer logFile.Close()
	}

	// start applications services
	application := app.New(log, cfg)
	application.MustRun(cfg)

	// мониторинг сигналов ОС для корректного прерывания/завершения процесса
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-stop

	log.Warn("Forced application shutdown...")
	application.Stop()
	log.Info("Application has shutdown")
}
