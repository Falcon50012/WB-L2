package wgetm

import (
	"crypto/md5"
	"fmt"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func urlToLocalPath(rawURL, contentType string) string {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {

		h := fmt.Sprintf("%x", md5.Sum([]byte(rawURL)))
		return filepath.Join("unknown", h+".bin")
	}

	host := sanitizeName(u.Host)
	p := u.Path

	switch {
	case p == "" || p == "/":
		p = "/index.html"
	case strings.HasSuffix(p, "/"):
		p = p + "index.html"
	}

	segments := strings.Split(p, "/")
	for i, seg := range segments {
		segments[i] = sanitizeName(seg)
	}
	localPath := strings.Join(segments, string(filepath.Separator))

	base := filepath.Base(localPath)
	if ext := filepath.Ext(base); len(ext) <= 1 {
		if e := extForContentType(contentType); e != "" {
			localPath += e
		} else {
			localPath += ".html"
		}
	}

	if u.RawQuery != "" {
		h := fmt.Sprintf("%x", md5.Sum([]byte(u.RawQuery)))[:6]
		ext := filepath.Ext(localPath)
		localPath = strings.TrimSuffix(localPath, ext) + "_" + h + ext
	}

	return filepath.Join(host, localPath)
}

func saveFile(outputDir, localPath string, data []byte) error {
	abs := filepath.Join(outputDir, localPath)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return fmt.Errorf("make dir error %q: %w", abs, err)
	}
	if err := os.WriteFile(abs, data, 0o644); err != nil {
		return fmt.Errorf("write file error %q: %w", abs, err)
	}
	return nil
}

func sanitizeName(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case ':', '*', '?', '"', '<', '>', '|', '\\':
			b.WriteByte('_')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func extForContentType(ct string) string {
	if ct == "" {
		return ""
	}
	mt, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return ""
	}
	exts := map[string]string{
		"text/html":                ".html",
		"application/xhtml+xml":    ".html",
		"text/css":                 ".css",
		"application/javascript":   ".js",
		"text/javascript":          ".js",
		"application/json":         ".json",
		"image/jpeg":               ".jpg",
		"image/png":                ".png",
		"image/gif":                ".gif",
		"image/svg+xml":            ".svg",
		"image/webp":               ".webp",
		"image/x-icon":             ".ico",
		"image/vnd.microsoft.icon": ".ico",
		"font/woff":                ".woff",
		"font/woff2":               ".woff2",
		"font/ttf":                 ".ttf",
		"application/xml":          ".xml",
		"text/xml":                 ".xml",
		"application/pdf":          ".pdf",
	}
	return exts[mt]
}
