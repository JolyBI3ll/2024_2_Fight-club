package domain

import (
	"context"
)

const STARS = 5
const SMILES = 5
const RATE = 10

type Survey struct {
	ID    int    `gorm:"primary_key;auto_increment;column:id" json:"id"`
	Title string `gorm:"type:text;size:3000;column:title" json:"title"`
}

type Question struct {
	ID       int    `gorm:"primary_key;auto_increment;column:id" json:"id"`
	Title    string `gorm:"type:text;size:3000;column:title" json:"title"`
	Type     string `gorm:"type:text;size:3000;column:type" json:"type"`
	SurveyId string `gorm:"column:surveyId;not null" json:"surveyId"`
	Survey   Survey `gorm:"foreignkey:surveyId;references:ID" json:"-"`
}

type Answer struct {
	ID         int      `gorm:"primary_key;auto_increment;column:id" json:"id"`
	QuestionId int      `gorm:"column:questionId" json:"questionId"`
	UserId     string   `gorm:"column:userId" json:"userId"`
	Value      int      `gorm:"column:value" json:"value"`
	Question   Question `gorm:"foreignkey:questionId;references:ID" json:"-"`
}

type SurveyResponse struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Questions []Question `json:"questions"`
}

type PostSurvey struct {
	QuestionId int `json:"questionId"`
	Value      int `json:"value"`
}

//type GetStatictics struct {
//	Avg           float32     `json:"avg"`
//	AnswerNumbers map[int]int `json:"answerNumbers"`
//	QuestionId
//}

type CSATRepository interface {
	GetSurvey(ctx context.Context, surveyId int) (survey SurveyResponse, err error)
	PostSurvey(ctx context.Context, answers []PostSurvey, userId string) (err error)
}
