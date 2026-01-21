package common

func Transform[I, O any](in []I, conv func(in I) O) []O {
	out := make([]O, 0)
	for i := range in {
		out = append(out, conv(in[i]))
	}
	return out
}
