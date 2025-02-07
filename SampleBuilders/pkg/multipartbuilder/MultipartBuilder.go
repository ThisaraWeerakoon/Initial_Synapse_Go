package multipartbuilder

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

// FilePart represents a file's data.
type FilePart struct {
	Filename string
	Content  []byte
}

// MultipartData holds form fields and file parts.
type MultipartData struct {
	Fields map[string][]string
	Files  map[string][]FilePart
}

// MultipartBuilder defines methods for processing multipart/form-data.
type MultipartBuilder interface {
	// Parse extracts multipart data from an HTTP request.
	Parse(r *http.Request) (*MultipartData, error)
	// Build creates a multipart payload from internal data.
	Build(data *MultipartData) (contentType string, body []byte, err error)
}

// DefaultMultipartBuilder is an implementation using Go’s standard libraries.
type DefaultMultipartBuilder struct{}

// Parse uses r.ParseMultipartForm to extract form values and converts file headers
// into our FilePart type.
func (b *DefaultMultipartBuilder) Parse(r *http.Request) (*MultipartData, error) {
	// Parse the multipart form with a 10MB memory threshold.
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, err
	}
	md := &MultipartData{
		Fields: r.MultipartForm.Value,
		Files:  make(map[string][]FilePart),
	}
	// Convert each *multipart.FileHeader into our FilePart.
	for key, fhs := range r.MultipartForm.File {
		for _, fh := range fhs {
			file, err := fh.Open()
			if err != nil {
				return nil, err
			}
			content, err := ioutil.ReadAll(file)
			file.Close()
			if err != nil {
				return nil, err
			}
			md.Files[key] = append(md.Files[key], FilePart{
				Filename: fh.Filename,
				Content:  content,
			})
		}
	}
	return md, nil
}

// Build constructs a multipart/form-data payload using a multipart.Writer.
func (b *DefaultMultipartBuilder) Build(data *MultipartData) (string, []byte, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	// Write form fields.
	for key, vals := range data.Fields {
		for _, val := range vals {
			if err := writer.WriteField(key, val); err != nil {
				return "", nil, err
			}
		}
	}
	// Write file parts.
	for key, parts := range data.Files {
		for _, part := range parts {
			fw, err := writer.CreateFormFile(key, part.Filename)
			if err != nil {
				return "", nil, err
			}
			if _, err := fw.Write(part.Content); err != nil {
				return "", nil, err
			}
		}
	}
	// Close the writer to finalize the boundary.
	if err := writer.Close(); err != nil {
		return "", nil, err
	}
	return writer.FormDataContentType(), buf.Bytes(), nil
}

// simulateIncomingMultipartRequest creates a dummy HTTP request with multipart/form-data.
func simulateIncomingMultipartRequest() *http.Request {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	// Add form fields.
	writer.WriteField("username", "testuser")
	writer.WriteField("email", "test@example.com")
	// Add a file part.
	part, err := writer.CreateFormFile("profile", "profile.txt")
	if err != nil {
		log.Fatalf("Error creating form file: %v", err)
	}
	io.Copy(part, strings.NewReader("This is the file content."))
	writer.Close()
	req, err := http.NewRequest("POST", "http://example.com/upload", &buf)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

