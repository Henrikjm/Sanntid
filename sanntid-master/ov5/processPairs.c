#include <errno.h>
#include <sys/socket.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <netinet/in.h>
#include <strings.h>
#include <fcntl.h>
#include <unistd.h>

#define PORT 9930
#define BUFLEN 512

void err(char* str) {
	perror(str);
	exit(1);
}

typedef enum {
	E_BECAME_PRI,
	E_BECAME_SEC,
	E_PRIMARY_FAIL,
	E_FAIL
} event_t;

typedef enum {
	S_INVALID,
	S_PRIMARY,
	S_SECONDARY
} state_t;

typedef struct {
	int sock_fd;
	struct sockaddr_in my_addr;
} network_state_t;

event_t init_network(network_state_t* net_state){
	if(net_state == NULL) 
		err("pointer error, init_network");
	
	if ((net_state->sock_fd = socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP))==-1) {
		err("socket");
	}
	else printf("socket() successful\n");

	bzero(&net_state->my_addr, sizeof(net_state->my_addr));
	net_state->my_addr.sin_family = AF_INET;
	net_state->my_addr.sin_port = htons(PORT);
	net_state->my_addr.sin_addr.s_addr = htonl(INADDR_ANY);
	
	int dummy = 1;
	setsockopt(net_state->sock_fd, SOL_SOCKET, SO_REUSEADDR, &dummy, sizeof(int));//Allow binding to a UDP port already in use
	if (bind(net_state->sock_fd, (struct sockaddr*) &net_state->my_addr, sizeof(net_state->my_addr))==-1){
		if(errno == EADDRINUSE){
			printf("became secondary %d\n",errno);
			return E_BECAME_SEC;
		}else{
			err("bind, init_network");
		}
	}else{
		printf("Server : bind() successful\n");
	}
	// Wait a random time for a "keep alive-signal", then become the master
	clock_t start, end;
	start = clock();
	clock_t wait_time = CLOCKS_PER_SEC + 2*(rand() % (CLOCKS_PER_SEC));
	printf("Wait time: %f\n",(double)wait_time/CLOCKS_PER_SEC);
	int flags = fcntl(net_state->sock_fd, F_GETFL, 0);
	int read_length;
	fcntl(net_state->sock_fd, F_SETFL, flags | O_NONBLOCK);
	char buf[BUFLEN];
	int slen = sizeof(struct sockaddr_in);
	do {
		end = clock();
		read_length = recvfrom(net_state->sock_fd, buf, BUFLEN, 0, (struct sockaddr*) &net_state->my_addr, &slen);
		if(read_length > 0){
			printf("Recieved %d bytes\n", read_length);
			return E_BECAME_SEC;
		}
//			err("Recieve fail");
	}
	while(wait_time > (end - start));
	return E_BECAME_PRI;
}

event_t count_and_print(int* shared_var, network_state_t* net_state) {
	char buf[BUFLEN];
	clock_t start, end;
	start = clock();
	printf("%i\n",CLOCKS_PER_SEC);
	while(1){
		end = clock();
		//printf("%d\n",end-start);
		if (((end-start)/CLOCKS_PER_SEC) >= 1){
			(*shared_var)++;
			printf("Shared var: %i\n",*shared_var);
			start = end;
		}
		int slen = sizeof(struct sockaddr_in);
		sprintf(buf,"%d",*shared_var);
		if (sendto(net_state->sock_fd, buf, BUFLEN, 0, (struct sockaddr*) &net_state->my_addr, slen)==-1)
			err("sendto()");
	}
	return E_FAIL;
}

event_t listen_to_primary(int* shared_var, network_state_t* net_state) {
	int flags = fcntl(net_state->sock_fd, F_GETFL, 0);
	int read_length;
//	fcntl(net_state->sock_fd, F_SETFL, flags | O_NONBLOCK);
	int start, end = 0;
	int wait_time = CLOCKS_PER_SEC/1e5;
	char buf[BUFLEN];
	while(1) {
		//printf("I is now slave :(\n");
		int slen = sizeof(struct sockaddr_in);
		read_length = recvfrom(net_state->sock_fd, buf, BUFLEN, 0, (struct sockaddr*) &net_state->my_addr, &slen);
		if(read_length > 0){
			*shared_var = atoi(buf);
		}
		else {
			printf("Revieced %d bytes, I is master\n", read_length);
			return E_PRIMARY_FAIL;
		}
		// Wait a while
		start = end = clock();
		while(wait_time > (end - start)) end = clock();
	}
}

int main(void) {
	char buf[BUFLEN];
	struct sockaddr_in my_addr;
	int shared_var;
	state_t current_state = S_INVALID;
	event_t event;
	network_state_t net_state;
	srand(time(NULL));

	while(1) {
		switch(current_state) {
			case S_INVALID:
				event = init_network(&net_state);
				switch(event) {
					case E_BECAME_PRI:
						current_state = S_PRIMARY;
						shared_var = 0;
						break;
					case E_BECAME_SEC:
						current_state = S_SECONDARY;
						break;
					case E_FAIL:
					default:
						err("invalid event, S_INVALID");
						break;
				}
				break;
			case S_PRIMARY:
				event = count_and_print(&shared_var, &net_state);
				switch(event) {
					case E_FAIL:
					default:
						err("invalid event, S_PRIMARY");
						break;
				}
				break;
			case S_SECONDARY:
				event = listen_to_primary(&shared_var, &net_state);
				switch(event) {
					case E_PRIMARY_FAIL:
						current_state = S_PRIMARY;
						break;
					case E_FAIL:
					default:
						err("invalid event, S_SECONDARY");
						break;
				}
				break;
			default:
				err("invalid state");
				break;
		}
	}

	return 0;
}
