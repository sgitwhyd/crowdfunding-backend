package payment

import "be-bwastartup/campaign"

type Transaction struct {
	ID       int
	Amount   int
	Campaign campaign.Campaign
}