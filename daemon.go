package main

import (
	"os"
	"os/exec"
)

func runDaemon(profile string, usermode, systemd bool) {
	cmd := exec.Command("./bedrock_server")
	cmd.Env = append(cmd.Env, "disable_stdout=1", "profile="+profile, "LD_PRELOAD=./ModLoader.so", "LD_LIBRARY_PATH=.:./mods:./libs", "XDG_CACHE_HOME=./cache")
	if usermode {
		cmd.Env = append(cmd.Env, "user_dbus=1")
	}
	cmd.Dir, _ = os.Getwd()
	if systemd {
		cmd.Wait()
	}
}
