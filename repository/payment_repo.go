package repository

type PaymentRepo interface {
	GetRevenue() (float64, error)
}

func (r *Repo) GetRevenue() (float64, error) {
	var revenue float64
	res := r.DB.Raw("SELECT sum(p.payment_amount) as revenue FROM payments as p").First(&revenue)
	if res.Error != nil {
		return 0, res.Error
	}

	return revenue, nil
}
