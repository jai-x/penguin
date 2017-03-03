var current = new Object();

$(document).ready(function() {
  link_form_override();
  // Loop
  setInterval(function() {
    get_list();
    update_dl();
    update_np();
    update_pl();

  }, 1000);
});

// Ajax override of form
function link_form_override() {
  $("#queue").submit(function(event) {
    var formData = {
      "video_link": $("input[name=video_link]").val(),
      "ajax": true
    }
    // Set form button to Submitting...
    $("#queuebutton").val("Submitting...");
    // The ajax request
    $.ajax({
      type: "POST",
      url: "/ajax/queue",
      data: formData,
    })
    // When done change button back and clear form
    .done(function(data) {
      $("#queuebutton").val("Go");
      $("#queueinput").val("");
      // Notify user from response data
      $("#queue").notify(data.Message, data.Type);
    });
    // Stop the normal button behaviour
    event.preventDefault();
  });
}

function update_np() {
  // Verify nowPlaying exists
  if (current && current.nowPlaying) {
    // Check if ID is empty
    if (current.nowPlaying.ID == "") {
      $("#nowplaying").text("Nothing Playing");
    } else {
      // .text() will auto-escape html
      $("#nowplaying").text(
        current.nowPlaying.Title + " uploaded by " + current.nowPlaying.Uploader
      );
    }
  }
}

function update_pl() {
  // verify playlist exists
  if(current && current.playlist) {
    for (var b = 0; b < current.playlist.length; b++) {
      // Verify current bucket is not empty
      if (current.playlist[b]) {
        var bucketHTML = "";
        for (var i = 0; i < current.playlist[b].length; i++) {
          bucketHTML += format_video(current.playlist[b][i]);
        }
        var bucketDOMID = "#bucket" + b;
        $(bucketDOMID).html(bucketHTML);
      } else {
        // If bucket is empty, ensure it displays as empty
        var bucketDOMID = "#bucket" + b;
        $(bucketDOMID).html("");
      }
    }
  }
}

function update_dl() {
  var downloadingHTML = "";
  // Verify downloading exists
  if (current.downloading && current.downloading.length > 0) {
    downloadingHTML += "<h5>Downloading...</h5>";
    for (var i = 0; i < current.downloading.length; i++) {
      downloadingHTML += "<p>" + current.downloading[i] + "</p>";
    }
  }
  $("#downloading").html(downloadingHTML);
}

function format_video(vid) {
  var out = `<div class="row">`;
    // Title
    out += `<div class="small-7 columns">`;
    out += vid.Title;
    out += "</div>";
    // Remove button 
    out += `<div class="small-2 columns">`;
    // Button if user is uploader
    if (vid.Uploader == current.userAlias) {
      out += `<form action="/remove" method="POST">`;
      out += `<input type="hidden" name="video_id" value="` + vid.ID + `">`;
      out += `<input type="submit" value="Remove" class="alert button tiny remove-button">`;
      out += `</form>`;
    }
    out += `</div>`
    // Alias of uploader
    out += `<div class="small-3 columns">` + escapeHtml(vid.Uploader) + `</div>`;
  out += `</div>`;
  return out;
}

function get_list() {
  var res = $.getJSON("/playlist", function(result) {
    current.nowPlaying = result.NowPlaying;
    current.playlist = result.Playlist;
    current.userAlias = result.UserAlias;
    current.downloading = result.Downloading;
  });

  res.fail(function() {
    console.log("Cannot auto update: Lost connection to server");
  });
}

function escapeHtml(string) {
  var entityMap = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#39;',
    '/': '&#x2F;',
    '`': '&#x60;',
    '=': '&#x3D;'
  };
  return String(string).replace(/[&<>"'`=/]/g, function (s) {
    return entityMap[s];
  });
}