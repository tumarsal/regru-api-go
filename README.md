# REG.RU API 2.0 Go Client

Go-клиент для работы с API REG.RU (Рег.API 2.0) — домены, хостинг, DNS, счета, SSL-сертификаты.

## Установка

```bash
go get github.com/example/regru-api-go
```

### CLI (`regru`)

```bash
# из исходников репозитория
make install          # → $(go env GOPATH)/bin/regru
# или
go install ./cmd/regru/

cp .env.example .env  # username, password
```

Конфигурация: файл `.env`, переменные `REGRU_*` или флаги `--username`, `--password`. Базовый URL API по умолчанию: `https://api.reg.ru`.

## Быстрый старт

```go
package main

import (
    "context"
    "fmt"
    "log"

    regru "github.com/example/regru-api-go"
)

func main() {
    ctx := context.Background()

    // Создание клиента
    client := regru.NewClient("username", "password")

    // Проверка доступности домена
    result, err := client.Domain().Check(ctx, "example.ru", false)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Status: %s\n", result.Result)

    // Получение баланса
    balance, err := client.User().GetBalance(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Balance: %.2f %s\n", balance.Balance, balance.Currency)
}
```

## Аутентификация

### По логину и паролю (простой способ)

```go
client := regru.NewClient("username", "password")
```

### По сигнатуре (безопасный способ)

Требует SSL-сертификат, загруженный в настройках API REG.RU:

```go
privateKey, err := os.ReadFile("/path/to/private.key")
if err != nil {
    log.Fatal(err)
}

client, err := regru.ClientWithSignature("username", privateKey)
if err != nil {
    log.Fatal(err)
}
```

## Документация API

Полная документация REG.API 2.0 находится в файле [`docs/reg-api-docs.md`](docs/reg-api-docs.md).

Оригинал: https://www.reg.ru/reseller/api2doc

## Поддерживаемые функции

### Общие функции
- `Nop` — проверка доступности API
- `ResellerNop` — проверка доступности API для партнеров
- `GetUserID` — получение ID пользователя
- `GetServiceID` — получение ID услуги

### Пользователь (`User`)
- `user/nop` — проверка доступности
- `user/create` — регистрация пользователя (партнеры)
- `user/get_statistics` — статистика
- `user/get_balance` — баланс
- `user/set_reseller_url` / `user/get_reseller_url` — URL редиректа
- `user/get_persons` — список персон

### Счета (`Bill`)
- `bill/nop` — проверка доступности
- `bill/get_not_payed` — неоплаченные счета
- `bill/get_for_period` — счета за период
- `bill/change_pay_type` — изменение способа оплаты
- `bill/delete` — удаление счета

### Услуги (`Service`)
- `service/nop` — проверка доступности
- `service/get_prices` — цены на услуги
- `service/get_servtype_details` — детали типа услуги
- `service/create` — заказ услуги
- `service/check_create` — валидация параметров
- `service/delete` — удаление услуги
- `service/get_info` — информация об услугах
- `service/get_list` — список услуг
- `service/get_folders` — папки услуги
- `service/get_details` — детали услуги
- `service/renew` — продление
- `service/get_bills` — счета услуги
- `service/set_autorenew_flag` — автопродление
- `service/suspend` / `service/resume` — приостановка/возобновление
- `service/upgrade` — смена тарифа
- `service/partcontrol_grant` / `service/partcontrol_revoke` — частичное управление

### Домены (`Domain`)
- `domain/nop` — проверка доступности
- `domain/get_prices` — цены доменов
- `domain/get_suggest` — подбор имени
- `domain/get_premium_prices` — цены премиум-доменов
- `domain/get_deleted` — список удаленных доменов
- `domain/check` — проверка занятости
- `domain/create` — регистрация домена
- `domain/transfer` — перенос домена
- `domain/get_transfer_status` — статус переноса
- `domain/set_new_authinfo` — новый код авторизации
- `domain/cancel_transfer` — отмена переноса
- `domain/get_rereg_data` — данные об освобождающихся доменах
- `domain/set_rereg_bids` / `domain/get_user_rereg_bids` — ставки на домены
- `domain/get_docs_upload_uri` — загрузка документов
- `domain/update_contacts` — обновление контактов
- `domain/update_private_person_flag` — флаг приватности WHOIS
- `domain/register_ns` / `domain/delete_ns` — управление NS
- `domain/get_nss` / `domain/update_nss` — список/обновление NS
- `domain/delegate` / `domain/undelegate` — делегирование
- `domain/transfer_to_another_account` — передача домена
- `domain/look_at_entering_list` — входящие передачи

### DNS-зона (`Zone`)

Документация REG.API: [zone/add_aaaa](https://www.reg.ru/reseller/api2doc), [zone/get_resource_records](https://www.reg.ru/reseller/api2doc).

- `zone/nop` — проверка доступности
- `zone/get_resource_records` — список записей (`GetResourceRecords`, `FilterRecords`)
- `zone/add_alias` / `zone/add_aaaa` / `zone/add_cname` / `zone/add_mx` / `zone/add_ns` / `zone/add_txt` / `zone/add_srv` / `zone/add_caa` / `zone/add_https` — добавление (`AddRecord`)
- `zone/remove_record` — удаление (`RemoveRecord`)
- `zone/clear` — очистка зоны

Хелперы: `Add_Alias`, `Add_AAAA`, `Add_CNAME`, `Add_MX`, `Add_NS`, `Add_TXT`, `Add_SRV`.

#### CLI: просмотр и добавление записей

```bash
# все записи зоны
regru zone list example.ru

# с фильтрами
regru zone list example.ru --type AAAA
regru zone list example.ru --subdomain www
regru zone list example.ru --type A --content 1.2.3.4

# добавление (тип → соответствующий zone/add_* на стороне API)
regru zone add example.ru --type A --subdomain www --content 1.2.3.4
regru zone add example.ru --type AAAA --subdomain @ --content 2001:db8::1
regru zone add example.ru --type CNAME --subdomain mail --canonical-name mx10.example.ru
regru zone add example.ru --type MX --subdomain @ --content mail.example.ru --priority 10
regru zone add example.ru --type NS --subdomain tt --dns-server ns1.example.ru --record-number 10
regru zone add example.ru --type TXT --subdomain @ --text "v=spf1 ~all"
regru zone add example.ru --type SRV --service _sip._udp --target sip.example.ru --port 5060 --priority 0
regru zone add example.ru --type CAA --subdomain @ --tag issuewild --value ca.example.com --flags 0
regru zone add example.ru --type HTTPS --subdomain @ --target . --value "alpn=h3" --priority 1
```

```go
records, err := client.Zone().GetResourceRecords(ctx, "example.ru")
filtered := regru.FilterRecords(records, regru.RecordFilter{Type: "AAAA", Subdomain: "www"})

err = client.Zone().Add_AAAA(ctx, "example.ru", "www", "2001:db8::1")
```

### Папки (`Folder`)
- `folder/nop` — проверка доступности
- `folder/get_list` — список папок
- `folder/create` / `folder/rename` / `folder/delete`
- `folder/add_services` / `folder/remove_services` / `folder/replace_services`
- `folder/get_services` / `folder/get_not_in_folder_services`
- `folder/move_services`

### Магазин доменов (`Shop`)
- `shop/nop` — проверка доступности
- `shop/get_info` — информация о лоте
- `shop/enable` / `shop/disable` — включение/выключение продажи
- `shop/get_categories` — категории
- `shop/get_suggested_tags` — популярные теги

## Настройки клиента

```go
// Кастомный HTTP клиент с таймаутом
httpClient := &http.Client{Timeout: 60 * time.Second}
client := regru.NewClient("username", "password", regru.WithHTTPClient(httpClient))

// Кастомный endpoint (для тестирования)
client := regru.NewClient("username", "password", regru.WithBaseURL("https://api.reg.ru"))

// Формат ответа (json, yaml, xml, plain)
client := regru.NewClient("username", "password", regru.WithOutputFormat(regru.OutputJSON))
```

## Обработка ошибок

```go
result, err := client.Domain().Check(ctx, "example.ru", false)
if err != nil {
    // err может содержать:
    // - Ошибки HTTP (статус != 200)
    // - Ошибки сети
    // - Ошибки API (error_code из ответа REG.RU)
    log.Printf("Error: %v", err)
    return
}
```

## Константы

```go
// Способы оплаты
regru.PayTypePrepay // предоплата (по умолчанию)
regru.PayTypeBank   // банковский перевод
regru.PayTypePBank  // безналичный перевод
regru.PayTypeYaCard // ЮMoney

// Валюты
regru.CurrencyRUR // рубли
regru.CurrencyUSD // доллары
regru.CurrencyEUR // евро
regru.CurrencyUAH // гривны

// Типы услуг
regru.ServTypeDomain          // домен
regru.ServTypeHostingISP      // хостинг ISPmanager
regru.ServTypeWebFwd          // веб-перенаправление
regru.ServTypeParking         // парковка
regru.ServTypeSSLCertificate  // SSL-сертификат
regru.ServTypeVPS             // VPS
regru.ServTypeDedicatedServer // выделенный сервер
```

## Тестовый доступ

Для отладки используйте логин/пароль `test`/`test` — запросы выполняются без реальных действий и списания средств.

```go
client := regru.NewClient("test", "test")
err := client.Nop(ctx)
```

## Лицензия

MIT
