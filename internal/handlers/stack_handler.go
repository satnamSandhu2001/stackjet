package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/satnamSandhu2001/stackjet/internal/cli/stack"
	"github.com/satnamSandhu2001/stackjet/internal/dto"
	"github.com/satnamSandhu2001/stackjet/internal/services"
	"github.com/satnamSandhu2001/stackjet/pkg"
	"github.com/satnamSandhu2001/stackjet/pkg/API"
	"github.com/satnamSandhu2001/stackjet/pkg/helpers"
)

type StackHandler struct {
	service services.StackService
}

func NewStackHandler(service *services.StackService) *StackHandler {
	return &StackHandler{
		service: *service,
	}
}

func (h *StackHandler) CreateNewStack(c *gin.Context) {
	var body dto.Stack_CreateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		errors := pkg.TagValidationErrors(err, &body)
		API.ValidationsErrors(c, errors)
		return
	}

	// validate start command
	if body.Commands.Start != "" {
		switch body.Type {
		case "nodejs":
			if err := helpers.ValidateNodeStartCommand(body.Commands.Start); err != nil {
				API.Error(c, err.Error())
				return
			}
		}
	}
	if body.Commands.Start == "" {
		switch body.Type {
		case "nodejs":
			body.Commands.Start = "npm start"
		}
	}

	// handle stream
	stream := API.NewStreamWriter(c)
	if stream == nil {
		API.InternalServerError(c, "Streaming not supported", nil)
		return
	}

	logger := helpers.NewMultiLogger(stream)
	err := stack.CreateNewStack(logger, c.Request.Context(), h.service, &body)

	if err != nil {
		stream.WriteString("__ERROR__:" + err.Error() + "\n")
		return
	}

	stream.WriteString("__SUCCESS__\n")
}
