package gpio

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

//MakeSequence build an IO sequence
func MakeSequence(start State, durations ...time.Duration) Sequence {
	return Sequence{
		Start:     start,
		Durations: durations,
	}
}

//Sequence a sequence of GPIO state
type Sequence struct {
	Start     State
	Durations []time.Duration
}

//EndState final state of a sequence
func (s Sequence) EndState() State {
	if len(s.Durations)%2 == 0 {
		return s.Start
	}
	return swapState(s.Start)
}

//NextState start state for the next sequence
func (s Sequence) NextState() State {
	return swapState(s.EndState())
}

//ToSequence load a sequence from a string of the form:
//  high 100 1000 30000
func ToSequence(s string) (seq Sequence, err error) {
	r := strings.NewReader(s)
	var stateStr string
	_, err = fmt.Fscanf(r, "%s", &stateStr)
	if err != nil {
		return
	}

	seq.Start, err = ParseState(stateStr)
	if err != nil {
		return
	}

	for {
		var micros int
		_, err = fmt.Fscanf(r, "%d", &micros)
		if err == io.EOF {
			err = nil
			return
		}
		if err != nil {
			return
		}

		seq.Durations = append(seq.Durations, time.Duration(micros)*time.Microsecond)
	}
}

//String convert to a string representation
func (s Sequence) String() string {
	dStrs := make([]string, len(s.Durations))
	for i, d := range s.Durations {
		dStrs[i] = strconv.FormatInt(int64(d)/1000, 10)
	}

	return fmt.Sprintf("%s %s", States[s.Start], strings.Join(dStrs, " "))
}

func swapState(state State) State {
	if state == Low {
		return High
	}
	return Low
}

//Execute execute a sequence of output states
func Execute(pin Pin, seq Sequence) {
	pin.Output()

	state := seq.Start
	Set(pin, state)
	for _, d := range seq.Durations {
		time.Sleep(d)
		state = swapState(state)
		Set(pin, state)
	}
}

//Monitor read a sequence of input states
func Monitor(pin Pin, start State, timeout time.Duration) (Sequence, error) {
	pin.Output()

	//wait for the start state, if we're not there already
	// if _, ok := WaitChange(pin, start, timeout); !ok {
	//   return Sequence{}, fmt.Errorf("initial state %s not met", States[start])
	// }
	start = pin.Read()

	durations := []time.Duration{}

	state := start
	for {
		state = swapState(state)
		d, ok := WaitChange(pin, state, timeout)
		if !ok {
			break
		}

		durations = append(durations, d)
	}

	return Sequence{
		Start:     start,
		Durations: durations,
	}, nil
}
