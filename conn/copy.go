package conn

import (
	"io"
)

func CoreCopy(dst io.Writer, src io.Reader) (written int64, err error) {
	var buffer [8192]byte
	buf := buffer[:]

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}

func CopyLink(dst, src io.ReadWriteCloser) {
	go func() {
		defer src.Close()
		CoreCopy(src, dst)
	}()
	defer dst.Close()
	CoreCopy(dst, src)
}
