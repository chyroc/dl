package cmd_helper

func Assert(err error) {
	if err != nil {
		panic(err)
	}
}
