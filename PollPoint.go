package onvif

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/beevik/etree"
	"github.com/use-go/onvif/event"
	"github.com/use-go/onvif/gosoap"
	"github.com/use-go/onvif/networking"
	"github.com/use-go/onvif/sdk"
)

type PullPoint struct {
	address         string
	terminationTime time.Time
	httpClient      *http.Client
	username        string
	password        string
}

type PullPointParams struct {
	Address         string
	TerminationTime string
	Username        string
	Password        string
	HTTPTimeout     time.Duration
}

func NewPullPoint(params PullPointParams) (*PullPoint, error) {
	// parse time from 2023-10-30T05:27:55Z
	parsedTime, err := time.Parse(time.RFC3339, params.TerminationTime)
	if err != nil {
		return nil, err
	}

	if params.HTTPTimeout == 0 {
		params.HTTPTimeout = time.Second * 60
	}

	// setup timeouts for client
	httpClient := &http.Client{
		Transport: &http.Transport{},
		Timeout:   params.HTTPTimeout,
	}

	return &PullPoint{
		address:         params.Address,
		terminationTime: parsedTime,
		httpClient:      httpClient,
		username:        params.Username,
		password:        params.Password,
	}, nil
}

func (pullPoint *PullPoint) GetAddress() string {
	return pullPoint.address
}

func (pullPoint *PullPoint) CallMethod(method interface{}) (*http.Response, error) {
	output, err := xml.MarshalIndent(method, "  ", "    ")
	if err != nil {
		return nil, err
	}

	soap, err := pullPoint.buildMethodSOAP(string(output))
	if err != nil {
		return nil, err
	}

	soap.AddWSSecurity(pullPoint.username, pullPoint.password)

	endpoint := pullPoint.GetAddress()

	return networking.SendSoap(context.TODO(), pullPoint.httpClient, endpoint, soap.String())
}

func (pullPoint *PullPoint) PullMessages(ctx context.Context, request event.PullMessages) (event.PullMessagesResponse, error) {
	type Envelope struct {
		Header struct{}
		Body   struct {
			PullMessagesResponse event.PullMessagesResponse
		}
	}
	var reply Envelope
	if httpReply, err := pullPoint.CallMethod(request); err != nil {
		return reply.Body.PullMessagesResponse, fmt.Errorf("call: %w", err)
	} else {
		err = sdk.ReadAndParse(ctx, httpReply, &reply, "PullMessages")
		if err != nil {
			return reply.Body.PullMessagesResponse, fmt.Errorf("parse: %w", err)
		}
		return reply.Body.PullMessagesResponse, nil
	}
}

func (pullPoint *PullPoint) buildMethodSOAP(msg string) (gosoap.SoapMessage, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(msg); err != nil {
		return "", err
	}
	element := doc.Root()

	soap := gosoap.NewEmptySOAP()
	soap.AddBodyContent(element)

	return soap, nil
}
