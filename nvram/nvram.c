// Prototypes
unsigned int Setup();
char *Get(unsigned int gOptionsRef, char *key);
void Set(unsigned int gOptionsRef, char *key, char *value);
//void Teardown(unsigned int gOptionsRef);
//int Delete(char *name, char **error, unsigned int gOptionsRef);

#include <stdio.h>
#include <IOKit/IOKitLib.h>
#include <IOKit/IOKitKeys.h>
#include <CoreFoundation/CoreFoundation.h>
#include <err.h>
#include <mach/mach_error.h>

unsigned int Setup() {
    mach_port_t mainPort;

    kern_return_t result = IOMainPort(bootstrap_port, &mainPort);
    if (result != KERN_SUCCESS) {
        errx(1, "Error getting the IOMainPort: %s", mach_error_string(result));
    }
    io_registry_entry_t gOptionsRef = IORegistryEntryFromPath(mainPort, "IODeviceTree:/options");
    if (gOptionsRef == 0) {
        errx(1, "nvram is not supported on this system");
    }

    return gOptionsRef;
}

char *Get(unsigned int gOptionsRef, char *key) {
    CFStringRef keyRef = CFStringCreateWithCString(kCFAllocatorDefault, key, kCFStringEncodingUTF8);
    if (keyRef == 0) {
        errx(1, "Error creating CFString for key %s", key);
    }

    CFTypeRef valueRef = IORegistryEntryCreateCFProperty(gOptionsRef, keyRef, 0, 0);
    if (valueRef == 0) {
        errx(1, "Error key is not set");
    }
    return (char *) CFDataGetBytePtr(valueRef);
}

void Set(unsigned int gOptionsRef, char *key, char *value) {
    CFStringRef keyRef = CFStringCreateWithCString(kCFAllocatorDefault, key,kCFStringEncodingUTF8);
    if (keyRef == 0) {
        errx(1, "Error creating CFString for key %s", key);
    }
    CFDataRef valueRef = CFDataCreateWithBytesNoCopy(kCFAllocatorDefault, (const UInt8 *)value, strlen(value), kCFAllocatorNull);
    kern_return_t result = IORegistryEntrySetCFProperty(gOptionsRef, keyRef, valueRef);
    if (result != KERN_SUCCESS) {
        errx(1, "Error could not write value: %s", mach_error_string(result));
    }
}