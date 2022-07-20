package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
)

var (
	serviceVersion = os.Getenv("SERVICE_VERSION")
	serviceName    = os.Getenv("SERVICE_NAME")
	servicePort    = os.Getenv("SERVICE_PORT")
	cosmosBinding  = os.Getenv("COSMOS_BINDING")

	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Lshortfile)
}

func main() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	Info.Printf("serviceName: %s", serviceName)
	Info.Printf("servicePort: %s", servicePort)
	Info.Printf("cosmosBinding: %s", cosmosBinding)

	port := fmt.Sprintf(":%s", servicePort)

	s, err := daprd.NewService(port)
	if err != nil {
		Error.Printf("%s: failed create the server: %v", serviceName, err)
	}

	// add invocation handler
	if err := s.AddServiceInvocationHandler("/getPhotos", InvocationHandler); err != nil {
		Error.Printf("%s: error adding service invocation handler: %v", serviceName, err)
	}

	// start the service
	if err := s.Start(); err != nil {
		Error.Printf("%s: server failed to start: %v", serviceName, err)
	}
}

func InvocationHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	Info.Printf("%s - ContentType:%s, Verb:%s, QueryString:%s, Data:%+v", serviceName, in.ContentType, in.Verb, in.QueryString, string(in.Data))

	type Query struct {
		Collection string `json:"collection"`
		Album      string `json:"album"`
	}

	q := &Query{}
	json.Unmarshal(in.Data, &q)

	var jsonStr = `{
		"filter": {
			"AND": [
				{
					"EQ": { "collection": "{{.collection}}" }
				},
				{
					"EQ": { "album": "{{.album}}" }
				}
			]
		}
	}`

	t, _ := template.New("text").Parse(jsonStr)
	o := bytes.Buffer{}
	t.Execute(&o, map[string]interface{}{"album": q.Album, "collection": q.Collection})

	Info.Printf("querying database: %s", serviceName)
	r, err := queryData(ctx, o.String())
	bytArr, err := json.Marshal(r)

	Info.Printf("query result: %s", r.Results)

	out = &common.Content{
		Data:        bytArr,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}

	return out, nil
}

func queryData(ctx context.Context, query string) (r *dapr.QueryResponse, err error) {
	client, err := dapr.NewClient()
	if err != nil {
		Error.Printf("%s: failed to create client: %v", serviceName, err)
		return nil, err
	}

	Info.Printf("query: %s", query)

	r, err = client.QueryStateAlpha1(ctx, cosmosBinding, query, nil)
	if err != nil {
		Error.Printf("%s: failed to query state: %v", serviceName, err)
		return nil, err
	}

	return r, nil
}
