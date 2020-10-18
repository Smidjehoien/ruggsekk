package internal

import (
	"context"
	"fmt"
	"github.com/ruckstack/ruckstack/common/ui"
	"github.com/ruckstack/ruckstack/server/daemon/internal/containerd"
	"github.com/ruckstack/ruckstack/server/daemon/internal/etcd"
	"github.com/ruckstack/ruckstack/server/daemon/internal/k3s"
	"github.com/ruckstack/ruckstack/server/daemon/internal/webserver"
	"github.com/ruckstack/ruckstack/server/internal/environment"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
)

func Start() error {
	//serverReady = make(chan bool, 1)
	//

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	packageConfig := environment.PackageConfig

	fmt.Printf("Starting %s version %s\n", packageConfig.Name, packageConfig.Version)

	if err := os.MkdirAll(filepath.Join(environment.ServerHome, "logs"), 0755); err != nil {
		return err
	}
	serverLog, err := os.OpenFile(filepath.Join(environment.ServerHome, "logs", "server.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer serverLog.Close()
	log.SetOutput(serverLog)

	fmt.Printf("    Server log: %s\n", serverLog.Name())
	fmt.Printf("    K3S log: %s\n", filepath.Join(environment.ServerHome, "logs", "k3s.log"))

	if err := ioutil.WriteFile(filepath.Join(environment.ServerHome, "data", "server.pid"), []byte(strconv.Itoa(os.Getpid())), 0644); err != nil {
		return err
	}

	log.Println("Starting server components...")

	if err := webserver.StartWebserver(ctx); err != nil {
		return fmt.Errorf("error starting webserver: %s", err)
	}

	if err := containerd.Start(ctx); err != nil {
		return fmt.Errorf("error starting containerd: %s", err)
	}

	if err := etcd.Start(ctx); err != nil {
		return fmt.Errorf("error starting etcd: %s", err)
	}

	if err := k3s.Start(ctx); err != nil {
		return fmt.Errorf("error starting k8s server: %s", err)
	}

	//go monitor.StartMonitor()

	//
	//if err := agent.Start(); err != nil {
	//	ui.Fatalf("error starting k3s agent: %s", err)
	//}
	//ui.Println("K3s agent started successfully")

	select {
	case <-ctx.Done():
		ui.Println("Server shutting down...")
		ui.VPrintf("Shutdown reason: ", ctx.Err())
		return nil
	}
}

func Stop() {
	k3s.Stop()
}
