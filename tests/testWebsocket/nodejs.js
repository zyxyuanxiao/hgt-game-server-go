var WSMessage;
var wsmessage;
var buffer;
const GameMessage = require('./protos/protobuf/GameMessage_pb.js')
let message = GameMessage.Message()
message.setProtocol(101)
meesage.setCode(1001)
message.setData(null)
let bytes = message.serializeBinary()
// protobuf.load("protos/protobuf/GameMessage.proto", function (err, root) {
//     console.log("123123213213")
//     if (err) throw err;
//     WSMessage = root.lookup("GameMessage.Message");
//     wsmessage = WSMessage.create({ protocol: 1, code: 1, data: null });
//     buffer = WSMessage.encode(wsmessage).finish();
// });
// var websocket = new WebSocket("ws://111.230.34.248:9503/wss");
// var websocket = new WebSocket("wss://api.sunanzhi.com/wss");
var websocket = new WebSocket("ws://localhost:8899/wss?Authorization=42c217ce-3278-47bc-affa-d074516c1692");
websocket.onopen = function () {
    console.log("websocket open");
    document.getElementById("recv").innerHTML = "Connected";
}
websocket.inclose = function () {
    console.log('websocket close');
}
websocket.onmessage = function (e) {
    console.log(e.data);
    document.getElementById("recv").innerHTML = e.data;
}
document.getElementById("sendBtn").onclick = function () {
    var txt = document.getElementById("sendTxt").value;
    websocket.send(bytes);
}