$(function() {
  var socket = null;
  var msgBox = $("#chatbox textarea");
  var messages = $("#messages");
  var timeRegex = /.*([0-9]{2}:[0-9]{2}:[0-9]{2}).*/g;
  $("#chatbox").submit( function() {
    if (!msgBox.val()) return false;
    if (!socket) {
      alert("Error: There is no socket connection.");
      return false;
    }
    socket.send(JSON.stringify({"Message": msgBox.val()}));
    console.log("Message sent to socket.");
    msgBox.val("");
    return false;
  });

  if (!window["WebSocket"]) {
    alert("Error: Your browser does not support web sockets.")
  } else {
    //API with the server
    var addr = $("#address").val();
    socket = new WebSocket(addr);
    console.log("WebSocket connection performed");
    socket.onclose = function() {
      alert("Connection has been closed.");
    }
    socket.onmessage = function(e) {
      var msg = JSON.parse(e.data);
      var time = timeRegex.exec(msg.When)[1];
      messages.append(
        $("<li>").append(
          $("<span>").text(time + " - "),
          $("<strong>").text(msg.Name + ": "),
          $("<span>").text(msg.Message)
        )
      );
    }
  }
});
