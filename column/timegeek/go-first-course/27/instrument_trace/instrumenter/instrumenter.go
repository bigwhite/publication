package instrumenter

type Instrumenter interface {
	Instrument(string) ([]byte, error)
}
