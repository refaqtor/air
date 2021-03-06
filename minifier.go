package air

import (
	"bytes"
	"errors"
	"image/jpeg"
	"image/png"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
	"github.com/tdewolff/minify/xml"
)

type (
	// Minifier is used to provide a `Minify()` method for an `Air` instance for minifies a
	// content by a MIME type.
	Minifier interface {
		// Init initializes the `Minifier`. It will be called in the `Air#Serve()`.
		Init() error

		// Minify minifies the b by the mimeType.
		Minify(mimeType string, b []byte) ([]byte, error)
	}

	// minifier implements the `Minifier`.
	minifier struct {
		pngEncoder *png.Encoder
		m          *minify.M
	}
)

// newMinifier returns a pointer of a new instance of the `minifier`.
func newMinifier() *minifier {
	return &minifier{
		pngEncoder: &png.Encoder{
			CompressionLevel: png.BestCompression,
		},
		m: minify.New(),
	}
}

// Init implements the `Minifier#Init()`.
func (m *minifier) Init() error {
	m.m.Add(MIMETextHTML, &html.Minifier{})

	m.m.Add(MIMETextCSS, &css.Minifier{
		Decimals: -1,
	})

	m.m.Add(MIMETextJavaScript, &js.Minifier{})

	m.m.Add(MIMEApplicationJSON, &json.Minifier{})

	m.m.Add(MIMETextXML, &xml.Minifier{})

	m.m.Add(MIMEImageSVGXML, &svg.Minifier{
		Decimals: -1,
	})

	return nil
}

// Minify implements the `Minifier#Minify()`.
func (m *minifier) Minify(mimeType string, b []byte) ([]byte, error) {
	switch mimeType {
	case MIMEImageJPEG:
		return m.minifyJPEG(b)
	case MIMEImagePNG:
		return m.minifyPNG(b)
	case MIMETextHTML, MIMETextCSS, MIMETextJavaScript,
		MIMEApplicationJSON, MIMETextXML, MIMEImageSVGXML:
		return m.minifyOthers(mimeType, b)
	}
	return nil, errors.New("unsupported mime type")
}

// minifyJPEG minifies the b by using the "image/jpeg".
func (m *minifier) minifyJPEG(b []byte) ([]byte, error) {
	img, err := jpeg.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if err := jpeg.Encode(buf, img, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// minifyPNG minifies the b by using the "image/png".
func (m *minifier) minifyPNG(b []byte) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if err := m.pngEncoder.Encode(buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// minifyOthers minifies the b by the mimeType by using the "github.com/tdewolff/minify".
func (m *minifier) minifyOthers(mimeType string, b []byte) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := m.m.Minify(mimeType, buf, bytes.NewReader(b)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
