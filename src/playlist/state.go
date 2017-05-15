package playlist

var (
	stateChan = make(chan bool)
)

func stateChange() {
	select {
	case stateChan <- true:
		break
	default:
		break
	}
}

func WaitForChange() {
	<-stateChan
}
