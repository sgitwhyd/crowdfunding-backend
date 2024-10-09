package transaction

type GetTransactionByIDInput struct {
	ID int `uri:"id" binding:"required"`
}