#define _XOPEN_SOURCE 600
#include <stdio.h>
#include <pthread.h>
#include <time.h>
#include <stdlib.h>
#include <stdbool.h>
#include <unistd.h>

#define N_THREADS 3

// Make 3 mutexes, and three variables that indicate failure or not. Barrier for atomic action sync
bool failure[N_THREADS];
pthread_mutex_t failure_mutex[N_THREADS];
pthread_mutex_t sync_mutex;
pthread_barrier_t barrier;

//bool sync(void);

void* thread(void *arg){
	int id = (int)arg;
	int safe_state = 0;
	int new_value = 0;
	int temp = 0;
	bool calc_fail = false;
	clock_t start = 0, end = 0;
	clock_t wait_time = 4*CLOCKS_PER_SEC; // One second, ish
	for (;;) {
		// sync(); // Start point
		pthread_barrier_wait(&barrier);
		start = clock(); 
		safe_state = new_value; // Store old state (recovery)
		temp = new_value + 1; // New calculation
		pthread_mutex_lock(&failure_mutex[id]);
		failure[id] = (rand() % 10) == 0;
		pthread_mutex_unlock(&failure_mutex[id]);
		// sync atomic action (end point)
		pthread_barrier_wait(&barrier);
		pthread_mutex_lock(&sync_mutex);
		bool fail = failure[0] || failure[1] || failure[2];
		pthread_mutex_unlock(&sync_mutex);
		if(fail) {
			//someone failed
			new_value = safe_state;
		} else {
			// everyone OK
			new_value = temp;	
		}
		if(id == 0) printf("----------------------------\n");
		sleep(1);
		printf("Thread id: %d, value: %d, failure: %d\n", id, new_value, failure[id]);
	}
}
/*
bool sync(void) {
	// read all failure variables
	static int count = 0;
	static bool release = false;
	pthread_mutex_lock(&sync_mutex);
	count = count + 1;
	release = false;
	pthread_mutex_unlock(&sync_mutex);
	while(1){ //wait for all threads to sync
		pthread_mutex_lock(&sync_mutex);
		if((count == N_THREADS) || (release)) {
			//compare and release
			count = 0;
			release = true;
			pthread_mutex_unlock(&sync_mutex);
			return failure[0] || failure[1] || failure[2];		
		}
		pthread_mutex_unlock(&sync_mutex);
	}
}
*/
int main(void) {
	pthread_t thr[N_THREADS];
	srand(time(NULL));

	if(pthread_barrier_init(&barrier, NULL, N_THREADS))
	{
		printf("Could not create a barrier\n");
		return -1;
	}

	for(int i = 0; i < N_THREADS; i++) {
		if(pthread_create(&thr[i], NULL, &thread, (void*)i)) {
			printf("Error in pthread create %d\n", i);
			return 1;
		}
	}
	
	for(int i = 0; i < N_THREADS; i++) {
		if(pthread_join(thr[i], NULL)) {
			printf("Error joining %d\n", i);
		}
	}

	return 0;
}
