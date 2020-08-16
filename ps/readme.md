# Powershell в Golang

Довольно часто в своей работе приходится использовать скрипты PowerShell в Golang, в частности, запросы WMI, LDAP, события Windows и т.д.
В настоящий момент мне известны только 2 библиотеки, в которых имеется нужный мне функционал:

* [KnicKnic/go-powershell](https://github.com/KnicKnic/go-powershell). Цель этого проекта - дать возможность быстро писать код golang и взаимодействовать с windows через powershell, а не использовать exec. Используется .Net подобное API для взаимодействия с PowerShell. Этот проект имеет зависимость от native-powershell(проект c++ / cli, который позволяет взаимодействовать с powershell через интерфейс DLL C).

* [40a/go-powershell](https://github.com/40a/go-powershell). Этот пакет портирован Golang из [jPowershell](https://github.com/profesorfalken/jPowerShell) и позволяет запускать как локальные, так и удаленнное выполнение скриптов. Однако последние коммиты были 4 года назад, да и сам код сейчас компилируется только после многочисленных правок исходников.

В целом в Golang имеются достаточно возможностей реализовать выполнение скриптов PowerShell без привлечения внешних зависимостей.

Запустим "powershell.exe" с аргументом "/help" и выведем на экран справку по данной команде:

```golang
func TestGetPSHelp(t *testing.T) {
    cmd := exec.Command("powershell.exe", "-Help")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
    }
}
```

Из вывода выберу только несколько вариантов запуска скриптов:

```bash
PowerShell -PSConsoleFile SqlSnapIn.Psc1
PowerShell -version 2.0 -NoLogo -InputFormat text -OutputFormat XML
PowerShell -ConfigurationName AdminRoles
PowerShell -Command {Get-EventLog -LogName security}
PowerShell -Command "& {Get-EventLog -LogName security}"

#To use the -EncodedCommand parameter:
$command = 'dir "c:\program files" '
$bytes = [System.Text.Encoding]::Unicode.GetBytes ($command)
$encodedCommand = [Convert]::ToBase64String($bytes)
powershell.exe -encodedCommand $encodedCommand
```

Здесь мы видим, как можно выполнить скрипт Powershell, используя командную строку. Я остановлюсь только на работе запуска с опцией "-encodedCommand", которая позволяет использовать все возможные скрипты PowerShell.

Для начала необходимо конвертировать наш скрипт в кодировку UTF16  LittleIndian.
Можно воспользоваться дополнительными пакетами от Google.

>"golang.org/x/text/encoding/unicode"
>
>"golang.org/x/text/transform"

Для кодирования напишем отдельную функцию:

```golang
func encodeUtf16Le(utf8 []byte) []byte {
    utf16le, _, err := transform.Bytes(
        unicode.UTF16(
            unicode.LittleEndian,
            unicode.IgnoreBOM,
        ).NewEncoder(),
        utf8,
    )
    if err != nil {
        panic err
    }
    return utf16le
}
```

Но можно и самим написать кодирование с использованием стандартных библиотек:

```golang
func encodeUtf16Le(utf8 []byte) []byte {
    uint16s := utf16.Encode([]rune(string(utf8)))
    data := make([]byte, 2*len(uint16s))
    for index, value := range uint16s {
        binary.LittleEndian.PutUint16(data[index*2:], value)
    }
    return data
}
```

Теперь нам остается только полученный массив байт перевести в Base64-строку.

Напишем тест запуска в PowerShell с помощью кодированного скрипта:

```golang
func TestRunPSScript(t *testing.T) {
    script := []byte(`dir "c:\program files" `)
    cmd := exec.Command(
        "powershell.exe",
        "-EncodedCommand",
        base64.StdEncoding.EncodeToString(
            encodeUtf8ToUtf16Le(script),
        ),
    )
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
    }
}
```

Скрипт выполняется, и мы видим результат его работы:

``` bash
#< CLIXML


    Directory: C:\


Mode                LastWriteTime         Length Name                                                                  
----                -------------         ------ ----                                                                  
d-----       28.05.2020     20:28                dev                                                                   
d-----       09.05.2020     22:30                Intel                                                                 
d-----       16.05.2020      2:04                PerfLogs                                                              
d-r---       30.05.2020     11:07                Program Files                                                         
d-r---       11.03.2020     22:07                Program Files (x86)                                                   
d-r---       12.03.2020      7:56                Users                                                                 
d-----       02.08.2020     13:37                Windows                                                               
d-----       02.08.2020     13:24                ���                                                                   


<Objs Version="1.1.0.1" xmlns="http://schemas.microsoft.com/powershell/2004/04"><Obj S="progress" RefId="0"><TN RefId="0"><T>System.Management.Automation.PSCustomObject</T><T>System.Object</T></TN><MS><I64 N="SourceId">1</I64><PR N="Record"><AV>Preparing modules for first use.</AV><AI>0</AI><Nil /><PI>-1</PI><PC>-1</PC><T>Completed</T><SR>-1</SR><SD> </SD></PR></MS></Obj></Objs>
```

Первое, что бросается в глаза - это сломанная кодировка. Имя папки в кирилице выводтся в виде ���.
Оказывается, по умолчанию вывод происходит в кодировке Codepage866.
Давайте исправим это.
Создадим функцию, которая декодирует вывод:

```golang
func decodeCp866(cp866 []byte) ([]byte, error) {
    return charmap.CodePage866.NewDecoder().Bytes(cp866)
}
```

Давайте немного переделаем вывод на экран. Создадим структуру, которая реализуется интерфейс io.Writer. За базовую стрктуру возьмем "bytes.Buffer". И переопределим 2 метода вывода:

* в виде массива байт
* в виде строки.

``` golang
type Capturer struct {
    bytes.Buffer
}

func (c *Capturer) Bytes() []byte {
    return decodeCp866(c.Buffer.Bytes())
}

func (c *Capturer) String() string {
    return string(c.Bytes())
}

```

Изменим TestRunPSScript:

``` golang

func TestRunPSScript(t *testing.T) {
    script := []byte(`dir "c:\" `)
    cmd := exec.Command(
        "powershell.exe",
        "-EncodedCommand",
        base64.StdEncoding.EncodeToString(
            encodeUtf8ToUtf16Le(script),
        ),
    )

    outCapturer := &Capturer{}
    cmd.Stdout = outCapturer

    errCapturer := &Capturer{}
    cmd.Stderr = errCapturer

    err := cmd.Run()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
    }

    fmt.Println(outCapturer.String())
    fmt.Println(errCapturer.String())
}

```

Теперь мы видим кириллицу в правильном отображении:

``` bash



    Directory: C:\


Mode                LastWriteTime         Length Name                                                                  
----                -------------         ------ ----                                                                  
d-----       28.05.2020     20:28                dev                                                                   
d-----       09.05.2020     22:30                Intel                                                                 
d-----       16.05.2020      2:04                PerfLogs                                                              
d-r---       30.05.2020     11:07                Program Files                                                         
d-r---       11.03.2020     22:07                Program Files (x86)                                                   
d-r---       12.03.2020      7:56                Users                                                                 
d-----       02.08.2020     13:37                Windows                                                               
d-----       02.08.2020     13:24                Имя                                                                   



#< CLIXML
<Objs Version="1.1.0.1" xmlns="http://schemas.microsoft.com/powershell/2004/04"><Obj S="progress" RefId="0"><TN RefId="0"><T>System.Management.Automation.PSCustomObject</T><T>System.Object</T></TN><MS><I64 N="SourceId">1</I64><PR N="Record"><AV>Preparing modules for first use.</AV><AI>0</AI><Nil /><PI>-1</PI><PC>-1</PC><T>Completed</T><SR>-1</SR><SD> </SD></PR></MS></Obj></Objs>


```
