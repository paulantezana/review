package models

import "time"

// Teacher struct
type Teacher struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DNI       string    `json:"dni" gorm:"type:varchar(15); not null; unique"`
	LastName  string    `json:"last_name"`
	FirstName string    `json:"first_name"`

	BirthDate time.Time `json:"birth_date"`
	Gender    string    `json:"gender"`

	//Country string `json:"country"`
	//Department string `json:"department"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	//CivilStatus string `json:"civil_status"`
	WorkConditions string    `json:"work_conditions"`
	EducationLevel string    `json:"education_level"`
	AdmissionDate  time.Time `json:"admission_date"`
	RetirementDate time.Time `json:"retirement_date"`
	//YearsDecency uint `json:"years_decency"`
	//TeachingMonths uint `json:"teaching_months"`
	Specialty string `json:"specialty"`

	UserID uint `json:"user_id"`

	// Only in REQUEST OR RESPONSE JSON
	Type      string `json:"type" gorm:"-"`
	ProgramID uint   `json:"program_id" gorm:"-"`

	TeacherPrograms []TeacherProgram `json:"teacher_programs"`
}
