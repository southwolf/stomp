package stomp

import (
	"bufio"
	"io"
)

// slices used to write frames
var (
	colonSlice   = []byte{58}     // colon ':'
	crlfSlice    = []byte{13, 10} // CR-LF
	newlineSlice = []byte{10}     // newline (LF)
	nullSlice    = []byte{0}      // null character
)

// Writes STOMP frames to an underlying io.Writer.
type Writer struct {
	writer *bufio.Writer
}

// Creates a new Writer object, which writes to an underlying io.Writer.
func NewWriter(writer io.Writer) *Writer {
	return NewWriterSize(writer, 4096)
}

func NewWriterSize(writer io.Writer, bufferSize int) *Writer {
	return &Writer{writer: bufio.NewWriterSize(writer, bufferSize)}
}

// Write the contents of a frame to the underlying io.Writer.
func (w *Writer) Write(f *Frame) error {
	var err error

	if f == nil {
		// nil frame means send a heart-beat LF
		_, err = w.writer.Write(newlineSlice)
		if err != nil {
			return err
		}
	} else {
		_, err = w.writer.Write([]byte(f.Command))
		if err != nil {
			return err
		}

		_, err = w.writer.Write(newlineSlice)
		if err != nil {
			return err
		}

		for key, values := range f.Header {
			for _, value := range values {
				_, err = w.writer.Write([]byte(key))
				if err != nil {
					return err
				}
				_, err = w.writer.Write(colonSlice)
				if err != nil {
					return err
				}
				_, err = w.writer.Write([]byte(value))
				if err != nil {
					return err
				}
				_, err = w.writer.Write(newlineSlice)
				if err != nil {
					return err
				}
			}
		}

		_, err = w.writer.Write(newlineSlice)
		if err != nil {
			return err
		}

		if len(f.Body) > 0 {
			_, err = w.writer.Write(f.Body)
			if err != nil {
				return err
			}
		}

		// write the final null (0) byte	
		_, err = w.writer.Write(nullSlice)
		if err != nil {
			return err
		}
	}

	err = w.writer.Flush()
	return err
}
