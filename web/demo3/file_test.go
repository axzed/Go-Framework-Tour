package web

import (
	"testing"
)

func TestFileDownloader_Handle(t *testing.T) {
	s := NewHTTPServer()
	s.Get("/download", (&FileDownloader{
		// 下载的文件所在目录
		Dir: "./testdata/download",
	}).Handle())
	// 在浏览器里面输入 localhost:8081/download?file=test.txt
	s.Start(":8081")
}

func TestStaticResourceHandler_Handle(t *testing.T) {
	s := NewHTTPServer()
	handler := NewStaticResourceHandler(199, WithStaticResourceDir("./testdata/img"))
	s.Get("/img/:file", handler.Handle)
	// 在浏览器里面输入 localhost:8081/download?file=test.txt
	s.Start(":8081")
}
