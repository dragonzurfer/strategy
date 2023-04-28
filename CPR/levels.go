package cpr

type CPRLevelType string

const (
	Resistance      CPRLevelType = "RESISTANCE"
	Support         CPRLevelType = "SUPPORT"
	CentralPivot    CPRLevelType = "CENTRAL PIVOT"
	BottomPivot     CPRLevelType = "BOTTOM PIVOT"
	TopPivot        CPRLevelType = "TOP PIVOT"
	InitBalanceHigh CPRLevelType = "INITIAL BALANCE HIGH"
	InitBalanceLow  CPRLevelType = "INITIAL BALANCE LOW"
	PDH             CPRLevelType = "PREVIOUS DAY HIGH"
	PDL             CPRLevelType = "PREVIOUS DAY LOW"
)

type CPRLevel struct {
	Price float64
	Type  CPRLevelType
}

func (c *CPRLevel) HasCrossed(reversalPrice float64, closingPrice float64) bool {

	if reversalPrice > c.Price && c.Price > closingPrice {
		return true
	}

	if reversalPrice < c.Price && c.Price < closingPrice {
		return true
	}
	return false
}
