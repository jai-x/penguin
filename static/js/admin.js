var current = new Object();

$(document).ready(function() {
  console.log("Autoupdate page started...");
  setInterval(function() {
    get_list();
    update_np();
    update_pl();

  }, 2000);
});

function update_np() {
  // Verify nowPlaying exists
  if (current && current.nowPlaying) {
    // Check if ID is empty
    if (current.nowPlaying.ID == "") {
      $("#nowplaying").html("Nothing Playing");
    } else {
      $("#nowplaying").html(
        current.nowPlaying.Title + " uploaded by " + escapeHtml(current.nowPlaying.Uploader) +
        `<form action="/admin/kill">` +
          `<input type="submit" value="Stop current video" class="button expanded alert">` +
        `</form>`
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

function format_video(vid) {
  var out = `<div class="row">`;
  // Title
  out += `<div class="small-7 columns">`;
  out += vid.Title;
  out += "</div>";
  // Remove button as admin
  out += `<div class="small-2 columns">`;
  out += `<form action="/admin/remove" method="POST">`;
  out += `<input type="hidden" name="video_id" value="" + vid.ID + "">`;
  out += `<input type="submit" value="Remove" class="alert button small">`;
  out += `</form>`;
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