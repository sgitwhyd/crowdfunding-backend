package transaction

type service struct {
	repository Repositoy
}

type Service interface {
	GetTransactionByCampaignID(input GetTransactionByIDInput) ([]Transaction, error)
	GetTransactionsUser(userID int) ([]Transaction, error)
}

func NewService(repository Repositoy) *service {
	return &service{repository}
}

func (s *service) GetTransactionByCampaignID(input GetTransactionByIDInput) ([]Transaction, error) {

	transactions, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}

	if len(transactions) == 0 {
		return []Transaction{}, nil
	}

	return transactions, nil
}

func (s *service) GetTransactionsUser(userID int) ([]Transaction, error) {

	transactions, err := s.repository.GetByUserID(userID)
	if err != nil {
		return transactions, err
	}

	if len(transactions) == 0 {
		return []Transaction{}, nil
	}

	return transactions, nil
}