/*
 * main.c
 *
 *  Created on: Jan 6, 2012
 *      Author: mladen
 */
#include <unistd.h>
#include <stdio.h>
#include <malloc.h>

typedef struct _DUMMY
{
	char szSomething[1000];
}DUMMY, *PDUMMY;

int main(int argc, char** argv)
{
	while(1)
	{
		PDUMMY* p;
		p = malloc(sizeof(DUMMY));
		printf("--Main Tick--\n");
		sleep(1);
	}

	return 0;
}

