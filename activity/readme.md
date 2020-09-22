# Журналирование действий пользователя (User activity)

## Основная функция

Для записи события необходимо использовать функцию __Log__:

```golang
activity.Log(activity Activity, data ...interface{})
```
В качестве первого аргумента передается комбинация идентификатора события и его статус в случае неуспешного действия.

Второй аргумент необязательный: это кортеж из примитивных типов и также могут быть использованы структуры __activity.Param__.

__Примеры вызова функции:__

```golang
// успешная авторизация пользователя __user_012345__ с __192.168.0.12__ адреса.
activity.Log(activity.Auth, activity.NewParam("UserID", "user_012345"), activity.NewParam("Ip", "192.168.0.12"))

// неудачная авторизация пользователя __user_012345__ с __192.168.0.12__ адреса - Пользователь заблокирован.
activity.Log(activity.Auth | activity.AuthUserLocked, activity.NewParam("UserID", "user_012345"), activity.NewParam("Ip", "192.168.0.12"))
```

## Структура

Логирование состоит из 2 компонентов:

1. __Писателя событий__ (interface handler). Записывает их в файл или работает как syslog-клиент,
3. __Шаблонизатора событий__ (interface formatter). В зависимости от писателя нормализует событие к формату, соответствующему RFC5424.

Используется "ленивая инициализация" логгирования. По умолчанию логирование выключено.

Настройка логгирования возможна с помощью следующих флагов:

```golang
func main() {
	cli.App = cli.Command{
		Name: "Trade Finances",
		Flags: []*cli.Flag{
            // обработчик событий
            cli.NewFlag("activity_handler", 0, "user activity: 0 - none, 1 - syslog 2 - file. Default: 0"),
            // имя файла лога (только для файлового обработчика)
			cli.NewFlag("activity_filename", "activity.log", "filename of activity only for file handler. Default: activity.log"),
            // идентификатор лога
            cli.NewFlag("activity_tag", "cbg", "user activity tag. Default: cbg"),
		},
		Before: func(c *cli.Command) error {
			activity.TheConfig = &activity.Config{
				c.Int("activity_handler"),
				c.String("activity_filename"),
				c.String("activity_tag"),
			}
			return nil
		},
	}
	cli.RunAndExit(os.Args)
}
```
Флаги используются для создания конфига логгера.

В первом обращения к функции к функции __Log__ происходит инициализация логгера и его запуск (Потокобезопасно. Используется __sync.Once__). Если логирование не настроено, то функция работает как _stab_.

## Настройка журнала логирования ##

Независимо от того, какой используется обработчик событий, в журнале логирования они отображаюся единообразно в формате RFC5424. Например,

```txt
<14>1 2020-09-21T18:38:15.756002+03:00 SF314-57 cbg 602032 - -  [UserID=user_012345 Ip=192.168.0.12] Authentication OK
<11>1 2020-09-21T18:38:15.756171+03:00 SF314-57 cbg 602032 - -  [UserID=user_012345 Ip=192.168.0.12] Authentication failed: user is locked

2020-09-21T18:38:15.756171+03:00 Authentication failed: user is locked
```

### Настройка /etc/rsyslog.conf при использовании обработчика __syslog__

```txt
# события сохраняются в локальный файл
:syslogtag, contains, "cbg" /var/log/cbg.log;RSYSLOG_SyslogProtocol23Format
# события отправляются на сервер syslog по UDP 
:syslogtag, contains, "cbg" @0.0.0.0:514;RSYSLOG_SyslogProtocol23Format
```

Затем необходимо перезапустить сервис __rsyslog__ для применения настроек:

```sh
$ sudo service rsyslog restart
```

### Настройка /etc/rsyslog.conf при записи в локальный файл 

Проблемы:
1. Для создания файлов и записи логов в директории __/var/log/*__ требуются дополнительные привилегии, поэтому запись в файл производится в ~/ каталоге
2. Необходимо отслеживать размер файла и делать ротацию (реализовано с использованием пакета __"gopkg.in/natefinch/lumberjack.v2"__)
3. Не удалось подключить файл и мониторить его изменения в __rsyslog__ как сделано в __rsyslogtest__.

```txt
#https://www.rsyslog.com/doc/v8-stable/configuration/templates.html
template(name="RSYSLOG_SyslogProtocol23Format" type="string"
     string="<%PRI%>1 %TIMESTAMP:::date-rfc3339% %HOSTNAME% %APP-NAME% %PROCID% %MSGID% %STRUCTURED-DATA% %msg%\n")

input(
    type="imfile"
    File="/home/apotapov/go/src/log/cbgfile.log"
    Tag="cbg:"
    Facility="user"
    Severity="info"
    PersistStateInterval="1"
    reopenOnTruncate="on"
    freshStartTail="on"

    ruleset="my-test-file-logging"
)

ruleset(name="my-test-file-logging") {
    action(
        type="omfwd"
        template="RSYSLOG_SyslogProtocol23Format"
        queue.type="LinkedList"
        queue.filename="fileq1"
        queue.saveonshutdown="on"
        action.resumeRetryCount="-1"
        action.ResumeInterval="10"
        Target="127.0.0.1"
        Port="514"
        Protocol="udp"
        #timeout="5"
    )
}
```