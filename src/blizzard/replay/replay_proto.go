package replay

import "blizzard/blizzval"

func readPlayer(v blizzval.Value) *Player {
	out := Player{}
	out.Raw = v
	m := v.(map[int]blizzval.Value)
	out.Name = m[0].(string)
	out.Team = m[5].(int64)
	out.Handicap = m[6].(int64)
	out.IsWinner = m[8].(int64)
	out.Character = m[10].(string)
	return &out
}
func readDetails(v blizzval.Value) *Details {
	out := Details{}
	out.Raw = v
	m := v.(map[int]blizzval.Value)
	if f, ok := m[0]; ok {
		s := f.([]blizzval.Value)
		out.Players = make([]*Player, len(s))
		for i := 0; i < len(s); i++ {
			out.Players[i] = readPlayer(s[i])
		}
	}
	out.Map = m[1].(string)
	out.TimeStamp = m[5].(int64)
	out.UTCOffset = m[6].(int64)
	return &out
}
