package handlers

import (
	"fmt"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/satnamSandhu2001/stackjet/internal/dto"
	"github.com/satnamSandhu2001/stackjet/internal/services"
	"github.com/satnamSandhu2001/stackjet/pkg"
	"github.com/satnamSandhu2001/stackjet/pkg/API"
	"github.com/satnamSandhu2001/stackjet/pkg/commands"
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
	if !slices.Contains(pkg.Config().VALID_STACKS, body.Type) {
		API.Error(c, "Invalid stack type")
		return

	}
	if body.Port < 1024 {
		API.Error(c, "Invalid port number")
		return
	}
	if isFree := commands.IsPortFree(body.Port); !isFree {
		API.Error(c, fmt.Sprintf("Port %d is in use or blocked", body.Port))
		return
	}

	if err := h.service.CreateStack(c.Request.Context(), &body); err != nil {
		API.InternalServerError(c, "Failed to create stack", err)
		return
	}

	API.Success(c, "Stack created successfully", nil)

}
