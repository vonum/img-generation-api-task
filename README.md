# Task
Implement a service that does the following:

1. Accepts image uploads through an HTTP API (in a standard *multipart/form-data* encoding) by POST requests on URL `/upload`.
   - The size of the upload should be limited to `8192 Kilobytes`. The service shouldn’t read more data 
     from the client into memory. If the limit is exceeded the service should respond with `413 Request entity too large`.
   - The images are expected to be in `JPEG` format in the multipart body-part named `image`.
   - The file type should be checked before handing the image over for processing.
   - If it’s not a JPEG image then the client should receive a *`400 Bad Request`* error.
2. Generates the ID for the image in `UUIDv4` format.
3. The client should already receive a `200 Ok` response with the ID for the image as soon as the upload finishes and
   the image is accepted for further asynchronous processing.
   The image-ID is to be returned in the response body in JSON format in the following form: `{"image_id":"<uuid-v4>"}`.
4. Perform image downscaling by calling the provided function `image.Rescale`.
5. Save the resulting image on the local filesystem. The path to the image file should be `<base path>/<uuid>.jpg` (`base path` will be supplied in a environment variable named `BASE_PATH`).
6. Because the image processing is quite memory- and CPU-intensive, no more than N images should be processed in parallel
   (the number N should be configurable at start-up time but default to number of processors).
   If there are already N images being processed when a new request comes in and no slot becomes available within
   *100 ms* – the client should receive a `429 Too Many Requests` HTTP-error and his upload should not be accepted for processing.

The whole task should be implemented using only Go's standard library plus google's uuid-library.
You can use Go's generic type parameters, but we suggest to limit it to places where you can easily argue about the benefits.

This task should not be seen like a coding quiz exercise where just strict requirements have to be met. 
We would like to see your engineering skills. It is not only important for us whether you are able to solve it, but also how you do that.

There are only very few tests provided in the skeleton project. It's probably useful to complete them and add more
testing in general, especially for the asynchronous part.

The service doesn't have to be production ready in every aspect.
You may leave out certain aspects like tracing. 
Please do point out what you would add in order to make this service ready for production.

As a bonus: implement graceful shutdown to make sure every accepted request is processed before terminating the program.

## Legal
- Source of testimage_big.jpg https://www.flickr.com/photos/34745138@N00/4484817402, upscaled
- Source of testimage_small.jpg https://www.flickr.com/photos/91501748@N07/20709325971

The set of files that includes this README file is part of an candidate assessment test from tutti.ch.
The candidate agreed to never share these files or derivations with any other parties outside tutti.ch.
 
