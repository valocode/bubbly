package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/valocode/bubbly/env"
)

const (
	defaultFormName = "file"
	formatJSON      = "json"
	formatZip       = "zip"
	contentTypeZip  = "application/zip"
)

type RunError struct {
	Code    int
	Message string
}

type WorkerRun struct {
	RemoteInput RemoteInput
	Name        string // name of the `run` resource
}

type RemoteInput struct {
	Data     []byte // the raw input data
	Filename string // name and path of the file
	Format   string // json or zip
}

func ProcessRunData(bCtx *env.BubblyContext, c echo.Context) (WorkerRun, error) {
	contentType := c.Request().Header.Get(echo.HeaderContentType)
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		// if no media has been provided, assume that this is a simple
		// remote run trigger without input data
		if err.Error() == "mime: no media type" {
			return WorkerRun{}, nil
		}
		return WorkerRun{}, err
	}

	switch mediaType {
	case echo.MIMEMultipartForm:
		bCtx.Logger.Debug().Str("type", echo.MIMEMultipartForm).Msg("identified valid content")

		form, err := c.MultipartForm()

		if err != nil {
			return WorkerRun{}, err
		}

		wr, err := handleMultipartForm(bCtx, form)

		if err != nil {
			return WorkerRun{}, err
		}
		return wr, nil
	default:
		return WorkerRun{}, http.ErrNotMultipart
	}
}

// TODO: support multi-file upload
func handleMultipartForm(bCtx *env.BubblyContext, form *multipart.Form) (WorkerRun, error) {
	files := form.File[defaultFormName]

	switch len(files) {
	case 0:
		// an unsupported form part has been used when sending the data
		if len(form.File) != 0 {
			return WorkerRun{}, errors.New(fmt.Sprintf(`invalid form name: use "%s"`, defaultFormName))
		}
		return WorkerRun{}, errors.New(fmt.Sprintf(`no files provided`))
	case 1:
		bCtx.Logger.Debug().Str("name", files[0].Filename).Msg("identified valid content")
	default:
		return WorkerRun{}, errors.New("unsupported: more than one file included in POST")
	}

	var wr WorkerRun

	file := files[0]

	contentType := file.Header.Get("Content-Type")
	switch contentType {
	case echo.MIMEApplicationJSON:
		bCtx.Logger.Debug().Str("type", echo.MIMEApplicationJSON).Msg("identified valid content")

		jsonBytes, err := handleJSONFile(*file)
		if err != nil {
			return WorkerRun{}, err
		}

		wr = WorkerRun{
			RemoteInput: RemoteInput{
				Filename: file.Filename,
				Data:     jsonBytes,
				Format:   formatJSON,
			},
		}

	case contentTypeZip:
		bCtx.Logger.Debug().Str("type", contentTypeZip).Msg("identified valid content")
		// Source
		src, err := file.Open()

		if err != nil {
			return WorkerRun{}, err
		}
		defer src.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, src); err != nil {
			return WorkerRun{}, fmt.Errorf("failed to parse .zip file into bytes: %w", err)
		}

		bCtx.Logger.Debug().Bytes("zip", buf.Bytes()).Msg("successfully parsed .zip file into bytes")

		wr = WorkerRun{
			RemoteInput: RemoteInput{
				Filename: file.Filename,
				Data:     buf.Bytes(),
				Format:   formatZip,
			},
		}

	default:
		rString := fmt.Sprintf(`Unsupported content type "%s". Options: "%s", "%s"`, contentType, echo.MIMEApplicationJSON, contentTypeZip)
		return WorkerRun{}, errors.New(rString)
	}

	return wr, nil
}

func handleJSONFile(file multipart.FileHeader) ([]byte, error) {
	// Source
	src, err := file.Open()

	if err != nil {
		return nil, err
	}
	defer src.Close()
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
