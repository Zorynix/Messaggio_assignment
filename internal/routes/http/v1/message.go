package v1

import (
	routeerrs "messagio_testsuite/internal/routes/http/v1/route_errors"
	"messagio_testsuite/internal/service"
	serviceerrs "messagio_testsuite/internal/service/service_errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type messageRoutes struct {
	messageService service.Message
}

func NewMessageRoutes(g *echo.Group, messageService service.Message) {
	r := &messageRoutes{
		messageService: messageService,
	}

	g.POST("/create", r.create)
	g.GET("/messages", r.getAll)
	g.GET("/messages/:id", r.getByID)
	g.GET("/messages/stats", r.getStats)
	g.PUT("/messages/:id/process", r.markAsProcessed)
}

func (r *messageRoutes) create(c echo.Context) error {
	type request struct {
		Message string `json:"message" validate:"required"`
	}
	var req request
	if err := c.Bind(&req); err != nil {
		routeerrs.NewErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(&req); err != nil {
		routeerrs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := r.messageService.CreateMessage(c.Request().Context(), req.Message)
	if err != nil {
		routeerrs.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}
	type response struct {
		Id uuid.UUID `json:"id"`
	}

	return c.JSON(http.StatusCreated, response{
		Id: id,
	})
}

func (r *messageRoutes) getAll(c echo.Context) error {
	messages, err := r.messageService.GetMessages(c.Request().Context())
	if err != nil {
		routeerrs.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, messages)
}

func (r *messageRoutes) getByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		routeerrs.NewErrorResponse(c, http.StatusBadRequest, "invalid id format")
		return err
	}

	message, err := r.messageService.GetMessageById(c.Request().Context(), id)
	if err != nil {
		if err == serviceerrs.ErrMessageNotFound {
			routeerrs.NewErrorResponse(c, http.StatusNotFound, err.Error())
			return err
		}
		routeerrs.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, message)
}

func (r *messageRoutes) getStats(c echo.Context) error {
	count, err := r.messageService.GetProcessedMessagesStats(c.Request().Context())
	if err != nil {
		routeerrs.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		ProcessedMessages int `json:"processed_messages"`
	}

	return c.JSON(http.StatusOK, response{
		ProcessedMessages: count,
	})
}

func (r *messageRoutes) markAsProcessed(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		routeerrs.NewErrorResponse(c, http.StatusBadRequest, "invalid id format")
		return err
	}

	err = r.messageService.MarkMessageAsProcessed(c.Request().Context(), id)
	if err != nil {
		if err == serviceerrs.ErrMessageNotFound {
			routeerrs.NewErrorResponse(c, http.StatusNotFound, err.Error())
			return err
		}
		routeerrs.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Message string `json:"message"`
	}

	return c.JSON(http.StatusOK, response{
		Message: "Message successfully marked as processed",
	})
}
