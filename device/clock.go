package device

var clock chan uint
var cycle uint

func WaitClockTick() {
	<-clock
}

func WaitClockTicks(n uint) {
	for n > 0 {
		<-clock
		n--
	}
}

func ClockTick() {
	cycle++
	clock <- cycle
}
