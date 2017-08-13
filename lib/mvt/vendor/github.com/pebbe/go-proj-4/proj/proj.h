#include <proj_api.h>

char *transform(projPJ srcdefn, projPJ dstdefn, long point_count, double *x, double *y, double *z);
char *fwd(projPJ src, double *lng, double *lat);
char *inv(projPJ dst, double *lng, double *lat);
char *get_err(void);
