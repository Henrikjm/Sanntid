// gcc 4.7.2 +
// gcc -std=gnu99 -Wall -g -o helloworld_c helloworld_c.c -lpthread

#include <pthread.h>
#include <stdio.h>

int i = 0;
pthread_mutex_t mtx;

// Note the return type: void*
void* adder(){
   
    for(int x = 0; x < 999999; x++){
        pthread_mutex_lock(&mtx);
        i++;
        pthread_mutex_unlock(&mtx);
    }

    return NULL;
}

void* adder2(){
    
    for(int x = 0; x < 1000000; x++){
        pthread_mutex_lock(&mtx);
        i--;
        pthread_mutex_unlock(&mtx);

    }
   
    return NULL;
}

int main(){
 
    
    pthread_mutex_init(&mtx, NULL);

    pthread_t adder_thr;
    pthread_t adder_thr2;

    pthread_create(&adder_thr, NULL, adder, NULL);
    pthread_create(&adder_thr2, NULL, adder2, NULL);
    
    pthread_join(adder_thr, NULL);
    pthread_join(adder_thr2, NULL);

    pthread_mutex_destroy(&mtx);


    printf("Done: %i\n", i);
    return 0;
    
}