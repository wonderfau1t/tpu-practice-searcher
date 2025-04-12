package models

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

// Статус (хранятся все статусы, как пользователей, так и компаний (пока на модерации))
type Status struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

// Роли
type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`

	Users []User `gorm:"foreignKey:RoleID"`
}

// Пользователь
type User struct {
	ID          int64          `gorm:"primaryKey"`
	Username    string         `gorm:"unique"`
	PhoneNumber sql.NullString `gorm:"unique"`
	CreatedAt   time.Time
	Description sql.NullString
	StatusID    uint
	RoleID      uint

	Status Status
	Role   Role
}

// Школы
type School struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	ShortName string

	Courses []Course `gorm:"foreignKey:SchoolID"`
}

// Направления
type Course struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	SchoolID uint   `json:"-"`

	School    School    `json:"-"`
	Vacancies []Vacancy `gorm:"many2many:vacancy_courses" json:"-"`
}

// Модератор (человек от ТПУ), который ответственен за направления школы
type Moderator struct {
	UserID   int64 `gorm:"primaryKey"`
	SchoolID uint  `gorm:"primaryKey"`

	User   User   `gorm:"foreignKey:UserID"`
	School School `gorm:"foreignKey:SchoolID"`
}

// Компания
type Company struct {
	gorm.Model
	Name         string `gorm:"unique"`
	Description  sql.NullString
	Link         sql.NullString
	StatusID     uint
	HeadHrID     int64
	AcceptedByID *int64

	Status     Status `gorm:"foreignKey:StatusID"`
	HeadHr     User   `gorm:"foreignKey:HeadHrID"`
	AcceptedBy User   `gorm:"foreignKey:AcceptedByID"`
}

// HR каждой компании
type HrManager struct {
	UserID    int64 `gorm:"primaryKey"`
	CompanyID uint  `gorm:"primaryKey"`

	User    User    `gorm:"foreignKey:UserID"`
	Company Company `gorm:"foreignKey:CompanyID"`
}

// Формат проведения вакансии
type Format struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique" json:"name"`

	Vacancies []Vacancy `gorm:"foreignKey:FormatID" json:"-"`
}

// Категории (IT, нефтегаз и тд)
type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique" json:"name"`

	Vacancies []Vacancy `gorm:"foreignKey:CategoryID" json:"-"`
}

type PaymentForAccommodation struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique" json:"name"`
}

type FarePayment struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique" json:"name"`
}

// Вакансия
type Vacancy struct {
	gorm.Model
	Name                           string
	CompanyID                      uint
	HrID                           int64
	StatusID                       uint
	FormatID                       uint
	CategoryID                     uint
	Courses                        []Course           `gorm:"many2many:vacancy_courses"`
	Description                    VacancyDescription `gorm:"foreignKey:VacancyID"`
	Keywords                       []VacancyKeywords  `gorm:"foreignKey:VacancyID"`
	NumberOfResponses              int                `gorm:"default:0"`
	DeadlineAt                     string
	PaymentForAccommodationID      uint
	PaymentForAccommodationDetails sql.NullString
	FarePaymentID                  uint
	FarePaymentDetails             sql.NullString

	Company                 Company                 `gorm:"foreignKey:CompanyID"`
	Hr                      User                    `gorm:"foreignKey:HrID"`
	Status                  Status                  `gorm:"foreignKey:StatusID"`
	Format                  Format                  `gorm:"foreignKey:FormatID"`
	Category                Category                `gorm:"foreignKey:CategoryID"`
	PaymentForAccommodation PaymentForAccommodation `gorm:"foreignKey:PaymentForAccommodationID"`
	FarePayment             FarePayment             `gorm:"foreignKey:FarePaymentID"`
}

// Подробное описание вакансии
type VacancyDescription struct {
	VacancyID      uint `gorm:"primaryKey"`
	Workplace      sql.NullString
	Position       sql.NullString
	Salary         sql.NullInt64
	Requirements   sql.NullString
	Food           sql.NullString
	Conditions     sql.NullString
	AdditionalInfo sql.NullString
}

type VacancyKeywords struct {
	VacancyID uint   `gorm:"primaryKey"`
	Keyword   string `gorm:"primaryKey"`
}
