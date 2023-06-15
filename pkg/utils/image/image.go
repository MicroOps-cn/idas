package image

import (
	"bytes"
	"context"
	"image/png"
	"io"
	"unicode/utf8"

	"github.com/MicroOps-cn/idas/pkg/errors"
	avatar "github.com/disintegration/letteravatar"
)

func GenerateAvatar(_ context.Context, content string) (reader io.Reader, err error) {
	firstLetter, _ := utf8.DecodeRuneInString(content)

	img, err := avatar.Draw(128, firstLetter, nil)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	if err = png.Encode(buf, img); err != nil {
		return nil, errors.WithServerError(500, err, "failed to write buffer")
	}
	return buf, nil
}
