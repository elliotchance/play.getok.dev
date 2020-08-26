package main

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/elliotchance/ok/compiler"
	"github.com/elliotchance/ok/parser"
	"github.com/elliotchance/ok/vm"
)

func errorResponse(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       err.Error(),
	}, nil
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	resp := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       buf.String(),
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
