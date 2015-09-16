.PHONY: all bin genfiles blizzval

all: bin

bin: genfiles
	go install blizzard/hots blizzard/mpq/mpqtool

genfiles: src/blizzard/replay/replay_proto.go src/blizzard/replay/typeinfo.go

src/blizzard/replay/replay_proto.go: src/blizzard/replay/replay.go src/blizzard/blizzval/gen/main.go
	go run src/blizzard/blizzval/gen/main.go -in src/blizzard/replay/replay.go > $@
	go fmt $@

src/blizzard/replay/typeinfo.go: src/blizzard/replay/gen.py
	python src/blizzard/replay/gen.py > $@
	go fmt $@
