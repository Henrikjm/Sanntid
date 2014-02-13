#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <pthread.h>

typedef void* ThreadParam;
typedef void* (*Thread)(void*);

void *print_number(int* num) {
    while (1) {
        printf("Tallet er: %d\n", *num);
        sleep(1);
    }
	return NULL;
}

void start_thread(int tall) {
	int* a = malloc(sizeof(tall));
	*a = tall;
	printf("%d",*a);
	printf("starting to print");
	pthread_t t;
//	void* value_ptr;
	pthread_create(&t, NULL, (Thread)print_number, a);
//	pthread_join(t, &value_ptr);
//	pthread_exit(value_ptr);
	printf("done");
}


int main() {
	start_thread(123);
	start_thread(321);
	pthread_exit(NULL);
}

