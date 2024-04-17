#define BOOL int
#define TRUE 1
#define FALSE 0
typedef struct{
    char *errMsg;
    char *data;
    BOOL ok;
} RSResult;