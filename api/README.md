# Oracle Api

## Format

### Time and Duration

The time and duration has to be of `ISO8601` format UTC
Examples :

```
time: 2020-05-12T08:00:00Z //UTC
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
  "publicKey":"02d7e8908aa101d0f7d3565fff11629d3b8fe0a7c431ad336e07de062df5053d6a"
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
    "startDate": "2020-01-01T00:00:00Z",
    "frequency": "PT1H",
    "range": "P10DT"
  }
  ```
- GET `/asset/<asset id>/announcement/<time ISO8601>` to get an announcement for an asset at a requested date (generated lazily). The api will return an announcement corresponding to the next publication of the requested date (depending on oracle configuration)
  example :

  ```
  GET /asset/btcusd/announcement/2021-01-14T07:21:00Z
  200  OK
  ```

  ```json
  {
   "announcementSignature":"",
   "oraclePublicKey":"ce4b7ad2b45de01f0897aa716f67b4c2f596e54506431e693f898712fe7e9bf3",
   "oracleEvent":{
      "nonces":[
         "74558fffd4ef133cb923c066bcc5dd56477bede5da9f3793cb882ce38cc7ef34",
         "5027602119090d44a3466347b4435ae3d5e2729ba0d2d62c1e1733bd980bdf20",
         "3589ebd06f32507757c39acf7f902ad92228a9b82d9120878149a00df1c8dbc5",
         "33fc8f4427198c944ea80e25283b63b241bf5a485f442d0119a1766fb4e34180",
         "a33a3e8678998b8f64edaea7059b23cbf1f9eeda1c300235eff70f54289c6aee",
         "a8f91ae839e117212de79f5dc314d810cef29390304700d0661b0ff3e4166a02",
         "3f526a4995314f3d17c374cdd78c5881fa057ca051e01aa4a88d35bd7a6700ab",
         "8e39e8e08e0adae364f5019d2d0b8fe5cab0a272316426ea9b6d3c5191c43b47",
         "705378b3e483d519755034e93927d326ffe45339c42ab3562f93a8a5b9ef4bb0",
         "afc8b8832150bbde74c926aeb6c5e3f64822fa775005588ee6845755e3f979c3",
         "85adf9e307f45828d4d5743f8cab31339887fd95c562261b84e6107aec5a7828",
         "4b3754062f9ab2cfb7448e939dec05a7749497ed57816710cf0121572f05c115",
         "0d461cf8d8f1a6c8abbabcb66a79dba85948e58f193a0f9d77e7ec995ccfd6f7",
         "d46b37ff056fc07faa2bb388980dcb422c555e3f8829f1d3e5df6c869709ec24",
         "75a0f85d2c448125512846d9cdd045f06a5a700f93073c7c2e782e1183314522",
         "6c0d902798f28db6b2e09859945cc6b0f67423cfb45a6804324ed99507928957",
         "da919395fcacc192dd6f6f98972eed2eb5f18cda5b3ca6f64e4c65a1710193f9",
         "669fa16bf4987ce52fdc8125bf3de8634f2c30db3cecd64d4c02a069a77f8196",
         "e569275f82a94df1591fe63529d947039c47a74be8fb24de5e19e2d768892571",
         "f7e88b898541a1daebbebb3357ff29b551915b9965350d5d62ade00771935038"
      ],
      "eventMaturity":"2021-01-14T07:21:00Z",
      "eventDescriptor":{
         "base":2,
         "isSigned":false,
         "unit":"",
         "precision":0
      },
      "eventId":"btcusd1610608860"
   }
}
```

- GET `/asset/<asset id>/attestation/<time ISO8601>` to get an attestation for an asset at a requested date (generated lazily). The api will return an attestation corresponding to the next publication of the requested date (depending on oracle configuration). if the publication date has not happened yet, an Bad Request Error will be returned.
  example :
  ```
  GET /asset/btcusd/attestation/2021-01-14T07:21:00Z
  200  OK
  ```
  ```json
  {
   "eventId":"btcusd1610608860",
   "signatures":[
      "74558fffd4ef133cb923c066bcc5dd56477bede5da9f3793cb882ce38cc7ef34821c9114cfb8f159a934452331c22c2f7a413c938de25f791321fed90334238e",
      "5027602119090d44a3466347b4435ae3d5e2729ba0d2d62c1e1733bd980bdf2017cb6e085ba26bcfb446a8d40ff09fc7e3ca12c78c5127d64376fa250cfc6c60",
      "3589ebd06f32507757c39acf7f902ad92228a9b82d9120878149a00df1c8dbc5fa50085692c65704453ec05fd76532c2071a664724f15f203ee1a963661d4e8e",
      "33fc8f4427198c944ea80e25283b63b241bf5a485f442d0119a1766fb4e341804bc43586968bb74678b951206743c22ed2f125a4e6397958f76ab8369a59cee3",
      "a33a3e8678998b8f64edaea7059b23cbf1f9eeda1c300235eff70f54289c6aeedf95d843ed39baff6bb8d0c9d9a3c891a3ec255c6202d4e81bad79f6dbd18dc7",
      "a8f91ae839e117212de79f5dc314d810cef29390304700d0661b0ff3e4166a02f5ac530d6bc0f2a651ddb3f896239f5382ee8751cfbe8772f84d56f23ebbfbdf",
      "3f526a4995314f3d17c374cdd78c5881fa057ca051e01aa4a88d35bd7a6700ab4fd4f2a167501e9f122fec06c92a64107a2058b817ea62e472f04a1bf381a053",
      "8e39e8e08e0adae364f5019d2d0b8fe5cab0a272316426ea9b6d3c5191c43b47b081c3b89ec3f4a59804aa63b42076bbdbdbad66b7b3a74386659a101958c8ae",
      "705378b3e483d519755034e93927d326ffe45339c42ab3562f93a8a5b9ef4bb0017c39e060404eac75da56bddc8cdb09603629d3f0039f18e30a09596349e209",
      "afc8b8832150bbde74c926aeb6c5e3f64822fa775005588ee6845755e3f979c36c82200f34545edc2001eac6de827e1d33ea9f00e152984b962a6eb294048ab0",
      "85adf9e307f45828d4d5743f8cab31339887fd95c562261b84e6107aec5a782830ba4df1bb57cdbeb3aded71f1f740b4cc40597e1f7ea3ad50ce595733c04dd0",
      "4b3754062f9ab2cfb7448e939dec05a7749497ed57816710cf0121572f05c11534188aff6bb2bbb61d35551d00427b448c914f701fa8f6ce9255793c42cc4592",
      "0d461cf8d8f1a6c8abbabcb66a79dba85948e58f193a0f9d77e7ec995ccfd6f7b18864be1e15e0e0405cb0a59adc77cc9d372d6c63af20938ed6e8fb2946e1e6",
      "d46b37ff056fc07faa2bb388980dcb422c555e3f8829f1d3e5df6c869709ec2418a3e5a9013df911322f14167850d7b44201dd067bf84bf1dd248514a3d2a2af",
      "75a0f85d2c448125512846d9cdd045f06a5a700f93073c7c2e782e118331452274d7f2f0e2f374586684a3783d1cb202d93ef756aad67ab5eaf06e2984924963",
      "6c0d902798f28db6b2e09859945cc6b0f67423cfb45a6804324ed995079289579ba3a081258edc8bb86fa880f60bfc662eb1607fd165e53ed7a9e02e7e6443c7",
      "da919395fcacc192dd6f6f98972eed2eb5f18cda5b3ca6f64e4c65a1710193f98727cb7c12317a204dc83c37c0e5df72440026411819a8eee630461df2ea89a7",
      "669fa16bf4987ce52fdc8125bf3de8634f2c30db3cecd64d4c02a069a77f8196f342aecae1e8263b621fbf3969f2dd6b3fa390c91e2d0b018038715f88862ed6",
      "e569275f82a94df1591fe63529d947039c47a74be8fb24de5e19e2d768892571ccbfd8414f64ad518a8b8148c08a3dea14787e817bebb061ef89e7ce3c80ef9e",
      "f7e88b898541a1daebbebb3357ff29b551915b9965350d5d62ade0077193503862c7f15a557de46e3f546304dc409cf7a1a2cae40efb737fabc240d7d56ba8d1"
   ],
   "values":[
      "0",
      "0",
      "0",
      "0",
      "1",
      "0",
      "0",
      "1",
      "0",
      "1",
      "0",
      "1",
      "0",
      "1",
      "1",
      "0",
      "1",
      "1",
      "1",
      "1"
   ]
}
```
