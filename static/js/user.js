var domain;

$(document).ready(function() {
	domain = window.location.href.substring(0, window.location.href.lastIndexOf("/"));
	// Form overrides to callback functions
    $("#queue").submit(ajax_queue);
    $("#upload").submit(ajax_upload);
	// Ajax polling to update page
    playlist_refresh();
});

var progress_bar_html = `
	<div class="progress" role="progressbar" id="progress-outer" tabindex="0" aria-valuenow="0" aria-valuemin="0" aria-valuemax="100">
		<div class="progress-meter" id="meter" style="width: 0%"></div>
	</div>
`;

function ajax_queue(event) {
    // Stop the normal button behaviour
    event.preventDefault();
    // Set form button to Submitting...
    $("#queuebutton").val("Submitting...");
    $("#queuebutton").attr("disabled", true);
    // Get form data
    var formData = {
		"video_link": $("input[name=video_link]").val(),
    }
    // The ajax request
    $.ajax({
		type: "POST",
		url: domain + "/ajax/queue",
		data: formData,
    })
    .done(function(data) {
        // Reset form button and link input
    	$("#queuebutton").attr("disabled", false);
        $("#queuebutton").val("Go");
        $("#queueinput").val("");
        // Notify user from response data
        $("#queue").notify(data.Response, data.Type);
    });
}

function ajax_upload(event) {
    // Stop normal button behaviour
    event.preventDefault();
    // Set form button to Submitting...
    $("#filebutton").val("Submitting...");
    $("#filebutton").attr("disabled", true);
	// Add progress bar to DOM
    $("#vid-input").append(progress_bar_html);
    // The ajax request
    $.ajax({
		type: "POST",
		url: domain + "/ajax/upload",
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
                }, false);
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
		// Remove progress bar
		$("#progress-outer").remove();
        // Notify user from response data
        $("#upload").notify(data.Response, data.Type);
    });
}

function playlist_refresh() {
    $("#main").load("/ajax/playlist", function(response, status) {
        // Set timeout dependant on success or error of previous request
        const timeout = status === "success" ? 2000 : 10000;
        window.setTimeout(function() {
            playlist_refresh();
        }, timeout);
    });
}
