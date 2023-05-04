package charmaps

var (
	All      = Combine(BlocksBasic, Blocks, BlocksAdvanced, ASCII, Misc)
	CharMaps = map[string][]Pattern{
		"blocks":       Blocks,
		"blocks_basic": BlocksBasic,
		"blocks_adv":   BlocksAdvanced,
		"ascii":        ASCII,
		"misc":         Misc,
	}
)

func Combine(patterns ...[]Pattern) []Pattern {
	var pat []Pattern
	for i := range patterns {
		pat = append(pat, patterns[i]...)
	}
	return pat
}
