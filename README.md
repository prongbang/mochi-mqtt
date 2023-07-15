# mochi-mqtt

## Support 

### MQTT

- mqtt://host:1883

### Websocket

- ws://host:8083

### Stats

- http://host:8080

## Custom Auth Hook with JWT

### Using MQTT.js

- Install

```shell
npm install mqtt --save
```

- Example

```js
const mqtt = require("mqtt");
const jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c";
const client = mqtt.connect({
    host: "host",
    port: 1883, // or 8083
    username: jwt,
});

client.on("connect", function () {
  client.subscribe("presence", function (err) {
    if (!err) {
      client.publish("presence", "Hello mqtt");
    }
  });
});

client.on("message", function (topic, message) {
  // message is Buffer
  console.log(message.toString());
  client.end();
});
```

