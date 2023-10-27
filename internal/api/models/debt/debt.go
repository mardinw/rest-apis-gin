package debt

type DebtRemainder struct {
	ID            int64
	RetailerID    int64
	TransactionID int64
	Date          int64
	PaymentID     int64
	PaymentMethod string
	Status        string
}
