package handlers

import (
	"io"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/satnamSandhu2001/stackjet/internal/core/stack"
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
	var body dto.Stack_Create_Request
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

	// handle streaming
	logWriter := API.NewSSEWriter(c.Writer)

	err := stack.CreateNewStack(logWriter, c.Request.Context(), h.service, &body)

	if err != nil {
		logWriter.Write([]byte("__ERROR__: " + err.Error()))
	}
	logWriter.Close()
}

func (h *StackHandler) DeployStack(c *gin.Context) {
	var body dto.Stack_Deploy_Request
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		API.Error(c, "Invalid ID format")
		return
	}
	body.ID = id
	if err := c.ShouldBindJSON(&body); err != nil {
		errors := pkg.TagValidationErrors(err, &body)
		API.ValidationsErrors(c, errors)
		return
	}

	// Create log collector
	var logBuf strings.Builder
	sseWriter := API.NewSSEWriter(c.Writer)
	logWriter := io.MultiWriter(sseWriter, &logBuf)

	deploymentID, err := stack.DeployStack(logWriter, c.Request.Context(), h.service, &body)
	if err != nil {
		logWriter.Write([]byte("__ERROR__: " + err.Error()))
	}

	// Save logs to DB
	if deploymentID != 0 {
		h.service.CreateDeploymentLog(c.Request.Context(), &dto.DeploymentLog_Create_Request{DeploymentID: deploymentID, Log: logBuf.String()})
	}
	sseWriter.Close()
}

func (h *StackHandler) ListStacks(c *gin.Context) {
	stacks, err := h.service.GetStackList(c.Request.Context())
	if err != nil {
		API.Error(c, "failed to list stacks")
		return
	}
	API.Success(c, "success", stacks)
}
