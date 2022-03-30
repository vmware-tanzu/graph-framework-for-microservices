package helper

func GetCrdParentsMap() map[string][]string {
	return map[string][]string{
		"root.helloworld.com":   {""},
		"config.helloworld.com": {"root.helloworld.com"},
	}
}
