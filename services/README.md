# Services

Модуль `services` содержит набор сервисов, которые обеспечивают взаимодействие с системами поставщиков.

## Использование

Каждый сервис описывается в отдельной директории и должен имплементировать интерфейс `Service` или `ERPSystem` в случае, если это ERP система. <br/>

Регистрация сервисов в системе происходит при запуске метода `Init()` в модуле `initializer`

```go
package initializer

import (
	"MerlionScript/services/elektronmir"
)

func InitServices() {
	elektronmir.Init()
}
```

![alt text](/assets/services_schm.png)