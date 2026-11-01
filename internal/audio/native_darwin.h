#ifndef NATIVE_DARWIN_H
#define NATIVE_DARWIN_H

int applesoda_tap_start(void);
void applesoda_tap_stop(void);
int applesoda_tap_is_running(void);
float applesoda_tap_get_level(void);
int applesoda_tap_read_samples(float* dst, int max_samples);

#endif
