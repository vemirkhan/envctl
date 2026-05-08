package cmd

func init() {
	extraCommands = append(extraCommands, NewWatchCmd)
}
