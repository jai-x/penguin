$(document).ready(function() {
	admin_playlist_refresh();
});

function admin_playlist_refresh() {
  $("#main").load("/ajax/admin/playlist", function(response, status) {
	  // Set timeout dependant on success or error of previous request
	  const timeout = status === "success" ? 2000 : 10000;
	  window.setTimeout(function() {
		  playlist_refresh();
	  }, timeout);
  });
}
