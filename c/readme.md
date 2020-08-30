# cgo

В Go можно использовать C-код. Для этого есть существует утилита cgo, которая неявно подключается при компиляции проекта.

## Использование фрагмента кода, написанного на C в Go

Чтобы использовать C-код, необходимо импортировать псевдопакет "С".
Код на C пишется перед импортом в комментариях и его использует cgo.
Обратите внимание, что между комментарием cgo и оператором импорта не должно быть пустых строк.

Для начала, мы хотим сложить 2 целых числа с использованием C-кода

```golang

// // C - код
// #include <stdio.h>
// void add(int a, int b) {
//     printf("%d\n", a + b);
// }
import "C"

import "fmt"

// Add 2 numbers
func Add(a, b int) {
    fmt.Printf("%d + %d = ", a, b)
    C.add(C.int(a), C.int(b))
}

```

При использовании C в коде можно выделить 2 части: все что идет вначале в комментариях и строка импорта C, относится к утилите cgo, а остальное - к Go
Обратите внимание, что перед вызовом C-функции в Go используется обращение к псевдопакету C, а переменные в виде аргументов также приводятся к формату C.

## Использование библиотеки, написанной на C в Go

Немного усложним наш пример. Пусть теперь функция сложения у нас реализована в библиотеке C и у нас есть сама библиотека и ее заголовок.

Для начала давайте ее создадим.

```c

// addition.h
void add(int a, int b);

// addition.c
#include <stdio.h>
#include "addition.h"

void add(int a, int b) {
    printf("%d\n", a + b);
}

```

Скомпилируем библиотеку:

```bash

gcc -c addition.c
gcc -shared -o libaddition.so addition.o

```

А теперь в файле addition.go изменим комментарий для cgo следующим образом:

```golang

// // C - код
// #cgo CFLAGS: -I.
// #cgo LDFLAGS: -L. -laddition
// #include "addition.h"
import "C"

```

И этого достаточно для использования библиотеки addition.so в Go.

## Использование Go в программах, написанных на C

Для того, чтобы экспортировать функцию, написанную на Go с C-код, необходимо ее пометить для cgo специальным комментарием __export__.
Например, мы хотим, чтобы наша функцию складывала 2 числа и выводила результат в stdout.

```golang

//export Add2Numbers
func Add2Numbers(a C.int, b C.int) {
    fmt.Println(a + b)
}

```

При компиляции cgo с помощью этого комментария формирует Header-файл ___cgo_export.h__, с описанием функций, помеченных директивой __export__

Создадим C-файл со следующим содержанием - создаем поток, в котором будем складывать 2 числа:

```c
// addition.c
#include "_cgo_export.h"
#include <pthread.h>

void *myThreadFun(void *vargp){
    int a, b;
    for (size_t i = 0; i < 5; i++){
        Add2Numbers(a++, ++b);
    }
}

void adds() {
    pthread_t thread_id;
    printf("Before Thread\n");
    pthread_create(&thread_id, NULL, myThreadFun, NULL);
    pthread_join(thread_id, NULL);
    printf("After Thread\n");
    return;
}

```

Внесем правки в файл Go, теперь он у нас должен выглядить следующим образом:

```golang

/*
// C - код
#cgo LDFLAGS: -lpthread
#include <stdio.h>
extern void adds();
*/
import "C"

// Go-код
import "fmt"

//export Add2Numbers
func Add2Numbers(a C.int, b C.int) {
    fmt.Println(a + b)
}

// Add 2 numbers
func Add(a, b int) {
    C.adds()
}

```

В Go-коде возможно использовать и многострочные комментарии для cgo с помощью __/* ... */__.
