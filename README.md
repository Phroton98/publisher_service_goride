# Documentation

### URL => https://publisher-go-ride.herokuapp.com/

## For Driver

* Update Location
    * POST /api/location/:id
    * Request Body
    ```json
    {"x":3, "y":4, "available":true, "token":"5"}
    ```
    * Response OK (200)
    ```json
    {
        "status": "user 3 <id> location updated"
    }
    ```
    * Response Error (InternalServerError, BadRequest)
    ```json
    {
        "error_message": "message error"
    }
    ```
* Get Location
    * GET /api/location/:id
    * Response OK (200)
    ```json
    {
        "id":1,
        "token":"5",
        "x":3,
        "y":4,
        "available":true,
        "timestamp":1534595330
    }
    ```
    * Response Error (BadRequest)
    ```json
    {
        "error_message": "message error"
    }
    ```

## For Customer

* Create Order
    * POST /api/order 
    * Request Body
    ```json
    {"user_id":1, "user_token":"1", "x":3.0, "y":4.0, "origin":"jakarta", "destination":"bandung", "dest_x":2.0, "dest_y":2.0, "price":10000, "go_pay":false}
    ```
    * Response (StatusCreated)
    ```json
    {
        "message": "order created", 
        "order_id": "<id>"
    }
    ```
    * Response Error (BadRequest, dkk)
    ```json
    {
        "error_message": "message error"
    }
    ```
* Get Order
    * GET /api/order/:id
    * Response (OK)
    ```json
    {
        "ID": 3,
        "UserId": 1,
        "UserToken": "1",
        "DriverId": 0,
        "DriverToken": "",
        "Origin": "jakarta",
        "Destination": "bandung",
        "Price": 10000,
        "Status": "cancelled",
        "CreatedAt": "2018-08-18T13:14:30.31866Z",
        "GoPay": false
    }
    ```
    * Response Error (BadRequest, dkk)
    ```json
    {
        "error_message": "message error"
    }
    ```
* Cancel Order
    * DELETE /api/order/:id
    * Request Body
    ```json
    {"user_id":1, "user_token":"1"}
    ```
    * Response (200)
    ```json
    {
        "message": "order cancelled"
    }
    ```
    * Response Error (BadRequest, dkk)
    ```json
    {
        "error_message": "message error"
    }
    ```