package main

func runDaemon(profile string, usermode bool) {
	cmd := getCommand(profile, usermode)
	cmd.Wait()
}
