# Automatic Update Testing

This are the instructions on how to test the automatic update locally. 
Do this every time before publishing new version to se if  automatic update works.
 
#### 1. Serving new releases 

We will use python3 to serve our releases.

* Create a hierarchy like the one below somewhere on your disk:
```bash
├── download
│   └── vX.Y.Z # Replace with your new version
│       ├── lia-sdk-linux.zip # This is your new linux build
│       ├── lia-sdk-macos.zip # This is your new macos build
│       └── lia-sdk-windows.zip # This is your new windows build
└── latest
``` 
* In `latest` file paste `{"tag_name": "vX.Y.Z"}` again replacing `vX.Y.Z` with new version.
* Run the server with `python -m SimpleHTTPServer 5000` from the root of the hierarchy created above.

#### 2. Using local server in Lia-SDK

In the terminal where you will be running your update you need to export the base URL to your releases.

```bash
export RELEASES_BASE_URL="http://127.0.0.1:5000/"
```


#### 3. Run update

```bash
./lia update
```

The current Lia-SDK should now be updated to the new one. 
Test if everything works as expected.