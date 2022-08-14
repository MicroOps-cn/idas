/*
 Copyright © 2022 MicroOps-cn.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package buffer

import "io"

type PreReader interface {
	io.ReadCloser
	ReadBuf() []byte
	String() string
}

type preReader struct {
	io.ReadCloser
	buf    []byte
	bufIdx int
	maxLen int
}

func (r *preReader) Close() error {
	r.buf = nil
	r.bufIdx = 0
	r.maxLen = 0
	return r.ReadCloser.Close()
}

func (r preReader) ReadBuf() []byte {
	return r.buf
}

func (r preReader) String() string {
	return string(r.buf)
}

func (r *preReader) Read(p []byte) (n int, err error) {
	if r.bufIdx >= r.maxLen {
		return r.ReadCloser.Read(p)
	} else if len(r.buf) <= r.bufIdx {
		return 0, io.EOF
	}
	n = copy(p, r.buf[r.bufIdx:])
	r.bufIdx += n
	return
}

func NewPreReader(reader io.ReadCloser, maxLen int) (PreReader, error) {
	r := &preReader{ReadCloser: reader, maxLen: maxLen}
	for i := 0; i < maxLen; {
		p2 := make([]byte, maxLen-i)
		if n, err := reader.Read(p2); n != 0 {
			r.buf = append(r.buf, p2[:n]...)
			i += n
		} else if err == io.EOF {
			break
		} else if err != nil {
			return r, err
		}
	}
	return r, nil
}
