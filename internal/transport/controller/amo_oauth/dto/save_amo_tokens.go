package amo_oauth

type SaveAmoTokensQuery struct {
	Code    string `form:"code" binding:"required"`
	Referer string `form:"referer" binding:"required"`
}
