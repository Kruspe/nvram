// Prototypes
unsigned int Setup();
char *Get(unsigned int gOptionsRef, char *key, char **err);
void Set(unsigned int gOptionsRef, char *key, char *value, char **err);
void Teardown(unsigned int gOptionsRef);
void Delete(unsigned int gOptionsRef, char *name, char **error);

#include <stdio.h>
#include <IOKit/IOKitLib.h>
#include <IOKit/IOKitKeys.h>
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

void Teardown(unsigned int gOptionsRef) {
    IOObjectRelease(gOptionsRef);
}

char *Get(unsigned int gOptionsRef, char *key, char **err) {
    CFStringRef keyRef = CFStringCreateWithCString(kCFAllocatorDefault, key, kCFStringEncodingUTF8);
    if (keyRef == 0) {
        asprintf(err, "Error creating CFString for key %s", key);
        return "";
    }

    CFTypeRef valueRef = IORegistryEntryCreateCFProperty(gOptionsRef, keyRef, 0, 0);
    if (valueRef == 0) {
        asprintf(err, "key '%s' is not set", key);
        return "";
    }
    return (char *) CFDataGetBytePtr(valueRef);
}

void Set(unsigned int gOptionsRef, char *key, char *value, char **err) {
    CFStringRef keyRef = CFStringCreateWithCString(kCFAllocatorDefault, key, kCFStringEncodingUTF8);
    if (keyRef == 0) {
        asprintf(err, "Error creating CFString for key %s", key);
        return;
    }
    CFDataRef valueRef = CFDataCreateWithBytesNoCopy(kCFAllocatorDefault, (const UInt8 *) value, strlen(value),
                                                     kCFAllocatorNull);
    kern_return_t result = IORegistryEntrySetCFProperty(gOptionsRef, keyRef, valueRef);
    if (result != KERN_SUCCESS) {
        asprintf(err, "Error could not write value: %s", mach_error_string(result));
        return;
    }
}

void Delete(unsigned int gOptionsRef, char *key, char **err) {
    CFStringRef deleteKeyRef = CFStringCreateWithCString(kCFAllocatorDefault, kIONVRAMDeletePropertyKey,
                                                         kCFStringEncodingUTF8);
    if (deleteKeyRef == 0) {
        asprintf(err, "Error creating DeleteKey CFString");
        return;
    }
    CFTypeRef keyRef = CFStringCreateWithCString(kCFAllocatorDefault, key, kCFStringEncodingUTF8);
    if (keyRef == 0) {
        asprintf(err, "Error creating CFString for key %s", key);
        return;
    }
    kern_return_t result = IORegistryEntrySetCFProperty(gOptionsRef, deleteKeyRef, keyRef);
    if (result != KERN_SUCCESS) {
        asprintf(err, "Error during key deletion: %s", mach_error_string(result));
        return;
    }
}

