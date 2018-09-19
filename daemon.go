package main

func runDaemon(profile string, usermode, systemd bool) {
	cmd := getCommand(profile, usermode)
	if systemd {
		cmd.Wait()
	}
}
