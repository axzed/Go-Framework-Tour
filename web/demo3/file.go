package web

import (
	lru "github.com/hashicorp/golang-lru"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type FileUploader struct {
	FileField string
	// 比如说 DST 是一个目录
	Dst string
	// DstPathFunc 用于计算目标路径
	DstPathFunc func(fh *multipart.FileHeader) string
}

func (f *FileUploader) Handle() HandleFunc {
	return func(ctx *Context) {
		file, header, err := ctx.Req.FormFile(f.FileField)
		if err != nil {
			ctx.RespStatusCode = 500
			ctx.RespData = []byte("上传失败")
			return
		}
		dst, err := os.OpenFile(filepath.Join(f.DstPathFunc(header), header.Filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		if err != nil {
			ctx.RespStatusCode = 500
			ctx.RespData = []byte("上传失败")
			return
		}
		io.CopyBuffer(dst, file, nil)
	}
}

type FileDownloader struct {
	// 设计各种参数
	Dir string
}

func (f *FileDownloader) Handle() HandleFunc {

	// 你可以在这里
	return func(ctx *Context) {
		// file 也可以是多段的呀 /a/b/c.txt
		file, err := ctx.QueryValue("file")
		if err != nil {
			ctx.RespStatusCode = 500
			ctx.RespData = []byte("文件找不到")
			return
		}
		// file 可能是一些乱七八糟的东西
		// file =///////abc.txt
		path := filepath.Join(f.Dir, filepath.Clean(file))
		// 从完整路径里面拿到文件名
		fn := filepath.Base(path)
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")
		// 文件下载本质上，就是把文件写到这里
		http.ServeFile(ctx.Resp, ctx.Req, path)
	}
}

type StaticResourceHandlerOption func(handler *StaticResourceHandler)

func WithStaticResourceDir(dir string) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.dir = dir
	}
}

func WithMoreExtension(extMap map[string]string) StaticResourceHandlerOption {
	return func(h *StaticResourceHandler) {
		for ext, contentType := range extMap {
			h.extensionContentTypeMap[ext] = contentType
		}
	}
}

func WithMaxFileSize(fileSize int) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.maxFileSize = fileSize
	}
}

type StaticResourceHandler struct {
	dir                     string
	extensionContentTypeMap map[string]string

	// 我要缓存，但不是缓存所有的，并且要控制住内存消耗
	cache       *lru.Cache
	maxFileSize int
}

func NewStaticResourceHandler(maxSize int, opts ...StaticResourceHandlerOption) *StaticResourceHandler {
	c, _ := lru.New(maxSize)
	res := &StaticResourceHandler{
		cache:       c,
		maxFileSize: 50 * 1024 * 1024,
		extensionContentTypeMap: map[string]string{
			// 这里根据自己的需要不断添加
			"jpeg": "image/jpeg",
			"jpe":  "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"pdf":  "image/pdf",
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (s *StaticResourceHandler) Handler() HandleFunc {
	if s.maxFileSize <= 0 {
		s.maxFileSize = 50 * 1024 * 1024
	}
	return func(ctx *Context) {

	}
}

// Handle 假设我们请求 /static/come_on_baby.jpg
func (s *StaticResourceHandler) Handle(ctx *Context) {
	file, ok := ctx.PathParams["file"]
	if !ok {
		ctx.RespData = []byte("文件不存在")
		ctx.RespStatusCode = 500
		return
	}
	header := ctx.Resp.Header()

	itm, ok := s.cache.Get(file)
	if ok {
		ctx.RespStatusCode = 200
		ctx.RespData = itm.(*cacheItem).data
		header.Set("Content-Type", itm.(*cacheItem).contentType)
		return
	}
	path := filepath.Join(s.dir, filepath.Clean(file))
	f, err := os.Open(path)
	if err != nil {
		ctx.RespData = []byte("文件不存在")
		ctx.RespStatusCode = 500
		return
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		ctx.RespData = []byte("文件不存在")
		ctx.RespStatusCode = 500
		return
	}
	newItm := &cacheItem{
		data:        data,
		contentType: s.extensionContentTypeMap[filepath.Ext(file)],
	}
	// 大于 50M 我就不缓存
	if len(data) <= s.maxFileSize {
		s.cache.Add(file, newItm)
	}

	header.Set("Content-Type", newItm.contentType)
	ctx.RespStatusCode = 200
	ctx.RespData = data
}

type cacheItem struct {
	data        []byte
	contentType string
}
