# Анализ данных из JSON-файлов  

## Установка и запуск

В терминале:  
* ```git clone https://github.com/AleksMa/testsParser.git```  
* ```cd testsParser```   
Далее либо
* ```export GO111MODULE=on```
* ```go build```    
либо
* ```go get github.com/xeipuuv/gojsonschema```    
---
Запуск осуществляется командой   
```go run parser.go file1 file2 file3 [fileResult]```    
  
Данные из примеров:   
```go run parser.go Data/1.json Data/2.json Data/3.json Data/result.json```  




---  
Результат представлен в виде массива объектов, содержащих следующие свойства:
* **"name"**     - название теста;
* **"status"**   - статус теста (*OK* или *fail*);
* **"expected"** - ожидаемое значение проверки;
* **"actual"**   - реальное значение проверки;

## Тесты
Простейшее тестирование можно запустить посредством команды ```go test```  
Происходит запуск скрипта с данными из Data/ и сравнение с правильным результатом.

## Валидация
Входные данные валидируются по JSON Schema, используется пакет ***gojsonschema***.  
Его можно установить командой ```go get github.com/xeipuuv/gojsonschema```  
Пакет чувствителен к языку. Пожалуйста, не размещайте входные данные в директориях с названиями не на английском языке.