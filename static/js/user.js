var domain;

hide_link_options();

$(document).ready(function() {
	domain = window.location.href;
	if (domain[domain.length - 1] === "/") {
		// Remove trailing slash
		domain = domain.slice(0, -1);
	}

	// Show link options on focus
	$("#queueinput").focus(show_link_options);
	$("#close-queue-options").click(hide_link_options);

	// Form overrides to callback functions
    $("#queue-form").submit(ajax_queue);
    $("#upload-form").submit(ajax_upload);
	// Ajax polling to update page
    playlist_refresh();
});

function hide_link_options() {
	$("#queue-options").hide();
	$("#link-group").css("margin-bottom", 0);
	$("#queue-form").css("border-style", "none");
}

function show_link_options() {
	// Animation is pure swag
	$("#queue-options").show(100);
	// Reset to initial values in css file
	$("#link-group").css("margin-bottom", "");
	$("#queue-form").css("border-style", "");
}

function ajax_queue(event) {
	// Hide the extra options
	hide_link_options();
    // Stop the normal button behaviour
    event.preventDefault();
    // Set form button to Submitting...
    $("#queuebutton").val("Submitting...");
    $("#queuebutton").attr("disabled", true);
    // Get form data
    var formData = {
		"video_link": $("input[name=video_link]").val(),
		"download_subs": $("input[name=download_subs]").val(),
		"vid_offset": $("input[name=vid_offset]").val(),
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
		$("#sub-checkbox").prop("checked", false);
        $("#queuebutton").val("Go");
        $("#queueinput").val("");
        $("#vid-offset").val("");
        // Notify user from response data
        $("#queue-form").notify(data.Response, data.Type);
    });
}

function ajax_upload(event) {
    // Stop normal button behaviour
    event.preventDefault();
    // Set form button to Submitting...
    $("#filebutton").val("Submitting...");
    $("#filebutton").attr("disabled", true);
	// Show progress bar
    $("#progress-outer").show(100);
    // The ajax request
    $.ajax({
		type: "POST",
		url: domain + "/ajax/upload",
		data: new FormData($("#upload-form")[0]),
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
		// Remove and reset progress bar
		$("#progress-outer").css("display", "");
		$("#meter").css("width", "0%");
        // Notify user from response data
        $("#upload-form").notify(data.Response, data.Type);
    });
}

function playlist_refresh() {
    $("#main").load(domain + "/ajax/playlist", function(response, status) {
        // Set timeout dependant on success or error of previous request
        const timeout = status === "success" ? 2000 : 15000;
        window.setTimeout(function() {
            playlist_refresh();
        }, timeout);
    });
}
