# Bidder

*Bidder service*

## Dependencies:
* [gorilla/mux](https://github.com/gorilla/mux)
* [negroni](https://github.com/urfave/negroni)
* [xid](https://github.com/rs/xid)


### Steps to build
* #### compile
    * Linux  
        ```
        make build-linux
        ```
    * current-platform  
        ```
        make build
        ```    
* #### test
    ```
    make test
    ```
    
***
Output would be a compiled binary with name **bidder** which is powering the docker image.
    