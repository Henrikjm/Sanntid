#include<pthread.h>
#include<stdio.h>

volatile int i;
void* thread1(void* args){
	for(int j=0; j<100000; j++){	
		i=i+1;
	}
}

void* thread2(void* args){
	for(int j=0; j<100000; j++){
		i=i-1;
	}
}

int main(){
	for(int j=0; j<10;j++){
		i = 0;
		pthread_t t1;
		pthread_t t2;
		pthread_create(&t1, NULL, &thread1, NULL);
		pthread_create(&t2, NULL, &thread2, NULL);

		pthread_join(t1,NULL);
		pthread_join(t2,NULL);
		printf("%d\n", i);
	}
	return 0;
}
