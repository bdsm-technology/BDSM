package main

func runDaemon(profile string, usermode bool) {
	printPair("BDSM Version", version)
	cmd := getCommand(profile, usermode)
	cmd.Run()
	cmd.Wait()
}
