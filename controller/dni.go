package controller

import (
	"fmt"
	"github.com/paulantezana/review/models"
	"github.com/paulantezana/review/provider"
	"strings"
	"time"
)

func Dni(dni string) (models.Student, error) {
	student := models.Student{}
	student.DNI = dni

	// Query essalud
	container, err := provider.EssaludQuery(dni)

	// Validate
	if err == nil && len(container.DatosPerson) >= 1 {
		first := strings.ToUpper(container.DatosPerson[0].ApellidoPaterno)
		last := strings.ToUpper(container.DatosPerson[0].ApellidoMaterno)
		name := strings.Title(strings.ToLower(container.DatosPerson[0].Nombres))
		student.FullName = fmt.Sprintf("%s %s, %s", first, last, name)
		if container.DatosPerson[0].Sexo == "2" {
			student.Gender = "M"
		} else {
			student.Gender = "F"
		}
		layout := "02/01/2006"
		t, _ := time.Parse(layout, container.DatosPerson[0].FechaNacimiento)
		student.BirthDate = t
	} else {
		// Query jne
		data, _ := provider.JneQuery(student.DNI)
		student.FullName = fmt.Sprintf("%s %s, %s", strings.ToUpper(data.LastName), strings.ToUpper(data.FirstName), strings.Title(data.Name))
	}

	// Response data
	return student, nil
}
