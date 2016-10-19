#import "foundation.h"

NSColor* color_at_screen()
{
    CGPoint mouseLoc = CGEventGetLocation(CGEventCreate(nil));

    // Grab the display for said mouse location.
    uint32_t count = 0;
    CGDirectDisplayID displayForPoint;

    if (CGGetDisplaysWithPoint(mouseLoc, 1, &displayForPoint, &count) != kCGErrorSuccess)
    {
        NSLog(@"Oops.");
        return 0;
    }

    // Grab the color on said display at said mouse location.
    CGImageRef image = CGDisplayCreateImageForRect(displayForPoint, CGRectMake(mouseLoc.x, mouseLoc.y, 1, 1));
    NSBitmapImageRep* bitmap = [[NSBitmapImageRep alloc] initWithCGImage:image];
    CGImageRelease(image);
    NSColor* color = [bitmap colorAtX:0 y:0];
    [bitmap release];

    return color;
}

float color_red_component(NSColor* color)
{
    return ([color redComponent]);
}

float color_green_component(NSColor* color)
{
    return ([color greenComponent]);
}

float color_blue_component(NSColor* color)
{
    return ([color blueComponent]);
}