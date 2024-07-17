package v1

import (
	routeerrs "messagio_testsuite/internal/routes/http/v1/route_errors"
	"messagio_testsuite/internal/service"
	serviceerrs "messagio_testsuite/internal/service/service_errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type messageRoutes struct {
	messageService service.Message
}

func newMessageRoutes(g *echo.Group, messageService service.Message) {
	r := &messageRoutes{
		messageService: messageService,
	}

	g.POST("/create", r.create)
	g.GET("/messages", r.getAll)
	g.GET("/messages/:id", r.getByID)
}

func (r *messageRoutes) create(c echo.Context) error {
	id, err := r.messageService.CreateMessage(c.Request().Context())
	if err != nil {
		if err == serviceerrs.ErrMessageAlreadyExists {
			routeerrs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		routeerrs.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}
	type response struct {
		Id int `json:"id"`
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
	id, err := strconv.Atoi(idStr)
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
