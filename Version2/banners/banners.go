package banners

import (
	ascii "ascii-art-web-stylize/rendering"
	"bytes"
	"crypto/md5"
	_ "embed"
	"errors"
	"fmt"
	"io"
)

func MD5Sum(r io.Reader) (string, error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5.Sum(bytes)), nil
}

//go:embed standard.txt
var standard []byte

//go:embed shadow.txt
var shadow []byte

//go:embed thinkertoy.txt
var thinkertoy []byte

var ErrBanner = errors.New("banner intergrity is compromised")

type banner struct {
	md5Sum string
	body   []byte
}

func (b *banner) checkMD5Hash() error {
	md5Sum, err := MD5Sum(bytes.NewReader(b.body))
	if err != nil {
		return err
	}
	if md5Sum != b.md5Sum {
		return ErrBanner
	}
	return nil
}

var banners = map[string]banner{
	"standard":   {"ac85e83127e49ec42487f272d9b9db8b", standard},
	"shadow":     {"a49d5fcb0d5c59b2e77674aa3ab8bbb1", shadow},
	"thinkertoy": {"86d9947457f6a41a18cb98427e314ff8", thinkertoy},
}

func ParseTemplates() (map[string]*ascii.Template, error) {
	templates := make(map[string]*ascii.Template, len(banners))

	for name, banner := range banners {
		if err := banner.checkMD5Hash(); err != nil {
			return nil, fmt.Errorf("banner %q %w", name, err)
		}

		t, err := ascii.NewTemplate(bytes.NewReader(banner.body))
		if err != nil {
			return nil, fmt.Errorf("banner %q %w", name, err)
		}

		templates[name] = t
	}

	return templates, nil
}
