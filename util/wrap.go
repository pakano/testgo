package util

import "io"

type WrapReaderWriter struct {
	R       io.Reader
	W       io.Writer
	Trigger func()
}

func (wrap *WrapReaderWriter) Read(p []byte) (n int, err error) {
	n, err = wrap.R.Read(p)
	if wrap.Trigger != nil {
		wrap.Trigger()
	}
	return
}

func (wrap *WrapReaderWriter) Write(p []byte) (n int, err error) {
	n, err = wrap.W.Write(p)
	if wrap.Trigger != nil {
		wrap.Trigger()
	}
	return
}
