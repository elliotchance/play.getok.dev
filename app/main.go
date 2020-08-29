package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/elliotchance/ok/compiler"
	"github.com/elliotchance/ok/parser"
	"github.com/elliotchance/ok/vm"
	"github.com/gobuffalo/packr/v2"
)

func errorResponse(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       err.Error(),
	}, nil
}

func Run(request events.APIGatewayProxyRequest) (resp events.APIGatewayProxyResponse, _ error) {
	defer func() {
		if r := recover(); r != nil {
			resp = events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       fmt.Sprintf("%v", r),
			}
		}
	}()

	p := parser.ParseString(request.Body, "main.ok")
	if errs := p.Errors(); len(errs) > 0 {
		return errorResponse(errors.New(errs.String()))
	}

	pkg, err := compiler.CompileFile(p.File, p.Interfaces, p.Constants)
	if err != nil {
		return errorResponse(err)
	}

	m := vm.NewVM(pkg.Funcs, pkg.Tests, pkg.Interfaces, "play")
	buf := bytes.NewBuffer(nil)
	m.Stdout = buf
	err = m.Run()
	if err != nil {
		return errorResponse(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       buf.String(),
	}, nil
}

func Index(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	box := packr.New("pkg", "../static")
	index, err := box.FindString("index.html")
	if err != nil {
		return errorResponse(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       index,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}, nil
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.Path {
	case "/":
		return Index(request)

	case "/run":
		return Run(request)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Body:       request.Path,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
