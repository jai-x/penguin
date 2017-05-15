$(document).ready(function() {
	sse_admin_playlist();
});

function sse_admin_playlist() {
	if (typeof(EventSource) === "undefined") {
		console.log("SSE Not supported on this browser");
		return
	}
	var org = window.location.origin;
	var source = new EventSource(org + "/sse/admin");
	source.onopen = function (e) {
		console.log("SSE Connected");
	};
	source.onmessage = function (e) {
		console.log("Playlist update");
		$("#main").html(e.data);
	};
}
