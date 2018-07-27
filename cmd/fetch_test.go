package cmd_test

/*func TestCmdFetch(t *testing.T) {
	cases := []struct {
		url string
		name string
		customBotDir string
		isErr bool
		errMsg string
	}{
		{
			url: "https://github.com/liagame/java-bot/archive/master.zip",
			name: "birko",
			customBotDir: "",
			isErr: false,
			errMsg: "failed to fetch java bot and put into working dir",
		},
	}

	for i, c := range cases {
		if os.Getenv("BE_CRASHER") == "1" {
			Crasher()
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestCrasher")
		cmd.Env = append(os.Environ(), "BE_CRASHER=1")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		t.Fatalf("process ran with err %v, want exit status 1", err)
	}
}*/

/*func TestCmdFetch(t *testing.T) {
	url := "https://github.com/liagame/java-bot/archive/master.zip"
	name := "birko"
	customBotDir := ""
	//isErr := false
	//errMsg := "failed to fetch java bot and put into working dir"


	if os.Getenv("RUN_FUNC") == "1" {
		setupTmpConfigPaths()
		internal.FetchBot(url, name, customBotDir)

		return
	}

	exitStatus := getCmdStatus("TestCmdFetch")
	fmt.Println(exitStatus)
}

func getCmdStatus(funcName string) int {
	cmd := exec.Command(os.Args[0], "-test.run=" + funcName)
	cmd.Env = append(os.Environ(), "RUN_FUNC=1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		if status, ok := e.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 0
}

func setupTmpConfigPaths() {
	// Set tmp path to bots
	var err error
	config.PathToBots, err = ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	// Copy data to tmp path to bots
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var index int
	if runtime.GOOS == "windows" {
		index = strings.LastIndex(wd, "\\")
	} else {
		index = strings.LastIndex(wd, "/")
	}
	pathToData := filepath.Join(wd[:index], "build", "data")
	pathToTmpData := filepath.Join(config.PathToBots, "data")

	if err := copy.Dir(pathToData, pathToTmpData); err != nil {
		msg := fmt.Sprintf("failed to copy data to executable path %s", err)
		panic(msg)
	}
}
*/
