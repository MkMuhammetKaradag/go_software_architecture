/*


### 4. IoC Container (Inversion of Control Container - Kontrolün Tersine Çevrilmesi Konteyneri)

**Nedir?**

IoC (Inversion of Control - Kontrolün Tersine Çevrilmesi), bir yazılım tasarım prensibidir. Geleneksel programlamada, bir objenin
(veya modülün) kendisi ihtiyaç duyduğu bağımlılıkları (diğer objeleri) oluşturur ve yönetir. IoC'de ise, objenin bağımlılıklarını objenin kendisi değil,
 dışarıdan bir mekanizma (genellikle bir "container" veya "framework") sağlar.

**IoC Container (veya Dependency Injection Container - DI Container):**

Bu, IoC prensibini uygulayan bir yazılım aracıdır. Başka bir deyişle, bir IoC Container:

*   Uygulamanın objelerini oluşturmaktan ve onların yaşam döngülerini yönetmekten sorumludur.
*   Objelerin bağımlılıklarını çözümler ve onları objelere "enjekte eder" (Dependency Injection).
*   Bu sayede, objeler kendi bağımlılıklarını bulmak veya oluşturmak zorunda kalmazlar, bu sorumluluğu container üstlenir.

**Neden Önemli?**

IoC Container'lar, özellikle büyük ve karmaşık uygulamalarda yukarıda öğrendiğimiz SOLID prensiplerini (özellikle DIP ve OCP) uygulamanın pratik bir yolunu sunar.

*   **Bağımlılık Yönetimi:** Objelerin karmaşık bağımlılık zincirlerini otomatik olarak çözümler.
*   **Modülerlik ve Gevşek Bağlantı:** Objeleri kendi bağımlılıklarını aramaktan kurtararak daha bağımsız ve test edilebilir hale getirir.
*   **Test Edilebilirlik:** Testlerde gerçek implementasyonlar yerine mock objeleri enjekte etmeyi kolaylaştırır.
*   **Konfigürasyon Esnekliği:** Bağımlılıkların implementasyonlarını kolayca değiştirmeye olanak tanır (örneğin, test ortamında farklı bir veritabanı,
üretimde başka bir veritabanı).
*   **Boilerplate Kodu Azaltır:** Obje oluşturma ve bağımlılık enjeksiyonu için yazılması gereken tekrarlayan kodu azaltır.

**Go Dilinde Örnek (Basit Bir Yaklaşım):**

Go'da Java/Spring veya .NET/Autofac gibi tam teşekküllü IoC Container'lar standart olarak gelmez. Bunun nedeni,
 Go'nun interface'ler ve fonksiyonel seçenekler (`functional options`) gibi dil özelliklerinin,
 Dependency Injection'ı (DI) genellikle çok daha hafif ve el ile yapılabilir kılmasındandır.
 Genellikle, Go'da DI'yı constructor'lar aracılığıyla yaparız. Ancak,
 daha büyük projeler için kendi basit IoC Container'ımızı veya bir kütüphane (örneğin `go-cleanarch/di`, `uber-go/dig`) kullanabiliriz.

Şimdi, ilk `ErrorHandler` ve `Logger` örneğimizin basit bir IoC Container benzeri yapıyla nasıl daha da soyutlanabileceğini gösterelim.
 Kendi basit "Container" yapımızı oluşturacağız.



**Açıklama:**

*   `Container` adında basit bir struct tanımladık. Bu struct, `loggers` adında bir `map` ile `Logger` interface'ini uygulayan farklı implementasyonları saklar.
*   `RegisterLogger` metodu ile farklı `Logger` implementasyonlarını (konsol, dosya vb.) bir isimle container'a kaydediyoruz.
*   `GetLogger` metodu, kayıtlı bir `Logger` implementasyonunu isimle geri döndürür.
*   `GetErrorHandler` metodu, `ErrorHandler`'ın ihtiyacı olan `Logger` bağımlılığını container'dan alır ve `NewErrorHandler` fonksiyonunu kullanarak `ErrorHandler` objesini oluşturur. **Burada kontrol tersine dönmüştür:** `main` fonksiyonu veya `ErrorHandler` kendisi `Logger`'ı oluşturmuyor, bunun yerine container ona `Logger`'ı sağlıyor.
*   `main` fonksiyonunda, önce container'ı başlatıyor, sonra `Logger` implementasyonlarını kaydediyoruz. Sonra da `ErrorHandler` objelerini container üzerinden alıyoruz. `main` fonksiyonu artık `ConsoleLogger` veya `FileLogger`'ın nasıl oluşturulduğunu veya `ErrorHandler`'ın bu logger'ları nasıl kullandığını bilmek zorunda değil; bu sorumluluk container'a devredilmiştir.

*/

package main

import (
	"fmt"
	"log"
	"os"
)

// Logger interface'i (Soyutlama)
type Logger interface {
	Log(message string)
}

// ConsoleLogger struct'ı (Detay)
type ConsoleLogger struct{}

func (cl *ConsoleLogger) Log(message string) {
	fmt.Printf("Konsol Log: %s\n", message)
}

// FileLogger struct'ı (Detay)
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

// ErrorHandler üst seviye modül
type ErrorHandler struct {
	logger Logger // Logger interface'ine bağımlılık
}

// NewErrorHandler constructor fonksiyonu
func NewErrorHandler(logger Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

func (eh *ErrorHandler) HandleError(err error) {
	eh.logger.Log(fmt.Sprintf("Hata meydana geldi: %v", err))
	// Diğer hata işleme mantığı...
}

// ------------ Basit IoC Container Yapısı ------------

// Container struct'ı, bağımlılıkları saklar.
type Container struct {
	loggers map[string]Logger
}

// NewContainer, yeni bir Container başlatır.
func NewContainer() *Container {
	return &Container{
		loggers: make(map[string]Logger),
	}
}

// RegisterLogger, bir Logger implementasyonunu container'a kaydeder.
func (c *Container) RegisterLogger(name string, logger Logger) {
	c.loggers[name] = logger
}

// GetLogger, kayıtlı bir Logger'ı container'dan döner.
func (c *Container) GetLogger(name string) (Logger, error) {
	logger, ok := c.loggers[name]
	if !ok {
		return nil, fmt.Errorf("Logger bulunamadı: %s", name)
	}
	return logger, nil
}

// GetErrorHandler, bağımlılıklarını container'dan alarak bir ErrorHandler oluşturur.
// Bu metot, bağımlılık enjeksiyonunu (DI) simüle eder.
func (c *Container) GetErrorHandler(loggerName string) (*ErrorHandler, error) {
	logger, err := c.GetLogger(loggerName)
	if err != nil {
		return nil, err
	}
	return NewErrorHandler(logger), nil
}

func main() {
	fmt.Println("IoC Container Benzeri Yapı ile Örnek")

	// 1. Container'ı oluştur
	container := NewContainer()

	// 2. Bağımlılıkları container'a kaydet
	container.RegisterLogger("consoleLogger", &ConsoleLogger{})
	container.RegisterLogger("fileLogger", &FileLogger{fileName: "app_errors_container.log"})

	// 3. Container'dan objeleri ve onların bağımlılıklarını al
	// Artık main fonksiyonu, ErrorHandler'ın nasıl oluşturulduğunu veya hangi logger'ı kullandığını doğrudan yönetmek zorunda değil.
	// Bu sorumluluğu container'a devretti.

	errorHandlerWithConsole, err := container.GetErrorHandler("consoleLogger")
	if err != nil {
		log.Fatal(err)
	}
	errorHandlerWithConsole.HandleError(fmt.Errorf("container üzerinden konsola loglandı"))
	fmt.Println("--------------------")

	errorHandlerWithFile, err := container.GetErrorHandler("fileLogger")
	if err != nil {
		log.Fatal(err)
	}
	errorHandlerWithFile.HandleError(fmt.Errorf("container üzerinden dosyaya loglandı"))
	fmt.Println("--------------------")

	// Olmayan bir logger'ı deneme
	_, err = container.GetErrorHandler("databaseLogger")
	if err != nil {
		fmt.Println("Hata (beklendiği gibi):", err)
	}
}
