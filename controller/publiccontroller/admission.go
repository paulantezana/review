package publiccontroller

import (
    "github.com/labstack/echo"
    "github.com/paulantezana/review/utilities"
    "net/http"
)

func GetAdmissionExamResults(c echo.Context) error {
    // Get data request
    request := utilities.Request{}
    if err := c.Bind(&request); err != nil {
        return err
    }

    // Required params
    if request.SubsidiaryID == 0 {
        c.JSON(http.StatusOK,utilities.Response{Message: "EL parametro subsidiary_id es obligatorio"})
    }

    // Required params
    if request.SubsidiaryID == 0 {
        c.JSON(http.StatusOK,utilities.Response{Message: "EL parametro subsidiary_id es obligatorio"})
    }

    // Response data
    return c.JSON(http.StatusOK,utilities.Response{
        Success: true,
        Data: request,
    })
}