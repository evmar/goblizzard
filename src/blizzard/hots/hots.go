package main

import (
	"fmt"
	"log"
	"os"

	"blizzard/mpq"
	"blizzard/replay"
)

func walkGameEvents(mpqr *mpq.Reader) {
	r := replay.NewGameEventReader(mpqr)
	counts := map[string]int{}
	for event := r.Read(); event != nil; event = r.Read() {
		switch event := event.(type) {
		case *replay.BankFileEvent, *replay.BankSectionEvent, *replay.BankKeyEvent, *replay.BankSignatureEvent:
			// Per-user counter metadata.

		case *replay.UserFinishedLoadingSyncEvent:
			// Only for userid 16?

		case *replay.TriggerMouseMovedEvent:
			// Only for one player (?).
		case *replay.CameraUpdateEvent:
			// Camera position change.
		case *replay.CmdUpdateTargetPointEvent:
			// User move command.
		case *replay.CmdUpdateTargetUnitEvent:
			// User attack command.

		case *replay.UnitClickEvent:
			// ?
		case *replay.CommandManagerStateEvent:
			// ?
		case *replay.CmdEvent:
			// User command, maybe abilities?
		case *replay.TriggerSoundLengthSyncEvent, *replay.TriggerSoundOffsetEvent, *replay.TriggerSoundtrackDoneEvent:
			// Sound stuff.
		case *replay.SelectionDeltaEvent:
			// User selection change.
		case *replay.UserOptionsEvent:
			// User options settings.

		case *replay.TriggerTransmissionCompleteEvent, *replay.TriggerTransmissionOffsetEvent:
			// ?

		case *replay.TriggerDialogControlEvent:
			// ?

		case *replay.TriggerPingEvent:
			// Map ping.

		case *replay.TriggerChatMessageEvent:
			// Chat messages.

		case *replay.HeroTalentTreeSelectedEvent:
			// Selected a talent.

		case *replay.GameUserLeaveEvent:
			// End of game (?).

		default:
			log.Printf("%#v", event)
		}
		counts[fmt.Sprintf("%T", event)]++
	}

	NewHist(counts).Desc().Print()
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("need path")
	}
	path := os.Args[1]

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("%s", err)
	}
	mpqr := mpq.NewReader(f)

	// readReplayDetails(&r)
	// readReplayTrackerEvents(&r)
	walkGameEvents(mpqr)
}
