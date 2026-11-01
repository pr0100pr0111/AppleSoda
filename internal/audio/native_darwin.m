#import <Foundation/Foundation.h>
#import <CoreAudio/CoreAudio.h>
#import <CoreAudio/AudioHardwareTapping.h>
#import <CoreAudio/CATapDescription.h>
#include <pthread.h>
#include <math.h>
#include "native_darwin.h"

#define TAP_RING_BUFFER_SIZE 16384

static AudioObjectID g_tap_id = 0;
static AudioObjectID g_aggregate_device = 0;
static AudioDeviceIOProcID g_proc_id = NULL;
static CFStringRef g_tap_uid = NULL;

static float g_ring_buffer[TAP_RING_BUFFER_SIZE];
static int g_write_pos = 0;
static int g_sample_count = 0;
static float g_latest_level = 0.0f;
static int g_is_running = 0;

static pthread_mutex_t g_audio_mutex = PTHREAD_MUTEX_INITIALIZER;

static OSStatus audio_io_proc(AudioObjectID inDevice,
                             const AudioTimeStamp* inNow,
                             const AudioBufferList* inInputData,
                             const AudioTimeStamp* inInputTime,
                             AudioBufferList* outOutputData,
                             const AudioTimeStamp* inOutputTime,
                             void* inClientData) {
    if (!inInputData || inInputData->mNumberBuffers == 0) {
        return noErr;
    }

    float* samples = (float*)inInputData->mBuffers[0].mData;
    UInt32 numSamples = inInputData->mBuffers[0].mDataByteSize / sizeof(float);
    if (numSamples == 0 || !samples) {
        return noErr;
    }

    float sumSq = 0.0f;
    pthread_mutex_lock(&g_audio_mutex);
    for (UInt32 i = 0; i < numSamples; i++) {
        float s = samples[i];
        g_ring_buffer[g_write_pos] = s;
        g_write_pos = (g_write_pos + 1) % TAP_RING_BUFFER_SIZE;
        if (g_sample_count < TAP_RING_BUFFER_SIZE) {
            g_sample_count++;
        }
        sumSq += s * s;
    }
    g_latest_level = sqrtf(sumSq / (float)numSamples);
    pthread_mutex_unlock(&g_audio_mutex);

    return noErr;
}

int applesoda_tap_start(void) {
    pthread_mutex_lock(&g_audio_mutex);
    if (g_is_running) {
        pthread_mutex_unlock(&g_audio_mutex);
        return 0;
    }
    pthread_mutex_unlock(&g_audio_mutex);

    @autoreleasepool {
        CATapDescription *desc = [[CATapDescription alloc] initStereoGlobalTapButExcludeProcesses:@[]];
        if (!desc) {
            return -1;
        }
        desc.privateTap = YES;
        desc.muteBehavior = CATapUnmuted;

        AudioObjectID tapID = 0;
        OSStatus status = AudioHardwareCreateProcessTap(desc, &tapID);
        if (status != noErr) {
            return (int)status;
        }

        AudioObjectPropertyAddress addr = {
            .mSelector = kAudioTapPropertyUID,
            .mScope = kAudioObjectPropertyScopeGlobal,
            .mElement = kAudioObjectPropertyElementMain
        };
        CFStringRef tapUID = NULL;
        UInt32 size = sizeof(tapUID);
        status = AudioObjectGetPropertyData(tapID, &addr, 0, NULL, &size, &tapUID);
        if (status != noErr || !tapUID) {
            AudioHardwareDestroyProcessTap(tapID);
            return (int)status;
        }

        NSString *aggregateUID = [[NSUUID UUID] UUIDString];
        NSDictionary *tapDict = @{
            @kAudioSubTapUIDKey: (__bridge NSString*)tapUID
        };
        NSDictionary *aggregateDesc = @{
            @kAudioAggregateDeviceNameKey: @"AppleSodaTapDevice",
            @kAudioAggregateDeviceUIDKey: aggregateUID,
            @kAudioAggregateDeviceIsPrivateKey: @(YES),
            @kAudioAggregateDeviceTapListKey: @[ tapDict ]
        };

        AudioObjectID aggregateDevice = 0;
        status = AudioHardwareCreateAggregateDevice((__bridge CFDictionaryRef)aggregateDesc, &aggregateDevice);
        if (status != noErr) {
            CFRelease(tapUID);
            AudioHardwareDestroyProcessTap(tapID);
            return (int)status;
        }

        AudioDeviceIOProcID procID = NULL;
        status = AudioDeviceCreateIOProcID(aggregateDevice, audio_io_proc, NULL, &procID);
        if (status != noErr) {
            AudioHardwareDestroyAggregateDevice(aggregateDevice);
            CFRelease(tapUID);
            AudioHardwareDestroyProcessTap(tapID);
            return (int)status;
        }

        status = AudioDeviceStart(aggregateDevice, procID);
        if (status != noErr) {
            AudioDeviceDestroyIOProcID(aggregateDevice, procID);
            AudioHardwareDestroyAggregateDevice(aggregateDevice);
            CFRelease(tapUID);
            AudioHardwareDestroyProcessTap(tapID);
            return (int)status;
        }

        pthread_mutex_lock(&g_audio_mutex);
        g_tap_id = tapID;
        g_tap_uid = tapUID;
        g_aggregate_device = aggregateDevice;
        g_proc_id = procID;
        g_write_pos = 0;
        g_sample_count = 0;
        g_latest_level = 0.0f;
        g_is_running = 1;
        pthread_mutex_unlock(&g_audio_mutex);

        return 0;
    }
}

void applesoda_tap_stop(void) {
    pthread_mutex_lock(&g_audio_mutex);
    if (!g_is_running) {
        pthread_mutex_unlock(&g_audio_mutex);
        return;
    }

    AudioObjectID agg = g_aggregate_device;
    AudioDeviceIOProcID proc = g_proc_id;
    AudioObjectID tap = g_tap_id;
    CFStringRef uid = g_tap_uid;

    g_is_running = 0;
    g_aggregate_device = 0;
    g_proc_id = NULL;
    g_tap_id = 0;
    g_tap_uid = NULL;
    g_write_pos = 0;
    g_sample_count = 0;
    g_latest_level = 0.0f;
    pthread_mutex_unlock(&g_audio_mutex);

    if (agg && proc) {
        AudioDeviceStop(agg, proc);
        AudioDeviceDestroyIOProcID(agg, proc);
    }
    if (agg) {
        AudioHardwareDestroyAggregateDevice(agg);
    }
    if (uid) {
        CFRelease(uid);
    }
    if (tap) {
        AudioHardwareDestroyProcessTap(tap);
    }
}

int applesoda_tap_is_running(void) {
    pthread_mutex_lock(&g_audio_mutex);
    int r = g_is_running;
    pthread_mutex_unlock(&g_audio_mutex);
    return r;
}

float applesoda_tap_get_level(void) {
    pthread_mutex_lock(&g_audio_mutex);
    float lvl = g_latest_level;
    pthread_mutex_unlock(&g_audio_mutex);
    return lvl;
}

int applesoda_tap_read_samples(float* dst, int max_samples) {
    if (!dst || max_samples <= 0) {
        return 0;
    }

    pthread_mutex_lock(&g_audio_mutex);
    int n = max_samples;
    if (n > g_sample_count) {
        n = g_sample_count;
    }
    if (n == 0) {
        pthread_mutex_unlock(&g_audio_mutex);
        return 0;
    }

    int start = (g_write_pos - n + TAP_RING_BUFFER_SIZE) % TAP_RING_BUFFER_SIZE;
    for (int i = 0; i < n; i++) {
        int idx = (start + i) % TAP_RING_BUFFER_SIZE;
        dst[i] = g_ring_buffer[idx];
    }
    pthread_mutex_unlock(&g_audio_mutex);

    return n;
}
