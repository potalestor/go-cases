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