# Oracle Api

## Format

### Time and Duration

The time and duration has to be of `ISO8601` format  
Examples :

```
time: 2020-05-12T08:00:00Z
duration: P10DT (= 10 days)
```

## Routes

- GET `/oracle/publickey` to recover the oracle public key as a string  
  example :
  ```
  GET /oracle/publickey
  200  OK
  ```
  ```json
  {
  "public_key":"02d7e8908aa101d0f7d3565fff11629d3b8fe0a7c431ad336e07de062df5053d6a"
  }
  ```
- GET `/asset` will list available assets
  example :
  ```
  GET /asset
  200  OK
  ```
  ```json
  ["btcusd","btcjpy"]
  ```
- GET `/asset/<asset id>/config` will return the asset configuration
  example :
  ```
  GET /asset/btcusd/config
  200  OK
  ```
  ```json
  {
    "frequency": "PT1H",
    "range": "P10DT"
  }
  ```
- GET `/asset/<asset id>/rvalue/<time ISO8601>` to get an rvalue for an asset at a requested date (generated lazily). The api will return an rvalue corresponding to the next publication of the requested date (depending on oracle configuration)  
  example :

  ```
  GET /asset/btcusd/rvalue/2020-05-12T07:20:00.00Z
  200  OK
  ```

  ```json
  {
    "publish_date": "2020-05-12T08:00:00Z",
    "asset": "btcusd",
    "rvalue": "03dbdc72bab02979ca8af0d2d91a887ea245031aab78bc3edc2380e22f5deabe63"
  }
  ```

  the response can include the signature and value if the signature has been generated. In that case, the response will be the same kind as GET Signature api

- GET `/asset/<asset id>/signature/<time ISO8601>` to get a signature for an asset at a requested date (generated lazily). The api will return a signature corresponding to the next publication of the requested date (depending on oracle configuration). if the publication date has not happened yet, an error Bad Request Error will be sent.
  example :
  ```
  GET /asset/btcusd/signature/2020-05-12T07:20:00.00Z
  200  OK
  ```
  ```json
  {
    "publish_date": "2020-05-12T08:00:00Z",
    "asset": "btcusd",
    "rvalue": "03dbdc72bab02979ca8af0d2d91a887ea245031aab78bc3edc2380e22f5deabe63",
    "signature": "d3d54ab1f385739e931a91204c3a0c2f1482e7e6006a378a4aeae96599ebc990",
    "value": "220761900556055069"
  }
  ```
