package x_vendor

import (
"bytes"
"github.com/matttproud/golang_protobuf_extensions/pbutil"
"github.com/palantir/stacktrace"
"github.com/pkg/errors"
"io"
"fmt"
"github.com/liagame/lia-SDK/x_vendor/curves"
)

type ReplayData struct {
	Duration           float32
	GamerWinner        Winner
	Bot1RemainingUnits int
	Bot2RemainingUnits int
	Bot1CrashTime float32
	Bot2CrashTime float32
}

type Winner string

const (
	NONE  Winner = "none"
	BOT_1 Winner = "bot_1"
	BOT_2 Winner = "bot_2"
)

func GetReplayData(replayReader io.Reader) (*ReplayData, error) {
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(replayReader)
	if err != nil {
		return nil, errors.New("failed to read replay from io.Reader")
	}

	elements, err := parseReplayFile(buf.Bytes())
	if err != nil {
		return nil, errors.New("failed to parse replay to protobuf elements")
	}

	replayData, err := parseReplayElements(elements)
	if err != nil {
		return nil, errors.New("failed to parse curve elements")
	}
	return replayData, nil
}

func parseReplayFile(replay []byte) ([]curves.Element, error) {
	r := bytes.NewReader(replay)
	var elements []curves.Element

	for {
		el := &curves.Element{}
		_, err := pbutil.ReadDelimited(r, el)
		if err == io.EOF {
			break // End of file
		} else if err != nil {
			return nil, stacktrace.Propagate(err, "failed to parse replay file")
		}
		elements = append(elements, *el)
	}

	return elements, nil
}

func parseReplayElements(elements []curves.Element) (*ReplayData, error) {
	gameOverEvent, err := GetGameOverEvent(elements)
	if err != nil {
		return nil, fmt.Errorf("failed to get gameOverTime. Err %s", err)
	}

	winner, err := getWinner(elements)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to get winner")
	}

	remainingUnitsTeam1, remainingUnitsTeam2 := getRemainingUnits(elements)

	return &ReplayData{
		Duration:           gameOverEvent.T,
		GamerWinner:        winner,
		Bot1RemainingUnits: remainingUnitsTeam1,
		Bot2RemainingUnits: remainingUnitsTeam2,
		Bot1CrashTime: gameOverEvent.Bot1CrashTime,
		Bot2CrashTime: gameOverEvent.Bot2CrashTime,
	}, nil
}

func GetGameOverEvent(elements []curves.Element) (*curves.GameOverEvent, error) {
	for _, el := range elements {
		gameOverEvent := el.GetGameOverEvent()
		if gameOverEvent != nil {
			return gameOverEvent, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func getWinner(elements []curves.Element) (Winner, error) {
	for _, el := range elements {
		goe := el.GetGameOverEvent()
		if goe != nil {
			// Extract winner
			switch goe.Winner {
			case curves.GameOverEvent_TEAM_1:
				return BOT_1, nil
			case curves.GameOverEvent_TEAM_2:
				return BOT_2, nil
			}
		}
	}
	return BOT_1, fmt.Errorf("not found")
}

func getRemainingUnits(elements []curves.Element) (int, int) {
	// Get entities ids
	var eidsTeam1 []int32
	var eidsTeam2 []int32
	for _, el := range elements {
		e := el.GetCreateEntityEvent()
		if e != nil {
			switch e.Team {
			case "TEAM_1":
				eidsTeam1 = append(eidsTeam1, e.Eid)
			case "TEAM_2":
				eidsTeam2 = append(eidsTeam2, e.Eid)
			}
		}
	}

	// Get health curves ids
	var cidsTeam1 []int32
	var cidsTeam2 []int32
	for _, el := range elements {
		c := el.GetCurve()
		if c != nil && c.Type == "HEALTH" {
			if inArray(eidsTeam1, c.Eid) {
				cidsTeam1 = append(cidsTeam1, c.Id)
			} else if inArray(eidsTeam2, c.Eid) {
				cidsTeam2 = append(cidsTeam2, c.Id)
			}
		}
	}

	// Count curve ids that reached health 0
	nDeadTeam1 := 0
	nDeadTeam2 := 0
	for _, el := range elements {
		p := el.GetBasicPoint()
		if p != nil {
			if inArray(cidsTeam1, p.Cid) && p.X == 0 {
				nDeadTeam1++
			} else if inArray(cidsTeam2, p.Cid) && p.X == 0 {
				nDeadTeam2++
			}
		}
	}

	nAliveTeam1 := len(cidsTeam1) - nDeadTeam1
	nAliveTeam2 := len(cidsTeam2) - nDeadTeam2

	return nAliveTeam1, nAliveTeam2
}

func inArray(list []int32, value int32) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

