# tcp-tunnel

Утилита, позволяющая обращаться к ресурсу локальной сети без публичного IP адреса из внешней сети через сервер с публичным IP (например виртуалка в Яндекс Облаке).

# Как запускать

`go run tunnel-server/main.go --client-port <port> --tunnel-port <port>` запустить сервер на публичном сервере. `client-port` для подключения внешнего клиента, `tunnel-port` для подключения туннельного агента. Сервера запускаются на `0.0.0.0`

`go run tunnel-agent/main.go --tunnel-ip <ip> --tunnel-port <port> --local-ip <ip> --local-port <port>` запустить агента в локальной сети с доступом к нужному ресурсу. `tunnel-ip` и `tunnel-port` адресуют туннельный сервер, `local-ip` и `local-port` - локальный ресурс (например, по умолчанию это nginx на 80 порту localhost)

Полный пример:

1. На тачке в облаке с ip 1.2.3.4 запускаем сервер. `go run tunnel-server/main.go --client-port 9090 --tunnel-port 9999`
2. На локальном сервере запускаем агента. `go run tunnel-agent/main.go --tunnel-ip 1.2.3.4 --tunnel-port 9999 --local-ip 127.0.0.1 --local-port 80`
3. на том же локальном сервере запускаем nginx на 80 порту.
4. Из внешней сети можем обратиться к nginx по 1.2.3.4:9090