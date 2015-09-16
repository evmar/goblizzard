package replay

import (
	"fmt"
	"io"
)

type EventMeta struct {
	GameLoop int
	UserId   int
}
type Event interface {
	Meta() *EventMeta
}

func (e *EventMeta) Meta() *EventMeta {
	return e
}

// typeinfo 0 (int)

// typeinfo 1 (int)

// typeinfo 2 (int)

// typeinfo 3 (int)

// typeinfo 4 (int)

// typeinfo 5 (int)

// typeinfo 6 (int)

// typeinfo 7 (choice)
func decodeGameLoopDeltaAuto(r *bitReader) int32 {
	switch tag := readBits(r, 2); tag {
	case 0: // m_uint6
		return int32(int8(readBits(r, 6)))
	case 1: // m_uint14
		return int32(int16(readBits(r, 14)))
	case 2: // m_uint22
		return int32(int32(readBits(r, 22)))
	case 3: // m_uint32
		return int32(int32(readBits(r, 32)))
	default:
		panic(fmt.Errorf("unknown choice tag %d", tag))
	}
}

// typeinfo 8 (struct)
type Unknown8 struct {
	UserId int8 // 2
}

func decodeUnknown8(r *bitReader) *Unknown8 {
	out := &Unknown8{}
	out.UserId = int8(readBits(r, 5))
	return out
}

// typeinfo 9 (blob)
func decodeByteString_0_8(r *bitReader) string {
	n := readBits(r, 8)
	r.SyncToByte()
	buf := make([]byte, n)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

// typeinfo 10 (int)

// typeinfo 11 (struct)
type Version struct {
	Flags     int8  // 10
	Major     int8  // 10
	Minor     int8  // 10
	Revision  int8  // 10
	Build     int32 // 6
	BaseBuild int32 // 6
}

func decodeVersion(r *bitReader) *Version {
	out := &Version{}
	out.Flags = int8(readBits(r, 8))
	out.Major = int8(readBits(r, 8))
	out.Minor = int8(readBits(r, 8))
	out.Revision = int8(readBits(r, 8))
	out.Build = int32(readBits(r, 32))
	out.BaseBuild = int32(readBits(r, 32))
	return out
}

// typeinfo 12 (int)

// typeinfo 13 (bool)

// typeinfo 14 (array)
func decodeUnknown14(r *bitReader) []int8 {
	n := int(16)
	arr := make([]int8, n)
	for i := 0; i < n; i++ {
		arr[i] = int8(readBits(r, 8))
	}
	return arr
}

// typeinfo 15 (optional)
func decodeUnknown15(r *bitReader) *[]int8 {
	if readBits(r, 1) != 0 {
		ret := decodeUnknown14(r)
		return &ret
	}
	return nil
}

// typeinfo 16 (blob)
func decodeByteString_16_0(r *bitReader) string {
	n := 16
	r.SyncToByte()
	buf := make([]byte, n)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

// typeinfo 17 (struct)
type NgdpRootKey struct {
	DataDeprecated *[]int8 // 15
	Data           string  // 16
}

func decodeNgdpRootKey(r *bitReader) *NgdpRootKey {
	out := &NgdpRootKey{}
	out.DataDeprecated = decodeUnknown15(r)
	out.Data = decodeByteString_16_0(r)
	return out
}

// typeinfo 18 (struct)
type Header struct {
	Signature        string       // 9
	Version          *Version     // 11
	Type             int8         // 12
	ElapsedGameLoops int32        // 6
	UseScaledTime    bool         // 13
	NgdpRootKey      *NgdpRootKey // 17
	DataBuildNum     int32        // 6
}

func decodeHeader(r *bitReader) *Header {
	out := &Header{}
	out.Signature = decodeByteString_0_8(r)
	out.Version = decodeVersion(r)
	out.Type = int8(readBits(r, 3))
	out.ElapsedGameLoops = int32(readBits(r, 32))
	out.UseScaledTime = readBits(r, 1) != 0
	out.NgdpRootKey = decodeNgdpRootKey(r)
	out.DataBuildNum = int32(readBits(r, 32))
	return out
}

// TODO: Unknown19 (19)

// typeinfo 20 (blob)
func decodeByteString_0_7(r *bitReader) string {
	n := readBits(r, 7)
	r.SyncToByte()
	buf := make([]byte, n)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

// typeinfo 21 (int)

// typeinfo 22 (struct)
type Toon struct {
	Region int8 // 10
	// TODO: ProgramId TODO // 19
	Realm int32  // 6
	Name  string // 20
	Id    int64  // 21
}

func decodeToon(r *bitReader) *Toon {
	out := &Toon{}
	out.Region = int8(readBits(r, 8))
	panic("TODO") // decode ProgramId
	out.Realm = int32(readBits(r, 32))
	out.Name = decodeByteString_0_7(r)
	out.Id = int64(readBits(r, 64))
	return out
}

// typeinfo 23 (struct)
type Color struct {
	A int8 // 10
	R int8 // 10
	G int8 // 10
	B int8 // 10
}

func decodeColor(r *bitReader) *Color {
	out := &Color{}
	out.A = int8(readBits(r, 8))
	out.R = int8(readBits(r, 8))
	out.G = int8(readBits(r, 8))
	out.B = int8(readBits(r, 8))
	return out
}

// typeinfo 24 (int)

// typeinfo 25 (optional)
func decodeUnknown25(r *bitReader) *int8 {
	if readBits(r, 1) != 0 {
		ret := int8(readBits(r, 8))
		return &ret
	}
	return nil
}

// typeinfo 26 (struct)
type Unknown26 struct {
	Name             string // 9
	Toon             *Toon  // 22
	Race             string // 9
	Color            *Color // 23
	Control          int8   // 10
	TeamId           int8   // 1
	Handicap         int8   // 0
	Observe          int8   // 24
	Result           int8   // 24
	WorkingSetSlotId *int8  // 25
	Hero             string // 9
}

func decodeUnknown26(r *bitReader) *Unknown26 {
	out := &Unknown26{}
	out.Name = decodeByteString_0_8(r)
	out.Toon = decodeToon(r)
	out.Race = decodeByteString_0_8(r)
	out.Color = decodeColor(r)
	out.Control = int8(readBits(r, 8))
	out.TeamId = int8(readBits(r, 4))
	out.Handicap = int8(readBits(r, 7))
	out.Observe = int8(readBits(r, 2))
	out.Result = int8(readBits(r, 2))
	out.WorkingSetSlotId = decodeUnknown25(r)
	out.Hero = decodeByteString_0_8(r)
	return out
}

// typeinfo 27 (array)
func decodeUnknown27(r *bitReader) []*Unknown26 {
	n := int(readBits(r, 5))
	arr := make([]*Unknown26, n)
	for i := 0; i < n; i++ {
		arr[i] = decodeUnknown26(r)
	}
	return arr
}

// typeinfo 28 (optional)
func decodeUnknown28(r *bitReader) *[]*Unknown26 {
	if readBits(r, 1) != 0 {
		ret := decodeUnknown27(r)
		return &ret
	}
	return nil
}

// typeinfo 29 (blob)
func decodeByteString_0_10(r *bitReader) string {
	n := readBits(r, 10)
	r.SyncToByte()
	buf := make([]byte, n)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

// typeinfo 30 (blob)
func decodeByteString_0_11(r *bitReader) string {
	n := readBits(r, 11)
	r.SyncToByte()
	buf := make([]byte, n)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

// typeinfo 31 (struct)
type Thumbnail struct {
	File string // 30
}

func decodeThumbnail(r *bitReader) *Thumbnail {
	out := &Thumbnail{}
	out.File = decodeByteString_0_11(r)
	return out
}

// typeinfo 32 (optional)
func decodeUnknown32(r *bitReader) *bool {
	if readBits(r, 1) != 0 {
		ret := readBits(r, 1) != 0
		return &ret
	}
	return nil
}

// typeinfo 33 (int)

// typeinfo 34 (blob)
func decodeByteString_0_12(r *bitReader) string {
	n := readBits(r, 12)
	r.SyncToByte()
	buf := make([]byte, n)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

// typeinfo 35 (blob)
func decodeByteString_40_0(r *bitReader) string {
	n := 40
	r.SyncToByte()
	buf := make([]byte, n)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

// typeinfo 36 (array)
func decodeUnknown36(r *bitReader) []string {
	n := int(readBits(r, 6))
	arr := make([]string, n)
	for i := 0; i < n; i++ {
		arr[i] = decodeByteString_40_0(r)
	}
	return arr
}

// typeinfo 37 (optional)
func decodeUnknown37(r *bitReader) *[]string {
	if readBits(r, 1) != 0 {
		ret := decodeUnknown36(r)
		return &ret
	}
	return nil
}

// typeinfo 38 (array)
func decodeUnknown38(r *bitReader) []string {
	n := int(readBits(r, 6))
	arr := make([]string, n)
	for i := 0; i < n; i++ {
		arr[i] = decodeByteString_0_11(r)
	}
	return arr
}

// typeinfo 39 (optional)
func decodeUnknown39(r *bitReader) *[]string {
	if readBits(r, 1) != 0 {
		ret := decodeUnknown38(r)
		return &ret
	}
	return nil
}

// typeinfo 40 (struct)
type Unknown40 struct {
	PlayerList             *[]*Unknown26 // 28
	Title                  string        // 29
	Difficulty             string        // 9
	Thumbnail              *Thumbnail    // 31
	IsBlizzardMap          bool          // 13
	RestartAsTransitionMap *bool         // 32
	TimeUTC                int64         // 33
	TimeLocalOffset        int64         // 33
	Description            string        // 34
	ImageFilePath          string        // 30
	CampaignIndex          int8          // 10
	MapFileName            string        // 30
	CacheHandles           *[]string     // 37
	MiniSave               bool          // 13
	GameSpeed              int8          // 12
	DefaultDifficulty      int8          // 3
	ModPaths               *[]string     // 39
}

func decodeUnknown40(r *bitReader) *Unknown40 {
	out := &Unknown40{}
	out.PlayerList = decodeUnknown28(r)
	out.Title = decodeByteString_0_10(r)
	out.Difficulty = decodeByteString_0_8(r)
	out.Thumbnail = decodeThumbnail(r)
	out.IsBlizzardMap = readBits(r, 1) != 0
	out.RestartAsTransitionMap = decodeUnknown32(r)
	out.TimeUTC = int64(-9223372036854775808 + int64(readBits(r, 64)))
	out.TimeLocalOffset = int64(-9223372036854775808 + int64(readBits(r, 64)))
	out.Description = decodeByteString_0_12(r)
	out.ImageFilePath = decodeByteString_0_11(r)
	out.CampaignIndex = int8(readBits(r, 8))
	out.MapFileName = decodeByteString_0_11(r)
	out.CacheHandles = decodeUnknown37(r)
	out.MiniSave = readBits(r, 1) != 0
	out.GameSpeed = int8(readBits(r, 3))
	out.DefaultDifficulty = int8(readBits(r, 6))
	out.ModPaths = decodeUnknown39(r)
	return out
}

// typeinfo 41 (optional)
func decodeUnknown41(r *bitReader) *string {
	if readBits(r, 1) != 0 {
		ret := decodeByteString_0_8(r)
		return &ret
	}
	return nil
}

// typeinfo 42 (optional)
func decodeUnknown42(r *bitReader) *string {
	if readBits(r, 1) != 0 {
		ret := decodeByteString_40_0(r)
		return &ret
	}
	return nil
}

// typeinfo 43 (optional)
func decodeUnknown43(r *bitReader) *int32 {
	if readBits(r, 1) != 0 {
		ret := int32(readBits(r, 32))
		return &ret
	}
	return nil
}

// typeinfo 44 (struct)
// name hints: set(['RacePref', 'RacePreference'])
type RacePref struct {
	Race *int8 // 25
}

func decodeRacePref(r *bitReader) *RacePref {
	out := &RacePref{}
	out.Race = decodeUnknown25(r)
	return out
}

// typeinfo 45 (struct)
type TeamPreference struct {
	Team *int8 // 25
}

func decodeTeamPreference(r *bitReader) *TeamPreference {
	out := &TeamPreference{}
	out.Team = decodeUnknown25(r)
	return out
}

// typeinfo 46 (blob)
func decodeByteString_0_9(r *bitReader) string {
	n := readBits(r, 9)
	r.SyncToByte()
	buf := make([]byte, n)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

// typeinfo 47 (struct)
type Unknown47 struct {
	Name               string          // 9
	ClanTag            *string         // 41
	ClanLogo           *string         // 42
	HighestLeague      *int8           // 25
	CombinedRaceLevels *int32          // 43
	RandomSeed         int32           // 6
	RacePreference     *RacePref       // 44
	TeamPreference     *TeamPreference // 45
	TestMap            bool            // 13
	TestAuto           bool            // 13
	Examine            bool            // 13
	CustomInterface    bool            // 13
	TestType           int32           // 6
	Observe            int8            // 24
	Hero               string          // 46
	Skin               string          // 46
	Mount              string          // 46
	ToonHandle         string          // 20
}

func decodeUnknown47(r *bitReader) *Unknown47 {
	out := &Unknown47{}
	out.Name = decodeByteString_0_8(r)
	out.ClanTag = decodeUnknown41(r)
	out.ClanLogo = decodeUnknown42(r)
	out.HighestLeague = decodeUnknown25(r)
	out.CombinedRaceLevels = decodeUnknown43(r)
	out.RandomSeed = int32(readBits(r, 32))
	out.RacePreference = decodeRacePref(r)
	out.TeamPreference = decodeTeamPreference(r)
	out.TestMap = readBits(r, 1) != 0
	out.TestAuto = readBits(r, 1) != 0
	out.Examine = readBits(r, 1) != 0
	out.CustomInterface = readBits(r, 1) != 0
	out.TestType = int32(readBits(r, 32))
	out.Observe = int8(readBits(r, 2))
	out.Hero = decodeByteString_0_9(r)
	out.Skin = decodeByteString_0_9(r)
	out.Mount = decodeByteString_0_9(r)
	out.ToonHandle = decodeByteString_0_7(r)
	return out
}

// typeinfo 48 (array)
func decodeUnknown48(r *bitReader) []*Unknown47 {
	n := int(readBits(r, 5))
	arr := make([]*Unknown47, n)
	for i := 0; i < n; i++ {
		arr[i] = decodeUnknown47(r)
	}
	return arr
}

// typeinfo 49 (struct)
type GameOptions struct {
	LockTeams             bool  // 13
	TeamsTogether         bool  // 13
	AdvancedSharedControl bool  // 13
	RandomRaces           bool  // 13
	BattleNet             bool  // 13
	Amm                   bool  // 13
	Ranked                bool  // 13
	Competitive           bool  // 13
	Practice              bool  // 13
	Cooperative           bool  // 13
	NoVictoryOrDefeat     bool  // 13
	HeroDuplicatesAllowed bool  // 13
	Fog                   int8  // 24
	Observers             int8  // 24
	UserDifficulty        int8  // 24
	ClientDebugFlags      int64 // 21
}

func decodeGameOptions(r *bitReader) *GameOptions {
	out := &GameOptions{}
	out.LockTeams = readBits(r, 1) != 0
	out.TeamsTogether = readBits(r, 1) != 0
	out.AdvancedSharedControl = readBits(r, 1) != 0
	out.RandomRaces = readBits(r, 1) != 0
	out.BattleNet = readBits(r, 1) != 0
	out.Amm = readBits(r, 1) != 0
	out.Ranked = readBits(r, 1) != 0
	out.Competitive = readBits(r, 1) != 0
	out.Practice = readBits(r, 1) != 0
	out.Cooperative = readBits(r, 1) != 0
	out.NoVictoryOrDefeat = readBits(r, 1) != 0
	out.HeroDuplicatesAllowed = readBits(r, 1) != 0
	out.Fog = int8(readBits(r, 2))
	out.Observers = int8(readBits(r, 2))
	out.UserDifficulty = int8(readBits(r, 2))
	out.ClientDebugFlags = int64(readBits(r, 64))
	return out
}

// typeinfo 50 (int)

// typeinfo 51 (int)

// typeinfo 52 (bitarray)
func decodeUnknown52(r *bitReader) uint64 {
	panic("TODO")
}

// typeinfo 53 (bitarray)
func decodeUnknown53(r *bitReader) uint64 {
	panic("TODO")
}

// typeinfo 54 (bitarray)
func decodeUnknown54(r *bitReader) uint64 {
	panic("TODO")
}

// typeinfo 55 (bitarray)
func decodeUnknown55(r *bitReader) uint64 {
	panic("TODO")
}

// typeinfo 56 (struct)
type Unknown56 struct {
	AllowedColors       uint64 // 52
	AllowedRaces        uint64 // 53
	AllowedDifficulty   uint64 // 52
	AllowedControls     uint64 // 53
	AllowedObserveTypes uint64 // 54
	AllowedAIBuilds     uint64 // 55
}

func decodeUnknown56(r *bitReader) *Unknown56 {
	out := &Unknown56{}
	out.AllowedColors = decodeUnknown52(r)
	out.AllowedRaces = decodeUnknown53(r)
	out.AllowedDifficulty = decodeUnknown52(r)
	out.AllowedControls = decodeUnknown53(r)
	out.AllowedObserveTypes = decodeUnknown54(r)
	out.AllowedAIBuilds = decodeUnknown55(r)
	return out
}

// typeinfo 57 (array)
func decodeUnknown57(r *bitReader) []*Unknown56 {
	n := int(readBits(r, 5))
	arr := make([]*Unknown56, n)
	for i := 0; i < n; i++ {
		arr[i] = decodeUnknown56(r)
	}
	return arr
}

// typeinfo 58 (struct)
type GameDescription struct {
	RandomValue         int32        // 6
	GameCacheName       string       // 29
	GameOptions         *GameOptions // 49
	GameSpeed           int8         // 12
	GameType            int8         // 12
	MaxUsers            int8         // 2
	MaxObservers        int8         // 2
	MaxPlayers          int8         // 2
	MaxTeams            int64        // 50
	MaxColors           int8         // 3
	MaxRaces            int64        // 51
	MaxControls         int8         // 10
	MapSizeX            int8         // 10
	MapSizeY            int8         // 10
	MapFileSyncChecksum int32        // 6
	MapFileName         string       // 30
	MapAuthorName       string       // 9
	ModFileSyncChecksum int32        // 6
	SlotDescriptions    []*Unknown56 // 57
	DefaultDifficulty   int8         // 3
	DefaultAIBuild      int8         // 0
	CacheHandles        []string     // 36
	HasExtensionMod     bool         // 13
	IsBlizzardMap       bool         // 13
	IsPremadeFFA        bool         // 13
	IsCoopMode          bool         // 13
}

func decodeGameDescription(r *bitReader) *GameDescription {
	out := &GameDescription{}
	out.RandomValue = int32(readBits(r, 32))
	out.GameCacheName = decodeByteString_0_10(r)
	out.GameOptions = decodeGameOptions(r)
	out.GameSpeed = int8(readBits(r, 3))
	out.GameType = int8(readBits(r, 3))
	out.MaxUsers = int8(readBits(r, 5))
	out.MaxObservers = int8(readBits(r, 5))
	out.MaxPlayers = int8(readBits(r, 5))
	out.MaxTeams = int64(1 + int64(readBits(r, 4)))
	out.MaxColors = int8(readBits(r, 6))
	out.MaxRaces = int64(1 + int64(readBits(r, 8)))
	out.MaxControls = int8(readBits(r, 8))
	out.MapSizeX = int8(readBits(r, 8))
	out.MapSizeY = int8(readBits(r, 8))
	out.MapFileSyncChecksum = int32(readBits(r, 32))
	out.MapFileName = decodeByteString_0_11(r)
	out.MapAuthorName = decodeByteString_0_8(r)
	out.ModFileSyncChecksum = int32(readBits(r, 32))
	out.SlotDescriptions = decodeUnknown57(r)
	out.DefaultDifficulty = int8(readBits(r, 6))
	out.DefaultAIBuild = int8(readBits(r, 7))
	out.CacheHandles = decodeUnknown36(r)
	out.HasExtensionMod = readBits(r, 1) != 0
	out.IsBlizzardMap = readBits(r, 1) != 0
	out.IsPremadeFFA = readBits(r, 1) != 0
	out.IsCoopMode = readBits(r, 1) != 0
	return out
}

// typeinfo 59 (optional)
func decodeUnknown59(r *bitReader) *int8 {
	if readBits(r, 1) != 0 {
		ret := int8(readBits(r, 4))
		return &ret
	}
	return nil
}

// typeinfo 60 (optional)
func decodeUnknown60(r *bitReader) *int8 {
	if readBits(r, 1) != 0 {
		ret := int8(readBits(r, 5))
		return &ret
	}
	return nil
}

// typeinfo 61 (struct)
type ColorPref struct {
	Color *int8 // 60
}

func decodeColorPref(r *bitReader) *ColorPref {
	out := &ColorPref{}
	out.Color = decodeUnknown60(r)
	return out
}

// typeinfo 62 (array)
func decodeUnknown62(r *bitReader) []string {
	n := int(readBits(r, 4))
	arr := make([]string, n)
	for i := 0; i < n; i++ {
		arr[i] = decodeByteString_0_9(r)
	}
	return arr
}

// typeinfo 63 (array)
func decodeUnknown63(r *bitReader) []int32 {
	n := int(readBits(r, 17))
	arr := make([]int32, n)
	for i := 0; i < n; i++ {
		arr[i] = int32(readBits(r, 32))
	}
	return arr
}

// typeinfo 64 (array)
func decodeUnknown64(r *bitReader) []int32 {
	n := int(readBits(r, 9))
	arr := make([]int32, n)
	for i := 0; i < n; i++ {
		arr[i] = int32(readBits(r, 32))
	}
	return arr
}

// typeinfo 65 (struct)
type Unknown65 struct {
	Control            int8       // 10
	UserId             *int8      // 59
	TeamId             int8       // 1
	ColorPref          *ColorPref // 61
	RacePref           *RacePref  // 44
	Difficulty         int8       // 3
	AiBuild            int8       // 0
	Handicap           int8       // 0
	Observe            int8       // 24
	LogoIndex          int32      // 6
	Hero               string     // 46
	Skin               string     // 46
	Mount              string     // 46
	Artifacts          []string   // 62
	WorkingSetSlotId   *int8      // 25
	Rewards            []int32    // 63
	ToonHandle         string     // 20
	Licenses           []int32    // 64
	TandemLeaderUserId *int8      // 59
	Commander          string     // 46
}

func decodeUnknown65(r *bitReader) *Unknown65 {
	out := &Unknown65{}
	out.Control = int8(readBits(r, 8))
	out.UserId = decodeUnknown59(r)
	out.TeamId = int8(readBits(r, 4))
	out.ColorPref = decodeColorPref(r)
	out.RacePref = decodeRacePref(r)
	out.Difficulty = int8(readBits(r, 6))
	out.AiBuild = int8(readBits(r, 7))
	out.Handicap = int8(readBits(r, 7))
	out.Observe = int8(readBits(r, 2))
	out.LogoIndex = int32(readBits(r, 32))
	out.Hero = decodeByteString_0_9(r)
	out.Skin = decodeByteString_0_9(r)
	out.Mount = decodeByteString_0_9(r)
	out.Artifacts = decodeUnknown62(r)
	out.WorkingSetSlotId = decodeUnknown25(r)
	out.Rewards = decodeUnknown63(r)
	out.ToonHandle = decodeByteString_0_7(r)
	out.Licenses = decodeUnknown64(r)
	out.TandemLeaderUserId = decodeUnknown59(r)
	out.Commander = decodeByteString_0_9(r)
	return out
}

// typeinfo 66 (array)
func decodeUnknown66(r *bitReader) []*Unknown65 {
	n := int(readBits(r, 5))
	arr := make([]*Unknown65, n)
	for i := 0; i < n; i++ {
		arr[i] = decodeUnknown65(r)
	}
	return arr
}

// typeinfo 67 (struct)
type LobbyState struct {
	Phase             int8         // 12
	MaxUsers          int8         // 2
	MaxObservers      int8         // 2
	Slots             []*Unknown65 // 66
	RandomSeed        int32        // 6
	HostUserId        *int8        // 59
	IsSinglePlayer    bool         // 13
	GameDuration      int32        // 6
	DefaultDifficulty int8         // 3
	DefaultAIBuild    int8         // 0
}

func decodeLobbyState(r *bitReader) *LobbyState {
	out := &LobbyState{}
	out.Phase = int8(readBits(r, 3))
	out.MaxUsers = int8(readBits(r, 5))
	out.MaxObservers = int8(readBits(r, 5))
	out.Slots = decodeUnknown66(r)
	out.RandomSeed = int32(readBits(r, 32))
	out.HostUserId = decodeUnknown59(r)
	out.IsSinglePlayer = readBits(r, 1) != 0
	out.GameDuration = int32(readBits(r, 32))
	out.DefaultDifficulty = int8(readBits(r, 6))
	out.DefaultAIBuild = int8(readBits(r, 7))
	return out
}

// typeinfo 68 (struct)
type SyncLobbyState struct {
	UserInitialData []*Unknown47     // 48
	GameDescription *GameDescription // 58
	LobbyState      *LobbyState      // 67
}

func decodeSyncLobbyState(r *bitReader) *SyncLobbyState {
	out := &SyncLobbyState{}
	out.UserInitialData = decodeUnknown48(r)
	out.GameDescription = decodeGameDescription(r)
	out.LobbyState = decodeLobbyState(r)
	return out
}

// typeinfo 69 (struct)
type Unknown69 struct {
	SyncLobbyState *SyncLobbyState // 68
}

func decodeUnknown69(r *bitReader) *Unknown69 {
	out := &Unknown69{}
	out.SyncLobbyState = decodeSyncLobbyState(r)
	return out
}

// typeinfo 70 (struct)
type BankFileEvent struct {
	EventMeta
	Name string // 20
}

func decodeBankFileEvent(r *bitReader) *BankFileEvent {
	out := &BankFileEvent{}
	out.Name = decodeByteString_0_7(r)
	return out
}

// typeinfo 71 (blob)
func decodeByteString_0_6(r *bitReader) string {
	n := readBits(r, 6)
	r.SyncToByte()
	buf := make([]byte, n)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

// typeinfo 72 (struct)
type BankSectionEvent struct {
	EventMeta
	Name string // 71
}

func decodeBankSectionEvent(r *bitReader) *BankSectionEvent {
	out := &BankSectionEvent{}
	out.Name = decodeByteString_0_6(r)
	return out
}

// typeinfo 73 (struct)
type BankKeyEvent struct {
	EventMeta
	Name string // 71
	Type int32  // 6
	Data string // 20
}

func decodeBankKeyEvent(r *bitReader) *BankKeyEvent {
	out := &BankKeyEvent{}
	out.Name = decodeByteString_0_6(r)
	out.Type = int32(readBits(r, 32))
	out.Data = decodeByteString_0_7(r)
	return out
}

// typeinfo 74 (struct)
type BankValueEvent struct {
	EventMeta
	Type int32  // 6
	Name string // 71
	Data string // 34
}

func decodeBankValueEvent(r *bitReader) *BankValueEvent {
	out := &BankValueEvent{}
	out.Type = int32(readBits(r, 32))
	out.Name = decodeByteString_0_6(r)
	out.Data = decodeByteString_0_12(r)
	return out
}

// typeinfo 75 (array)
func decodeUnknown75(r *bitReader) []int8 {
	n := int(readBits(r, 5))
	arr := make([]int8, n)
	for i := 0; i < n; i++ {
		arr[i] = int8(readBits(r, 8))
	}
	return arr
}

// typeinfo 76 (struct)
type BankSignatureEvent struct {
	EventMeta
	Signature  []int8 // 75
	ToonHandle string // 20
}

func decodeBankSignatureEvent(r *bitReader) *BankSignatureEvent {
	out := &BankSignatureEvent{}
	out.Signature = decodeUnknown75(r)
	out.ToonHandle = decodeByteString_0_7(r)
	return out
}

// typeinfo 77 (struct)
type UserOptionsEvent struct {
	EventMeta
	GameFullyDownloaded      bool   // 13
	DevelopmentCheatsEnabled bool   // 13
	TestCheatsEnabled        bool   // 13
	MultiplayerCheatsEnabled bool   // 13
	SyncChecksummingEnabled  bool   // 13
	IsMapToMapTransition     bool   // 13
	StartingRally            bool   // 13
	DebugPauseEnabled        bool   // 13
	UseGalaxyAsserts         bool   // 13
	PlatformMac              bool   // 13
	CameraFollow             bool   // 13
	BaseBuildNum             int32  // 6
	BuildNum                 int32  // 6
	VersionFlags             int32  // 6
	HotkeyProfile            string // 46
}

func decodeUserOptionsEvent(r *bitReader) *UserOptionsEvent {
	out := &UserOptionsEvent{}
	out.GameFullyDownloaded = readBits(r, 1) != 0
	out.DevelopmentCheatsEnabled = readBits(r, 1) != 0
	out.TestCheatsEnabled = readBits(r, 1) != 0
	out.MultiplayerCheatsEnabled = readBits(r, 1) != 0
	out.SyncChecksummingEnabled = readBits(r, 1) != 0
	out.IsMapToMapTransition = readBits(r, 1) != 0
	out.StartingRally = readBits(r, 1) != 0
	out.DebugPauseEnabled = readBits(r, 1) != 0
	out.UseGalaxyAsserts = readBits(r, 1) != 0
	out.PlatformMac = readBits(r, 1) != 0
	out.CameraFollow = readBits(r, 1) != 0
	out.BaseBuildNum = int32(readBits(r, 32))
	out.BuildNum = int32(readBits(r, 32))
	out.VersionFlags = int32(readBits(r, 32))
	out.HotkeyProfile = decodeByteString_0_9(r)
	return out
}

// typeinfo 78 (struct)
// names: set(['ServerPingMessage', 'TriggerPlanetPanelDeathCompleteEvent', 'SaveGameDoneEvent', 'TriggerMercenaryPanelPurchaseEvent', 'TriggerVictoryPanelExitEvent', 'TriggerProfilerLoggingFinishedEvent', 'TriggerGameCreditsFinishedEvent', 'TriggerMovieFinishedEvent', 'UserFinishedLoadingSyncEvent', 'TriggerMovieStartedEvent', 'TriggerPlanetPanelReplayEvent', 'LoadGameDoneEvent', 'TriggerPlanetPanelCanceledEvent', 'TriggerBattleReportPanelExitEvent', 'TriggerResearchPanelExitEvent', 'TriggerResearchPanelPurchaseEvent', 'TriggerAbortMissionEvent', 'TriggerMercenaryPanelExitEvent', 'TriggerSkippedEvent', 'TriggerPlanetPanelBirthCompleteEvent', 'TriggerPurchaseExitEvent'])
type ServerPingMessage struct {
	EventMeta
}
type TriggerPlanetPanelDeathCompleteEvent struct {
	EventMeta
}
type SaveGameDoneEvent struct {
	EventMeta
}
type TriggerMercenaryPanelPurchaseEvent struct {
	EventMeta
}
type TriggerVictoryPanelExitEvent struct {
	EventMeta
}
type TriggerProfilerLoggingFinishedEvent struct {
	EventMeta
}
type TriggerGameCreditsFinishedEvent struct {
	EventMeta
}
type TriggerMovieFinishedEvent struct {
	EventMeta
}
type UserFinishedLoadingSyncEvent struct {
	EventMeta
}
type TriggerMovieStartedEvent struct {
	EventMeta
}
type TriggerPlanetPanelReplayEvent struct {
	EventMeta
}
type LoadGameDoneEvent struct {
	EventMeta
}
type TriggerPlanetPanelCanceledEvent struct {
	EventMeta
}
type TriggerBattleReportPanelExitEvent struct {
	EventMeta
}
type TriggerResearchPanelExitEvent struct {
	EventMeta
}
type TriggerResearchPanelPurchaseEvent struct {
	EventMeta
}
type TriggerAbortMissionEvent struct {
	EventMeta
}
type TriggerMercenaryPanelExitEvent struct {
	EventMeta
}
type TriggerSkippedEvent struct {
	EventMeta
}
type TriggerPlanetPanelBirthCompleteEvent struct {
	EventMeta
}
type TriggerPurchaseExitEvent struct {
	EventMeta
}

// typeinfo 79 (int)

// typeinfo 80 (struct)
type CameraTarget struct {
	X int16 // 79
	Y int16 // 79
}

func decodeCameraTarget(r *bitReader) *CameraTarget {
	out := &CameraTarget{}
	out.X = int16(readBits(r, 16))
	out.Y = int16(readBits(r, 16))
	return out
}

// typeinfo 81 (struct)
type CameraSaveEvent struct {
	EventMeta
	Which  int8          // 12
	Target *CameraTarget // 80
}

func decodeCameraSaveEvent(r *bitReader) *CameraSaveEvent {
	out := &CameraSaveEvent{}
	out.Which = int8(readBits(r, 3))
	out.Target = decodeCameraTarget(r)
	return out
}

// typeinfo 82 (struct)
type SaveGameEvent struct {
	EventMeta
	FileName    string // 30
	Automatic   bool   // 13
	Overwrite   bool   // 13
	Name        string // 9
	Description string // 29
}

func decodeSaveGameEvent(r *bitReader) *SaveGameEvent {
	out := &SaveGameEvent{}
	out.FileName = decodeByteString_0_11(r)
	out.Automatic = readBits(r, 1) != 0
	out.Overwrite = readBits(r, 1) != 0
	out.Name = decodeByteString_0_8(r)
	out.Description = decodeByteString_0_10(r)
	return out
}

// typeinfo 83 (struct)
type CommandManagerResetEvent struct {
	EventMeta
	Sequence int32 // 6
}

func decodeCommandManagerResetEvent(r *bitReader) *CommandManagerResetEvent {
	out := &CommandManagerResetEvent{}
	out.Sequence = int32(readBits(r, 32))
	return out
}

// typeinfo 84 (int)

// typeinfo 85 (struct)
type Point struct {
	X int64 // 84
	Y int64 // 84
}

func decodePoint(r *bitReader) *Point {
	out := &Point{}
	out.X = int64(-2147483648 + int64(readBits(r, 32)))
	out.Y = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 86 (struct)
type Data struct {
	Point     *Point // 85
	Time      int64  // 84
	Verb      string // 29
	Arguments string // 29
}

func decodeData(r *bitReader) *Data {
	out := &Data{}
	out.Point = decodePoint(r)
	out.Time = int64(-2147483648 + int64(readBits(r, 32)))
	out.Verb = decodeByteString_0_10(r)
	out.Arguments = decodeByteString_0_10(r)
	return out
}

// typeinfo 87 (struct)
type GameCheatEvent struct {
	EventMeta
	Data *Data // 86
}

func decodeGameCheatEvent(r *bitReader) *GameCheatEvent {
	out := &GameCheatEvent{}
	out.Data = decodeData(r)
	return out
}

// typeinfo 88 (int)

// typeinfo 89 (struct)
type Unknown89 struct {
	AbilLink     int16 // 79
	AbilCmdIndex int8  // 2
	AbilCmdData  *int8 // 25
}

func decodeUnknown89(r *bitReader) *Unknown89 {
	out := &Unknown89{}
	out.AbilLink = int16(readBits(r, 16))
	out.AbilCmdIndex = int8(readBits(r, 5))
	out.AbilCmdData = decodeUnknown25(r)
	return out
}

// typeinfo 90 (optional)
func decodeUnknown90(r *bitReader) **Unknown89 {
	if readBits(r, 1) != 0 {
		ret := decodeUnknown89(r)
		return &ret
	}
	return nil
}

// typeinfo 91 (null)

// typeinfo 92 (int)

// typeinfo 93 (struct)
// name hints: set(['PosWorld', 'SnapshotPoint', 'Target'])
type PosWorld struct {
	X int32 // 92
	Y int32 // 92
	Z int64 // 84
}

func decodePosWorld(r *bitReader) *PosWorld {
	out := &PosWorld{}
	out.X = int32(readBits(r, 20))
	out.Y = int32(readBits(r, 20))
	out.Z = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 94 (struct)
type Target struct {
	TargetUnitFlags         int16     // 79
	Timer                   int8      // 10
	Tag                     int32     // 6
	SnapshotUnitLink        int16     // 79
	SnapshotControlPlayerId *int8     // 59
	SnapshotUpkeepPlayerId  *int8     // 59
	SnapshotPoint           *PosWorld // 93
}

func decodeTarget(r *bitReader) *Target {
	out := &Target{}
	out.TargetUnitFlags = int16(readBits(r, 16))
	out.Timer = int8(readBits(r, 8))
	out.Tag = int32(readBits(r, 32))
	out.SnapshotUnitLink = int16(readBits(r, 16))
	out.SnapshotControlPlayerId = decodeUnknown59(r)
	out.SnapshotUpkeepPlayerId = decodeUnknown59(r)
	out.SnapshotPoint = decodePosWorld(r)
	return out
}

// typeinfo 95 (choice)
func decodeUnknown95(r *bitReader) interface{} {
	switch tag := readBits(r, 2); tag {
	case 0: // None
		return nil
	case 1: // TargetPoint
		return decodePosWorld(r)
	case 2: // TargetUnit
		return decodeTarget(r)
	case 3: // Data
		return int32(readBits(r, 32))
	default:
		panic(fmt.Errorf("unknown choice tag %d", tag))
	}
}

// typeinfo 96 (int)

// typeinfo 97 (struct)
type CmdEvent struct {
	EventMeta
	CmdFlags  int32       // 88
	Abil      **Unknown89 // 90
	Data      interface{} // 95
	Sequence  int64       // 96
	OtherUnit *int32      // 43
	UnitGroup *int32      // 43
}

func decodeCmdEvent(r *bitReader) *CmdEvent {
	out := &CmdEvent{}
	out.CmdFlags = int32(readBits(r, 23))
	out.Abil = decodeUnknown90(r)
	out.Data = decodeUnknown95(r)
	out.Sequence = int64(1 + int64(readBits(r, 32)))
	out.OtherUnit = decodeUnknown43(r)
	out.UnitGroup = decodeUnknown43(r)
	return out
}

// typeinfo 98 (int)

// typeinfo 99 (bitarray)
func decodeUnknown99(r *bitReader) uint64 {
	panic("TODO")
}

// typeinfo 100 (array)
func decodeUnknown100(r *bitReader) []int16 {
	n := int(readBits(r, 9))
	arr := make([]int16, n)
	for i := 0; i < n; i++ {
		arr[i] = int16(readBits(r, 9))
	}
	return arr
}

// typeinfo 101 (choice)
func decodeUnknown101(r *bitReader) interface{} {
	switch tag := readBits(r, 2); tag {
	case 0: // None
		return nil
	case 1: // Mask
		return decodeUnknown99(r)
	case 2: // OneIndices
		return decodeUnknown100(r)
	case 3: // ZeroIndices
		return decodeUnknown100(r)
	default:
		panic(fmt.Errorf("unknown choice tag %d", tag))
	}
}

// typeinfo 102 (struct)
type Unknown102 struct {
	UnitLink              int16 // 79
	SubgroupPriority      int8  // 10
	IntraSubgroupPriority int8  // 10
	Count                 int16 // 98
}

func decodeUnknown102(r *bitReader) *Unknown102 {
	out := &Unknown102{}
	out.UnitLink = int16(readBits(r, 16))
	out.SubgroupPriority = int8(readBits(r, 8))
	out.IntraSubgroupPriority = int8(readBits(r, 8))
	out.Count = int16(readBits(r, 9))
	return out
}

// typeinfo 103 (array)
func decodeUnknown103(r *bitReader) []*Unknown102 {
	n := int(readBits(r, 9))
	arr := make([]*Unknown102, n)
	for i := 0; i < n; i++ {
		arr[i] = decodeUnknown102(r)
	}
	return arr
}

// typeinfo 104 (struct)
type Delta struct {
	SubgroupIndex int16         // 98
	RemoveMask    interface{}   // 101
	AddSubgroups  []*Unknown102 // 103
	AddUnitTags   []int32       // 64
}

func decodeDelta(r *bitReader) *Delta {
	out := &Delta{}
	out.SubgroupIndex = int16(readBits(r, 9))
	out.RemoveMask = decodeUnknown101(r)
	out.AddSubgroups = decodeUnknown103(r)
	out.AddUnitTags = decodeUnknown64(r)
	return out
}

// typeinfo 105 (struct)
type SelectionDeltaEvent struct {
	EventMeta
	ControlGroupId int8   // 1
	Delta          *Delta // 104
}

func decodeSelectionDeltaEvent(r *bitReader) *SelectionDeltaEvent {
	out := &SelectionDeltaEvent{}
	out.ControlGroupId = int8(readBits(r, 4))
	out.Delta = decodeDelta(r)
	return out
}

// typeinfo 106 (struct)
type ControlGroupUpdateEvent struct {
	EventMeta
	ControlGroupIndex  int8        // 1
	ControlGroupUpdate int8        // 24
	Mask               interface{} // 101
}

func decodeControlGroupUpdateEvent(r *bitReader) *ControlGroupUpdateEvent {
	out := &ControlGroupUpdateEvent{}
	out.ControlGroupIndex = int8(readBits(r, 4))
	out.ControlGroupUpdate = int8(readBits(r, 2))
	out.Mask = decodeUnknown101(r)
	return out
}

// typeinfo 107 (struct)
type SelectionSyncData struct {
	Count                   int16 // 98
	SubgroupCount           int16 // 98
	ActiveSubgroupIndex     int16 // 98
	UnitTagsChecksum        int32 // 6
	SubgroupIndicesChecksum int32 // 6
	SubgroupsChecksum       int32 // 6
}

func decodeSelectionSyncData(r *bitReader) *SelectionSyncData {
	out := &SelectionSyncData{}
	out.Count = int16(readBits(r, 9))
	out.SubgroupCount = int16(readBits(r, 9))
	out.ActiveSubgroupIndex = int16(readBits(r, 9))
	out.UnitTagsChecksum = int32(readBits(r, 32))
	out.SubgroupIndicesChecksum = int32(readBits(r, 32))
	out.SubgroupsChecksum = int32(readBits(r, 32))
	return out
}

// typeinfo 108 (struct)
type SelectionSyncCheckEvent struct {
	EventMeta
	ControlGroupId    int8               // 1
	SelectionSyncData *SelectionSyncData // 107
}

func decodeSelectionSyncCheckEvent(r *bitReader) *SelectionSyncCheckEvent {
	out := &SelectionSyncCheckEvent{}
	out.ControlGroupId = int8(readBits(r, 4))
	out.SelectionSyncData = decodeSelectionSyncData(r)
	return out
}

// typeinfo 109 (array)
func decodeUnknown109(r *bitReader) []int64 {
	n := int(readBits(r, 3))
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(-2147483648 + int64(readBits(r, 32)))
	}
	return arr
}

// typeinfo 110 (struct)
type ResourceTradeEvent struct {
	EventMeta
	RecipientId int8    // 1
	Resources   []int64 // 109
}

func decodeResourceTradeEvent(r *bitReader) *ResourceTradeEvent {
	out := &ResourceTradeEvent{}
	out.RecipientId = int8(readBits(r, 4))
	out.Resources = decodeUnknown109(r)
	return out
}

// typeinfo 111 (struct)
type TriggerChatMessageEvent struct {
	EventMeta
	ChatMessage string // 29
}

func decodeTriggerChatMessageEvent(r *bitReader) *TriggerChatMessageEvent {
	out := &TriggerChatMessageEvent{}
	out.ChatMessage = decodeByteString_0_10(r)
	return out
}

// typeinfo 112 (int)

// typeinfo 113 (struct)
type TargetPoint struct {
	X int64 // 84
	Y int64 // 84
	Z int64 // 84
}

func decodeTargetPoint(r *bitReader) *TargetPoint {
	out := &TargetPoint{}
	out.X = int64(-2147483648 + int64(readBits(r, 32)))
	out.Y = int64(-2147483648 + int64(readBits(r, 32)))
	out.Z = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 114 (struct)
type AICommunicateEvent struct {
	EventMeta
	Beacon                            int64        // 112
	Ally                              int64        // 112
	Flags                             int64        // 112
	Build                             int64        // 112
	TargetUnitTag                     int32        // 6
	TargetUnitSnapshotUnitLink        int16        // 79
	TargetUnitSnapshotUpkeepPlayerId  int64        // 112
	TargetUnitSnapshotControlPlayerId int64        // 112
	TargetPoint                       *TargetPoint // 113
}

func decodeAICommunicateEvent(r *bitReader) *AICommunicateEvent {
	out := &AICommunicateEvent{}
	out.Beacon = int64(-128 + int64(readBits(r, 8)))
	out.Ally = int64(-128 + int64(readBits(r, 8)))
	out.Flags = int64(-128 + int64(readBits(r, 8)))
	out.Build = int64(-128 + int64(readBits(r, 8)))
	out.TargetUnitTag = int32(readBits(r, 32))
	out.TargetUnitSnapshotUnitLink = int16(readBits(r, 16))
	out.TargetUnitSnapshotUpkeepPlayerId = int64(-128 + int64(readBits(r, 8)))
	out.TargetUnitSnapshotControlPlayerId = int64(-128 + int64(readBits(r, 8)))
	out.TargetPoint = decodeTargetPoint(r)
	return out
}

// typeinfo 115 (struct)
type SetAbsoluteGameSpeedEvent struct {
	EventMeta
	Speed int8 // 12
}

func decodeSetAbsoluteGameSpeedEvent(r *bitReader) *SetAbsoluteGameSpeedEvent {
	out := &SetAbsoluteGameSpeedEvent{}
	out.Speed = int8(readBits(r, 3))
	return out
}

// typeinfo 116 (struct)
type AddAbsoluteGameSpeedEvent struct {
	EventMeta
	Delta int64 // 112
}

func decodeAddAbsoluteGameSpeedEvent(r *bitReader) *AddAbsoluteGameSpeedEvent {
	out := &AddAbsoluteGameSpeedEvent{}
	out.Delta = int64(-128 + int64(readBits(r, 8)))
	return out
}

// typeinfo 117 (struct)
type TriggerPingEvent struct {
	EventMeta
	Point         *Point // 85
	Unit          int32  // 6
	PingedMinimap bool   // 13
	Option        int64  // 84
}

func decodeTriggerPingEvent(r *bitReader) *TriggerPingEvent {
	out := &TriggerPingEvent{}
	out.Point = decodePoint(r)
	out.Unit = int32(readBits(r, 32))
	out.PingedMinimap = readBits(r, 1) != 0
	out.Option = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 118 (struct)
type BroadcastCheatEvent struct {
	EventMeta
	Verb      string // 29
	Arguments string // 29
}

func decodeBroadcastCheatEvent(r *bitReader) *BroadcastCheatEvent {
	out := &BroadcastCheatEvent{}
	out.Verb = decodeByteString_0_10(r)
	out.Arguments = decodeByteString_0_10(r)
	return out
}

// typeinfo 119 (struct)
type AllianceEvent struct {
	EventMeta
	Alliance int32 // 6
	Control  int32 // 6
}

func decodeAllianceEvent(r *bitReader) *AllianceEvent {
	out := &AllianceEvent{}
	out.Alliance = int32(readBits(r, 32))
	out.Control = int32(readBits(r, 32))
	return out
}

// typeinfo 120 (struct)
type UnitClickEvent struct {
	EventMeta
	UnitTag int32 // 6
}

func decodeUnitClickEvent(r *bitReader) *UnitClickEvent {
	out := &UnitClickEvent{}
	out.UnitTag = int32(readBits(r, 32))
	return out
}

// typeinfo 121 (struct)
type UnitHighlightEvent struct {
	EventMeta
	UnitTag int32 // 6
	Flags   int8  // 10
}

func decodeUnitHighlightEvent(r *bitReader) *UnitHighlightEvent {
	out := &UnitHighlightEvent{}
	out.UnitTag = int32(readBits(r, 32))
	out.Flags = int8(readBits(r, 8))
	return out
}

// typeinfo 122 (struct)
type TriggerReplySelectedEvent struct {
	EventMeta
	ConversationId int64 // 84
	ReplyId        int64 // 84
}

func decodeTriggerReplySelectedEvent(r *bitReader) *TriggerReplySelectedEvent {
	out := &TriggerReplySelectedEvent{}
	out.ConversationId = int64(-2147483648 + int64(readBits(r, 32)))
	out.ReplyId = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 123 (optional)
func decodeUnknown123(r *bitReader) *string {
	if readBits(r, 1) != 0 {
		ret := decodeByteString_0_7(r)
		return &ret
	}
	return nil
}

// typeinfo 124 (struct)
type Unknown124 struct {
	GameUserId int8    // 1
	Observe    int8    // 24
	Name       string  // 9
	ToonHandle *string // 123
	ClanTag    *string // 41
	ClanLogo   *string // 42
}

func decodeUnknown124(r *bitReader) *Unknown124 {
	out := &Unknown124{}
	out.GameUserId = int8(readBits(r, 4))
	out.Observe = int8(readBits(r, 2))
	out.Name = decodeByteString_0_8(r)
	out.ToonHandle = decodeUnknown123(r)
	out.ClanTag = decodeUnknown41(r)
	out.ClanLogo = decodeUnknown42(r)
	return out
}

// typeinfo 125 (array)
func decodeUnknown125(r *bitReader) []*Unknown124 {
	n := int(readBits(r, 5))
	arr := make([]*Unknown124, n)
	for i := 0; i < n; i++ {
		arr[i] = decodeUnknown124(r)
	}
	return arr
}

// typeinfo 126 (int)

// typeinfo 127 (struct)
type HijackReplayGameEvent struct {
	EventMeta
	UserInfos []*Unknown124 // 125
	Method    int8          // 126
}

func decodeHijackReplayGameEvent(r *bitReader) *HijackReplayGameEvent {
	out := &HijackReplayGameEvent{}
	out.UserInfos = decodeUnknown125(r)
	out.Method = int8(readBits(r, 1))
	return out
}

// typeinfo 128 (struct)
// names: set(['TriggerPurchaseMadeEvent', 'TriggerPurchasePanelSelectedPurchaseItemChangedEvent'])
type TriggerPurchaseMadeEvent struct {
	EventMeta
	PurchaseItemId int64 // 84
}
type TriggerPurchasePanelSelectedPurchaseItemChangedEvent struct {
	EventMeta
	PurchaseItemId int64 // 84
}

func decodeTriggerPurchaseMadeEvent(r *bitReader) *TriggerPurchaseMadeEvent {
	out := &TriggerPurchaseMadeEvent{}
	out.PurchaseItemId = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}
func decodeTriggerPurchasePanelSelectedPurchaseItemChangedEvent(r *bitReader) *TriggerPurchasePanelSelectedPurchaseItemChangedEvent {
	out := &TriggerPurchasePanelSelectedPurchaseItemChangedEvent{}
	out.PurchaseItemId = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 129 (struct)
// names: set(['TriggerVictoryPanelPlayMissionAgainEvent', 'TriggerPlanetMissionLaunchedEvent'])
type TriggerVictoryPanelPlayMissionAgainEvent struct {
	EventMeta
	DifficultyLevel int64 // 84
}
type TriggerPlanetMissionLaunchedEvent struct {
	EventMeta
	DifficultyLevel int64 // 84
}

func decodeTriggerVictoryPanelPlayMissionAgainEvent(r *bitReader) *TriggerVictoryPanelPlayMissionAgainEvent {
	out := &TriggerVictoryPanelPlayMissionAgainEvent{}
	out.DifficultyLevel = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}
func decodeTriggerPlanetMissionLaunchedEvent(r *bitReader) *TriggerPlanetMissionLaunchedEvent {
	out := &TriggerPlanetMissionLaunchedEvent{}
	out.DifficultyLevel = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 130 (choice)
func decodeUnknown130(r *bitReader) interface{} {
	switch tag := readBits(r, 3); tag {
	case 0: // None
		return nil
	case 1: // Checked
		return readBits(r, 1) != 0
	case 2: // ValueChanged
		return int32(readBits(r, 32))
	case 3: // SelectionChanged
		return int64(-2147483648 + int64(readBits(r, 32)))
	case 4: // TextChanged
		return decodeByteString_0_11(r)
	case 5: // MouseButton
		return int32(readBits(r, 32))
	default:
		panic(fmt.Errorf("unknown choice tag %d", tag))
	}
}

// typeinfo 131 (struct)
type TriggerDialogControlEvent struct {
	EventMeta
	ControlId int64       // 84
	EventType int64       // 84
	EventData interface{} // 130
}

func decodeTriggerDialogControlEvent(r *bitReader) *TriggerDialogControlEvent {
	out := &TriggerDialogControlEvent{}
	out.ControlId = int64(-2147483648 + int64(readBits(r, 32)))
	out.EventType = int64(-2147483648 + int64(readBits(r, 32)))
	out.EventData = decodeUnknown130(r)
	return out
}

// typeinfo 132 (struct)
type TriggerSoundLengthQueryEvent struct {
	EventMeta
	SoundHash int32 // 6
	Length    int32 // 6
}

func decodeTriggerSoundLengthQueryEvent(r *bitReader) *TriggerSoundLengthQueryEvent {
	out := &TriggerSoundLengthQueryEvent{}
	out.SoundHash = int32(readBits(r, 32))
	out.Length = int32(readBits(r, 32))
	return out
}

// typeinfo 133 (array)
func decodeUnknown133(r *bitReader) []int32 {
	n := int(readBits(r, 7))
	arr := make([]int32, n)
	for i := 0; i < n; i++ {
		arr[i] = int32(readBits(r, 32))
	}
	return arr
}

// typeinfo 134 (struct)
type SyncInfo struct {
	SoundHash []int32 // 133
	Length    []int32 // 133
}

func decodeSyncInfo(r *bitReader) *SyncInfo {
	out := &SyncInfo{}
	out.SoundHash = decodeUnknown133(r)
	out.Length = decodeUnknown133(r)
	return out
}

// typeinfo 135 (struct)
type TriggerSoundLengthSyncEvent struct {
	EventMeta
	SyncInfo *SyncInfo // 134
}

func decodeTriggerSoundLengthSyncEvent(r *bitReader) *TriggerSoundLengthSyncEvent {
	out := &TriggerSoundLengthSyncEvent{}
	out.SyncInfo = decodeSyncInfo(r)
	return out
}

// typeinfo 136 (struct)
type TriggerAnimLengthQueryByNameEvent struct {
	EventMeta
	QueryId        int16 // 79
	LengthMs       int32 // 6
	FinishGameLoop int32 // 6
}

func decodeTriggerAnimLengthQueryByNameEvent(r *bitReader) *TriggerAnimLengthQueryByNameEvent {
	out := &TriggerAnimLengthQueryByNameEvent{}
	out.QueryId = int16(readBits(r, 16))
	out.LengthMs = int32(readBits(r, 32))
	out.FinishGameLoop = int32(readBits(r, 32))
	return out
}

// typeinfo 137 (struct)
type TriggerAnimLengthQueryByPropsEvent struct {
	EventMeta
	QueryId  int16 // 79
	LengthMs int32 // 6
}

func decodeTriggerAnimLengthQueryByPropsEvent(r *bitReader) *TriggerAnimLengthQueryByPropsEvent {
	out := &TriggerAnimLengthQueryByPropsEvent{}
	out.QueryId = int16(readBits(r, 16))
	out.LengthMs = int32(readBits(r, 32))
	return out
}

// typeinfo 138 (struct)
type TriggerAnimOffsetEvent struct {
	EventMeta
	AnimWaitQueryId int16 // 79
}

func decodeTriggerAnimOffsetEvent(r *bitReader) *TriggerAnimOffsetEvent {
	out := &TriggerAnimOffsetEvent{}
	out.AnimWaitQueryId = int16(readBits(r, 16))
	return out
}

// typeinfo 139 (struct)
type TriggerSoundOffsetEvent struct {
	EventMeta
	Sound int32 // 6
}

func decodeTriggerSoundOffsetEvent(r *bitReader) *TriggerSoundOffsetEvent {
	out := &TriggerSoundOffsetEvent{}
	out.Sound = int32(readBits(r, 32))
	return out
}

// typeinfo 140 (struct)
type TriggerTransmissionOffsetEvent struct {
	EventMeta
	TransmissionId int64 // 84
	Thread         int32 // 6
}

func decodeTriggerTransmissionOffsetEvent(r *bitReader) *TriggerTransmissionOffsetEvent {
	out := &TriggerTransmissionOffsetEvent{}
	out.TransmissionId = int64(-2147483648 + int64(readBits(r, 32)))
	out.Thread = int32(readBits(r, 32))
	return out
}

// typeinfo 141 (struct)
type TriggerTransmissionCompleteEvent struct {
	EventMeta
	TransmissionId int64 // 84
}

func decodeTriggerTransmissionCompleteEvent(r *bitReader) *TriggerTransmissionCompleteEvent {
	out := &TriggerTransmissionCompleteEvent{}
	out.TransmissionId = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 142 (optional)
func decodeUnknown142(r *bitReader) **CameraTarget {
	if readBits(r, 1) != 0 {
		ret := decodeCameraTarget(r)
		return &ret
	}
	return nil
}

// typeinfo 143 (optional)
func decodeUnknown143(r *bitReader) *int16 {
	if readBits(r, 1) != 0 {
		ret := int16(readBits(r, 16))
		return &ret
	}
	return nil
}

// typeinfo 144 (optional)
func decodeUnknown144(r *bitReader) *int64 {
	if readBits(r, 1) != 0 {
		ret := int64(-128 + int64(readBits(r, 8)))
		return &ret
	}
	return nil
}

// typeinfo 145 (struct)
type CameraUpdateEvent struct {
	EventMeta
	Target   **CameraTarget // 142
	Distance *int16         // 143
	Pitch    *int16         // 143
	Yaw      *int16         // 143
	Reason   *int64         // 144
	Follow   bool           // 13
}

func decodeCameraUpdateEvent(r *bitReader) *CameraUpdateEvent {
	out := &CameraUpdateEvent{}
	out.Target = decodeUnknown142(r)
	out.Distance = decodeUnknown143(r)
	out.Pitch = decodeUnknown143(r)
	out.Yaw = decodeUnknown143(r)
	out.Reason = decodeUnknown144(r)
	out.Follow = readBits(r, 1) != 0
	return out
}

// typeinfo 146 (struct)
type TriggerConversationSkippedEvent struct {
	EventMeta
	SkipType int8 // 126
}

func decodeTriggerConversationSkippedEvent(r *bitReader) *TriggerConversationSkippedEvent {
	out := &TriggerConversationSkippedEvent{}
	out.SkipType = int8(readBits(r, 1))
	return out
}

// typeinfo 147 (int)

// typeinfo 148 (struct)
type PosUI struct {
	X int16 // 147
	Y int16 // 147
}

func decodePosUI(r *bitReader) *PosUI {
	out := &PosUI{}
	out.X = int16(readBits(r, 11))
	out.Y = int16(readBits(r, 11))
	return out
}

// typeinfo 149 (struct)
type TriggerMouseClickedEvent struct {
	EventMeta
	Button   int32     // 6
	Down     bool      // 13
	PosUI    *PosUI    // 148
	PosWorld *PosWorld // 93
	Flags    int64     // 112
}

func decodeTriggerMouseClickedEvent(r *bitReader) *TriggerMouseClickedEvent {
	out := &TriggerMouseClickedEvent{}
	out.Button = int32(readBits(r, 32))
	out.Down = readBits(r, 1) != 0
	out.PosUI = decodePosUI(r)
	out.PosWorld = decodePosWorld(r)
	out.Flags = int64(-128 + int64(readBits(r, 8)))
	return out
}

// typeinfo 150 (struct)
type TriggerMouseMovedEvent struct {
	EventMeta
	PosUI    *PosUI    // 148
	PosWorld *PosWorld // 93
	Flags    int64     // 112
}

func decodeTriggerMouseMovedEvent(r *bitReader) *TriggerMouseMovedEvent {
	out := &TriggerMouseMovedEvent{}
	out.PosUI = decodePosUI(r)
	out.PosWorld = decodePosWorld(r)
	out.Flags = int64(-128 + int64(readBits(r, 8)))
	return out
}

// typeinfo 151 (struct)
type AchievementAwardedEvent struct {
	EventMeta
	AchievementLink int16 // 79
}

func decodeAchievementAwardedEvent(r *bitReader) *AchievementAwardedEvent {
	out := &AchievementAwardedEvent{}
	out.AchievementLink = int16(readBits(r, 16))
	return out
}

// typeinfo 152 (struct)
type TriggerHotkeyPressedEvent struct {
	EventMeta
	Hotkey int32 // 6
	Down   bool  // 13
}

func decodeTriggerHotkeyPressedEvent(r *bitReader) *TriggerHotkeyPressedEvent {
	out := &TriggerHotkeyPressedEvent{}
	out.Hotkey = int32(readBits(r, 32))
	out.Down = readBits(r, 1) != 0
	return out
}

// typeinfo 153 (struct)
type TriggerTargetModeUpdateEvent struct {
	EventMeta
	AbilLink     int16 // 79
	AbilCmdIndex int8  // 2
	State        int64 // 112
}

func decodeTriggerTargetModeUpdateEvent(r *bitReader) *TriggerTargetModeUpdateEvent {
	out := &TriggerTargetModeUpdateEvent{}
	out.AbilLink = int16(readBits(r, 16))
	out.AbilCmdIndex = int8(readBits(r, 5))
	out.State = int64(-128 + int64(readBits(r, 8)))
	return out
}

// typeinfo 154 (struct)
type TriggerSoundtrackDoneEvent struct {
	EventMeta
	Soundtrack int32 // 6
}

func decodeTriggerSoundtrackDoneEvent(r *bitReader) *TriggerSoundtrackDoneEvent {
	out := &TriggerSoundtrackDoneEvent{}
	out.Soundtrack = int32(readBits(r, 32))
	return out
}

// typeinfo 155 (struct)
type TriggerPlanetMissionSelectedEvent struct {
	EventMeta
	PlanetId int64 // 84
}

func decodeTriggerPlanetMissionSelectedEvent(r *bitReader) *TriggerPlanetMissionSelectedEvent {
	out := &TriggerPlanetMissionSelectedEvent{}
	out.PlanetId = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 156 (struct)
type TriggerKeyPressedEvent struct {
	EventMeta
	Key   int64 // 112
	Flags int64 // 112
}

func decodeTriggerKeyPressedEvent(r *bitReader) *TriggerKeyPressedEvent {
	out := &TriggerKeyPressedEvent{}
	out.Key = int64(-128 + int64(readBits(r, 8)))
	out.Flags = int64(-128 + int64(readBits(r, 8)))
	return out
}

// typeinfo 157 (struct)
type ResourceRequestEvent struct {
	EventMeta
	Resources []int64 // 109
}

func decodeResourceRequestEvent(r *bitReader) *ResourceRequestEvent {
	out := &ResourceRequestEvent{}
	out.Resources = decodeUnknown109(r)
	return out
}

// typeinfo 158 (struct)
type ResourceRequestFulfillEvent struct {
	EventMeta
	FulfillRequestId int64 // 84
}

func decodeResourceRequestFulfillEvent(r *bitReader) *ResourceRequestFulfillEvent {
	out := &ResourceRequestFulfillEvent{}
	out.FulfillRequestId = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 159 (struct)
type ResourceRequestCancelEvent struct {
	EventMeta
	CancelRequestId int64 // 84
}

func decodeResourceRequestCancelEvent(r *bitReader) *ResourceRequestCancelEvent {
	out := &ResourceRequestCancelEvent{}
	out.CancelRequestId = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 160 (struct)
type TriggerResearchPanelSelectionChangedEvent struct {
	EventMeta
	ResearchItemId int64 // 84
}

func decodeTriggerResearchPanelSelectionChangedEvent(r *bitReader) *TriggerResearchPanelSelectionChangedEvent {
	out := &TriggerResearchPanelSelectionChangedEvent{}
	out.ResearchItemId = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 161 (struct)
type TriggerMercenaryPanelSelectionChangedEvent struct {
	EventMeta
	MercenaryId int64 // 84
}

func decodeTriggerMercenaryPanelSelectionChangedEvent(r *bitReader) *TriggerMercenaryPanelSelectionChangedEvent {
	out := &TriggerMercenaryPanelSelectionChangedEvent{}
	out.MercenaryId = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 162 (struct)
type TriggerBattleReportPanelPlayMissionEvent struct {
	EventMeta
	BattleReportId  int64 // 84
	DifficultyLevel int64 // 84
}

func decodeTriggerBattleReportPanelPlayMissionEvent(r *bitReader) *TriggerBattleReportPanelPlayMissionEvent {
	out := &TriggerBattleReportPanelPlayMissionEvent{}
	out.BattleReportId = int64(-2147483648 + int64(readBits(r, 32)))
	out.DifficultyLevel = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 163 (struct)
// names: set(['TriggerBattleReportPanelPlaySceneEvent', 'TriggerBattleReportPanelSelectionChangedEvent'])
type TriggerBattleReportPanelPlaySceneEvent struct {
	EventMeta
	BattleReportId int64 // 84
}
type TriggerBattleReportPanelSelectionChangedEvent struct {
	EventMeta
	BattleReportId int64 // 84
}

func decodeTriggerBattleReportPanelPlaySceneEvent(r *bitReader) *TriggerBattleReportPanelPlaySceneEvent {
	out := &TriggerBattleReportPanelPlaySceneEvent{}
	out.BattleReportId = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}
func decodeTriggerBattleReportPanelSelectionChangedEvent(r *bitReader) *TriggerBattleReportPanelSelectionChangedEvent {
	out := &TriggerBattleReportPanelSelectionChangedEvent{}
	out.BattleReportId = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 164 (int)

// typeinfo 165 (struct)
type DecrementGameTimeRemainingEvent struct {
	EventMeta
	DecrementMs int32 // 164
}

func decodeDecrementGameTimeRemainingEvent(r *bitReader) *DecrementGameTimeRemainingEvent {
	out := &DecrementGameTimeRemainingEvent{}
	out.DecrementMs = int32(readBits(r, 19))
	return out
}

// typeinfo 166 (struct)
type TriggerPortraitLoadedEvent struct {
	EventMeta
	PortraitId int64 // 84
}

func decodeTriggerPortraitLoadedEvent(r *bitReader) *TriggerPortraitLoadedEvent {
	out := &TriggerPortraitLoadedEvent{}
	out.PortraitId = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 167 (struct)
type TriggerMovieFunctionEvent struct {
	EventMeta
	FunctionName string // 20
}

func decodeTriggerMovieFunctionEvent(r *bitReader) *TriggerMovieFunctionEvent {
	out := &TriggerMovieFunctionEvent{}
	out.FunctionName = decodeByteString_0_7(r)
	return out
}

// typeinfo 168 (struct)
type TriggerCustomDialogDismissedEvent struct {
	EventMeta
	Result int64 // 84
}

func decodeTriggerCustomDialogDismissedEvent(r *bitReader) *TriggerCustomDialogDismissedEvent {
	out := &TriggerCustomDialogDismissedEvent{}
	out.Result = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 169 (struct)
type TriggerGameMenuItemSelectedEvent struct {
	EventMeta
	GameMenuItemIndex int64 // 84
}

func decodeTriggerGameMenuItemSelectedEvent(r *bitReader) *TriggerGameMenuItemSelectedEvent {
	out := &TriggerGameMenuItemSelectedEvent{}
	out.GameMenuItemIndex = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 170 (struct)
type TriggerPurchasePanelSelectedPurchaseCategoryChangedEvent struct {
	EventMeta
	PurchaseCategoryId int64 // 84
}

func decodeTriggerPurchasePanelSelectedPurchaseCategoryChangedEvent(r *bitReader) *TriggerPurchasePanelSelectedPurchaseCategoryChangedEvent {
	out := &TriggerPurchasePanelSelectedPurchaseCategoryChangedEvent{}
	out.PurchaseCategoryId = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 171 (struct)
type TriggerButtonPressedEvent struct {
	EventMeta
	Button int16 // 79
}

func decodeTriggerButtonPressedEvent(r *bitReader) *TriggerButtonPressedEvent {
	out := &TriggerButtonPressedEvent{}
	out.Button = int16(readBits(r, 16))
	return out
}

// typeinfo 172 (struct)
type TriggerCutsceneBookmarkFiredEvent struct {
	EventMeta
	CutsceneId   int64  // 84
	BookmarkName string // 20
}

func decodeTriggerCutsceneBookmarkFiredEvent(r *bitReader) *TriggerCutsceneBookmarkFiredEvent {
	out := &TriggerCutsceneBookmarkFiredEvent{}
	out.CutsceneId = int64(-2147483648 + int64(readBits(r, 32)))
	out.BookmarkName = decodeByteString_0_7(r)
	return out
}

// typeinfo 173 (struct)
type TriggerCutsceneEndSceneFiredEvent struct {
	EventMeta
	CutsceneId int64 // 84
}

func decodeTriggerCutsceneEndSceneFiredEvent(r *bitReader) *TriggerCutsceneEndSceneFiredEvent {
	out := &TriggerCutsceneEndSceneFiredEvent{}
	out.CutsceneId = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 174 (struct)
type TriggerCutsceneConversationLineEvent struct {
	EventMeta
	CutsceneId          int64  // 84
	ConversationLine    string // 20
	AltConversationLine string // 20
}

func decodeTriggerCutsceneConversationLineEvent(r *bitReader) *TriggerCutsceneConversationLineEvent {
	out := &TriggerCutsceneConversationLineEvent{}
	out.CutsceneId = int64(-2147483648 + int64(readBits(r, 32)))
	out.ConversationLine = decodeByteString_0_7(r)
	out.AltConversationLine = decodeByteString_0_7(r)
	return out
}

// typeinfo 175 (struct)
type TriggerCutsceneConversationLineMissingEvent struct {
	EventMeta
	CutsceneId       int64  // 84
	ConversationLine string // 20
}

func decodeTriggerCutsceneConversationLineMissingEvent(r *bitReader) *TriggerCutsceneConversationLineMissingEvent {
	out := &TriggerCutsceneConversationLineMissingEvent{}
	out.CutsceneId = int64(-2147483648 + int64(readBits(r, 32)))
	out.ConversationLine = decodeByteString_0_7(r)
	return out
}

// typeinfo 176 (struct)
type GameUserLeaveEvent struct {
	EventMeta
	LeaveReason int8 // 1
}

func decodeGameUserLeaveEvent(r *bitReader) *GameUserLeaveEvent {
	out := &GameUserLeaveEvent{}
	out.LeaveReason = int8(readBits(r, 4))
	return out
}

// typeinfo 177 (struct)
type GameUserJoinEvent struct {
	EventMeta
	Observe               int8    // 24
	Name                  string  // 9
	ToonHandle            *string // 123
	ClanTag               *string // 41
	ClanLogo              *string // 42
	Hijack                bool    // 13
	HijackCloneGameUserId *int8   // 59
}

func decodeGameUserJoinEvent(r *bitReader) *GameUserJoinEvent {
	out := &GameUserJoinEvent{}
	out.Observe = int8(readBits(r, 2))
	out.Name = decodeByteString_0_8(r)
	out.ToonHandle = decodeUnknown123(r)
	out.ClanTag = decodeUnknown41(r)
	out.ClanLogo = decodeUnknown42(r)
	out.Hijack = readBits(r, 1) != 0
	out.HijackCloneGameUserId = decodeUnknown59(r)
	return out
}

// typeinfo 178 (optional)
func decodeUnknown178(r *bitReader) *int64 {
	if readBits(r, 1) != 0 {
		ret := int64(1 + int64(readBits(r, 32)))
		return &ret
	}
	return nil
}

// typeinfo 179 (struct)
type CommandManagerStateEvent struct {
	EventMeta
	State    int8   // 24
	Sequence *int64 // 178
}

func decodeCommandManagerStateEvent(r *bitReader) *CommandManagerStateEvent {
	out := &CommandManagerStateEvent{}
	out.State = int8(readBits(r, 2))
	out.Sequence = decodeUnknown178(r)
	return out
}

// typeinfo 180 (struct)
type CmdUpdateTargetPointEvent struct {
	EventMeta
	Target *PosWorld // 93
}

func decodeCmdUpdateTargetPointEvent(r *bitReader) *CmdUpdateTargetPointEvent {
	out := &CmdUpdateTargetPointEvent{}
	out.Target = decodePosWorld(r)
	return out
}

// typeinfo 181 (struct)
type CmdUpdateTargetUnitEvent struct {
	EventMeta
	Target *Target // 94
}

func decodeCmdUpdateTargetUnitEvent(r *bitReader) *CmdUpdateTargetUnitEvent {
	out := &CmdUpdateTargetUnitEvent{}
	out.Target = decodeTarget(r)
	return out
}

// typeinfo 182 (struct)
type CatalogModifyEvent struct {
	EventMeta
	Catalog int8   // 10
	Entry   int16  // 79
	Field   string // 9
	Value   string // 9
}

func decodeCatalogModifyEvent(r *bitReader) *CatalogModifyEvent {
	out := &CatalogModifyEvent{}
	out.Catalog = int8(readBits(r, 8))
	out.Entry = int16(readBits(r, 16))
	out.Field = decodeByteString_0_8(r)
	out.Value = decodeByteString_0_8(r)
	return out
}

// typeinfo 183 (struct)
type HeroTalentTreeSelectedEvent struct {
	EventMeta
	Index int32 // 6
}

func decodeHeroTalentTreeSelectedEvent(r *bitReader) *HeroTalentTreeSelectedEvent {
	out := &HeroTalentTreeSelectedEvent{}
	out.Index = int32(readBits(r, 32))
	return out
}

// typeinfo 184 (struct)
type HeroTalentTreeSelectionPanelToggledEvent struct {
	EventMeta
	Shown bool // 13
}

func decodeHeroTalentTreeSelectionPanelToggledEvent(r *bitReader) *HeroTalentTreeSelectionPanelToggledEvent {
	out := &HeroTalentTreeSelectionPanelToggledEvent{}
	out.Shown = readBits(r, 1) != 0
	return out
}

// typeinfo 185 (struct)
type ChatMessage struct {
	EventMeta
	Recipient int8   // 12
	String    string // 30
}

func decodeChatMessage(r *bitReader) *ChatMessage {
	out := &ChatMessage{}
	out.Recipient = int8(readBits(r, 3))
	out.String = decodeByteString_0_11(r)
	return out
}

// typeinfo 186 (struct)
type PingMessage struct {
	EventMeta
	Recipient int8   // 12
	Point     *Point // 85
}

func decodePingMessage(r *bitReader) *PingMessage {
	out := &PingMessage{}
	out.Recipient = int8(readBits(r, 3))
	out.Point = decodePoint(r)
	return out
}

// typeinfo 187 (struct)
type LoadingProgressMessage struct {
	EventMeta
	Progress int64 // 84
}

func decodeLoadingProgressMessage(r *bitReader) *LoadingProgressMessage {
	out := &LoadingProgressMessage{}
	out.Progress = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 188 (struct)
type ReconnectNotifyMessage struct {
	EventMeta
	Status int8 // 24
}

func decodeReconnectNotifyMessage(r *bitReader) *ReconnectNotifyMessage {
	out := &ReconnectNotifyMessage{}
	out.Status = int8(readBits(r, 2))
	return out
}

// typeinfo 189 (struct)
type Stats struct {
	ScoreValueMineralsCurrent                  int64 // 84
	ScoreValueVespeneCurrent                   int64 // 84
	ScoreValueMineralsCollectionRate           int64 // 84
	ScoreValueVespeneCollectionRate            int64 // 84
	ScoreValueWorkersActiveCount               int64 // 84
	ScoreValueMineralsUsedInProgressArmy       int64 // 84
	ScoreValueMineralsUsedInProgressEconomy    int64 // 84
	ScoreValueMineralsUsedInProgressTechnology int64 // 84
	ScoreValueVespeneUsedInProgressArmy        int64 // 84
	ScoreValueVespeneUsedInProgressEconomy     int64 // 84
	ScoreValueVespeneUsedInProgressTechnology  int64 // 84
	ScoreValueMineralsUsedCurrentArmy          int64 // 84
	ScoreValueMineralsUsedCurrentEconomy       int64 // 84
	ScoreValueMineralsUsedCurrentTechnology    int64 // 84
	ScoreValueVespeneUsedCurrentArmy           int64 // 84
	ScoreValueVespeneUsedCurrentEconomy        int64 // 84
	ScoreValueVespeneUsedCurrentTechnology     int64 // 84
	ScoreValueMineralsLostArmy                 int64 // 84
	ScoreValueMineralsLostEconomy              int64 // 84
	ScoreValueMineralsLostTechnology           int64 // 84
	ScoreValueVespeneLostArmy                  int64 // 84
	ScoreValueVespeneLostEconomy               int64 // 84
	ScoreValueVespeneLostTechnology            int64 // 84
	ScoreValueMineralsKilledArmy               int64 // 84
	ScoreValueMineralsKilledEconomy            int64 // 84
	ScoreValueMineralsKilledTechnology         int64 // 84
	ScoreValueVespeneKilledArmy                int64 // 84
	ScoreValueVespeneKilledEconomy             int64 // 84
	ScoreValueVespeneKilledTechnology          int64 // 84
	ScoreValueFoodUsed                         int64 // 84
	ScoreValueFoodMade                         int64 // 84
	ScoreValueMineralsUsedActiveForces         int64 // 84
	ScoreValueVespeneUsedActiveForces          int64 // 84
	ScoreValueMineralsFriendlyFireArmy         int64 // 84
	ScoreValueMineralsFriendlyFireEconomy      int64 // 84
	ScoreValueMineralsFriendlyFireTechnology   int64 // 84
	ScoreValueVespeneFriendlyFireArmy          int64 // 84
	ScoreValueVespeneFriendlyFireEconomy       int64 // 84
	ScoreValueVespeneFriendlyFireTechnology    int64 // 84
}

func decodeStats(r *bitReader) *Stats {
	out := &Stats{}
	out.ScoreValueMineralsCurrent = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneCurrent = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueMineralsCollectionRate = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneCollectionRate = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueWorkersActiveCount = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueMineralsUsedInProgressArmy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueMineralsUsedInProgressEconomy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueMineralsUsedInProgressTechnology = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneUsedInProgressArmy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneUsedInProgressEconomy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneUsedInProgressTechnology = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueMineralsUsedCurrentArmy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueMineralsUsedCurrentEconomy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueMineralsUsedCurrentTechnology = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneUsedCurrentArmy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneUsedCurrentEconomy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneUsedCurrentTechnology = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueMineralsLostArmy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueMineralsLostEconomy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueMineralsLostTechnology = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneLostArmy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneLostEconomy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneLostTechnology = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueMineralsKilledArmy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueMineralsKilledEconomy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueMineralsKilledTechnology = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneKilledArmy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneKilledEconomy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneKilledTechnology = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueFoodUsed = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueFoodMade = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueMineralsUsedActiveForces = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneUsedActiveForces = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueMineralsFriendlyFireArmy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueMineralsFriendlyFireEconomy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueMineralsFriendlyFireTechnology = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneFriendlyFireArmy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneFriendlyFireEconomy = int64(-2147483648 + int64(readBits(r, 32)))
	out.ScoreValueVespeneFriendlyFireTechnology = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 190 (struct)
type PlayerStatsEvent struct {
	EventMeta
	PlayerId int8   // 1
	Stats    *Stats // 189
}

func decodePlayerStatsEvent(r *bitReader) *PlayerStatsEvent {
	out := &PlayerStatsEvent{}
	out.PlayerId = int8(readBits(r, 4))
	out.Stats = decodeStats(r)
	return out
}

// typeinfo 191 (struct)
// names: set(['UnitBornEvent', 'UnitInitEvent'])
type UnitBornEvent struct {
	EventMeta
	UnitTagIndex    int32  // 6
	UnitTagRecycle  int32  // 6
	UnitTypeName    string // 29
	ControlPlayerId int8   // 1
	UpkeepPlayerId  int8   // 1
	X               int8   // 10
	Y               int8   // 10
}
type UnitInitEvent struct {
	EventMeta
	UnitTagIndex    int32  // 6
	UnitTagRecycle  int32  // 6
	UnitTypeName    string // 29
	ControlPlayerId int8   // 1
	UpkeepPlayerId  int8   // 1
	X               int8   // 10
	Y               int8   // 10
}

func decodeUnitBornEvent(r *bitReader) *UnitBornEvent {
	out := &UnitBornEvent{}
	out.UnitTagIndex = int32(readBits(r, 32))
	out.UnitTagRecycle = int32(readBits(r, 32))
	out.UnitTypeName = decodeByteString_0_10(r)
	out.ControlPlayerId = int8(readBits(r, 4))
	out.UpkeepPlayerId = int8(readBits(r, 4))
	out.X = int8(readBits(r, 8))
	out.Y = int8(readBits(r, 8))
	return out
}
func decodeUnitInitEvent(r *bitReader) *UnitInitEvent {
	out := &UnitInitEvent{}
	out.UnitTagIndex = int32(readBits(r, 32))
	out.UnitTagRecycle = int32(readBits(r, 32))
	out.UnitTypeName = decodeByteString_0_10(r)
	out.ControlPlayerId = int8(readBits(r, 4))
	out.UpkeepPlayerId = int8(readBits(r, 4))
	out.X = int8(readBits(r, 8))
	out.Y = int8(readBits(r, 8))
	return out
}

// typeinfo 192 (struct)
type UnitDiedEvent struct {
	EventMeta
	UnitTagIndex         int32  // 6
	UnitTagRecycle       int32  // 6
	KillerPlayerId       *int8  // 59
	X                    int8   // 10
	Y                    int8   // 10
	KillerUnitTagIndex   *int32 // 43
	KillerUnitTagRecycle *int32 // 43
}

func decodeUnitDiedEvent(r *bitReader) *UnitDiedEvent {
	out := &UnitDiedEvent{}
	out.UnitTagIndex = int32(readBits(r, 32))
	out.UnitTagRecycle = int32(readBits(r, 32))
	out.KillerPlayerId = decodeUnknown59(r)
	out.X = int8(readBits(r, 8))
	out.Y = int8(readBits(r, 8))
	out.KillerUnitTagIndex = decodeUnknown43(r)
	out.KillerUnitTagRecycle = decodeUnknown43(r)
	return out
}

// typeinfo 193 (struct)
type UnitOwnerChangeEvent struct {
	EventMeta
	UnitTagIndex    int32 // 6
	UnitTagRecycle  int32 // 6
	ControlPlayerId int8  // 1
	UpkeepPlayerId  int8  // 1
}

func decodeUnitOwnerChangeEvent(r *bitReader) *UnitOwnerChangeEvent {
	out := &UnitOwnerChangeEvent{}
	out.UnitTagIndex = int32(readBits(r, 32))
	out.UnitTagRecycle = int32(readBits(r, 32))
	out.ControlPlayerId = int8(readBits(r, 4))
	out.UpkeepPlayerId = int8(readBits(r, 4))
	return out
}

// typeinfo 194 (struct)
type UnitTypeChangeEvent struct {
	EventMeta
	UnitTagIndex   int32  // 6
	UnitTagRecycle int32  // 6
	UnitTypeName   string // 29
}

func decodeUnitTypeChangeEvent(r *bitReader) *UnitTypeChangeEvent {
	out := &UnitTypeChangeEvent{}
	out.UnitTagIndex = int32(readBits(r, 32))
	out.UnitTagRecycle = int32(readBits(r, 32))
	out.UnitTypeName = decodeByteString_0_10(r)
	return out
}

// typeinfo 195 (struct)
type UpgradeEvent struct {
	EventMeta
	PlayerId        int8   // 1
	UpgradeTypeName string // 29
	Count           int64  // 84
}

func decodeUpgradeEvent(r *bitReader) *UpgradeEvent {
	out := &UpgradeEvent{}
	out.PlayerId = int8(readBits(r, 4))
	out.UpgradeTypeName = decodeByteString_0_10(r)
	out.Count = int64(-2147483648 + int64(readBits(r, 32)))
	return out
}

// typeinfo 196 (struct)
type UnitDoneEvent struct {
	EventMeta
	UnitTagIndex   int32 // 6
	UnitTagRecycle int32 // 6
}

func decodeUnitDoneEvent(r *bitReader) *UnitDoneEvent {
	out := &UnitDoneEvent{}
	out.UnitTagIndex = int32(readBits(r, 32))
	out.UnitTagRecycle = int32(readBits(r, 32))
	return out
}

// typeinfo 197 (array)
func decodeUnknown197(r *bitReader) []int64 {
	n := int(readBits(r, 10))
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(-2147483648 + int64(readBits(r, 32)))
	}
	return arr
}

// typeinfo 198 (struct)
type UnitPositionsEvent struct {
	EventMeta
	FirstUnitIndex int32   // 6
	Items          []int64 // 197
}

func decodeUnitPositionsEvent(r *bitReader) *UnitPositionsEvent {
	out := &UnitPositionsEvent{}
	out.FirstUnitIndex = int32(readBits(r, 32))
	out.Items = decodeUnknown197(r)
	return out
}

// typeinfo 199 (struct)
type PlayerSetupEvent struct {
	EventMeta
	PlayerId int8   // 1
	Type     int32  // 6
	UserId   *int32 // 43
	SlotId   *int32 // 43
}

func decodePlayerSetupEvent(r *bitReader) *PlayerSetupEvent {
	out := &PlayerSetupEvent{}
	out.PlayerId = int8(readBits(r, 4))
	out.Type = int32(readBits(r, 32))
	out.UserId = decodeUnknown43(r)
	out.SlotId = decodeUnknown43(r)
	return out
}
func readGameEvent(r *bitReader, typ int) Event {
	switch typ {
	case 5:
		return &UserFinishedLoadingSyncEvent{}
	case 7:
		return decodeUserOptionsEvent(r)
	case 9:
		return decodeBankFileEvent(r)
	case 10:
		return decodeBankSectionEvent(r)
	case 11:
		return decodeBankKeyEvent(r)
	case 12:
		return decodeBankValueEvent(r)
	case 13:
		return decodeBankSignatureEvent(r)
	case 14:
		return decodeCameraSaveEvent(r)
	case 21:
		return decodeSaveGameEvent(r)
	case 22:
		return &SaveGameDoneEvent{}
	case 23:
		return &LoadGameDoneEvent{}
	case 25:
		return decodeCommandManagerResetEvent(r)
	case 26:
		return decodeGameCheatEvent(r)
	case 27:
		return decodeCmdEvent(r)
	case 28:
		return decodeSelectionDeltaEvent(r)
	case 29:
		return decodeControlGroupUpdateEvent(r)
	case 30:
		return decodeSelectionSyncCheckEvent(r)
	case 31:
		return decodeResourceTradeEvent(r)
	case 32:
		return decodeTriggerChatMessageEvent(r)
	case 33:
		return decodeAICommunicateEvent(r)
	case 34:
		return decodeSetAbsoluteGameSpeedEvent(r)
	case 35:
		return decodeAddAbsoluteGameSpeedEvent(r)
	case 36:
		return decodeTriggerPingEvent(r)
	case 37:
		return decodeBroadcastCheatEvent(r)
	case 38:
		return decodeAllianceEvent(r)
	case 39:
		return decodeUnitClickEvent(r)
	case 40:
		return decodeUnitHighlightEvent(r)
	case 41:
		return decodeTriggerReplySelectedEvent(r)
	case 43:
		return decodeHijackReplayGameEvent(r)
	case 44:
		return &TriggerSkippedEvent{}
	case 45:
		return decodeTriggerSoundLengthQueryEvent(r)
	case 46:
		return decodeTriggerSoundOffsetEvent(r)
	case 47:
		return decodeTriggerTransmissionOffsetEvent(r)
	case 48:
		return decodeTriggerTransmissionCompleteEvent(r)
	case 49:
		return decodeCameraUpdateEvent(r)
	case 50:
		return &TriggerAbortMissionEvent{}
	case 51:
		return decodeTriggerPurchaseMadeEvent(r)
	case 52:
		return &TriggerPurchaseExitEvent{}
	case 53:
		return decodeTriggerPlanetMissionLaunchedEvent(r)
	case 54:
		return &TriggerPlanetPanelCanceledEvent{}
	case 55:
		return decodeTriggerDialogControlEvent(r)
	case 56:
		return decodeTriggerSoundLengthSyncEvent(r)
	case 57:
		return decodeTriggerConversationSkippedEvent(r)
	case 58:
		return decodeTriggerMouseClickedEvent(r)
	case 59:
		return decodeTriggerMouseMovedEvent(r)
	case 60:
		return decodeAchievementAwardedEvent(r)
	case 61:
		return decodeTriggerHotkeyPressedEvent(r)
	case 62:
		return decodeTriggerTargetModeUpdateEvent(r)
	case 63:
		return &TriggerPlanetPanelReplayEvent{}
	case 64:
		return decodeTriggerSoundtrackDoneEvent(r)
	case 65:
		return decodeTriggerPlanetMissionSelectedEvent(r)
	case 66:
		return decodeTriggerKeyPressedEvent(r)
	case 67:
		return decodeTriggerMovieFunctionEvent(r)
	case 68:
		return &TriggerPlanetPanelBirthCompleteEvent{}
	case 69:
		return &TriggerPlanetPanelDeathCompleteEvent{}
	case 70:
		return decodeResourceRequestEvent(r)
	case 71:
		return decodeResourceRequestFulfillEvent(r)
	case 72:
		return decodeResourceRequestCancelEvent(r)
	case 73:
		return &TriggerResearchPanelExitEvent{}
	case 74:
		return &TriggerResearchPanelPurchaseEvent{}
	case 75:
		return decodeTriggerResearchPanelSelectionChangedEvent(r)
	case 77:
		return &TriggerMercenaryPanelExitEvent{}
	case 78:
		return &TriggerMercenaryPanelPurchaseEvent{}
	case 79:
		return decodeTriggerMercenaryPanelSelectionChangedEvent(r)
	case 80:
		return &TriggerVictoryPanelExitEvent{}
	case 81:
		return &TriggerBattleReportPanelExitEvent{}
	case 82:
		return decodeTriggerBattleReportPanelPlayMissionEvent(r)
	case 83:
		return decodeTriggerBattleReportPanelPlaySceneEvent(r)
	case 84:
		return decodeTriggerBattleReportPanelSelectionChangedEvent(r)
	case 85:
		return decodeTriggerVictoryPanelPlayMissionAgainEvent(r)
	case 86:
		return &TriggerMovieStartedEvent{}
	case 87:
		return &TriggerMovieFinishedEvent{}
	case 88:
		return decodeDecrementGameTimeRemainingEvent(r)
	case 89:
		return decodeTriggerPortraitLoadedEvent(r)
	case 90:
		return decodeTriggerCustomDialogDismissedEvent(r)
	case 91:
		return decodeTriggerGameMenuItemSelectedEvent(r)
	case 93:
		return decodeTriggerPurchasePanelSelectedPurchaseItemChangedEvent(r)
	case 94:
		return decodeTriggerPurchasePanelSelectedPurchaseCategoryChangedEvent(r)
	case 95:
		return decodeTriggerButtonPressedEvent(r)
	case 96:
		return &TriggerGameCreditsFinishedEvent{}
	case 97:
		return decodeTriggerCutsceneBookmarkFiredEvent(r)
	case 98:
		return decodeTriggerCutsceneEndSceneFiredEvent(r)
	case 99:
		return decodeTriggerCutsceneConversationLineEvent(r)
	case 100:
		return decodeTriggerCutsceneConversationLineMissingEvent(r)
	case 101:
		return decodeGameUserLeaveEvent(r)
	case 102:
		return decodeGameUserJoinEvent(r)
	case 103:
		return decodeCommandManagerStateEvent(r)
	case 104:
		return decodeCmdUpdateTargetPointEvent(r)
	case 105:
		return decodeCmdUpdateTargetUnitEvent(r)
	case 106:
		return decodeTriggerAnimLengthQueryByNameEvent(r)
	case 107:
		return decodeTriggerAnimLengthQueryByPropsEvent(r)
	case 108:
		return decodeTriggerAnimOffsetEvent(r)
	case 109:
		return decodeCatalogModifyEvent(r)
	case 110:
		return decodeHeroTalentTreeSelectedEvent(r)
	case 111:
		return &TriggerProfilerLoggingFinishedEvent{}
	case 112:
		return decodeHeroTalentTreeSelectionPanelToggledEvent(r)
	default:
		panic(fmt.Errorf("unknown event type %d", typ))
	}
}
func readMessageEvent(r *bitReader, typ int) Event {
	switch typ {
	case 0:
		return decodeChatMessage(r)
	case 1:
		return decodePingMessage(r)
	case 2:
		return decodeLoadingProgressMessage(r)
	case 3:
		return &ServerPingMessage{}
	case 4:
		return decodeReconnectNotifyMessage(r)
	default:
		panic(fmt.Errorf("unknown event type %d", typ))
	}
}
func readTrackerEvent(r *bitReader, typ int) Event {
	switch typ {
	case 0:
		return decodePlayerStatsEvent(r)
	case 1:
		return decodeUnitBornEvent(r)
	case 2:
		return decodeUnitDiedEvent(r)
	case 3:
		return decodeUnitOwnerChangeEvent(r)
	case 4:
		return decodeUnitTypeChangeEvent(r)
	case 5:
		return decodeUpgradeEvent(r)
	case 6:
		return decodeUnitInitEvent(r)
	case 7:
		return decodeUnitDoneEvent(r)
	case 8:
		return decodeUnitPositionsEvent(r)
	case 9:
		return decodePlayerSetupEvent(r)
	default:
		panic(fmt.Errorf("unknown event type %d", typ))
	}
}
