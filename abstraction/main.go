/*
1. Abstraction (Soyutlama)
Nedir?
Soyutlama, karmaşık sistemlerin detaylarını gizleyerek, kullanıcının veya başka bir yazılım parçasının sadece ilgili fonksiyonelliğe odaklanmasını sağlayan bir prensiptir. Bir nesnenin "ne yaptığına" odaklanırken, "nasıl yaptığı" detayını gizler. Bu, kodun daha okunabilir, sürdürülebilir ve esnek olmasını sağlar. Go dilinde soyutlamayı genellikle interface'ler aracılığıyla yaparız.
Neden Önemli?
Karmaşıklığı Azaltır: Kullanıcının sadece bilmesi gereken şeyleri görmesini sağlar.
Değişiklikleri Yönetir: İç detaylar değişse bile, dış arayüz sabit kalabilir, böylece bağımlı kodların etkilenmesini engeller.
Esneklik Sağlar: Farklı implementasyonların aynı arayüzü kullanarak birbirinin yerine geçebilmesini sağlar.
Go Dilinde Örnek:
Bir ödeme sistemi düşünelim. Kredi kartı, PayPal gibi farklı ödeme yöntemleri olabilir. Her bir ödeme yöntemi farklı çalışsa da, sonuçta hepsi bir ödeme işlemini gerçekleştirir.
*/
package main

import "fmt"

// PaymentMethod interface'i, tüm ödeme yöntemlerinin sahip olması gereken davranışı soyutlar.
type PaymentMethod interface {
	Pay(amount float64) string
}

// CreditCard struct'ı, PaymentMethod interface'ini uygular.
type CreditCard struct {
	CardNumber string
	ExpiryDate string
	CVV        string
}

// Pay metodu CreditCard için ödeme işlemini gerçekleştirir.
func (cc *CreditCard) Pay(amount float64) string {
	return fmt.Sprintf("Kredi kartı ile %.2f TL ödendi. Kart No: %s", amount, cc.CardNumber)
}

// PayPal struct'ı, PaymentMethod interface'ini uygular.
type PayPal struct {
	Email    string
	Password string
}

// Pay metodu PayPal için ödeme işlemini gerçekleştirir.
func (pp *PayPal) Pay(amount float64) string {
	return fmt.Sprintf("PayPal ile %.2f TL ödendi. E-posta: %s", amount, pp.Email)
}

// ProcessPayment fonksiyonu, PaymentMethod interface'ini kabul eder.
// Bu sayede farklı ödeme yöntemleri için aynı fonksiyonu kullanabiliriz.
func ProcessPayment(method PaymentMethod, amount float64) {
	fmt.Println(method.Pay(amount))
}

func main() {
	creditCard := &CreditCard{
		CardNumber: "1234-5678-9012-3456",
		ExpiryDate: "12/25",
		CVV:        "123",
	}

	payPal := &PayPal{
		Email:    "example@example.com",
		Password: "mysecretpassword",
	}

	fmt.Println("Alışveriş başladı...")
	ProcessPayment(creditCard, 100.50) // Kredi kartı ile ödeme
	ProcessPayment(payPal, 50.00)      // PayPal ile ödeme
	fmt.Println("Alışveriş bitti.")
}
