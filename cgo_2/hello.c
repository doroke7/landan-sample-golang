#include <stdio.h>

#include "hello.h"

void hello_from_file(const char* name) {
	printf("Hello, %s! (來自 hello.c)\n", name);
	fflush(stdout);
}
