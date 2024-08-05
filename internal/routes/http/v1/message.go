package v1

import (
	"errors"
	routeerrs "messagio_testsuite/internal/routes/http/v1/route_errors"
	"messagio_testsuite/internal/service"
	serviceerrs "messagio_testsuite/internal/service/service_errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MessageRoutes struct {
	MessageService service.Message
}

func NewMessageRoutes(g *echo.Group, MessageService service.Message) {
	r := &MessageRoutes{
		MessageService: MessageService,
	}

	g.POST("/create", r.Create)
	g.GET("/messages", r.GetAll)
	g.GET("/messages/:id", r.GetByID)
	g.GET("/messages/stats", r.GetStats)
	g.PUT("/messages/:id/process", r.MarkAsProcessed)
}

func (r *MessageRoutes) Create(c echo.Context) error {
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

	id, err := r.MessageService.CreateMessage(c.Request().Context(), req.Message)
	if err != nil {
		if errors.Is(err, serviceerrs.ErrCannotCreateMessage) {
			routeerrs.NewErrorResponse(c, http.StatusInternalServerError, "failed to create message")
		} else if errors.Is(err, serviceerrs.ErrCannotProduceMessage) {
			routeerrs.NewErrorResponse(c, http.StatusInternalServerError, "failed to produce message to Kafka")
		} else {
			routeerrs.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return err
	}

	type response struct {
		Id uuid.UUID `json:"id"`
	}

	return c.JSON(http.StatusCreated, response{
		Id: id,
	})
}

func (r *MessageRoutes) GetAll(c echo.Context) error {
	messages, err := r.MessageService.GetMessages(c.Request().Context())
	if err != nil {
		routeerrs.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, messages)
}

func (r *MessageRoutes) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		routeerrs.NewErrorResponse(c, http.StatusBadRequest, "invalid id format")
		return err
	}

	message, err := r.MessageService.GetMessageById(c.Request().Context(), id)
	if err != nil {
		if err == serviceerrs.ErrMessageNotFound {
			routeerrs.NewErrorResponse(c, http.StatusNotFound, "message not found")
			return err
		}
		routeerrs.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, message)
}

func (r *MessageRoutes) GetStats(c echo.Context) error {
	count, err := r.MessageService.GetProcessedMessagesStats(c.Request().Context())
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

func (r *MessageRoutes) MarkAsProcessed(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		routeerrs.NewErrorResponse(c, http.StatusBadRequest, "invalid id format")
		return err
	}

	err = r.MessageService.MarkMessageAsProcessed(c.Request().Context(), id)
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
