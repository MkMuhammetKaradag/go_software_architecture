/*

2. Open-Closed Principle (Açık-Kapalı Prensibi - OCP)
Nedir?
Open-Closed Prensibi, Bertrand Meyer tarafından ortaya atılmış ve yazılım tasarımında çok temel bir prensiptir. Kısaca şöyle der:
"Yazılım varlıkları (sınıflar, modüller, fonksiyonlar vb.) geliştirmeye açık, ancak değiştirmeye kapalı olmalıdır."
Bu ne anlama geliyor?
Geliştirmeye Açık (Open for Extension): Bir modülün davranışını, yeni işlevsellik ekleyerek genişletebilmeliyiz. Yani, mevcut kodun üzerine inşa edebiliriz.
Değiştirmeye Kapalı (Closed for Modification): Bir modülün davranışını değiştirmek için mevcut kaynak kodunu değiştirmememiz gerekir. Yani, mevcut,
çalışan kodu bozma veya hata ekleme riskini azaltırız.
Neden Önemli?
Sürdürülebilirliği Artırır: Mevcut kodu değiştirmek yerine yeni kod eklemek, sistemin daha az hata içermesini sağlar.
Esnekliği Artırır: Yeni özellikler eklemek daha kolay hale gelir.
Bakım Maliyetini Azaltır: Değişikliklerin potansiyel yan etkilerini azaltarak bakım maliyetlerini düşürür.
Go Dilinde Örnek:
İlk konumuz olan "Abstraction" ile zaten OCP'ye doğru bir adım atmıştık. Interface'ler, Go'da OCP'yi uygulamanın temel yollarından biridir.
Bir mağaza sisteminde indirimleri hesapladığımızı düşünelim. Başlangıçta sadece sabit bir indirimimiz olsun, sonra yeni indirim türleri eklemek isteyelim
 (Öğrenci indirimi, Kara Cuma indirimi vb.).


*/

package main

import "fmt"

// DiscountStrategy interface'i, tüm indirim stratejilerinin sahip olması gereken davranışı soyutlar.
type DiscountStrategy interface {
	ApplyDiscount(price float64) float64
}

// NoDiscount struct'ı, DiscountStrategy interface'ini uygular.
type NoDiscount struct{}

func (nd NoDiscount) ApplyDiscount(price float64) float64 {
	return price
}

// FixedDiscount struct'ı, DiscountStrategy interface'ini uygular.
type FixedDiscount struct{}

func (fd FixedDiscount) ApplyDiscount(price float64) float64 {
	return price * 0.90 // %10 sabit indirim
}

// StudentDiscount struct'ı, DiscountStrategy interface'ini uygular.
type StudentDiscount struct{}

func (sd StudentDiscount) ApplyDiscount(price float64) float64 {
	return price * 0.85 // %15 öğrenci indirimi
}

// BlackFridayDiscount struct'ı, DiscountStrategy interface'ini uygular. (Yeni eklenen indirim)
type BlackFridayDiscount struct{}

func (bfd BlackFridayDiscount) ApplyDiscount(price float64) float64 {
	return price * 0.70 // %30 Kara Cuma indirimi
}

// Product (Önceki örnekten, fiyat bilgisi için kullanılıyor)
type Product struct {
	Name  string
	Price float64
}

// DiscountCalculator struct'ı, bir DiscountStrategy alır.
// Bu sayede, farklı indirim stratejileri ile çalışabiliriz.
type DiscountCalculator struct {
	strategy DiscountStrategy
}

// NewDiscountCalculator, bir indirim stratejisi ile DiscountCalculator oluşturur.
func NewDiscountCalculator(strategy DiscountStrategy) *DiscountCalculator {
	return &DiscountCalculator{strategy: strategy}
}

// CalculateDiscount metodu, mevcut stratejiyi kullanarak indirimi uygular.
func (dc *DiscountCalculator) CalculateDiscount(productPrice float64) float64 {
	return dc.strategy.ApplyDiscount(productPrice)
}

func main() {
	laptop := Product{Name: "Laptop", Price: 1500.00}

	// Sabit indirim ile hesaplama
	fixedDiscountCalc := NewDiscountCalculator(FixedDiscount{})
	priceAfterFixedDiscount := fixedDiscountCalc.CalculateDiscount(laptop.Price)
	fmt.Printf("%s (Sabit İndirim sonrası): %.2f TL\n", laptop.Name, priceAfterFixedDiscount)

	// Öğrenci indirimi ile hesaplama
	studentDiscountCalc := NewDiscountCalculator(StudentDiscount{})
	priceAfterStudentDiscount := studentDiscountCalc.CalculateDiscount(laptop.Price)
	fmt.Printf("%s (Öğrenci İndirimi sonrası): %.2f TL\n", laptop.Name, priceAfterStudentDiscount)

	// Kara Cuma indirimi ile hesaplama (Yeni bir indirim türü ekledik!)
	blackFridayDiscountCalc := NewDiscountCalculator(BlackFridayDiscount{})
	priceAfterBlackFridayDiscount := blackFridayDiscountCalc.CalculateDiscount(laptop.Price)
	fmt.Printf("%s (Kara Cuma İndirimi sonrası): %.2f TL\n", laptop.Name, priceAfterBlackFridayDiscount)

	// İndirimsiz hesaplama
	noDiscountCalc := NewDiscountCalculator(NoDiscount{})
	priceNoDiscount := noDiscountCalc.CalculateDiscount(laptop.Price)
	fmt.Printf("%s (İndirimsiz): %.2f TL\n", laptop.Name, priceNoDiscount)
}







/*  --- IGNORE ---

package main

import "fmt"

type Product struct {
	Name  string
	Price float64
}

// DiscountCalculator fonksiyonu, OCP'ye aykırı bir şekilde tasarlanmıştır.
// Yeni indirim türleri eklendiğinde bu fonksiyonun içini değiştirmemiz gerekir.
func CalculateDiscountBad(product Price, discountType string) float64 {
	switch discountType {
	case "no_discount":
		return product.Price
	case "fixed_discount":
		return product.Price * 0.90 // %10 sabit indirim
	case "student_discount":
		return product.Price * 0.85 // %15 öğrenci indirimi
	// Yarın yeni bir indirim türü geldiğinde buraya yeni bir 'case' eklemeliyiz.
	// Bu, mevcut çalışan kodu değiştirmek demektir.
	default:
		return product.Price
	}
}

// Bu kısım sadece gösterim amaçlıdır, ana fonksiyonun devamı için gerekli değil.
type Price interface {
	GetPrice() float64
}

func (p Product) GetPrice() float64 {
	return p.Price
}

*/
