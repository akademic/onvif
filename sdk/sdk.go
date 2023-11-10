package sdk

import (
	"context"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/juju/errors"
)

func ReadAndParse(ctx context.Context, httpReply *http.Response, reply interface{}, tag string) error {
	defaultLogger.Debug("action=%s msg=\"%s\" status=%s", tag, httpReply.Status, httpReply.StatusCode)

	b, err := ioutil.ReadAll(httpReply.Body)
	if err != nil {
		return errors.Annotate(err, "read")
	}

	httpReply.Body.Close()

	err = xml.Unmarshal(b, reply)
	return errors.Annotate(err, "decode")
}
