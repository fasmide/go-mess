go-mess provides a rudimentary API for Festo's MES4 by querying directly into its database. 

It doesn't work, it seems to race against the MES application in the database, resulting in what appears to be locking errors, but this requires further investigation.

## Running
`go-mess` takes in database path as its first argument, but will default to `FestoMES.accdb` when omitted:

```
E:\MES4>go-mess.exe
2020/09/01 09:29:17 opening FestoMES.accdb
2020/09/01 09:29:17 listening on :8000
```

## Using
There are three endpoints to this API, `active`, `previous` and `changes` . 

`/active`

This endpoint returns a JSON representation of all ongoing orders in the database:

```
$ http localhost:8000/active
HTTP/1.1 200 OK
Content-Length: 1123
Content-Type: text/plain; charset=utf-8
Date: Tue, 01 Sep 2020 07:34:00 GMT

[
    {
        "CNo": 0,
        "Enabled": true,
        "End": null,
        "ONo": 5241,
        "PlanedEnd": "2020-08-31T13:15:18+02:00",
        "PlanedStart": "2020-08-31T13:14:31+02:00",
        "Positions": [
            {
                "Carrier": null,
                "End": null,
                "Error": false,
                "MainOPos": 0,
                "ONo": 5241,
                "OPos": 1,
                "OpNo": 200,
                "PNo": 1211,
                "Part": {
                    "BasePallet": 25,
                    "Description": "Product No Fuse",
                    "LotSize": 0,
                    "MrpType": 1,
                    "PNo": 1211,
                    "Picture": "Pictures\\TransferFactory\\TF-Part-no-Fuses-Complete.png",
                    "SafetyStock": 0,
                    "Type": 3,
                    "WPNo": 4
                },
                "PlanedEnd": "2020-08-31T13:15:18+02:00",
                "PlanedStart": "2020-08-31T13:14:31+02:00",
                "Resource": {
                    "Automatic": false,
                    "DefaultBrowser": true,
                    "Description": "application module stacking magazine",
                    "IP": "172.20.3.1",
                    "ParallelProcessing": false,
                    "Picture": "Pictures\\TransferFactory\\ModulMagazinNeu2014.png",
                    "PlcType": 1,
                    "ResourceID": 65,
                    "ResourceName": "AM-MAG-IO",
                    "TopologyType": 1,
                    "WebPage": "http://192.168.0.6:8080/webvisu.htm"
                },
                "ResourceID": 65,
                "Start": null,
                "State": {
                    "Description": " not started yet",
                    "Short": "IDLE",
                    "State": 0
                },
                "StateID": 0,
                "StepNo": 10,
                "WONo": 0,
                "WPNo": 4
            }
        ],
        "Release": "2020-08-31T13:14:58+02:00",
        "Start": null,
        "State": {
            "Description": " not started yet",
            "Short": "IDLE",
            "State": 0
        },
        "StateID": 0
    }
]

```

`/previous` 

Requires the user to pass in ID's of the previous orders in question:

```
$ http "localhost:8000/previous?id=5224&id=5196"
HTTP/1.1 200 OK
Content-Length: 486
Content-Type: text/plain; charset=utf-8
Date: Tue, 01 Sep 2020 07:36:57 GMT

[
    {
        "CNo": 10004,
        "Enabled": true,
        "End": "2020-08-26T11:08:11+02:00",
        "ONo": 5224,
        "PlanedEnd": "2020-08-25T11:35:45+02:00",
        "PlanedStart": "2020-08-25T11:35:32+02:00",
        "Release": "2020-08-25T11:34:45+02:00",
        "Start": "2020-08-26T11:03:50+02:00",
        "State": 100
    },
    {
        "CNo": 10001,
        "Enabled": true,
        "End": "2020-08-25T10:51:05+02:00",
        "ONo": 5196,
        "PlanedEnd": "2020-08-25T10:29:14+02:00",
        "PlanedStart": "2020-08-25T10:27:57+02:00",
        "Release": "2020-08-25T10:28:54+02:00",
        "Start": "2020-08-25T10:28:49+02:00",
        "State": 100
    }
]

```

`/changes`

This endpoint has not been tested and may not work. It is supposed to accept WebSocket connections and emit changes to the active orders, but your mileage may vary.