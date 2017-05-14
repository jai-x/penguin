$(document).ready(function() {
  link_form_override();
  upload_form_override();
  playlist_refresh();
  /* Loop
  setInterval(function() {
    update_playlist();
  }, 1500);*/
});

function playlist_refresh() {
  var loc = window.location;
  var ws = new WebSocket("ws://" + loc.host + "/ajax/socket");
  ws.onopen = function (event) {
    console.log("Playlist socket connected");
  }
  ws.onerror = function (error) {
    console.log(error);
  }
  ws.onclose = function (event) {
    console.log("Playlist socket closed");
  }
  ws.onmessage = function (event) {
    console.log("Getting message: " + event.data);
  }
}

// Ajax override of form
function link_form_override() {
  $("#queue").submit(function(event) {
    // Stop the normal button behaviour
    event.preventDefault();
    // Set form button to Submitting...
    $("#queuebutton").val("Submitting...");
    // Get form data
    var formData = {
      "video_link": $("input[name=video_link]").val(),
    }
    // The ajax request
    $.ajax({
      type: "POST",
      url: "/ajax/queue",
      data: formData,
    })
    .done(function(data) {
      // Reset form button and link input
      $("#queuebutton").val("Go");
      $("#queueinput").val("");
      // Notify user from response data
      $("#queue").notify(data.Response, data.Type);
    });
  });
}

function upload_form_override() {
  $("#upload").submit(function(event) {
    // Stop normal button behaviour
    event.preventDefault();
    // Set form button to Submitting...
    $("#filebutton").val("Submitting...");
    $("#filebutton").attr("disabled", true);
    // The ajax request
    $.ajax({
      type: "POST",
      url: "/ajax/upload",
      data: new FormData($("#upload")[0]),
      cache: false,
      processData: false,
      contentType: false,
      // Custom XMLHttpRequest for progress bar
      xhr: function() {
          var myXhr = $.ajaxSettings.xhr();
          if (myXhr.upload) {
              // For handling the progress of the upload
              myXhr.upload.addEventListener('progress', function(e) {
                  if (e.lengthComputable) {
                    var percent = Math.round((e.loaded / e.total) * 100) + "%";
                    $("#meter").css("width", percent)
                  }
              } , false);
          }
          return myXhr;
      }
    })
    // Function on done
    .done(function(data) {
      // Reset form button, file input and progress bar
      $("#filebutton").attr("disabled", false);
      $("#filebutton").val("Go");
      $("#fileinput").val("");
      $("#meter").css("width", "0%");
      // Notify user from response data
      $("#upload").notify(data.Response, data.Type);
    });
  });
}

function update_playlist() {
  $("#main").load("/ajax/playlist");
}
