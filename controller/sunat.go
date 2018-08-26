package controller

import (
	"github.com/labstack/echo"
	"github.com/paulantezana/review/utilities"
	"io/ioutil"
	"net/http"
	"strings"
)

func Sunat(c echo.Context) error {
	url := "https://ruc.com.pe/api/v1/ruc"
	payload := strings.NewReader("{\r\n  \"token\": \"989b7941-f931-4983-988e-fa4f2e2be451-edc33212-f1cd-4c8e-9369-3fa671ff5155\",\r\n  \"ruc\": \"10739757555\"\r\n}")

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	return c.JSON(http.StatusOK, utilities.Response{
		Success: true,
		Data:    body,
	})
}
