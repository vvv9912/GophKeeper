package main

import (
	"GophKeeper/internal/Agent/app"
)

func main() {
	//создаём шаблон сертификата
	//cert := &x509.Certificate{
	//	// указываем уникальный номер сертификата
	//	SerialNumber: big.NewInt(1658),
	//	// заполняем базовую информацию о владельце сертификата
	//	Subject: pkix.Name{
	//		Organization: []string{"TestSert"},
	//		Country:      []string{"RU"},
	//	},
	//	// разрешаем использование сертификата для 127.0.0.1 и ::1
	//	IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
	//	// сертификат верен, начиная со времени создания
	//	NotBefore: time.Now(),
	//	// время жизни сертификата — 10 лет
	//	NotAfter:     time.Now().AddDate(10, 0, 0),
	//	SubjectKeyId: []byte{1, 2, 3, 4, 6},
	//	// устанавливаем использование ключа для цифровой подписи,
	//	// а также клиентской и серверной авторизации
	//	ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	//	KeyUsage:    x509.KeyUsageDigitalSignature,
	//	DNSNames:    []string{"localhost"},
	//}
	//
	//// создаём новый приватный RSA-ключ длиной 4096 бит
	//// обратите внимание, что для генерации ключа и сертификата
	//// используется rand.Reader в качестве источника случайных данных
	//privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// создаём сертификат x.509
	//certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// кодируем сертификат и ключ в формате PEM, который
	//// используется для хранения и обмена криптографическими ключами
	////var certPEM bytes.Buffer
	//
	//certOut, err := os.Create("cert.pem")
	//if err != nil {
	//	fmt.Println("Ошибка при создании файла сертификата:", err)
	//	return
	//}
	//pem.Encode(certOut, &pem.Block{
	//	Type:  "CERTIFICATE",
	//	Bytes: certBytes,
	//})
	//keyOut, err := os.OpenFile("key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	////var privateKeyPEM bytes.Buffer
	//if err != nil {
	//	fmt.Println("Ошибка при создании файла сертификата:", err)
	//	return
	//}
	//pem.Encode(keyOut, &pem.Block{
	//	Type:  "RSA PRIVATE KEY",
	//	Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	//})
	app.Run()
}
