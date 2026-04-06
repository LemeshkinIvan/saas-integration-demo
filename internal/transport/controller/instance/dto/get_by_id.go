package instance

type GetByIDUri struct {
	ID int64 `uri:"id" binding:"required"`
}
