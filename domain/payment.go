package domain

type PaymentCreateRequest struct {
	AdId   string `json:"adId"`
	Amount string `json:"amount"`
}

// PaymentRequest структура для запроса на создание платежа
type PaymentRequest struct {
	PatternID string `json:"patternId"` // ID шаблона (p2p — перевод на кошелек)
	To        string `json:"to"`        // Номер кошелька получателя
	Amount    string `json:"amountDue"` // Сумма платежа
	Comment   string `json:"comment"`   // Комментарий к платежу
}

// PaymentResponse структура для обработки ответа от API
type PaymentResponse struct {
	Status       string `json:"status"`     // Статус платежа
	Error        string `json:"error"`      // Ошибки, если есть
	PaymentID    string `json:"payment_id"` // ID платежа
	Confirmation struct {
		ConfirmationURL string `json:"confirmation_url"` // Ссылка на оплату
	} `json:"money_source"`
}
