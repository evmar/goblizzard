// Package replay implements decoders for replays (Starcraft, Heroes
// of the Storm etc.).
package replay

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"blizzard/blizzval"
	"blizzard/mpq"
)

type Player struct {
	Raw  blizzval.Value
	Name string "0"
	// bnet region "1"
	// starcraft race "2"
	// color "3"
	Team     int64 "5"
	Handicap int64 "6"
	// unknown "7"
	IsWinner int64 "8"
	// unknown "9"
	Character string "10"
}

type Details struct {
	Raw     blizzval.Value
	Players []*Player "0"
	Map     string    "1"
	// 2 unknown empty string
	// 3 preview image(?)
	// 4 always 1 ?
	TimeStamp int64 "5"
	UTCOffset int64 "6"
	// And a bunch more mystery fields...
}

func readReplayDetails(r *mpq.Reader) {
	fr := r.OpenFile("replay.details")
	e := blizzval.Read(bufio.NewReader(fr))

	// printTrackerEvent(e, 0)
	// fmt.Printf("\n")

	d := readDetails(e)
	log.Printf("%#v", d)
	// for _, p := range d.Players {
	// 	log.Printf("%#v", p)
	// }
}

type lr struct {
	r io.Reader
}

func (r *lr) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	log.Printf("read %d %x", n, buf[:n])
	return n, err
}

type TrackerEventType int

//go:generate stringer -type=TrackerEventType
const (
	EventPlayerStats   TrackerEventType = 0
	EventUnitBorn      TrackerEventType = 1
	EventUnitDied      TrackerEventType = 2
	EventOwnerChange   TrackerEventType = 3
	EventTypeChange    TrackerEventType = 4
	EventUnknown1      TrackerEventType = 5
	EventUnitPositions TrackerEventType = 8
	EventPlayerSetup   TrackerEventType = 9
)

type TrackerEvent struct {
	Frame int64
	Type  TrackerEventType
	Val   blizzval.Value
}

func readReplayTrackerEvents(mpqr *mpq.Reader) {
	r := mpqr.OpenFile("replay.tracker.events")

	for i := 0; i < 10000; i++ {
		te := &TrackerEvent{}
		var buf [3]byte
		_, err := r.Read(buf[:3])
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}
		if string(buf[:3]) != "\x03\x00\x09" &&
			string(buf[:3]) != "\x03\x02\x09" {
			panic(fmt.Errorf("unexpected event head %x", buf[:3]))
		}

		te.Frame = blizzval.ReadVarInt(r)

		_, err = r.Read(buf[:1])
		if err != nil {
			panic(err)
		}
		if string(buf[:1]) != "\x09" {
			panic("unexpected event trailer")
		}

		te.Type = TrackerEventType(blizzval.ReadVarInt(r))
		te.Val = blizzval.Read(r)
		log.Printf("%d %+v", i, te)
	}
}

func readBits(r *bitReader, n int) uint64 {
	bits, err := r.ReadBits(n)
	if err != nil {
		panic(err)
	}
	return bits
}

func decodeGameLoopDelta(r *bitReader) (int, error) {
	tag, err := r.ReadBits(2)
	if err != nil {
		return 0, err
	}
	switch tag {
	case 0: // m_uint6
		return int(readBits(r, 6)), nil
	case 1: // m_uint14
		return int(readBits(r, 14)), nil
	case 2: // m_uint22
		return int(readBits(r, 22)), nil
	case 3: // m_uint32
		return int(readBits(r, 32)), nil
	default:
		panic(fmt.Errorf("unknown choice tag %d", tag))
	}
}

type GameEventReader struct {
	r        *bitReader
	gameLoop int
}

func NewGameEventReader(mpqr *mpq.Reader) *GameEventReader {
	fr := mpqr.OpenFile("replay.game.events")
	r := newBitReader(bufio.NewReader(fr))

	return &GameEventReader{r: r}
}

func (r *GameEventReader) Read() Event {
	delta, err := decodeGameLoopDelta(r.r)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		panic(err)
	}
	r.gameLoop += delta
	userId := int(readBits(r.r, 5))

	eventId := int(readBits(r.r, 7))
	event := readGameEvent(r.r, eventId)
	r.r.SyncToByte()

	meta := event.Meta()
	meta.GameLoop = r.gameLoop
	meta.UserId = userId

	return event
}
