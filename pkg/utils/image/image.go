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

package image

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"io"

	"github.com/golang/freetype"
	"golang.org/x/exp/utf8string"
	"golang.org/x/image/font"

	"github.com/MicroOps-cn/idas/pkg/errors"
)

func GenerateAvatar(ctx context.Context, content string) (reader io.Reader, err error) {
	buf := &bytes.Buffer{}
	img := image.NewNRGBA(image.Rect(0, 0, 128, 128))
	c := freetype.NewContext()
	tfont, err := LoadSystemFonts(ctx, "PingFangSC", "PingFangSC-Regular", "Microsoft YaHei", "STXihei", "华文细黑", "Georgia", "Times New Roman", "serif")
	if err != nil {
		return nil, errors.WithServerError(500, err, "failed to load font file")
	}
	if tfont == nil {
		return nil, errors.NewServerError(500, "failed to load font file")
	}

	for x := 0; x < img.Rect.Dx(); x++ {
		for y := 0; y < img.Rect.Dx(); y++ {
			img.Set(x, y, color.RGBA{R: 51, G: 112, B: 255, A: 255})
		}
	}

	c.SetDPI(72)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetFont(tfont)
	c.SetFontSize(float64(img.Rect.Dx()) * 2 / 3)
	c.SetSrc(image.NewUniform(image.White))
	c.SetHinting(font.HintingNone)
	pt := freetype.Pt(img.Rect.Dx()/5, img.Rect.Dy()-img.Rect.Dy()/5)
	_, err = c.DrawString(utf8string.NewString(content).Slice(0, 1), pt)
	if err != nil {
		return nil, errors.WithServerError(500, err, "failed to draw text")
	}
	if err = png.Encode(buf, img); err != nil {
		return nil, errors.WithServerError(500, err, "failed to write buffer")
	}
	return buf, err
}
