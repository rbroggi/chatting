$(function() {
  var socket = null;
  var msgBox = $("#chatbox textarea");
  var messages = $("#messages");
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
      console.log("Raw time JSON field" + msg.When)
      var timeRegex = /.*([0-9]{2}:[0-9]{2}:[0-9]{2}).*/g;
      var time = timeRegex.exec(msg.When)[1];
      messages.append(
        $("<li>").append(
          $("<img>").attr("title", msg.Name).css({
            width:50,
            verticalAlign:"middle"
          }).attr("src", msg.AvatarURL),
          $("<span>").text(time + " - "),
          $("<strong>").text(msg.Name + ": "),
          $("<span>").text(msg.Message)
        )
      );
    }
  }
});
