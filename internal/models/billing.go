package models

// BillingUser represents a user in the billing system (table 'users')
type BillingUser struct {
	ID        int64   `json:"id" db:"id"`
	Contract  string  `json:"contract" db:"contract"`
	Name      string  `json:"name" db:"fio"`
	Phone     string  `json:"phone" db:"telefon"`
	Balance   float64 `json:"balance" db:"balance"`
	PlanID    *int64  `json:"plan_id" db:"paket"`
	Status    string  `json:"status" db:"state"` // 'on' or 'off'
	TimePay   int     `json:"time_pay" db:"t_pay"`
	StartDay  string  `json:"start_day" db:"start_day"`
	Services  string  `json:"services" db:"srvs"`
	Address   string  `json:"address" db:"address"`
	Group     int     `json:"group" db:"grp"`
}

// BillingPlan represents a tariff plan (table 'plans2')
type BillingPlan struct {
	ID    int64   `json:"id" db:"id"`
	Price float64 `json:"price" db:"price"`
	Name  string  `json:"name" db:"name"` // Assuming 'name' exists, though not used in Python
}

// BillingPay represents a payment record (table 'pays')
type BillingPay struct {
	ID      int64   `json:"id" db:"id"`
	UserID  int64   `json:"user_id" db:"mid"`
	Amount  float64 `json:"amount" db:"cash"`
	Time    float64 `json:"time" db:"time"` // Unix timestamp
	Admin   string  `json:"admin" db:"admin"`
	Reason  string  `json:"reason" db:"reason"`
	Comment string  `json:"comment" db:"coment"`
	Bonus   string  `json:"bonus" db:"bonus"` // 'y' or 'n'
	Flag    string  `json:"flag" db:"flag"`   // 't' for temporary
}

// BillingService represents a service (table 'services' - inferred, might not be used directly if 'srvs' in users is enough)
// Keeping it for now if needed for detailed service info
type BillingService struct {
	ID        int64  `json:"id" db:"id"`
	UserID    int64  `json:"user_id" db:"user_id"`
	Name      string `json:"name" db:"name"`
	Status    string `json:"status" db:"status"`
	IPAddress string `json:"ip_address" db:"ip_address"`
}

