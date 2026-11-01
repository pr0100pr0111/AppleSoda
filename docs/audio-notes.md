# CoreAudio Process Tap & Audio Capture Notes

## 1. Native Capture Overview (macOS 14.2+)

AppleSoda captures system audio directly using CoreAudio's **Process Tap** API introduced in macOS 14.2 (`CoreAudio/AudioHardwareTapping.h` and `CoreAudio/CATapDescription.h`).

### Key Components

1. **`CATapDescription`:**
   - Initialized via `[[CATapDescription alloc] initStereoGlobalTapButExcludeProcesses:@[]]`.
   - Global stereo tap captures all audio output streams mixed down to stereo.
   - Configured with `privateTap = YES` (visible only to AppleSoda) and `muteBehavior = CATapUnmuted` (normal hardware playback continues uninterrupted).

2. **`AudioHardwareCreateProcessTap`:**
   - Creates a new tap object (`AudioObjectID`).
   - Retrieves the unique tap UID via `kAudioTapPropertyUID` (`'tuid'`).

3. **Private Aggregate Device:**
   - A private aggregate device (`AudioHardwareCreateAggregateDevice`) is dynamically created containing the tap under the `kAudioAggregateDeviceTapListKey` (`"taps"`).
   - This routes the tapped PCM stream into standard CoreAudio input stream architecture.

4. **IOProc Callback (`AudioDeviceCreateIOProcID`):**
   - High-priority realtime audio callback (`audio_io_proc`) delivers buffers of `float32` PCM audio samples.
   - Samples are fed into an internal lock-protected circular buffer (16,384 samples) without memory allocations in the audio thread.
   - Root Mean Square (RMS) signal level is calculated for fast peak/level monitoring.

5. **Clean Teardown:**
   - On shutdown, `AudioDeviceStop`, `AudioDeviceDestroyIOProcID`, `AudioHardwareDestroyAggregateDevice`, and `AudioHardwareDestroyProcessTap` release all CoreAudio resources.

## 2. Permissions (TCC) & System Audio Recording

Capturing system audio on macOS requires explicit permission:
- **macOS 14.2 - macOS 14.x:** `System Settings → Privacy & Security → System Audio Recording`
- **macOS 15.0+ & macOS 26+:** `System Settings → Privacy & Security → Screen & System Audio Recording`

When AppleSoda or a terminal emulator (Terminal, iTerm2, Alacritty, Ghostty) invokes `AudioHardwareCreateProcessTap` for the first time:
- If permitted: tap initialization succeeds (`noErr = 0`).
- If permission is denied or pending: CoreAudio returns an error code (such as `-1743` or `1852797029`).

### Graceful Fallback
If native tap fails for any reason (permission denied, running in headless SSH, non-macOS, or compiled with `CGO_ENABLED=0`), AppleSoda does not crash:
- It automatically falls back to `simulated` mode.
- The status bar displays the active mode (`[simulated]`).

## 3. Testing the Spike

A standalone CLI spike is provided to test audio capture and signal level monitoring:

```bash
# Test native tap
go run ./cmd/tapspike native

# Test simulated source
go run ./cmd/tapspike simulated
```
