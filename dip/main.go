/*
3. Dependency Inversion Principle (Bağımlılıkların Tersine Çevrilmesi Prensibi - DIP)
Nedir?
Robert C. Martin (Uncle Bob) tarafından popüler hale getirilen DIP, SOLID prensiplerinin D harfidir ve kısaca şöyle der:
"1. Üst seviye modüller alt seviye modüllere bağımlı olmamalıdır. Her ikisi de soyutlamalara bağımlı olmalıdır."
"2. Soyutlamalar detaylara bağımlı olmamalıdır. Detaylar soyutlamalara bağımlı olmalıdır."
Biraz karışık görünebilir, ama temel fikir şudur:
Üst Seviye Modül: Uygulamanın iş mantığını içeren modüllerdir (örneğin, bir iş akışını yöneten bir servis).
Alt Seviye Modül: Detayları ve yardımcı araçları içeren modüllerdir (örneğin, bir veritabanı sürücüsü, bir dosya yazıcısı).
DIP'nin amacı, üst seviye modüllerin alt seviye modüllerin somut implementasyonlarına (yani, doğrudan kodlarına) bağımlı olmasını engellemektir. Bunun yerine, her ikisi de soyutlamalara (Go'da interface'ler) bağımlı olmalıdır.
Neden Önemli?
Esneklik ve Genişletilebilirlik: Alt seviye modüller değiştiğinde üst seviye modüllerin etkilenmesini engeller.
Test Edilebilirlik: Bağımlılıkları kolayca mock'layarak (sahte nesnelerle değiştirerek) birim testlerini çok daha kolay hale getirir.
Modülerlik: Kod parçalarını daha bağımsız hale getirir, böylece farklı bağlamlarda yeniden kullanılabilirler.
Bakım Kolaylığı: Kod tabanındaki değişikliklerin etkisini izole eder.
*/
package main

import (
	"fmt"
	"log"
	"os"
)

// Dosya işlemleri için

// Logger interface'i (Soyutlama)
// Üst seviye modül (ErrorHandler) ve alt seviye modüller (ConsoleLogger, FileLogger)
// bu soyutlamaya bağımlı olacak.
type Logger interface {
	Log(message string)
}

// ConsoleLogger struct'ı, Logger interface'ini uygular. (Detay)
type ConsoleLogger struct{}

func (cl *ConsoleLogger) Log(message string) {
	fmt.Printf("Konsol Log: %s\n", message)
}

// FileLogger struct'ı, Logger interface'ini uygular. (Detay)
type FileLogger struct {
	fileName string
}

func (fl *FileLogger) Log(message string) {
	file, err := os.OpenFile(fl.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Dosyaya log yazarken hata oluştu: %v", err)
		return
	}
	defer file.Close()
	logger := log.New(file, "FILE: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println(message)
	fmt.Printf("Dosyaya log yazıldı: %s\n", message)
}

// ErrorHandler üst seviye modüldür.
// Artık doğrudan somut bir loglayıcıya değil, Logger interface'ine bağımlıdır.
// Bu, bağımlılıkların tersine çevrilmesidir.
type ErrorHandler struct {
	logger Logger // Logger interface'ine bağımlılık
}

// NewErrorHandler constructor fonksiyonu, ErrorHandler'ı bir Logger interface'i ile başlatır.
func NewErrorHandler(logger Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

func (eh *ErrorHandler) HandleError(err error) {
	eh.logger.Log(fmt.Sprintf("Hata meydana geldi: %v", err))
	// Diğer hata işleme mantığı...
}

func main() {
	fmt.Println("DIP ile Loglama Örneği")

	// ConsoleLogger ile ErrorHandler oluşturma
	consoleLogger := &ConsoleLogger{}
	errorHandlerWithConsole := NewErrorHandler(consoleLogger)
	errorHandlerWithConsole.HandleError(fmt.Errorf("konsola yazılacak bir hata"))
	fmt.Println("--------------------")

	// FileLogger ile ErrorHandler oluşturma
	// "app_errors.log" dosyasına yazacak.
	fileLogger := &FileLogger{fileName: "app_errors.log"}
	errorHandlerWithFile := NewErrorHandler(fileLogger)
	errorHandlerWithFile.HandleError(fmt.Errorf("dosyaya yazılacak başka bir hata"))
	fmt.Println("--------------------")

	// Gelecekte eklenecek bir DatabaseLogger veya CloudLogger da bu şekilde kullanılabilir.
}

/*  --- IGNORE ---
type ConsoleLogger struct{}

func (cl *ConsoleLogger) Log(message string) {
	fmt.Printf("Konsol Log: %s\n", message)
}

// ErrorHandler üst seviye modüldür ve doğrudan ConsoleLogger'a bağımlıdır.
// Bu, DIP'ye aykırıdır.
type ErrorHandler struct {
	logger ConsoleLogger // Doğrudan somut bir bağımlılık
}

func (eh *ErrorHandler) HandleError(err error) {
	eh.logger.Log(fmt.Sprintf("Hata meydana geldi: %v", err))
	// Diğer hata işleme mantığı...
}

func main() {
	// Doğrudan somut bir logger oluşturup ErrorHandler'a veriyoruz.
	errorHandler := ErrorHandler{logger: ConsoleLogger{}}
	errorHandler.HandleError(fmt.Errorf("bir test hatası"))
}
*/
