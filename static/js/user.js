$(document).ready(function() {
  link_form_override();
  // Loop
  setInterval(function() {
    update_playlist();
  }, 1500);
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

function update_playlist() {
  $("#main").load("/ajax/playlist");
}