package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void SetAspectRatio(void) {
    dispatch_async(dispatch_get_main_queue(), ^{
        NSArray* windows = [NSApp windows];
        for (NSWindow* window in windows) {
            // Apply to the main app window (usually the one with a title or is visible)
            // We skip small utility windows if any
            if ([window isVisible] && (window.styleMask & NSWindowStyleMaskTitled)) {
                [window setContentAspectRatio:NSMakeSize(1.0, 1.0)];
                // Optional: Force a resize to snap to ratio immediately if needed
                // [window setFrame:[window frame] display:YES];
                NSLog(@"Set aspect ratio for window: %@", [window title]);
            }
        }
    });
}
*/
import "C"
import (
	"time"
)

func SetWindowAspectRatio() {
	// Try immediately
	C.SetAspectRatio()

	// Retry after a short delay to ensure window is fully initialized
	go func() {
		time.Sleep(500 * time.Millisecond)
		C.SetAspectRatio()
	}()
}
