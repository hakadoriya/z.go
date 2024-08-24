package cliz

func (c *Command) hasOnlyHiddenSubCommands() bool {
	for _, subcmd := range c.SubCommands {
		if !subcmd.Hidden {
			return false
		}
	}

	return true
}
