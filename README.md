### Анализ данных из трех JSON-файлов

Исходя из предложенных примеров, общим ключом логов было выбрано время теста.  

По умолчанию входные данные лежат в файлах ***1.json, 2.json, 3.json***.  
Выходные данные записываются в файл ***result.json***. 

Результирующие данные представлены в виде массива объектов, содержащих следующие свойства:
* **"name"**     - название теста;
* **"status"**   - статус теста (*OK* или *fail*);
* **"expected"** - ожидаемое значение проверки;
* **"actual"**   - реальное значение проверки;

