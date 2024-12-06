# Service Specification
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
