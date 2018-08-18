# Documentation

## URL => https://publisher-go-ride.herokuapp.com/

## For Driver

* Update Location
   * POST /api/location/:id
   * Request Body
   ```json
   {"x":3, "y":4, "available":true, "token":"5"}
   ```
   * Response OK (200)
   ```json
   {"status": "user 3 location updated"}
   ```
   * Response Error (InternalServerError, BadRequest)
   ```json
   {"error_message": "message error"}
   ```
* Get Location
   * GET /api/location/:id