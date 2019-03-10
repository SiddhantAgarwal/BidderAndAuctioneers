# Bidders & Auctioneer

*Basic http web services written in golang running on docker*

### Steps to run
* #### Build binaries
    ```
    cd Auctioner/
    make build-linux
    ```
    ```
    cd Bidder/
    make build-linux
    ```
* #### Build service images
    ```
    docker-compose -f service-build.yml build
    ```
* #### Up the services
    ```
    docker-compose up
    ```
    
### Auctioner Exposed API's
* /adplacement
    ```
    POST /adplacement HTTP/1.1
    Host: <HOST>:<PORT>
    Content-Type: application/json
    cache-control: no-cache
    {
    	"ad_placement_id": <AD_PLACEMENT_ID>
    }------WebKitFormBoundary7MA4YWxkTrZu0gW--
    ```

### Auctioner Exposed API's
* /adrequest
    ```
    POST /adrequest HTTP/1.1
    Host: <HOST>:<PORT>
    Content-Type: application/json
    cache-control: no-cache
    {
    	"ad_placement_id": <AD_PLACEMENT_ID>
    }------WebKitFormBoundary7MA4YWxkTrZu0gW--
    ```
    
***
The docker images are built on **scratch** they don't contain any linux libs to make the image sizes small.
The binary of all services contains liked libs, so the images just contain the executable.

The bidders name according to the service name are there in .env
example: BIDDERS=bidder1,bidder2,bidder3,bidder4

to add a bidder add a bidder service in docker-compose.yml and in .env to register it as a bidder