#include <stdio.h>
#include <unistd.h>
#include <pthread.h>

volatile int k = 0;

void* thread1(void* args){
    for(unsigned int i = 0; i < 1e3; i++){
        k = k + 1;
    }
}

void* thread2(void* args){
    for(unsigned int i = 0; i < 1e3; i++){
        k = k - 1;
    }
}

int main(){
    pthread_t t1, t2;
    for(int i = 0; i < 10; i++){
        k = 0;
        pthread_create(&t2, NULL, &thread1, NULL);
        pthread_create(&t1, NULL, &thread2, NULL);
        pthread_join(t2, NULL);
        pthread_join(t1, NULL);

        printf("k: %d\n", k);
    }

    return 0;
}
