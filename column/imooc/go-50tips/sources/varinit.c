#include <stdio.h>

static int cnt;

void f() {
    int n;
    printf("local n = %d\n", n);

    if (cnt > 5) {
        return;
    }

    cnt++;
    f();
}

int main() {
    f();
    return 0;
}
