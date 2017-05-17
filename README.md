Secure Sun{Rise,Set}
===================

This tool aims to provide a watchable camera streams located on Earth at:

 - sunrise
 - sunset

Inner working
-------------------
Example for: 
**T** = 60
**B** = 5
Silling sun in buffer:
1. Select a camera using the appropriate method
2. Record and encode for **T** seconds to make a video segment[^encodersegment]
3. Add recorded segment to a buffer[^ringbuffer] of **T** x **B** size
4. Go to 1 until the buffer is full
5. Wait for buffer size to be less 
6. Go to 1
[^ringbuffer]: [Circular buffer](https://en.wikipedia.org/wiki/Circular_buffer)

[^encodersegment]: Lossless quality, muxer: Transport Stream, encoder: todo fixed settings ?



Seeding sun to RTMP server:
1. Pop segment from the buffer
2. Feed video segment to streamer encoder realtime[^encoderstreamer]
3. (optionnal) Add shiny effect to transition
4. Go to 1 until the sun consume itself

[^encoderstreamer]: todo check if recorder stream can be feed to encoder without problem 

IDEA
-------
Camera selection proposal:


TODO
--------

 * Choose how to select a cam
 * Module to record the cam for T seconds
 * Module stream the cam
 * Module to ingest / broadcast streams
 * HTML page example to watch the sunrise or sunset !

Useful links: 
 - https://www.insecam.org
