package stats

type TestResults struct {
	Queue  []string
	InMsgs chan string
}

func InitTest() *TestResults {
	t := TestResults{}
	t.Queue = []string{}
	t.InMsgs = make(chan string)
	go t.Listen()
	return &t
}

func (t *TestResults) Listen() {
	for {
		t.AddResult(<-t.InMsgs)
	}
}

func (t *TestResults) AddResult(s string) {
	print(s)
	t.Queue = append(t.Queue, s)
}
